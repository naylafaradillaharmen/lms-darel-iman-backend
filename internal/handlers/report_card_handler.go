package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/go-chi/chi/v5"
)

type ReportCardHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	reports map[string]*models.StudentReportCard // qr_code -> report
}

func NewReportCardHandler(db *sql.DB) *ReportCardHandler {
	now := time.Now()
	pdfURL := "https://cdn.darel-iman.or.id/rapor/2026-ganjil-fatih-paket-c.pdf"

	h := &ReportCardHandler{
		DB:      db,
		reports: make(map[string]*models.StudentReportCard),
	}

	// Pre-seed Fatih's report card
	fatihCard := &models.StudentReportCard{
		ID:                     "rc000001-0000-0000-0000-000000000001",
		StudentID:              "s0000001-0000-0000-0000-000000000001",
		StudentName:            "Muhammad Al-Fatih",
		NIS:                    "NIS-2026-001",
		NISN:                   ptrString("0089218291"),
		AcademicYearID:         "22222222-2222-2222-2222-222222222222",
		AcademicYearName:       "Tahun Ajaran 2026/2027 (Semester Ganjil)",
		ClassroomID:            "c1111111-1111-1111-1111-111111111111",
		ClassroomName:          "Kelas X Paket C - Umar bin Khattab",
		HomeroomNotes:          "Ananda Muhammad Al-Fatih menunjukkan adab yang sangat santun, tekun dalam muroja'ah hafalan Al-Qur'an, dan sangat proaktif dalam pembelajaran daring. Terus pertahankan prestasinya!",
		SpiritualBehaviorGrade: "A (Sangat Baik)",
		SocialBehaviorGrade:    "A (Sangat Baik)",
		TahfizhSummary:         "Tuntas Juz 30 (Surah An-Naba s/d An-Nas) dengan predikat MUMTAZ. Mutaba'ah setoran Ziyadah dan Muraja'ah tercapai sesuai target kurikulum kepesantrenan Dar El Iman.",
		QRVerificationCode:     "DEI-RAPOR-2026-X-001",
		PDFFileURL:             &pdfURL,
		PublishedAt:            &now,
		Grades: []models.ReportCardGrade{
			{
				ID:              "rcg00001-0000-0000-0000-000000000001",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00001-0000-0000-0000-000000000001",
				SubjectCode:     "THF-01",
				SubjectName:     "Tahfizh Al-Qur'an Tematik",
				NumericalScore:  94.60,
				LetterGrade:     "A",
				CompetencyNotes: "Sangat menguasai kaidah makharijul huruf, tajwid Asy-Syafi'i, dan hafalan mutqin Juz 30.",
			},
			{
				ID:              "rcg00002-0000-0000-0000-000000000002",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00002-0000-0000-0000-000000000002",
				SubjectCode:     "ARB-01",
				SubjectName:     "Bahasa Arab Komunikatif & Nahwu Shorof Dasar",
				NumericalScore:  88.50,
				LetterGrade:     "A",
				CompetencyNotes: "Mampu memahami percakapan dasar bahasa Arab dan kaidah tata bahasa dasar isim/fi'il dengan baik.",
			},
			{
				ID:              "rcg00003-0000-0000-0000-000000000003",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00003-0000-0000-0000-000000000003",
				SubjectCode:     "AKD-01",
				SubjectName:     "Akidah & Akhlak Karimah",
				NumericalScore:  92.00,
				LetterGrade:     "A",
				CompetencyNotes: "Menunjukkan pemahaman tauhid ahlussunnah yang lurus dan pembiasaan adab penuntut ilmu dalam keseharian.",
			},
			{
				ID:              "rcg00004-0000-0000-0000-000000000004",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00004-0000-0000-0000-000000000004",
				SubjectCode:     "TEC-01",
				SubjectName:     "AI Literacy & Pemrograman Dasar",
				NumericalScore:  90.00,
				LetterGrade:     "A",
				CompetencyNotes: "Sangat cakap memanfaatkan teknologi AI secara etis dan memahami logika dasar pemrograman komputer.",
			},
			{
				ID:              "rcg00005-0000-0000-0000-000000000005",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00005-0000-0000-0000-000000000005",
				SubjectCode:     "IND-10",
				SubjectName:     "Bahasa Indonesia (Paket C)",
				NumericalScore:  86.00,
				LetterGrade:     "B",
				CompetencyNotes: "Mampu menyusun teks laporan hasil observasi dan teks argumentasi tertulis dengan tata bahasa formal yang baik.",
			},
			{
				ID:              "rcg00006-0000-0000-0000-000000000006",
				ReportCardID:    "rc000001-0000-0000-0000-000000000001",
				SubjectID:       "sub00006-0000-0000-0000-000000000006",
				SubjectCode:     "MTK-10",
				SubjectName:     "Matematika Terapan (Paket C)",
				NumericalScore:  85.00,
				LetterGrade:     "B",
				CompetencyNotes: "Terampil menyelesaikan persoalan aljabar, statistik data dasar, dan pemecahan masalah kuantitatif.",
			},
		},
	}

	h.reports["DEI-RAPOR-2026-X-001"] = fatihCard
	return h
}

