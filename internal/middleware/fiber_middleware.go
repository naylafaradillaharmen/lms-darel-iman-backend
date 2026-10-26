package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// FiberClaims merepresentasikan payload JWT yang terstandarisasi untuk Fiber
type FiberClaims struct {
	UserID   string `json:"sub"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

// FiberAuthMiddleware memvalidasi Bearer Token JWT pada request HTTP Fiber
func FiberAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    "UNAUTHORIZED",
				"message": "Header Authorization tidak ditemukan. Harap sertakan token Bearer.",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "Format token tidak valid. Gunakan format 'Bearer <token>'.",
			})
		}

		tokenString := parts[1]
		claims := &FiberClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak diharapkan: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    "TOKEN_EXPIRED_OR_INVALID",
				"message": "Token sesi tidak valid atau telah kedaluwarsa. Silakan login kembali.",
			})
		}

		// Hitung kategori peran
		category := DetermineRoleCategory(claims.RoleName)

		// Simpan klaim ke context locals Fiber
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role_name", claims.RoleName)
		c.Locals("role_category", category)

		return c.Next()
	}
}

// DetermineRoleCategory memetakan peran ke 7 kategori utama PKBM
func DetermineRoleCategory(roleName string) string {
	switch roleName {
	case "SUPER_ADMIN", "DEWAN_PEMBINA", "DIREKTUR_PKBM", "TIM_MUTU":
		return "PIMPINAN"
	case "MANAJER_AKADEMIK", "MANAJER_KESISWAAN", "MANAJER_KEUANGAN", "MANAJER_TEKNOLOGI":
		return "MANAJEMEN"
	case "WALI_KELAS", "TUTOR", "GURU":
		return "PENDIDIK"
	case "PROKTOR_UJIAN", "OPERATOR", "ADMIN":
		return "TEKNIS"
	case "SISWA", "SANTRI", "STUDENT":
		return "SANTRI"
	case "WALI_SANTRI", "ORANG_TUA":
		return "WALI_SANTRI"
	case "GUEST_PPDB":
		return "GUEST"
	default:
		return "UMUM"
	}
}

// FiberRequireRole membatasi endpoint hanya untuk peran-peran tertentu (Super Admin selalu diizinkan)
func FiberRequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentRole, _ := c.Locals("role_name").(string)
		if currentRole == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"code":    "SESSION_NOT_FOUND",
				"message": "Sesi autentikasi tidak ditemukan. Harap login terlebih dahulu.",
			})
		}

		// Super Admin selalu memiliki akses penuh (Bypass)
		if currentRole == "SUPER_ADMIN" {
			return c.Next()
		}

		for _, role := range allowedRoles {
			if currentRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success":        false,
			"code":           "ROLE_FORBIDDEN",
			"message":        fmt.Sprintf("Akses ditolak: Peran '%s' tidak memiliki hak akses untuk aksi ini.", currentRole),
			"current_role":   currentRole,
			"required_roles": allowedRoles,
		})
	}
}

// FiberRequireCategory membatasi endpoint berdasarkan kategori (PIMPINAN, MANAJEMEN, PENDIDIK, SANTRI)
func FiberRequireCategory(allowedCategories ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentCategory, _ := c.Locals("role_category").(string)
		currentRole, _ := c.Locals("role_name").(string)

		if currentRole == "SUPER_ADMIN" {
			return c.Next()
		}

		for _, cat := range allowedCategories {
			if currentCategory == cat {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success":            false,
			"code":               "CATEGORY_FORBIDDEN",
			"message":            fmt.Sprintf("Akses ditolak: Kategori akun '%s' tidak berwenang pada endpoint ini.", currentCategory),
			"current_category":   currentCategory,
			"allowed_categories": allowedCategories,
		})
	}
}
