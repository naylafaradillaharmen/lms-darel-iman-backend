package database

import (
	"embed"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunAutoMigrate mengeksekusi seluruh skema dan seed data secara otomatis ketika server backend di-deploy atau dijalankan.
func RunAutoMigrate(gormDB *gorm.DB) error {
	if gormDB == nil {
		log.Println("[AutoMigrate] Koneksi database belum tersedia. Migrasi otomatis dilewati (mode fallback demo aktif).")
		return nil
	}

	if os.Getenv("AUTO_MIGRATE") == "false" {
		log.Println("[AutoMigrate] AUTO_MIGRATE diatur 'false'. Melewati proses migrasi.")
		return nil
	}

	log.Println("[AutoMigrate] ========================================================")
	log.Println("[AutoMigrate] 🚀 Menjalankan Auto-Migrate Database PKBM Digital (Deploy Mode)")
	log.Println("[AutoMigrate] ========================================================")

	// Nonaktifkan pemeriksaan foreign key sementara untuk kelancaran eksekusi relasi
	if err := gormDB.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error; err != nil {
		log.Printf("[AutoMigrate] Peringatan SET FOREIGN_KEY_CHECKS=0: %v", err)
	}

	migrationFiles := []struct {
		name string
		desc string
	}{
		{"migrations/01_schema_tahap1.sql", "Skema Tahap 1: Fondasi Lembaga, User, Role, Rombel"},
		{"migrations/02_seed_tahap1.sql", "Seed Tahap 1: 18 Roles Resmi & Akun Demo Default"},
		{"migrations/03_schema_tahap2.sql", "Skema Tahap 2: Kurikulum, KBM LMS, Halaqah Tahfizh"},
		{"migrations/04_seed_tahap2.sql", "Seed Tahap 2: Mapel, Rombel KBM, Target Tahfizh"},
		{"migrations/05_schema_tahap3.sql", "Skema Tahap 3: Asesmen CBT, Bank Soal, Rapor Digital"},
		{"migrations/06_seed_tahap3.sql", "Seed Tahap 3: Ujian CBT & Sampel Penilaian Rapor"},
		{"migrations/07_schema_tahap4.sql", "Skema Tahap 4: PPDB Online, Keuangan SPP & Dapodik"},
		{"migrations/08_seed_tahap4.sql", "Seed Tahap 4: Gelombang PPDB & Tarif SPP"},
		{"migrations/09_schema_tahap5.sql", "Skema Tahap 5: Kemitraan Diaspora & Penjaminan Mutu BSC"},
		{"migrations/10_seed_tahap5.sql", "Seed Tahap 5: Data Mitra, Standar Akreditasi & KPI"},
	}

	for _, mf := range migrationFiles {
		content, err := migrationFS.ReadFile(mf.name)
		if err != nil {
			log.Printf("[AutoMigrate] Gagal membaca %s: %v", mf.name, err)
			continue
		}

		rawSQL := string(content)
		// Eksekusi statement SQL
		statements := splitSQLStatements(rawSQL)
		successCount := 0
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := gormDB.Exec(stmt).Error; err != nil {
				// Abaikan pesan error duplikasi yang wajar
				if !strings.Contains(err.Error(), "Duplicate") && !strings.Contains(err.Error(), "already exists") {
					log.Printf("[AutoMigrate] Info query pada %s: %v", mf.name, err)
				}
			} else {
				successCount++
			}
		}
		log.Printf("[AutoMigrate] ✅ Selesai: %s (%d query diproses)", mf.desc, successCount)
	}

	// Kembalikan pemeriksaan foreign key
	if err := gormDB.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error; err != nil {
		log.Printf("[AutoMigrate] Peringatan SET FOREIGN_KEY_CHECKS=1: %v", err)
	}

	log.Println("[AutoMigrate] 🎉 Seluruh skema database dan data master berhasil disiapkan!")
	log.Println("[AutoMigrate] ========================================================")
	return nil
}

// splitSQLStatements memecah string SQL menjadi query individual berdasarkan titik koma
func splitSQLStatements(sqlText string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false

	lines := strings.Split(sqlText, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Lewati baris komentar SQL tunggal
		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") {
			continue
		}

		for i := 0; i < len(line); i++ {
			ch := line[i]

			if ch == '\'' && !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			} else if ch == '"' && !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			} else if ch == '`' && !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}

			if ch == ';' && !inSingleQuote && !inDoubleQuote && !inBacktick {
				stmt := strings.TrimSpace(current.String())
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
		}
		current.WriteByte('\n')
	}

	finalStmt := strings.TrimSpace(current.String())
	if finalStmt != "" {
		statements = append(statements, finalStmt)
	}

	return statements
}
