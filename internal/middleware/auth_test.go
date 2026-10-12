package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "secret-test-darel-iman-2026"

func generateTestToken(userID, username, role string) string {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	mw := AuthMiddleware(testSecret)
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*Claims)
		if !ok || claims == nil {
			t.Errorf("Konteks user tidak ditemukan")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized"))
	})

	// 1. Tanpa header Authorization -> Harus 401
	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	rec := httptest.NewRecorder()
	mw(dummyHandler).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Ekspektasi 401 Unauthorized, tapi didapat %d", rec.Code)
	}

	// 2. Token tidak valid -> Harus 401
	reqInvalid := httptest.NewRequest("GET", "/api/v1/test", nil)
	reqInvalid.Header.Set("Authorization", "Bearer token-palsu-ngawur")
	recInvalid := httptest.NewRecorder()
	mw(dummyHandler).ServeHTTP(recInvalid, reqInvalid)
	if recInvalid.Code != http.StatusUnauthorized {
		t.Errorf("Ekspektasi 401 Unauthorized untuk token palsu, tapi didapat %d", recInvalid.Code)
	}

	// 3. Token valid -> Harus 200 OK
	validToken := generateTestToken("u001", "santri.fatih", "SISWA")
	reqValid := httptest.NewRequest("GET", "/api/v1/test", nil)
	reqValid.Header.Set("Authorization", "Bearer "+validToken)
	recValid := httptest.NewRecorder()
	mw(dummyHandler).ServeHTTP(recValid, reqValid)
	if recValid.Code != http.StatusOK {
		t.Errorf("Ekspektasi 200 OK untuk token valid, tapi didapat %d", recValid.Code)
	}
}

func TestRequireRole(t *testing.T) {
	adminRoleMw := RequireRole("OPERATOR")
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access Granted"))
	})

	// 1. User SISWA mencoba akses endpoint OPERATOR -> Harus 403 Forbidden
	ctxStudent := context.WithValue(context.Background(), UserContextKey, &Claims{
		UserID:   "s001",
		Username: "santri",
		Role:     "SISWA",
	})
	reqStudent := httptest.NewRequest("GET", "/api/v1/admin/resource", nil).WithContext(ctxStudent)
	recStudent := httptest.NewRecorder()
	adminRoleMw(dummyHandler).ServeHTTP(recStudent, reqStudent)
	if recStudent.Code != http.StatusForbidden {
		t.Errorf("Ekspektasi 403 Forbidden untuk SISWA membuka resource admin, didapat %d", recStudent.Code)
	}

	// 2. User OPERATOR mencoba akses endpoint OPERATOR -> Harus 200 OK
	ctxOperator := context.WithValue(context.Background(), UserContextKey, &Claims{
		UserID:   "op001",
		Username: "operator",
		Role:     "OPERATOR",
	})
	reqOperator := httptest.NewRequest("GET", "/api/v1/admin/resource", nil).WithContext(ctxOperator)
	recOperator := httptest.NewRecorder()
	adminRoleMw(dummyHandler).ServeHTTP(recOperator, reqOperator)
	if recOperator.Code != http.StatusOK {
		t.Errorf("Ekspektasi 200 OK untuk OPERATOR, didapat %d", recOperator.Code)
	}

	// 3. SUPER_ADMIN selalu boleh (bypass) -> Harus 200 OK
	ctxSuper := context.WithValue(context.Background(), UserContextKey, &Claims{
		UserID:   "sa001",
		Username: "superadmin",
		Role:     "SUPER_ADMIN",
	})
	reqSuper := httptest.NewRequest("GET", "/api/v1/admin/resource", nil).WithContext(ctxSuper)
	recSuper := httptest.NewRecorder()
	adminRoleMw(dummyHandler).ServeHTTP(recSuper, reqSuper)
	if recSuper.Code != http.StatusOK {
		t.Errorf("Ekspektasi 200 OK untuk SUPER_ADMIN, didapat %d", recSuper.Code)
	}
}
