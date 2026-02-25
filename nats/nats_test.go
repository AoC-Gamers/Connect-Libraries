package connectnats

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

const customServiceName = "custom-service"

func TestDefaultConfigAndEnv(t *testing.T) {
	t.Setenv("NATS_URL", "")
	t.Setenv("NATS_CLIENT_ID", "")
	defaultConfig := DefaultConfig()
	if defaultConfig.URL != "nats://localhost:4222" {
		t.Fatalf("unexpected default URL: %s", defaultConfig.URL)
	}
	if defaultConfig.ClientID != "connect-service" {
		t.Fatalf("unexpected default client id: %s", defaultConfig.ClientID)
	}

	t.Setenv("NATS_URL", "nats://custom:4222")
	t.Setenv("NATS_CLIENT_ID", customServiceName)
	overridden := DefaultConfig()
	if overridden.URL != "nats://custom:4222" {
		t.Fatalf("expected env URL override, got %s", overridden.URL)
	}
	if overridden.ClientID != customServiceName {
		t.Fatalf("expected env client id override, got %s", overridden.ClientID)
	}

	t.Setenv("ENV_TEST_VALUE", "configured")
	if value := getEnvOrDefault("ENV_TEST_VALUE", "fallback"); value != "configured" {
		t.Fatalf("expected configured value, got %s", value)
	}
	t.Setenv("ENV_TEST_VALUE", "")
	if value := getEnvOrDefault("ENV_TEST_VALUE", "fallback"); value != "fallback" {
		t.Fatalf("expected fallback value, got %s", value)
	}
}

func TestConnectAndEnsureStreamsErrors(t *testing.T) {
	if err := EnsureStreams(nil, nil); err == nil {
		t.Fatal("expected nil connection error in EnsureStreams")
	}

	cfg := &Config{
		URL:           "nats://127.0.0.1:1",
		ClientID:      "connect-test",
		ReconnectWait: 1,
		MaxReconnects: 0,
		Timeout:       50 * time.Millisecond,
	}

	if _, err := Connect(cfg); err == nil {
		t.Fatal("expected connection error for invalid endpoint")
	}
	if conn := MustConnect(cfg); conn != nil {
		t.Fatal("expected nil connection from MustConnect on failure")
	}

	if _, err := ConnectWithJetStream(cfg, nil); err == nil {
		t.Fatal("expected ConnectWithJetStream to fail with invalid endpoint")
	}
	if conn := MustConnectWithJetStream(cfg, nil); conn != nil {
		t.Fatal("expected nil from MustConnectWithJetStream on failure")
	}
}

func TestPublisherErrorsAndDefaults(t *testing.T) {
	logger := zerolog.Nop()
	if _, err := NewPublisher(nil, logger); err == nil {
		t.Fatal("expected error when creating publisher with nil conn")
	}

	publisher := &Publisher{log: logger}
	if err := publisher.PublishCore("subject", map[string]string{"x": "y"}); err == nil {
		t.Fatal("expected error when core connection is nil")
	}
	if err := publisher.PublishJetStream("subject", map[string]string{"x": "y"}); err == nil {
		t.Fatal("expected error when jetstream is nil")
	}
	if err := publisher.Publish("subject", map[string]string{"x": "y"}); err == nil {
		t.Fatal("expected Publish to fail when jetstream is nil")
	}

	custom := &Publisher{serviceName: customServiceName}
	WithServiceName("new-service")(custom)
	if custom.serviceName != "new-service" {
		t.Fatalf("expected service name override, got %s", custom.serviceName)
	}
}

func TestDefaultStreamConfigs(t *testing.T) {
	configs := DefaultStreamConfigs()
	if len(configs) != 2 {
		t.Fatalf("expected 2 default stream configs, got %d", len(configs))
	}
	if configs[0].Name == "" || len(configs[0].Subjects) == 0 {
		t.Fatal("expected first stream to have name and subjects")
	}
	if configs[1].Name == "" || len(configs[1].Subjects) == 0 {
		t.Fatal("expected second stream to have name and subjects")
	}
}

func TestCreateTLSConfig(t *testing.T) {
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")
	caFile := filepath.Join(tempDir, "ca.crt")

	certPEM, keyPEM := generateSelfSignedCertPEM(t)
	if err := os.WriteFile(certFile, certPEM, 0600); err != nil {
		t.Fatalf("failed writing cert file: %v", err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatalf("failed writing key file: %v", err)
	}
	if err := os.WriteFile(caFile, certPEM, 0600); err != nil {
		t.Fatalf("failed writing ca file: %v", err)
	}

	tlsConfig, err := createTLSConfig(certFile, keyFile, caFile)
	if err != nil {
		t.Fatalf("expected valid tls config, got %v", err)
	}
	if tlsConfig == nil || len(tlsConfig.Certificates) != 1 {
		t.Fatal("expected one certificate in tls config")
	}
	if tlsConfig.MinVersion == 0 {
		t.Fatal("expected tls min version to be configured")
	}

	if _, err = createTLSConfig(certFile, keyFile, filepath.Join(tempDir, "missing-ca.crt")); err == nil {
		t.Fatal("expected error when CA file is missing")
	}
}

func generateSelfSignedCertPEM(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyDER})

	return certPEM, keyPEM
}
