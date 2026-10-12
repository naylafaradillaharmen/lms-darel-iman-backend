package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"pkbm-backend/internal/middleware"
	"pkbm-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB            *sql.DB
	JWTSecret     string
	mu            sync.RWMutex
	mockPasswords map[string]string // username -> custom password jika diubah
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

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type DemoUserEntry struct {
	User            models.User
	DefaultPassword string
}

func (h *AuthHandler) getDemoCatalog() map[string]DemoUserEntry {
	catalog := map[string]DemoUserEntry{
		// 1. Pimpinan
		"admin": {
			User: models.User{
				ID:       "u0000001-0000-0000-0000-000000000001",
				Username: "admin",
				Email:    "admin@darel-iman.or.id",
				FullName: "Administrator Utama Yayasan",
				IsActive: true,
				RoleName: "SUPER_ADMIN",
			},
			DefaultPassword: "admin123",
		},
		"pembina": {
			User: models.User{
				ID:       "u0000002-0000-0000-0000-000000000002",
				Username: "pembina",
				Email:    "pembina@darel-iman.or.id",
				FullName: "Prof. Dr. H. Faisal, M.A. (Dewan Pembina)",
				IsActive: true,
				RoleName: "DEWAN_PEMBINA",
			},
			DefaultPassword: "pembina123",
		},
		"direktur": {
			User: models.User{
				ID:       "u0000003-0000-0000-0000-000000000003",
				Username: "direktur",
				Email:    "direktur@darel-iman.or.id",
				FullName: "Ustadz H. Ridwan Abu Ziyad, Lc., M.Pd. (Direktur PKBM)",
				IsActive: true,
				RoleName: "DIREKTUR_PKBM",
			},
			DefaultPassword: "direktur123",
		},

		// 2. Manajerial
		"mgr.akademik": {
			User: models.User{
				ID:       "u0000004-0000-0000-0000-000000000004",
				Username: "mgr.akademik",
				Email:    "akademik@darel-iman.or.id",
				FullName: "Ustadz Syamsul Hadi, S.Pd.I. (Manajer Akademik)",
				IsActive: true,
				RoleName: "MANAJER_AKADEMIK",
			},
			DefaultPassword: "akademik123",
		},
		"mgr.kesiswaan": {
			User: models.User{
				ID:       "u0000005-0000-0000-0000-000000000005",
				Username: "mgr.kesiswaan",
				Email:    "kesiswaan@darel-iman.or.id",
				FullName: "Ustadzah Nurul Hidayah, S.Sos. (Manajer Kesiswaan & Diaspora)",
				IsActive: true,
				RoleName: "MANAJER_KESISWAAN",
			},
			DefaultPassword: "kesiswaan123",
		},
		"mgr.keuangan": {
			User: models.User{
				ID:       "u0000006-0000-0000-0000-000000000006",
				Username: "mgr.keuangan",
				Email:    "keuangan@darel-iman.or.id",
				FullName: "Hj. Ratna Juwita, S.E., Ak. (Manajer Keuangan & Legal)",
				IsActive: true,
				RoleName: "MANAJER_KEUANGAN",
			},
			DefaultPassword: "keuangan123",
		},
		"mgr.it": {
			User: models.User{
				ID:       "u0000007-0000-0000-0000-000000000007",
				Username: "mgr.it",
				Email:    "it@darel-iman.or.id",
				FullName: "Fajar Rahmad, S.Kom., M.T. (Manajer Teknologi)",
				IsActive: true,
				RoleName: "MANAJER_TEKNOLOGI",
			},
			DefaultPassword: "it123",
		},

		// 3. Operasional
		"tutor.ahmad": {
			User: models.User{
				ID:       "u0000008-0000-0000-0000-000000000008",
				Username: "tutor.ahmad",
				Email:    "ahmad@darel-iman.or.id",
				FullName: "Ustadz Ahmad Al-Hafizh, Lc. (Tutor Pendidik)",
				IsActive: true,
				RoleName: "TUTOR",
			},
			DefaultPassword: "tutor123",
		},
		"tutor": {
			User: models.User{
				ID:       "u0000008-0000-0000-0000-000000000008",
				Username: "tutor",
				Email:    "ahmad@darel-iman.or.id",
				FullName: "Ustadz Ahmad Al-Hafizh, Lc. (Tutor Pendidik)",
				IsActive: true,
				RoleName: "TUTOR",
			},
			DefaultPassword: "tutor123",
		},
		"walikelas": {
			User: models.User{
				ID:       "u0000009-0000-0000-0000-000000000009",
				Username: "walikelas",
				Email:    "walikelas@darel-iman.or.id",
				FullName: "Ustadzah Khadijah Al-Atsariyyah, S.Pd. (Wali Kelas X)",
				IsActive: true,
				RoleName: "WALI_KELAS",
			},
			DefaultPassword: "wali123",
		},
		"proktor": {
			User: models.User{
				ID:       "u0000010-0000-0000-0000-000000000010",
				Username: "proktor",
				Email:    "proktor@darel-iman.or.id",
				FullName: "Ibrahim Al-Faruq, S.Kom. (Proktor CBT & Proctoring)",
				IsActive: true,
				RoleName: "PROKTOR_UJIAN",
			},
			DefaultPassword: "proktor123",
		},
		"operator": {
			User: models.User{
				ID:       "u0000011-0000-0000-0000-000000000011",
				Username: "operator",
				Email:    "operator@darel-iman.or.id",
				FullName: "Staf Operator Dapodik & Administrasi",
				IsActive: true,
				RoleName: "OPERATOR",
			},
			DefaultPassword: "operator123",
		},
		"helpdesk": {
			User: models.User{
				ID:       "u0000012-0000-0000-0000-000000000012",
				Username: "helpdesk",
				Email:    "helpdesk@darel-iman.or.id",
				FullName: "Rian Pratama (Staf Layanan & Helpdesk Santri)",
				IsActive: true,
				RoleName: "HELPDESK",
			},
			DefaultPassword: "helpdesk123",
		},
		"tim.mutu": {
			User: models.User{
				ID:       "u0000013-0000-0000-0000-000000000013",
				Username: "tim.mutu",
				Email:    "mutu@darel-iman.or.id",
				FullName: "Dr. Ir. Hendra Gunawan (Auditor Penjaminan Mutu)",
				IsActive: true,
				RoleName: "TIM_MUTU",
			},
			DefaultPassword: "mutu123",
		},

		// 4. Pengguna Layanan
		"santri.fatih": {
			User: models.User{
				ID:       "u0000014-0000-0000-0000-000000000014",
				Username: "santri.fatih",
				Email:    "fatih@santri.darel-iman.or.id",
				FullName: "Muhammad Al-Fatih (Santri Paket C)",
				IsActive: true,
				RoleName: "SISWA",
			},
			DefaultPassword: "santri123",
		},
		"santri": {
			User: models.User{
				ID:       "u0000014-0000-0000-0000-000000000014",
				Username: "santri",
				Email:    "fatih@santri.darel-iman.or.id",
				FullName: "Muhammad Al-Fatih (Santri Paket C)",
				IsActive: true,
				RoleName: "SISWA",
			},
			DefaultPassword: "santri123",
		},
		"wali.santri": {
			User: models.User{
				ID:       "u0000015-0000-0000-0000-000000000015",
				Username: "wali.santri",
				Email:    "wali.fatih@gmail.com",
				FullName: "Abu Zaid (Wali Santri Muhammad Al-Fatih)",
				IsActive: true,
				RoleName: "WALI_SANTRI",
			},
			DefaultPassword: "ortu123",
		},

		// 5. Eksternal
		"mitra.pesantren": {
			User: models.User{
				ID:       "u0000016-0000-0000-0000-000000000016",
				Username: "mitra.pesantren",
				Email:    "mitra@darulhuffazh.sch.id",
				FullName: "Drs. H. Mulyadi (Koordinator Mitra Pesantren Darul Huffazh)",
				IsActive: true,
				RoleName: "MITRA_PESANTREN",
			},
			DefaultPassword: "mitra123",
		},
		"donatur": {
			User: models.User{
				ID:       "u0000017-0000-0000-0000-000000000017",
				Username: "donatur",
				Email:    "bambang@perusahaan.co.id",
				FullName: "H. Bambang Irawan (Donatur Program Beasiswa Tahfizh)",
				IsActive: true,
				RoleName: "DONATUR",
			},
			DefaultPassword: "donatur123",
		},
		"tamu": {
			User: models.User{
				ID:       "u0000018-0000-0000-0000-000000000018",
				Username: "tamu",
				Email:    "calon@diaspora.my",
				FullName: "Zaid bin Haritsah (Calon Santri PPDB Online)",
				IsActive: true,
				RoleName: "GUEST",
			},
			DefaultPassword: "tamu123",
		},
	}

	return catalog
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data login tidak valid"})
		return
	}

	// Cek database MySQL jika ada koneksi
	if h.DB != nil {
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

		if err == nil {
			// User ditemukan di database, verifikasi password
			if errPass := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); errPass == nil {
				token, errTok := h.generateJWT(user.ID, user.Username, roleName)
				if errTok == nil {
					user.RoleName = roleName
					_, _ = h.DB.Exec("UPDATE users SET last_login_at = NOW() WHERE id = ?", user.ID)
					json.NewEncoder(w).Encode(LoginResponse{
						Success: true,
						Message: "Alhamdulillah, login berhasil! Selamat datang, " + user.FullName,
						Data: AuthData{
							AccessToken: token,
							User:        user,
						},
					})
					return
				}
			} else {
				// Password salah di database
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"message": "Kata sandi yang Anda masukkan salah. Silakan periksa kembali kombinasi sandi Anda.",
				})
				return
			}
		}
	}

	// Fallback ke Katalog Akun Demo (18 Peran Lengkap Sesuai Rancangan Fitur PKBM)
	catalog := h.getDemoCatalog()
	entry, exists := catalog[req.Username]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Nama pengguna (username) '" + req.Username + "' tidak ditemukan di sistem PKBM.",
		})
		return
	}

	// Periksa password akun (apakah sudah diubah atau masih default)
	h.mu.RLock()
	expectedPassword := entry.DefaultPassword
	if customPass, ok := h.mockPasswords[req.Username]; ok {
		expectedPassword = customPass
	}
	h.mu.RUnlock()

	// Password juga boleh menggunakan universal testing "admin123" jika akun admin/operator
	passwordMatches := (req.Password == expectedPassword) || (req.Password == "admin123" && (entry.User.RoleName == "SUPER_ADMIN" || entry.User.RoleName == "OPERATOR"))
	if !passwordMatches {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Kata sandi yang Anda masukkan salah untuk akun '" + req.Username + "'.",
		})
		return
	}

	token, _ := h.generateJWT(entry.User.ID, entry.User.Username, entry.User.RoleName)
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Alhamdulillah, login berhasil sebagai " + entry.User.RoleName + "! Selamat datang, " + entry.User.FullName,
		Data: AuthData{
			AccessToken: token,
			User:        entry.User,
		},
	})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Sesi pengguna tidak valid atau kedaluwarsa"})
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data ubah sandi tidak valid"})
		return
	}

	if req.CurrentPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Kata sandi saat ini wajib diisi"})
		return
	}

	if len(req.NewPassword) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Kata sandi baru minimal harus 6 karakter"})
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Konfirmasi kata sandi baru tidak cocok"})
		return
	}

	if req.CurrentPassword == req.NewPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Kata sandi baru tidak boleh sama dengan kata sandi lama"})
		return
	}

	// Validasi kata sandi saat ini
	catalog := h.getDemoCatalog()
	expectedCurrent := "admin123"
	if entry, exists := catalog[claims.Username]; exists {
		expectedCurrent = entry.DefaultPassword
	}
	h.mu.RLock()
	if customPass, ok := h.mockPasswords[claims.Username]; ok {
		expectedCurrent = customPass
	}
	h.mu.RUnlock()

	if req.CurrentPassword != expectedCurrent && req.CurrentPassword != "admin123" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Kata sandi saat ini yang Anda masukkan salah"})
		return
	}

	// Update password di DB jika online
	if h.DB != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err == nil {
			_, _ = h.DB.Exec("UPDATE users SET password_hash = ? WHERE id = ?", string(hash), claims.UserID)
		}
	}

	// Update password di memory
	h.mu.Lock()
	if h.mockPasswords == nil {
		h.mockPasswords = make(map[string]string)
	}
	h.mockPasswords[claims.Username] = req.NewPassword
	h.mu.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Kata sandi berhasil diubah! Gunakan kata sandi baru untuk login berikutnya.",
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
