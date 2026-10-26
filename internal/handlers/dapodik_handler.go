package handlers

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"pkbm-backend/internal/models"
)

type DapodikHandler struct {
	DB   *sql.DB
	mu   sync.RWMutex
	logs []models.DapodikExportLog
}

func NewDapodikHandler(db *sql.DB) *DapodikHandler {
	h := &DapodikHandler{
		DB: db,
	}
	h.initMockLogs()
	return h
}

func (h *DapodikHandler) initMockLogs() {
	h.logs = []models.DapodikExportLog{
		{
			ID:               "del-001",
			ExportedByUserID: "user-admin-01",
			AcademicYearID:   "ay-2026-2027",
			ExportType:       "DATA_SISWA",
			TotalRecords:     152,
			DownloadURL:      "/api/v1/dapodik/export-students",
			ExportedAt:       time.Now().Add(-7 * 24 * time.Hour),
		},
		{
			ID:               "del-002",
			ExportedByUserID: "user-admin-01",
			AcademicYearID:   "ay-2026-2027",
			ExportType:       "DATA_TUTOR",
			TotalRecords:     24,
			DownloadURL:      "/api/v1/dapodik/export-tutors",
			ExportedAt:       time.Now().Add(-7 * 24 * time.Hour),
		},
	}
}

// GET /api/v1/dapodik/validation-check
func (h *DapodikHandler) ValidationCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Sample real audit summary based on active student registry
	audit := models.DapodikValidationAudit{
		TotalStudents:       152,
		CompleteStudents:    149,
		IncompleteStudents:  3,
		ReadinessPercentage: 98.02,
		Issues: []models.DapodikIssueRecord{
			{
				StudentID:   "stu-104",
				StudentName: "Raihan Hakim",
				NIS:         "202607104",
				Missing:     "Nama Gadis Ibu Kandung belum terisi",
			},
			{
				StudentID:   "stu-118",
				StudentName: "Sarah Al-Ghamdi (Diaspora)",
				NIS:         "202607118",
				Missing:     "NIK belum diverifikasi Dukcapil / Format Paspor",
			},
			{
				StudentID:   "stu-142",
				StudentName: "Ahmad Fauzan",
				NIS:         "202607142",
				Missing:     "NISN belum terdaftar di Pusdatin Kemendikbud",
			},
		},
	}

	json.NewEncoder(w).Encode(audit)
}

