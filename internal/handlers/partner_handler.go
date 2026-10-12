package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PartnerHandler struct {
	DB              *sql.DB
	mu              sync.RWMutex
	partners        []models.Partner
	partnerStudents []models.PartnerStudent
}

func NewPartnerHandler(db *sql.DB) *PartnerHandler {
	h := &PartnerHandler{
		DB: db,
	}
	h.initMockData()
	return h
}

func (h *PartnerHandler) initMockData() {
	h.partners = []models.Partner{
		{
			ID:            "prt-001",
			Name:          "Pesantren Tahfizh Mitra Darul Qur'an",
			PartnerType:   "PESANTREN",
			Country:       "Indonesia",
			City:          "Bukittinggi",
			ContactPerson: "Ust. Lukman Hakim, Lc.",
			Phone:         "+628126789001",
			Email:         "pesantren.dq@gmail.com",
			IsActive:      true,
			StudentCount:  38,
			CreatedAt:     time.Now().Add(-120 * 24 * time.Hour),
		},
		{
			ID:            "prt-002",
			Name:          "Komunitas Muslim Indonesia Arab Saudi (KSA)",
			PartnerType:   "KOMUNITAS_DIASPORA",
			Country:       "Arab Saudi",
			City:          "Riyadh & Jeddah",
			ContactPerson: "Bpk. Ir. Hendra Saputra",
			Phone:         "+966509988776",
			Email:         "diaspora.ksa@darel-iman.or.id",
			IsActive:      true,
			StudentCount:  24,
			CreatedAt:     time.Now().Add(-90 * 24 * time.Hour),
		},
		{
			ID:            "prt-003",
			Name:          "Japan Islamic Indonesian Center (Tokyo & Osaka)",
			PartnerType:   "KOMUNITAS_DIASPORA",
			Country:       "Jepang",
			City:          "Tokyo",
			ContactPerson: "Ust. Dr. Fajar Nugroho",
			Phone:         "+818012345678",
			Email:         "jiic.tokyo@gmail.com",
			IsActive:      true,
			StudentCount:  12,
			CreatedAt:     time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			ID:            "prt-004",
			Name:          "Majelis Taklim Al-Muhajirin Malaysia",
			PartnerType:   "MAJELIS_TAKLIM",
			Country:       "Malaysia",
			City:          "Kuala Lumpur",
			ContactPerson: "Ibu Hj. Siti Nurhaliza",
			Phone:         "+60123456789",
			Email:         "mt.almuhajirin@kl.my",
			IsActive:      true,
			StudentCount:  16,
			CreatedAt:     time.Now().Add(-45 * 24 * time.Hour),
		},
	}

	h.partnerStudents = []models.PartnerStudent{
		{
			PartnerID:    "prt-002",
			StudentID:    "stu-005",
			StudentName:  "Khalid Al-Walid",
			NIS:          "202607005",
			ClassName:    "Paket C - Kelas 10 Abu Bakar",
			TahfizhJuz:   "Juz 30 (Lancar)",
			RegisteredAt: time.Now().Add(-40 * 24 * time.Hour),
		},
		{
			PartnerID:    "prt-002",
			StudentID:    "stu-118",
			StudentName:  "Sarah Al-Ghamdi",
			NIS:          "202607118",
			ClassName:    "Paket C - Kelas 10 Khadijah",
			TahfizhJuz:   "Juz 30 (Mumtaz)",
			RegisteredAt: time.Now().Add(-30 * 24 * time.Hour),
		},
		{
			PartnerID:    "prt-001",
			StudentID:    "stu-001",
			StudentName:  "Fatih Al-Ayyubi",
			NIS:          "202607001",
			ClassName:    "Paket C - Kelas 10 Abu Bakar",
			TahfizhJuz:   "Juz 30 & Juz 29",
			RegisteredAt: time.Now().Add(-50 * 24 * time.Hour),
		},
	}
}

