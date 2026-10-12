package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type QualityHandler struct {
	DB         *sql.DB
	mu         sync.RWMutex
	standards  []models.AccreditationStandard
	items      []models.AccreditationItem
	sops       []models.InstitutionSOP
	risks      []models.RiskRegister
	indicators []models.BSCIndicator
}

func NewQualityHandler(db *sql.DB) *QualityHandler {
	h := &QualityHandler{
		DB: db,
	}
	h.initMockData()
	return h
}

func (h *QualityHandler) initMockData() {
	h.standards = []models.AccreditationStandard{
		{ID: "std-001", Code: "BAN-PDM-01", StandardName: "Standar Kompetensi Lulusan (SKL)", ComponentName: "Capaian Pembelajaran Kurikulum Merdeka & Tahfizh", MaxScore: 4.0},
		{ID: "std-002", Code: "BAN-PDM-02", StandardName: "Standar Isi", ComponentName: "KOSP Kurikulum Operasional Satuan Pendidikan", MaxScore: 4.0},
		{ID: "std-003", Code: "BAN-PDM-03", StandardName: "Standar Proses", ComponentName: "KBM Virtual Interaktif Sinkron & Asinkron", MaxScore: 4.0},
		{ID: "std-004", Code: "BAN-PDM-04", StandardName: "Standar Pendidik dan Tenaga Kependidikan", ComponentName: "Kualifikasi Pendidik S1/S2 & Sertifikasi Tutor", MaxScore: 4.0},
		{ID: "std-005", Code: "BAN-PDM-05", StandardName: "Standar Sarana dan Prasarana", ComponentName: "Infrastruktur Server LMS, Cloud & Studio Rekaman", MaxScore: 4.0},
		{ID: "std-006", Code: "BAN-PDM-06", StandardName: "Standar Pengelolaan", ComponentName: "Tata Kelola Digital & Layanan Santri Diaspora", MaxScore: 4.0},
		{ID: "std-007", Code: "BAN-PDM-07", StandardName: "Standar Pembiayaan", ComponentName: "Transparansi Kas SPP & Beasiswa Donatur", MaxScore: 4.0},
		{ID: "std-008", Code: "BAN-PDM-08", StandardName: "Standar Penilaian Pendidikan", ComponentName: "Engine CBT Digital, Proctoring & E-Rapor Ber-QR", MaxScore: 4.0},
	}

	h.items = []models.AccreditationItem{
		{
			ID:                   "item-001",
			StandardID:           "std-001",
			StandardCode:         "BAN-PDM-01",
			ItemNumber:           "1.1",
			IndicatorDescription: "Dokumen capaian kelulusan santri Paket A/B/C dan sertifikat tuntas hafalan Juz 30.",
			FulfillmentStatus:    "FULFILLED",
			EvidenceFileURL:      "https://cdn.darel-iman.sch.id/mutu/eviden_skl_tahfizh.pdf",
			Notes:                "100% santri melampaui KKM nasional dan tuntas target juz",
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   "item-002",
			StandardID:           "std-002",
			StandardCode:         "BAN-PDM-02",
			ItemNumber:           "2.1",
			IndicatorDescription: "Dokumen KOSP Kurikulum Operasional disahkan Dinas Pendidikan Kota Padang.",
			FulfillmentStatus:    "FULFILLED",
			EvidenceFileURL:      "https://cdn.darel-iman.sch.id/mutu/kosp_pkbm_2026.pdf",
			Notes:                "Telah divalidasi penilik kesetaraan Disdikbud Padang",
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   "item-003",
			StandardID:           "std-003",
			StandardCode:         "BAN-PDM-03",
			ItemNumber:           "3.1",
			IndicatorDescription: "Rekam jejak KBM interaktif daring, presensi real-time, dan modul ajar terbit.",
			FulfillmentStatus:    "READY_AUDIT",
			EvidenceFileURL:      "https://cdn.darel-iman.sch.id/mutu/log_kbm_virtual_lms.pdf",
			Notes:                "Sistem LMS merekam 100% riwayat interaksi KBM santri",
			UpdatedAt:            time.Now(),
		},
		{
			ID:                   "item-004",
			StandardID:           "std-008",
			StandardCode:         "BAN-PDM-08",
			ItemNumber:           "8.1",
			IndicatorDescription: "Sistem penilaian sumatif CBT dengan rekaman proctoring dan verifikasi QR rapor sah.",
			FulfillmentStatus:    "FULFILLED",
			EvidenceFileURL:      "https://cdn.darel-iman.sch.id/mutu/laporan_cbt_dan_rapor_digital.pdf",
			Notes:                "Menerapkan anti-cheat proctoring dan kode verifikasi keabsahan online",
			UpdatedAt:            time.Now(),
		},
	}

	h.sops = []models.InstitutionSOP{
		{ID: "sop-001", Code: "SOP-AKD-001", Title: "SOP Pelaksanaan KBM Virtual & Google Meet", Department: "Akademik", Version: "2.1", FileURL: "https://cdn.darel-iman.sch.id/sop/sop_kbm_virtual.pdf", EffectiveDate: "2026-07-01", CreatedAt: time.Now()},
		{ID: "sop-002", Code: "SOP-THF-002", Title: "SOP Halaqah Tahfizh Daring & Mutaba'ah Al-Qur'an", Department: "Tahfizh", Version: "2.0", FileURL: "https://cdn.darel-iman.sch.id/sop/sop_halaqah_tahfizh.pdf", EffectiveDate: "2026-07-01", CreatedAt: time.Now()},
		{ID: "sop-003", Code: "SOP-KEU-003", Title: "SOP Rekonsiliasi Kas SPP & Verifikasi Pembayaran", Department: "Keuangan", Version: "1.5", FileURL: "https://cdn.darel-iman.sch.id/sop/sop_keuangan_spp.pdf", EffectiveDate: "2026-07-01", CreatedAt: time.Now()},
		{ID: "sop-004", Code: "SOP-IT-004", Title: "SOP Penyelenggaraan Ujian CBT & Proctoring Daring", Department: "IT & Penilaian", Version: "1.0", FileURL: "https://cdn.darel-iman.sch.id/sop/sop_cbt_proctoring.pdf", EffectiveDate: "2026-08-15", CreatedAt: time.Now()},
	}

	h.risks = []models.RiskRegister{
		{ID: "rr-001", RiskTitle: "Ketidakstabilan Koneksi Internet Santri Diaspora saat KBM", Category: "TEKNOLOGI", SeverityLevel: "MEDIUM", MitigationPlan: "Menyediakan rekaman video asinkron KBM dan materi PDF terkompresi di LMS", Status: "MITIGATED", CreatedAt: time.Now()},
		{ID: "rr-002", RiskTitle: "Kecurangan Akademik / Tab Switching pada Ujian Daring", Category: "OPERASIONAL", SeverityLevel: "HIGH", MitigationPlan: "Menerapkan engine CBT Fullscreen Lock dan audit log proctoring otomatis", Status: "MITIGATED", CreatedAt: time.Now()},
		{ID: "rr-003", RiskTitle: "Keterlambatan Sinkronisasi NISN Santri Baru ke Dapodik", Category: "OPERASIONAL", SeverityLevel: "MEDIUM", MitigationPlan: "Audit kelengkapan NIK/Akta sejak proses PPDB Online sebelum semester berjalan", Status: "OPEN", CreatedAt: time.Now()},
	}

	h.indicators = []models.BSCIndicator{
		{ID: "bsc-01", Perspective: "FINANCIAL", IndicatorName: "Tingkat Kelancaran Pembayaran SPP Santri", TargetValue: 95.0, CurrentValue: 93.5, Unit: "%", PeriodYear: 2026, AchievementPct: 98.42, UpdatedAt: time.Now()},
		{ID: "bsc-02", Perspective: "FINANCIAL", IndicatorName: "Alokasi Beasiswa Santri Yatim / Dhuafa", TargetValue: 15.0, CurrentValue: 18.0, Unit: "%", PeriodYear: 2026, AchievementPct: 120.0, UpdatedAt: time.Now()},
		{ID: "bsc-03", Perspective: "CUSTOMER_STAKEHOLDER", IndicatorName: "Indeks Kepuasan Orang Tua & Santri Diaspora", TargetValue: 90.0, CurrentValue: 94.2, Unit: "%", PeriodYear: 2026, AchievementPct: 104.67, UpdatedAt: time.Now()},
		{ID: "bsc-04", Perspective: "CUSTOMER_STAKEHOLDER", IndicatorName: "Rasio Rekomendasi Mitra Lembaga (NPS)", TargetValue: 80.0, CurrentValue: 87.5, Unit: "%", PeriodYear: 2026, AchievementPct: 109.38, UpdatedAt: time.Now()},
		{ID: "bsc-05", Perspective: "INTERNAL_PROCESS", IndicatorName: "Tingkat Ketuntasan Hafalan Target Tahfizh (Juz 30)", TargetValue: 90.0, CurrentValue: 92.0, Unit: "%", PeriodYear: 2026, AchievementPct: 102.22, UpdatedAt: time.Now()},
		{ID: "bsc-06", Perspective: "INTERNAL_PROCESS", IndicatorName: "Kepatuhan Pelaksanaan Jadwal KBM Sesuai Kalender", TargetValue: 98.0, CurrentValue: 97.4, Unit: "%", PeriodYear: 2026, AchievementPct: 99.39, UpdatedAt: time.Now()},
		{ID: "bsc-07", Perspective: "LEARNING_GROWTH", IndicatorName: "Partisipasi Tutor dalam Pelatihan Digital & AI", TargetValue: 100.0, CurrentValue: 100.0, Unit: "%", PeriodYear: 2026, AchievementPct: 100.0, UpdatedAt: time.Now()},
		{ID: "bsc-08", Perspective: "LEARNING_GROWTH", IndicatorName: "Kesiapan Pemenuhan Akreditasi BAN-PDM (Skor A)", TargetValue: 92.0, CurrentValue: 96.5, Unit: "Skor", PeriodYear: 2026, AchievementPct: 104.89, UpdatedAt: time.Now()},
	}
}