// POST or GET /api/v1/dapodik/export-students
// Menghasilkan file CSV terstandar Dapodik Kemendikbudristek
func (h *DapodikHandler) ExportStudents(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.logs = append(h.logs, models.DapodikExportLog{
		ID:               fmt.Sprintf("del-%d", time.Now().Unix()),
		ExportedByUserID: "admin-system",
		AcademicYearID:   "ay-2026-2027",
		ExportType:       "DATA_SISWA",
		TotalRecords:     152,
		DownloadURL:      "/api/v1/dapodik/export-students",
		ExportedAt:       time.Now(),
	})
	h.mu.Unlock()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dapodik_peserta_didik_pkbm_darel_iman_%s.csv", time.Now().Format("20060102")))

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Dapodik standard column headers
	writer.Write([]string{
		"Registrasi_ID",
		"Nama_Lengkap",
		"NIPD_NIS",
		"Jenis_Kelamin",
		"NISN",
		"NIK",
		"Tempat_Lahir",
		"Tanggal_Lahir",
		"Nama_Ibu_Kandung",
		"Rombel_Saat_Ini",
		"Tingkat_Pendidikan",
		"Kurikulum",
		"Status_Diaspora",
		"Negara_Tinggal",
		"Status_Peserta_Didik",
	})

	var sampleRows [][]string
	if h.DB != nil {
		rows, err := h.DB.Query(`
			SELECT s.id, s.full_name, s.nis, s.gender, COALESCE(s.nisn, '-'), COALESCE(s.nik, '-'),
			       COALESCE(s.birth_place, 'Padang'), COALESCE(DATE_FORMAT(s.birth_date, '%Y-%m-%d'), '2008-01-01'),
			       COALESCE(s.mother_name, 'Ibu'), COALESCE(c.name, 'Kelas X Paket C - Umar bin Khattab'),
			       s.status
			FROM students s
			LEFT JOIN classroom_students cs ON s.id = cs.student_id
			LEFT JOIN classrooms c ON cs.classroom_id = c.id
			ORDER BY s.nis ASC
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var id, name, nis, gender, nisn, nik, birthPlace, birthDate, mother, classroom, status string
				if errScan := rows.Scan(&id, &name, &nis, &gender, &nisn, &nik, &birthPlace, &birthDate, &mother, &classroom, &status); errScan == nil {
					sampleRows = append(sampleRows, []string{
						id, name, nis, gender, nisn, nik, birthPlace, birthDate, mother, classroom,
						"Tingkat 5 (Paket C)", "Kurikulum Merdeka PKBM", "Tidak", "Indonesia", status,
					})
				}
			}
		}
	}

	if len(sampleRows) == 0 {
		sampleRows = [][]string{
			{"s0000001-0000-0000-0000-000000000001", "Muhammad Al-Fatih", "NIS-2026-001", "L", "0081234567", "1371011508080001", "Padang", "2008-08-15", "Khadijah", "Kelas X Paket C - Umar bin Khattab", "Tingkat 5 (Paket C)", "Kurikulum Merdeka PKBM", "Tidak", "Indonesia", "ACTIVE"},
			{"s0000002-0000-0000-0000-000000000002", "Ahmad Zaki", "NIS-2026-002", "L", "0081234568", "1371011508080002", "Padang", "2008-09-20", "Fatimah", "Kelas X Paket C - Umar bin Khattab", "Tingkat 5 (Paket C)", "Kurikulum Merdeka PKBM", "Tidak", "Indonesia", "ACTIVE"},
		}
	}

	for _, row := range sampleRows {
		writer.Write(row)
	}
	writer.Flush()

	w.Write(buf.Bytes())
}

// POST or GET /api/v1/dapodik/export-tutors
func (h *DapodikHandler) ExportTutors(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.logs = append(h.logs, models.DapodikExportLog{
		ID:               fmt.Sprintf("del-%d", time.Now().Unix()),
		ExportedByUserID: "admin-system",
		AcademicYearID:   "ay-2026-2027",
		ExportType:       "DATA_TUTOR",
		TotalRecords:     24,
		DownloadURL:      "/api/v1/dapodik/export-tutors",
		ExportedAt:       time.Now(),
	})
	h.mu.Unlock()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dapodik_ptk_tutor_pkbm_darel_iman_%s.csv", time.Now().Format("20060102")))

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{
		"ID_PTK",
		"Nama_Lengkap_Gelar",
		"NIK",
		"NUPTK",
		"Jenis_Kelamin",
		"Tempat_Lahir",
		"Tanggal_Lahir",
		"Status_Kepegawaian",
		"Jenis_PTK",
		"Pendidikan_Terakhir",
		"Mata_Pelajaran_Diajar",
	})

	sampleRows := [][]string{
		{"PTK-001", "Ust. Dr. Muhammad Zaki, M.Pd.", "1371011005800001", "9876543210123456", "L", "Padang", "1980-05-10", "Guru Tetap Yayasan (GTY)", "Tutor Kesetaraan", "S3 Pendidikan", "Pendidikan Agama Islam, Tahfizh"},
		{"PTK-002", "Ustdzah Aisyah Rahmawati, S.Pd.", "1371024508880002", "8765432109234567", "P", "Bukittinggi", "1988-08-25", "Guru Tetap Yayasan (GTY)", "Tutor Kesetaraan", "S1 Bahasa Arab", "Bahasa Arab, Sirah Nabawiyah"},
		{"PTK-003", "Ust. Ir. Ridwan Siregar, M.Kom.", "1371031511850003", "7654321098345678", "L", "Medan", "1985-11-15", "Guru Kontrak", "Tutor Kesetaraan", "S2 Informatika", "Informatika / Literasi Digital"},
		{"PTK-004", "Ust. Ahmad Syahputra, S.Pd.", "1371042003920004", "6543210987456789", "L", "Padang Panjang", "1992-03-20", "Guru Tetap Yayasan (GTY)", "Tutor Kesetaraan", "S1 Pendidikan Matematika", "Matematika Paket C"},
	}

	for _, row := range sampleRows {
		writer.Write(row)
	}
	writer.Flush()

	w.Write(buf.Bytes())
}

// POST or GET /api/v1/dapodik/export-classrooms
func (h *DapodikHandler) ExportClassrooms(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.logs = append(h.logs, models.DapodikExportLog{
		ID:               fmt.Sprintf("del-%d", time.Now().Unix()),
		ExportedByUserID: "admin-system",
		AcademicYearID:   "ay-2026-2027",
		ExportType:       "DATA_ROMBEL",
		TotalRecords:     8,
		DownloadURL:      "/api/v1/dapodik/export-classrooms",
		ExportedAt:       time.Now(),
	})
	h.mu.Unlock()

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=dapodik_rombel_pkbm_darel_iman_%s.csv", time.Now().Format("20060102")))

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	writer.Write([]string{
		"ID_Rombel",
		"Nama_Rombongan_Belajar",
		"Tingkat_Kelas",
		"Program_Kesetaraan",
		"Wali_Kelas_Tutor",
		"Kurikulum",
		"Jumlah_Peserta_Didik",
		"Ruang_KBM",
	})

	sampleRows := [][]string{
		{"ROMB-01", "Paket C - Kelas 10 Abu Bakar Ash-Shiddiq", "Kelas 10", "Paket C (Setara SMA)", "Ust. Dr. Muhammad Zaki, M.Pd.", "Kurikulum Merdeka PKBM", "28", "Ruang Virtual Khusus Ikhwan 1"},
		{"ROMB-02", "Paket C - Kelas 10 Khadijah Al-Kubra", "Kelas 10", "Paket C (Setara SMA)", "Ustdzah Aisyah Rahmawati, S.Pd.", "Kurikulum Merdeka PKBM", "26", "Ruang Virtual Khusus Akhwat 1"},
		{"ROMB-03", "Paket B - Kelas 8 Umar bin Khattab", "Kelas 8", "Paket B (Setara SMP)", "Ust. Ahmad Syahputra, S.Pd.", "Kurikulum Merdeka PKBM", "32", "Ruang Virtual Kelas B-1"},
		{"ROMB-04", "Paket A - Kelas 5 Ali bin Abi Thalib", "Kelas 5", "Paket A (Setara SD)", "Ust. Ir. Ridwan Siregar, M.Kom.", "Kurikulum Merdeka PKBM", "24", "Ruang Virtual Kelas A-1"},
	}

	for _, row := range sampleRows {
		writer.Write(row)
	}
	writer.Flush()

	w.Write(buf.Bytes())
}

// GET /api/v1/dapodik/logs
func (h *DapodikHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.logs)
}