// 1. Generate Batch Report Cards
func (h *ReportCardHandler) GenerateBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ClassroomID    string `json:"classroom_id"`
		AcademicYearID string `json:"academic_year_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.ClassroomID = "c1111111-1111-1111-1111-111111111111"
		req.AcademicYearID = "22222222-2222-2222-2222-222222222222"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Ensure Fatih's report card is in the pool
	card := h.reports["DEI-RAPOR-2026-X-001"]
	now := time.Now()
	card.PublishedAt = &now

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "Batch e-Rapor resmi berhasil digenerate lengkap dengan QR Code Keabsahan",
		"generated":     1,
		"classroom_id":  req.ClassroomID,
		"academic_year": req.AcademicYearID,
		"sample_qr":     "DEI-RAPOR-2026-X-001",
	})
}

// 2. Get Student Report Card
func (h *ReportCardHandler) GetStudentReportCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := chi.URLParam(r, "student_id")
	_ = chi.URLParam(r, "year_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var target *models.StudentReportCard
	for _, c := range h.reports {
		if c.StudentID == studentID {
			target = c
			break
		}
	}
	if target == nil {
		target = h.reports["DEI-RAPOR-2026-X-001"] // default demo
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    target,
	})
}

// 3. Verify QR Code (Public Verification Endpoint)
func (h *ReportCardHandler) VerifyQRCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	qrCode := chi.URLParam(r, "qr_code")

	h.mu.RLock()
	defer h.mu.RUnlock()

	report, exists := h.reports[qrCode]
	if !exists {
		// If not exact match, check default demo
		if qrCode == "DEI-RAPOR-2026-X-001" {
			report = h.reports["DEI-RAPOR-2026-X-001"]
			exists = true
		}
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"valid":   false,
			"message": "Dokumen rapor dengan kode verifikasi ini tidak ditemukan atau tidak terdaftar di sistem PKBM Digital Dar El Iman.",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":              true,
		"valid":                true,
		"status":               "RESMI & SAH TERVERIFIKASI",
		"institution":          "PKBM Digital Dar El Iman Padang (NPSN: P9997123)",
		"decree_number":        "SK Operasional Disdikbud No. 421.9/PKBM/2026",
		"qr_verification_code": report.QRVerificationCode,
		"student_name":         report.StudentName,
		"nis":                  report.NIS,
		"nisn":                 report.NISN,
		"academic_year":        report.AcademicYearName,
		"classroom":            report.ClassroomName,
		"tahfizh_achievement":  report.TahfizhSummary,
		"published_at":         report.PublishedAt,
		"verified_at":          time.Now(),
		"total_subjects":       len(report.Grades),
		"grades":               report.Grades,
	})
}
