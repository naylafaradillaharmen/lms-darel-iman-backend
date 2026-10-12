package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"pkbm-backend/internal/middleware"
	"pkbm-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB        *sql.DB
	JWTSecret string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    AuthData     `json:"data"`
}

type AuthData struct {
	AccessToken string      `json:"access_token"`
	User        models.User `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Body permintaan tidak valid"})
		return
	}

	// Ambil data user beserta role
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.phone, u.avatar_url, u.is_active,
		       COALESCE(r.name, 'SISWA') as role_name
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE (u.username = ? OR u.email = ?) AND u.is_active = TRUE
		LIMIT 1
	`

	var user models.User
	var roleName string
	err := h.DB.QueryRow(query, req.Username, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.Phone, &user.AvatarURL, &user.IsActive, &roleName,
	)

	if err != nil {
		// Mock login fallback jika DB belum siap untuk memudahkan testing awal
		var mockUser *models.User
		if req.Username == "admin" && req.Password == "admin123" {
			mockUser = &models.User{
				ID:       "u0000001-0000-0000-0000-000000000001",
				Username: "admin",
				Email:    "admin@darel-iman.or.id",
				FullName: "Administrator Utama Yayasan",
				IsActive: true,
				RoleName: "SUPER_ADMIN",
			}
		} else if req.Username == "operator" && req.Password == "admin123" {
			mockUser = &models.User{
				ID:       "u0000002-0000-0000-0000-000000000002",
				Username: "operator",
				Email:    "operator@darel-iman.or.id",
				FullName: "Staf Operator Dapodik",
				IsActive: true,
				RoleName: "OPERATOR",
			}
		} else if (req.Username == "tutor.ahmad" || req.Username == "tutor") && (req.Password == "tutor123" || req.Password == "admin123") {
			mockUser = &models.User{
				ID:       "u0000003-0000-0000-0000-000000000003",
				Username: "tutor.ahmad",
				Email:    "ahmad@darel-iman.or.id",
				FullName: "Ustadz Ahmad Al-Hafizh, Lc.",
				IsActive: true,
				RoleName: "TUTOR",
			}
		} else if (req.Username == "santri.fatih" || req.Username == "santri") && (req.Password == "santri123" || req.Password == "admin123") {
			mockUser = &models.User{
				ID:       "u0000004-0000-0000-0000-000000000004",
				Username: "santri.fatih",
				Email:    "fatih@santri.darel-iman.or.id",
				FullName: "Muhammad Al-Fatih",
				IsActive: true,
				RoleName: "SISWA",
			}
		}

		if mockUser != nil {
			token, _ := h.generateJWT(mockUser.ID, mockUser.Username, mockUser.RoleName)
			json.NewEncoder(w).Encode(LoginResponse{
				Success: true,
				Message: "Login Berhasil (Mock Mode)",
				Data:    AuthData{AccessToken: token, User: *mockUser},
			})
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Username atau kata sandi salah"})
		return
	}

	user.RoleName = roleName

	// Verifikasi kata sandi
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Username atau kata sandi salah"})
		return
	}

	token, err := h.generateJWT(user.ID, user.Username, roleName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Gagal membuat sesi token"})
		return
	}

	// Update last_login_at
	_, _ = h.DB.Exec("UPDATE users SET last_login_at = NOW() WHERE id = ?", user.ID)

	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Login Berhasil",
		Data: AuthData{
			AccessToken: token,
			User:        user,
		},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Sesi tidak ditemukan"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		},
	})
}

func (h *AuthHandler) generateJWT(userID, username, role string) (string, error) {
	claims := middleware.Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pkbm-darel-iman",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.JWTSecret))
}
