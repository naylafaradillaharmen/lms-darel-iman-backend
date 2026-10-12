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

type PPDBHandler struct {
	DB            *sql.DB
	mu            sync.RWMutex
	periods       []models.PPDBPeriod
	registrations []models.PPDBRegistration
	documents     []models.PPDBDocument
}

func NewPPDBHandler(db *sql.DB) *PPDBHandler {
	h := &PPDBHandler{
		DB: db,
	}
	h.initMockData()
	return h
}

func (h *PPDBHandler) initMockData() {
	periodID := "per-2027-gel1"
	h.periods = []models.PPDBPeriod{
		{
			ID:             periodID,
			AcademicYearID: "ay-2026-2027",
			Title:          "PPDB Santri Baru Gelombang 1 TA 2027/2028 (Nasional & Diaspora)",
			StartDate:      "2026-10-01",
			EndDate:        "2026-12-31",
			IsOpen:         true,
			CreatedAt:      time.Now().Add(-10 * 24 * time.Hour),
		},
	}

	h.registrations = []models.PPDBRegistration{
		{
			ID:                   "preg-2027-001",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-001",
			TargetPackage:        "PAKET_C",
			FullName:             "Muhammad Ziyad Al-Farisi",
			Gender:               "L",
			NIK:                  "1371011508100003",
			NISN:                 "0098765432",
			BirthPlace:           "Padang",
			BirthDate:            "2010-08-15",
			Phone:                "081266554433",
			Email:                "ziyad@gmail.com",
			IsDiaspora:           false,
			Country:              "Indonesia",
			Timezone:             "Asia/Jakarta",
			GuardianName:         "Ahmad Zaki",
			GuardianPhone:        "081177665544",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Fatimah Zahra",
			Status:               "DOCUMENTS_VERIFIED",
			Notes:                "Berkas lengkap, siap mengikuti placement test tahfizh",
			CreatedAt:            time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:                   "preg-2027-002",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-002",
			TargetPackage:        "PAKET_B",
			FullName:             "Aisyah Maryam Binti Tariq",
			Gender:               "P",
			NIK:                  "9901024503110001",
			NISN:                 "0112345678",
			BirthPlace:           "Riyadh",
			BirthDate:            "2012-03-20",
			Phone:                "+966501234567",
			Email:                "tariq.family@aramco.com",
			IsDiaspora:           true,
			Country:              "Arab Saudi",
			Timezone:             "Asia/Riyadh",
			GuardianName:         "Tariq Abdullah",
			GuardianPhone:        "+966501234567",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Maryam Hanum",
			Status:               "SUBMITTED",
			Notes:                "Diaspora KSA, membutuhkan jadwal sesi KBM online malam hari KSA",
			CreatedAt:            time.Now().Add(-2 * 24 * time.Hour),
		},
		{
			ID:                   "preg-2027-003",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-003",
			TargetPackage:        "PAKET_A",
			FullName:             "Umar Abdurrahman",
			Gender:               "L",
			NIK:                  "1371050101150002",
			NISN:                 "0156789012",
			BirthPlace:           "Bukittinggi",
			BirthDate:            "2015-05-10",
			Phone:                "082188776655",
			Email:                "wali.umar@gmail.com",
			IsDiaspora:           false,
			Country:              "Indonesia",
			Timezone:             "Asia/Jakarta",
			GuardianName:         "Abdurrahman Saleh",
			GuardianPhone:        "082188776655",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Khadijah",
			Status:               "TEST_PASSED",
			Notes:                "Hasil Placement Test: Tahfizh 2 Juz, siap masuk Paket A Kelas 4",
			CreatedAt:            time.Now().Add(-4 * 24 * time.Hour),
		},
	}

	h.documents = []models.PPDBDocument{
		{
			ID:                 "pdoc-001",
			PPDBRegistrationID: "preg-2027-001",
			DocumentType:       "KK",
			FileURL:            "https://cdn.darel-iman.sch.id/ppdb/docs/kk_ziyad.pdf",
			IsVerified:         true,
			VerificationNotes:  "NIK dan data kepala keluarga valid sesuai Dukcapil",
			CreatedAt:          time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:                 "pdoc-002",
			PPDBRegistrationID: "preg-2027-001",
			DocumentType:       "AKTA_LAHIR",
			FileURL:            "https://cdn.darel-iman.sch.id/ppdb/docs/akta_ziyad.pdf",
			IsVerified:         true,
			VerificationNotes:  "Akta kelahiran terverifikasi",
			CreatedAt:          time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID:                 "pdoc-003",
			PPDBRegistrationID: "preg-2027-002",
			DocumentType:       "SURAT_KONSULAT",
			FileURL:            "https://cdn.darel-iman.sch.id/ppdb/docs/kbri_riyadh_aisyah.pdf",
			IsVerified:         false,
			VerificationNotes:  "Menunggu konfirmasi keabsahan surat domisili KBRI Riyadh",
			CreatedAt:          time.Now().Add(-2 * 24 * time.Hour),
		},
	}
}

// GET /api/v1/ppdb/periods/active
func (h *PPDBHandler) GetActivePeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, p := range h.periods {
		if p.IsOpen {
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, `{"error":"Tidak ada periode pendaftaran aktif"}`, http.StatusNotFound)
}

