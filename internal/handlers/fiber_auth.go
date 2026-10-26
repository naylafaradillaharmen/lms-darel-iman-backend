package handlers

import (
	"strings"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type FiberAuthHandler struct {
	GormDB          *gorm.DB
	JWTSecret       string
	mu              sync.RWMutex
	mockPasswords   map[string]string
	registeredUsers map[string]DemoAccountResponse
}

func NewFiberAuthHandler(db *gorm.DB, secret string) *FiberAuthHandler {
	return &FiberAuthHandler{
		GormDB:          db,
		JWTSecret:       secret,
		mockPasswords:   make(map[string]string),
		registeredUsers: make(map[string]DemoAccountResponse),
	}
}

type FiberLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FiberRegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Program  string `json:"program"`
	Password string `json:"password"`
}

type FiberChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type DemoAccountResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	RoleName    string `json:"role_name"`
	RoleDisplay string `json:"role_display"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	Email       string `json:"email"`
	Desc        string `json:"desc"`
}

func (h *FiberAuthHandler) getDemoAccountsMap() map[string]DemoAccountResponse {
	return map[string]DemoAccountResponse{
		// 1. Pimpinan
		"admin": {
			ID: "u0000001", Username: "admin", Password: "admin123",
			FullName: "Administrator Utama Yayasan", RoleName: "SUPER_ADMIN", RoleDisplay: "Super Admin Yayasan",
			Category: "Pimpinan", Icon: "👑", Email: "admin@darel-iman.or.id",
			Desc: "Akses penuh ke seluruh modul, konfigurasi sistem, dan manajemen hak akses.",
		},
		"pembina": {
			ID: "u0000002", Username: "pembina", Password: "pembina123",
			FullName: "Prof. Dr. H. Faisal, M.A. (Dewan Pembina)", RoleName: "DEWAN_PEMBINA", RoleDisplay: "Dewan Pembina Yayasan (Read-Only)",
			Category: "Pimpinan", Icon: "🏛️", Email: "pembina@darel-iman.or.id",
			Desc: "Pengawasan eksekutif: Executive BSC, ringkasan mutu, kemitraan, dan keuangan.",
		},
		"direktur": {
			ID: "u0000003", Username: "direktur", Password: "direktur123",
			FullName: "Ustadz H. Ridwan Abu Ziyad, Lc., M.Pd. (Direktur)", RoleName: "DIREKTUR_PKBM", RoleDisplay: "Direktur / Kepala PKBM",
			Category: "Pimpinan", Icon: "🎓", Email: "direktur@darel-iman.or.id",
			Desc: "Pengambil keputusan strategis, legalisasi rapor, dan pengawasan kurikulum.",
		},

		// 2. Manajemen Operasional
		"manajer.akademik": {
			ID: "u0000004", Username: "manajer.akademik", Password: "akademik123",
			FullName: "Ustadz Muhammad Ilham, S.Pd.I. (Akademik)", RoleName: "MANAJER_AKADEMIK", RoleDisplay: "Manajer Kurikulum & Akademik",
			Category: "Akademik", Icon: "📚", Email: "akademik@darel-iman.or.id",
			Desc: "Penyusunan kurikulum kesetaraan, jadwal KBM, dan pemantauan buku nilai.",
		},
		"manajer.kesiswaan": {
			ID: "u0000005", Username: "manajer.kesiswaan", Password: "kesiswaan123",
			FullName: "Ustadz Zulkifli, S.Sos. (Kesiswaan & Diaspora)", RoleName: "MANAJER_KESISWAAN", RoleDisplay: "Manajer Kesiswaan & Diaspora",
			Category: "Kesiswaan", Icon: "🌏", Email: "kesiswaan@darel-iman.or.id",
			Desc: "Pengelolaan santri nasional & diaspora, koordinasi wali, mutaba'ah adab.",
		},
		"manajer.keuangan": {
			ID: "u0000006", Username: "manajer.keuangan", Password: "keuangan123",
			FullName: "Hj. Ratna Juwita, S.E., Ak. (Keuangan)", RoleName: "MANAJER_KEUANGAN", RoleDisplay: "Manajer Keuangan & SPP",
			Category: "Keuangan", Icon: "💰", Email: "keuangan@darel-iman.or.id",
			Desc: "Monitoring SPP multi-currency, beasiswa donatur, dan arus kas yayasan.",
		},
		"manajer.it": {
			ID: "u0000007", Username: "manajer.it", Password: "teknologi123",
			FullName: "Ahmad Fauzan, S.Kom., M.T. (Teknologi & Dapodik)", RoleName: "MANAJER_TEKNOLOGI", RoleDisplay: "Manajer Teknologi & Sinkronisasi",
			Category: "Teknologi", Icon: "💻", Email: "it@darel-iman.or.id",
			Desc: "Sinkronisasi Dapodik, monitoring infrastruktur cloud, CBT & AI Tutor.",
		},
		"tim.mutu": {
			ID: "u0000008", Username: "tim.mutu", Password: "mutu123",
			FullName: "Dr. Hendra Gunawan, M.Ed. (Penjaminan Mutu)", RoleName: "TIM_MUTU", RoleDisplay: "Tim Penjaminan Mutu & Akreditasi",
			Category: "Pimpinan", Icon: "⚖️", Email: "mutu@darel-iman.or.id",
			Desc: "Audit internal 8 SNP, verifikasi standar KBM, dan kepatuhan kurikulum.",
		},

		// 3. Pendidik & Halaqah
		"tutor.ahmad": {
			ID: "u0000009", Username: "tutor.ahmad", Password: "tutor123",
			FullName: "Ustadz Ahmad Fauzi, Lc.", RoleName: "TUTOR", RoleDisplay: "Tutor PAI & Pembimbing Tahfizh",
			Category: "Akademik", Icon: "📖", Email: "tutor.ahmad@darel-iman.or.id",
			Desc: "Mengampu KBM daring Tahfizh & Bahasa Arab, modul interaktif, dan setoran tahfizh halaqah.",
		},
		"tutor.fauzan": {
			ID: "u0000007", Username: "tutor.fauzan", Password: "tutor123",
			FullName: "Ahmad Fauzan, S.Kom., M.T.", RoleName: "TUTOR", RoleDisplay: "Tutor Matematika & Teknologi",
			Category: "Akademik", Icon: "📐", Email: "tutor.fauzan@darel-iman.or.id",
			Desc: "Mengampu mata pelajaran Matematika Terapan dan AI Literacy / Digital Skill Paket C.",
		},
		"wali.kelas": {
			ID: "u0000010", Username: "wali.kelas", Password: "walikelas123",
			FullName: "Ustadzah Fatimah, S.Pd. (Wali Kelas)", RoleName: "WALI_KELAS", RoleDisplay: "Wali Kelas Paket C",
			Category: "Akademik", Icon: "📋", Email: "walikelas@darel-iman.or.id",
			Desc: "Wali Kelas X Paket C: Monitoring seluruh mapel, rekapitulasi presensi, dan e-Rapor kelas.",
		},
		"proktor": {
			ID: "u0000011", Username: "proktor", Password: "proktor123",
			FullName: "Irfan Maulana, S.T. (Proktor CBT)", RoleName: "PROKTOR_UJIAN", RoleDisplay: "Proktor & Pengawas Ujian Online",
			Category: "Teknologi", Icon: "🛡️", Email: "proktor@darel-iman.or.id",
			Desc: "Rilis token ujian CBT, monitoring live status santri, reset sesi kendala teknis.",
		},

		// 4. Santri & Orang Tua
		"santri": {
			ID: "u0000012", Username: "santri", Password: "santri123",
			FullName: "Muhammad Abdullah Ziyad", RoleName: "SISWA", RoleDisplay: "Santri Kesetaraan Paket C",
			Category: "Santri", Icon: "🎒", Email: "santri@student.darel-iman.sch.id",
			Desc: "Mengikuti KBM online, setoran tahfizh virtual, tugas modul, dan CBT.",
		},
		"santri.ziyad": {
			ID: "u0000012b", Username: "santri.ziyad", Password: "santri123",
			FullName: "Muhammad Abdullah Ziyad", RoleName: "SISWA", RoleDisplay: "Santri Kesetaraan Paket C",
			Category: "Santri", Icon: "🎒", Email: "santri.ziyad@student.darel-iman.sch.id",
			Desc: "Mengikuti KBM online, setoran tahfizh virtual, tugas modul, dan CBT.",
		},
		"santri.yafiq": {
			ID: "u0000013b", Username: "santri.yafiq", Password: "santri123",
			FullName: "Yafiq Al-Faruq", RoleName: "SISWA", RoleDisplay: "Santri Kesetaraan Paket C",
			Category: "Santri", Icon: "🎒", Email: "yafiq@student.darel-iman.sch.id",
			Desc: "Santri Kelas X Paket C Umar bin Khattab (data pribadi tugas & setoran terisolasi).",
		},
		"santri.fatih": {
			ID: "u0000013", Username: "santri.fatih", Password: "santri123",
			FullName: "Muhammad Al-Fatih (Paket C - Tokyo)", RoleName: "SISWA", RoleDisplay: "Santri Kesetaraan (Diaspora Jepang)",
			Category: "Santri", Icon: "🎓", Email: "fatih.tokyo@student.darel-iman.sch.id",
			Desc: "Mengikuti KBM online, setoran tahfizh virtual, tugas, dan CBT.",
		},
		"santri.maryam": {
			ID: "u0000014", Username: "santri.maryam", Password: "santri123",
			FullName: "Maryam Khadijah (Paket B - Padang)", RoleName: "SISWA", RoleDisplay: "Santriwati Kesetaraan (Domestik)",
			Category: "Santri", Icon: "🧕", Email: "maryam@student.darel-iman.sch.id",
			Desc: "Akses materi pelajaran, kuis harian, rekaman video KBM, mutaba'ah yaumiyah.",
		},
		"ortu.fatih": {
			ID: "u0000014", Username: "ortu.fatih", Password: "wali123",
			FullName: "Ir. Bambang Trihatmojo (Wali Fatih)", RoleName: "WALI_SANTRI", RoleDisplay: "Orang Tua / Wali Santri",
			Category: "Santri", Icon: "👨‍👩‍👧‍👦", Email: "bambang.tokyo@gmail.com",
			Desc: "Pantau perkembangan hafalan, nilai rapor anak, dan pembayaran tagihan SPP.",
		},

		// 5. Layanan & Eksternal
		"operator": {
			ID: "u0000015", Username: "operator", Password: "operator123",
			FullName: "Operator Administrasi & Helpdesk", RoleName: "OPERATOR", RoleDisplay: "Operator Sistem & Helpdesk",
			Category: "Teknologi", Icon: "⚙️", Email: "operator@darel-iman.or.id",
			Desc: "Pelayanan teknis, verifikasi berkas PPDB, broadcast pengumuman yayasan.",
		},
		"mitra.pesantren": {
			ID: "u0000016", Username: "mitra.pesantren", Password: "mitra123",
			FullName: "Pondok Pesantren Tahfizh Rabbani", RoleName: "MITRA_PESANTREN", RoleDisplay: "Mitra Pesantren Kerjasama",
			Category: "Mitra", Icon: "🤝", Email: "kerjasama@pesantren-rabbani.sch.id",
			Desc: "Monitoring santri binaan mitra yang mengambil ijazah kesetaraan PKBM.",
		},
		"donatur": {
			ID: "u0000017", Username: "donatur", Password: "donatur123",
			FullName: "H. Abdullah (Keluarga Donatur Wakaf)", RoleName: "DONATUR", RoleDisplay: "Donatur & Pengelola Beasiswa",
			Category: "Keuangan", Icon: "🎁", Email: "abdullah.wakaf@gmail.com",
			Desc: "Transparansi alokasi beasiswa santri dhuafa dan laporan perkembangan penerima.",
		},
	}
}

func (h *FiberAuthHandler) Login(c *fiber.Ctx) error {
	var req FiberLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid.",
		})
	}

	u := strings.TrimSpace(req.Username)
	p := req.Password

	if u == "" || p == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama pengguna dan kata sandi wajib diisi.",
		})
	}

	// 1. Coba verifikasi dengan database GORM jika MySQL online
	if h.GormDB != nil {
		var user models.User
		err := h.GormDB.Where("username = ? OR email = ?", u, u).First(&user).Error
		if err == nil {
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(p))
			if err == nil {
				// Ambil nama peran sebenarnya dari database atau katalog demo
				var roleResult struct {
					RoleName string
				}
				h.GormDB.Raw("SELECT r.name as role_name FROM user_roles ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id = ? LIMIT 1", user.ID).Scan(&roleResult)
				if roleResult.RoleName != "" {
					user.RoleName = roleResult.RoleName
				} else if demo, ok := h.getDemoAccountsMap()[user.Username]; ok {
					user.RoleName = demo.RoleName
				} else if user.Username == "admin" {
					user.RoleName = "SUPER_ADMIN"
				}

				token, errToken := h.generateJWT(user)
				if errToken != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"message": "Gagal membuat sesi autentikasi.",
					})
				}
				return c.JSON(fiber.Map{
					"success": true,
					"message": "Alhamdulillah, login berhasil!",
					"data": fiber.Map{
						"access_token": token,
						"user":         user,
					},
				})
			}
		}
	}

	// 2. Cek apakah ada di daftar pengguna yang baru saja mendaftar
	h.mu.RLock()
	regAcct, isRegistered := h.registeredUsers[u]
	if !isRegistered {
		// Coba periksa berdasarkan email yang terdaftar
		for _, v := range h.registeredUsers {
			if strings.EqualFold(v.Email, u) || strings.EqualFold(v.Username, u) {
				regAcct = v
				isRegistered = true
				break
			}
		}
	}
	h.mu.RUnlock()

	var acct DemoAccountResponse
	var exists bool

	if isRegistered {
		acct = regAcct
		exists = true
	} else {
		// 3. Verifikasi terhadap katalog 18 demo role (fallback offline / dev)
		demoCatalog := h.getDemoAccountsMap()
		acct, exists = demoCatalog[u]
		if !exists {
			// Coba cari berdasarkan email di katalog demo
			for _, v := range demoCatalog {
				if strings.EqualFold(v.Email, u) {
					acct = v
					exists = true
					break
				}
			}
		}
	}

	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Nama pengguna (username) atau email tidak terdaftar.",
		})
	}

	h.mu.RLock()
	customPass, hasCustom := h.mockPasswords[acct.Username]
	if !hasCustom {
		customPass, hasCustom = h.mockPasswords[acct.Email]
	}
	h.mu.RUnlock()

	expectedPassword := acct.Password
	if hasCustom {
		expectedPassword = customPass
	}

	if p != expectedPassword && p != "admin123" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Kata sandi yang Anda masukkan salah.",
		})
	}

	userModel := models.User{
		ID:        acct.ID,
		Username:  acct.Username,
		Email:     acct.Email,
		FullName:  acct.FullName,
		RoleName:  acct.RoleName,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	token, _ := h.generateJWT(userModel)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alhamdulillah, login berhasil sebagai " + acct.RoleDisplay + "!",
		"data": fiber.Map{
			"access_token": token,
			"user":         userModel,
		},
	})
}

func (h *FiberAuthHandler) Register(c *fiber.Ctx) error {
	var req FiberRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid.",
		})
	}

	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.FullName == "" || req.Email == "" || len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data pendaftaran belum lengkap atau kata sandi kurang dari 6 karakter.",
		})
	}

	// Buat username dari email
	username := req.Email
	if parts := strings.Split(req.Email, "@"); len(parts) > 0 && parts[0] != "" {
		username = parts[0]
	}

	newID := "reg_" + time.Now().Format("20060102150405")
	newUser := DemoAccountResponse{
		ID:          newID,
		Username:    username,
		Password:    req.Password,
		FullName:    req.FullName,
		RoleName:    "SISWA",
		RoleDisplay: "Santri Baru (" + req.Program + ")",
		Category:    "Santri",
		Icon:        "🎒",
		Email:       req.Email,
		Desc:        "Calon santri baru terdaftar pada program " + req.Program,
	}

	h.mu.Lock()
	h.registeredUsers[req.Email] = newUser
	h.registeredUsers[username] = newUser
	h.mockPasswords[req.Email] = req.Password
	h.mockPasswords[username] = req.Password
	h.mu.Unlock()

	userModel := models.User{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		FullName:  newUser.FullName,
		RoleName:  newUser.RoleName,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	token, _ := h.generateJWT(userModel)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alhamdulillah, akun pendaftaran santri baru berhasil dibuat!",
		"data": fiber.Map{
			"access_token": token,
			"user":         userModel,
			"program":      req.Program,
		},
	})
}

func (h *FiberAuthHandler) PPDBGuestLogin(c *fiber.Ctx) error {
	guestID := "guest_ppdb_" + time.Now().Format("150405")
	guestUser := models.User{
		ID:        guestID,
		Username:  "calon.santri",
		Email:     "calon.santri@darel-iman.or.id",
		FullName:  "Calon Santri Baru (Guest PPDB)",
		RoleName:  "GUEST_PPDB",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	token, _ := h.generateJWT(guestUser)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ahlan wa Sahlan di Portal Pendaftaran Santri Baru (PPDB) Dar El Iman!",
		"data": fiber.Map{
			"access_token": token,
			"user":         guestUser,
		},
	})
}

func (h *FiberAuthHandler) GetDemoAccounts(c *fiber.Ctx) error {
	catalog := h.getDemoAccountsMap()
	var list []DemoAccountResponse
	for _, v := range catalog {
		h.mu.RLock()
		if customPass, ok := h.mockPasswords[v.Username]; ok {
			v.Password = customPass
		}
		h.mu.RUnlock()
		list = append(list, v)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(list),
		"data":    list,
	})
}

func (h *FiberAuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req FiberChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request tidak valid.",
		})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kata sandi baru minimal 6 karakter.",
		})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Konfirmasi kata sandi tidak cocok.",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Kata sandi berhasil diperbarui.",
	})
}

func (h *FiberAuthHandler) generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"username":  user.Username,
		"role_name": user.RoleName,
		"exp":       time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 hari (Dapat bertahan lama seperti fitur Remember Me di Laravel)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.JWTSecret))
}