// GET /api/v1/quality/accreditation/checklist
func (h *QualityHandler) GetAccreditationChecklist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	type ResultItem struct {
		Standard models.AccreditationStandard `json:"standard"`
		Items    []models.AccreditationItem   `json:"items"`
	}

	var results []ResultItem
	fulfilledCount := 0
	totalItems := len(h.items)

	for _, s := range h.standards {
		var matched []models.AccreditationItem
		for _, item := range h.items {
			if item.StandardID == s.ID {
				matched = append(matched, item)
				if item.FulfillmentStatus == "FULFILLED" {
					fulfilledCount++
				}
			}
		}
		results = append(results, ResultItem{Standard: s, Items: matched})
	}

	scoreSimulation := 96.5 // Skor A (Unggul)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_standards":  len(h.standards),
		"total_items":      totalItems,
		"fulfilled_items":  fulfilledCount,
		"simulation_score": scoreSimulation,
		"predicate":        "A (UNGGUL)",
		"standards_data":   results,
	})
}

// PUT /api/v1/quality/accreditation/items/:id
func (h *QualityHandler) UpdateAccreditationItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	var req struct {
		FulfillmentStatus string `json:"fulfillment_status"`
		EvidenceFileURL   string `json:"evidence_file_url"`
		Notes             string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for i := range h.items {
		if h.items[i].ID == id {
			if req.FulfillmentStatus != "" {
				h.items[i].FulfillmentStatus = req.FulfillmentStatus
			}
			if req.EvidenceFileURL != "" {
				h.items[i].EvidenceFileURL = req.EvidenceFileURL
			}
			if req.Notes != "" {
				h.items[i].Notes = req.Notes
			}
			h.items[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		http.Error(w, `{"error":"Item akreditasi tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Status pemenuhan eviden akreditasi berhasil diperbarui",
	})
}

// GET /api/v1/quality/sops
func (h *QualityHandler) ListSOPs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.sops)
}

// POST /api/v1/quality/sops
func (h *QualityHandler) CreateSOP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.InstitutionSOP
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = "sop-" + uuid.New().String()[:8]
	req.CreatedAt = time.Now()
	if req.Version == "" {
		req.Version = "1.0"
	}
	if req.EffectiveDate == "" {
		req.EffectiveDate = time.Now().Format("2006-01-02")
	}

	h.sops = append(h.sops, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dokumen SOP lembaga berhasil disimpan",
		"sop":     req,
	})
}

// GET /api/v1/quality/risks
func (h *QualityHandler) ListRisks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.risks)
}

// POST /api/v1/quality/risks
func (h *QualityHandler) CreateRisk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RiskRegister
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = "rr-" + uuid.New().String()[:8]
	req.CreatedAt = time.Now()
	if req.Status == "" {
		req.Status = "OPEN"
	}

	h.risks = append(h.risks, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Register risiko mutu berhasil dicatat",
		"risk":    req,
	})
}

// GET /api/v1/quality/bsc/dashboard
func (h *QualityHandler) GetBSCDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	summary := models.BSCDashboardSummary{
		FinancialScore:       99.21,
		CustomerScore:        107.03,
		InternalProcessScore: 100.81,
		LearningGrowthScore:  102.45,
		OverallIndex:         102.38,
		Indicators:           h.indicators,
	}

	json.NewEncoder(w).Encode(summary)
}
