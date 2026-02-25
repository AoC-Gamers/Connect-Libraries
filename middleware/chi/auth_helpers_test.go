package chi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AoC-Gamers/connect-libraries/middleware/v2/authcontext"
	"github.com/AoC-Gamers/connect-libraries/middleware/v2/authjwt"
	"github.com/golang-jwt/jwt/v5"
)

const cookieTokenValue = "cookie-token"

type fakeResponder struct {
	unauthorizedCalled            bool
	tokenExpiredCalled            bool
	policyVersionMismatchCalled   bool
	insufficientPermissionsCalled bool
}

func (f *fakeResponder) Unauthorized(w http.ResponseWriter, detail string) {
	f.unauthorizedCalled = true
	w.WriteHeader(http.StatusUnauthorized)
}

func (f *fakeResponder) TokenExpired(w http.ResponseWriter) {
	f.tokenExpiredCalled = true
	w.WriteHeader(http.StatusUnauthorized)
}

func (f *fakeResponder) PolicyVersionMismatch(w http.ResponseWriter, tokenVersion, currentVersion int) {
	f.policyVersionMismatchCalled = true
	w.WriteHeader(http.StatusUnauthorized)
}

func (f *fakeResponder) InsufficientPermissions(w http.ResponseWriter, action string) {
	f.insufficientPermissionsCalled = true
	w.WriteHeader(http.StatusForbidden)
}

func TestContextHelpers(t *testing.T) {
	claims := &authjwt.Claims{
		SteamID:          "76561198000000123",
		Role:             "web_staff",
		AllowPermissions: 4,
	}

	ctx := context.WithValue(context.Background(), authcontext.SteamIDKey, claims.SteamID)
	ctx = context.WithValue(ctx, authcontext.RoleKey, claims.Role)
	ctx = context.WithValue(ctx, authcontext.ClaimsKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

	if steamID := GetSteamIDFromContext(req); steamID != claims.SteamID {
		t.Fatalf("expected steam id %s, got %s", claims.SteamID, steamID)
	}
	if role := GetRoleFromContext(req); role != claims.Role {
		t.Fatalf("expected role %s, got %s", claims.Role, role)
	}
	if GetClaimsFromContext(req) == nil {
		t.Fatal("expected claims in context")
	}
	if !HasRole(req, "web_staff") {
		t.Fatal("expected role match")
	}
	if !IsStaff(req) {
		t.Fatal("expected IsStaff true")
	}
	if IsAdmin(req) {
		t.Fatal("expected IsAdmin false")
	}
	if !HasPermission(req, 4) {
		t.Fatal("expected permission 4")
	}

	emptyReq := httptest.NewRequest(http.MethodGet, "/", nil)
	if GetSteamIDFromContext(emptyReq) != "" || GetRoleFromContext(emptyReq) != "" || GetClaimsFromContext(emptyReq) != nil {
		t.Fatal("expected empty values with no auth context")
	}
}

func TestExtractToken(t *testing.T) {
	reqHeader := httptest.NewRequest(http.MethodGet, "/?token=query-token", nil)
	reqHeader.Header.Set("Authorization", "Bearer header-token")
	reqHeader.AddCookie(&http.Cookie{Name: "token", Value: cookieTokenValue})
	if token := extractToken(reqHeader); token != "header-token" {
		t.Fatalf("expected header token, got %s", token)
	}

	reqQuery := httptest.NewRequest(http.MethodGet, "/?token=query-token", nil)
	if token := extractToken(reqQuery); token != "query-token" {
		t.Fatalf("expected query token, got %s", token)
	}

	reqCookie := httptest.NewRequest(http.MethodGet, "/", nil)
	reqCookie.AddCookie(&http.Cookie{Name: "token", Value: cookieTokenValue})
	if token := extractToken(reqCookie); token != cookieTokenValue {
		t.Fatalf("expected cookie token, got %s", token)
	}
}

func TestRequireRoleAndPermission(t *testing.T) {
	claims := &authjwt.Claims{Role: "web_admin", AllowPermissions: 8}
	ctx := context.WithValue(context.Background(), authcontext.ClaimsKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/admin", nil).WithContext(ctx)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	roleWrapped := RequireRole("web_admin")(next)
	rr := httptest.NewRecorder()
	roleWrapped.ServeHTTP(rr, req)
	if !called || rr.Code != http.StatusOK {
		t.Fatalf("expected role middleware to allow request, called=%v status=%d", called, rr.Code)
	}

	called = false
	permWrapped := RequirePermissionBitmask(8)(next)
	rr = httptest.NewRecorder()
	permWrapped.ServeHTTP(rr, req)
	if !called || rr.Code != http.StatusOK {
		t.Fatalf("expected permission middleware to allow request, called=%v status=%d", called, rr.Code)
	}
}

func TestRequireRoleAndPermissionErrors(t *testing.T) {
	responder := &fakeResponder{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	roleMiddleware := RequireRoleWithResponder(responder, "web_owner")(next)
	rr := httptest.NewRecorder()
	roleMiddleware.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if !responder.unauthorizedCalled {
		t.Fatal("expected unauthorized responder call for missing claims")
	}

	responder = &fakeResponder{}
	claims := &authjwt.Claims{Role: "web_user"}
	ctx := context.WithValue(context.Background(), authcontext.ClaimsKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	roleMiddleware = RequireRoleWithResponder(responder, "web_owner")(next)
	rr = httptest.NewRecorder()
	roleMiddleware.ServeHTTP(rr, req)
	if !responder.insufficientPermissionsCalled {
		t.Fatal("expected insufficient permissions call for invalid role")
	}

	responder = &fakeResponder{}
	permMiddleware := RequirePermissionBitmaskWithResponder(16, responder)(next)
	rr = httptest.NewRecorder()
	permMiddleware.ServeHTTP(rr, req)
	if !responder.insufficientPermissionsCalled {
		t.Fatal("expected insufficient permissions for missing bitmask permission")
	}
}

func TestRequireAndOptionalAuth(t *testing.T) {
	secret := "chi-jwt-secret"
	policyVersion := 3
	now := time.Now().Unix()
	token := mustSignChiToken(t, secret, jwt.MapClaims{
		"steamid":           "76561198000000999",
		"role":              "web_owner",
		"policy_version":    float64(policyVersion),
		"allow_permissions": float64(32),
		"deny_permissions":  float64(0),
		"iat":               float64(now),
		"exp":               float64(now + 3600),
	})

	config := authjwt.AuthConfig{SignerMaterial: secret, PolicyVersionGlobal: policyVersion}

	authCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCalled = true
		if GetSteamIDFromContext(r) == "" {
			t.Fatal("expected steam id in context")
		}
		if !IsAdmin(r) {
			t.Fatal("expected owner to be admin")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	RequireAuth(config)(next).ServeHTTP(rr, req)
	if !authCalled || rr.Code != http.StatusOK {
		t.Fatalf("expected authenticated request to pass, called=%v status=%d", authCalled, rr.Code)
	}

	optionalCalled := false
	optionalNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		optionalCalled = true
		w.WriteHeader(http.StatusOK)
	})
	rr = httptest.NewRecorder()
	OptionalAuth(config)(optionalNext).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/optional", nil))
	if !optionalCalled || rr.Code != http.StatusOK {
		t.Fatalf("expected optional auth without token to pass, called=%v status=%d", optionalCalled, rr.Code)
	}
}

func mustSignChiToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("could not sign chi token: %v", err)
	}
	return signed
}