// GET /api/v1/partners
func (h *PartnerHandler) ListPartners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.partners)
}

// POST /api/v1/partners
func (h *PartnerHandler) CreatePartner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.Partner
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.ContactPerson == "" {
		http.Error(w, `{"error":"Nama institusi dan nama narahubung wajib diisi"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = "prt-" + uuid.New().String()[:8]
	req.CreatedAt = time.Now()
	req.IsActive = true
	if req.Country == "" {
		req.Country = "Indonesia"
	}
	if req.PartnerType == "" {
		req.PartnerType = "PESANTREN"
	}

	h.partners = append(h.partners, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Mitra lembaga berhasil didaftarkan",
		"partner": req,
	})
}

// GET /api/v1/partners/:id/students
func (h *PartnerHandler) GetPartnerStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var res []models.PartnerStudent
	for _, ps := range h.partnerStudents {
		if id == "all" || ps.PartnerID == id {
			res = append(res, ps)
		}
	}

	json.NewEncoder(w).Encode(res)
}

// POST /api/v1/partners/:id/bulk-enroll
func (h *PartnerHandler) BulkEnrollPartnerStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partnerID := chi.URLParam(r, "id")

	var req struct {
		TargetPackage string   `json:"target_package"`
		StudentNames  []string `json:"student_names"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	count := 0
	for _, name := range req.StudentNames {
		if name == "" {
			continue
		}
		newStu := models.PartnerStudent{
			PartnerID:    partnerID,
			StudentID:    "stu-" + uuid.New().String()[:8],
			StudentName:  name,
			NIS:          fmt.Sprintf("202607%03d", len(h.partnerStudents)+10),
			ClassName:    fmt.Sprintf("Rombel %s Kolektif", req.TargetPackage),
			TahfizhJuz:   "Juz 30 (Baru Masuk)",
			RegisteredAt: time.Now(),
		}
		h.partnerStudents = append(h.partnerStudents, newStu)
		count++
	}

	// Update partner count
	for i := range h.partners {
		if h.partners[i].ID == partnerID {
			h.partners[i].StudentCount += count
			break
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        fmt.Sprintf("%d santri berhasil didaftarkan secara kolektif di bawah lembaga mitra", count),
		"enrolled_count": count,
		"partner_id":     partnerID,
	})
}

// GET /api/v1/partners/analytics/diaspora-map
func (h *PartnerHandler) GetDiasporaMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	summary := models.DiasporaMapSummary{
		TotalStudents:  152,
		NationalCount:  112,
		DiasporaCount:  40,
		CountriesCount: 6,
		Countries: []models.DiasporaCountryRecord{
			{CountryCode: "ID", CountryName: "Indonesia (WIB/WITA/WIT)", StudentCount: 112, Timezone: "UTC+7 to UTC+9", FlagEmoji: "🇮🇩"},
			{CountryCode: "SA", CountryName: "Arab Saudi (Riyadh/Makkah/Jeddah)", StudentCount: 22, Timezone: "UTC+3 (AST)", FlagEmoji: "🇸🇦"},
			{CountryCode: "JP", CountryName: "Jepang (Tokyo/Osaka/Nagoya)", StudentCount: 8, Timezone: "UTC+9 (JST)", FlagEmoji: "🇯🇵"},
			{CountryCode: "MY", CountryName: "Malaysia (Kuala Lumpur/Selangor)", StudentCount: 5, Timezone: "UTC+8 (MYT)", FlagEmoji: "🇲🇾"},
			{CountryCode: "AE", CountryName: "Uni Emirat Arab (Dubai/Abu Dhabi)", StudentCount: 3, Timezone: "UTC+4 (GST)", FlagEmoji: "🇦🇪"},
			{CountryCode: "GB", CountryName: "Inggris / Britania Raya (London)", StudentCount: 2, Timezone: "UTC+0 (BST)", FlagEmoji: "🇬🇧"},
		},
	}

	json.NewEncoder(w).Encode(summary)
}
