package settingsruntime

import "testing"

func TestCacheKey(t *testing.T) {
	got := CacheKey(" config ", "rate_limit.enabled ")
	want := "CONFIG:rate_limit.enabled"
	if got != want {
		t.Fatalf("unexpected cache key: got=%q want=%q", got, want)
	}
}

func TestNormalizeSpecKey(t *testing.T) {
	got, err := normalizeSpecKey(KeySpec{Entity: "config", Key: "x.y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "CONFIG:x.y" {
		t.Fatalf("unexpected normalized key: %q", got)
	}
}

func TestNormalizeSpecKey_Invalid(t *testing.T) {
	if _, err := normalizeSpecKey(KeySpec{Entity: "", Key: "x"}); err == nil {
		t.Fatal("expected error for empty entity")
	}
	if _, err := normalizeSpecKey(KeySpec{Entity: "CONFIG", Key: ""}); err == nil {
		t.Fatal("expected error for empty key")
	}
}