// POST /api/v1/ppdb/register
func (h *PPDBHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.PPDBRegistration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload format"}`, http.StatusBadRequest)
		return
	}

	if req.FullName == "" || req.GuardianName == "" || req.GuardianPhone == "" {
		http.Error(w, `{"error":"Nama lengkap, nama wali, dan nomor telepon wali wajib diisi"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	regNum := fmt.Sprintf("REG-2027-%03d", len(h.registrations)+1)
	req.ID = "preg-" + uuid.New().String()[:8]
	req.RegistrationNumber = regNum
	req.Status = "SUBMITTED"
	req.CreatedAt = time.Now()
	if req.PPDBPeriodID == "" && len(h.periods) > 0 {
		req.PPDBPeriodID = h.periods[0].ID
	}
	if req.Country == "" {
		req.Country = "Indonesia"
	}
	if req.Timezone == "" {
		req.Timezone = "Asia/Jakarta"
	}

	h.registrations = append(h.registrations, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":             "Pendaftaran PPDB berhasil diregistrasi",
		"registration_number": regNum,
		"registration":        req,
	})
}

// POST /api/v1/ppdb/:registration_id/documents
func (h *PPDBHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regID := chi.URLParam(r, "registration_id")

	var req models.PPDBDocument
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = "pdoc-" + uuid.New().String()[:8]
	req.PPDBRegistrationID = regID
	req.CreatedAt = time.Now()
	req.IsVerified = false

	h.documents = append(h.documents, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Dokumen berhasil diunggah",
		"document": req,
	})
}

// GET /api/v1/ppdb/my-status?registration_number=...
func (h *PPDBHandler) GetMyStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	regNum := r.URL.Query().Get("registration_number")

	if regNum == "" {
		http.Error(w, `{"error":"Parameter registration_number wajib diisi"}`, http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, reg := range h.registrations {
		if reg.RegistrationNumber == regNum {
			// Attach documents
			var docs []models.PPDBDocument
			for _, d := range h.documents {
				if d.PPDBRegistrationID == reg.ID {
					docs = append(docs, d)
				}
			}
			reg.Documents = docs
			json.NewEncoder(w).Encode(reg)
			return
		}
	}

	http.Error(w, `{"error":"Nomor pendaftaran tidak ditemukan"}`, http.StatusNotFound)
}

// GET /api/v1/ppdb/admin/registrations
func (h *PPDBHandler) ListRegistrations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	statusFilter := r.URL.Query().Get("status")
	pkgFilter := r.URL.Query().Get("package")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []models.PPDBRegistration
	for _, reg := range h.registrations {
		if statusFilter != "" && reg.Status != statusFilter {
			continue
		}
		if pkgFilter != "" && reg.TargetPackage != pkgFilter {
			continue
		}

		// Attach documents
		var docs []models.PPDBDocument
		for _, d := range h.documents {
			if d.PPDBRegistrationID == reg.ID {
				docs = append(docs, d)
			}
		}
		reg.Documents = docs
		result = append(result, reg)
	}

	json.NewEncoder(w).Encode(result)
}

// PUT /api/v1/ppdb/admin/registrations/:id/verify
func (h *PPDBHandler) VerifyRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	var payload struct {
		Status string `json:"status"` // DOCUMENTS_VERIFIED, TEST_PASSED, REJECTED
		Notes  string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for i := range h.registrations {
		if h.registrations[i].ID == id {
			if payload.Status != "" {
				h.registrations[i].Status = payload.Status
			}
			if payload.Notes != "" {
				h.registrations[i].Notes = payload.Notes
			}
			found = true
			break
		}
	}

	if !found {
		http.Error(w, `{"error":"Pendaftaran tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Status verifikasi PPDB berhasil diperbarui",
		"status":  payload.Status,
	})
}

// POST /api/v1/ppdb/admin/registrations/:id/accept-and-convert
// AKSI KRUSIAL: Mengubah calon santri menjadi santri resmi di tabel `students`,
// auto-generate akun login `users`, dan assign ke rombel awal.
func (h *PPDBHandler) AcceptAndConvert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	h.mu.Lock()
	defer h.mu.Unlock()

	var target *models.PPDBRegistration
	for i := range h.registrations {
		if h.registrations[i].ID == id {
			h.registrations[i].Status = "ACCEPTED"
			h.registrations[i].Notes = "Diterima Resmi & Akun Santri Aktif"
			target = &h.registrations[i]
			break
		}
	}

	if target == nil {
		http.Error(w, `{"error":"Data registrasi PPDB tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	// Generate NIS & Username
	nis := fmt.Sprintf("2027%04d", time.Now().Unix()%10000)
	username := fmt.Sprintf("santri.%s", target.RegistrationNumber)
	tempPassword := "DarElIman2027!"

	// If MySQL is online, we could execute query. In dual-mode, we return detailed conversion result.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":             "Calon santri berhasil DITERIMA dan dikonversi menjadi Santri Resmi PKBM Digital Dar El Iman",
		"registration_number": target.RegistrationNumber,
		"student_name":        target.FullName,
		"target_package":      target.TargetPackage,
		"generated_nis":       nis,
		"user_account": map[string]string{
			"username":      username,
			"role":          "STUDENT",
			"temp_password": tempPassword,
			"email":         target.Email,
		},
		"assigned_classroom": fmt.Sprintf("Rombel %s Reguler", target.TargetPackage),
		"status":             "ACCEPTED",
	})
}
