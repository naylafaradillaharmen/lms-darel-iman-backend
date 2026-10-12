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

type FinanceHandler struct {
	DB           *sql.DB
	mu           sync.RWMutex
	tariffs      []models.TuitionTariff
	invoices     []models.StudentInvoice
	transactions []models.PaymentTransaction
	scholarships []models.ScholarshipGrant
}

func NewFinanceHandler(db *sql.DB) *FinanceHandler {
	h := &FinanceHandler{
		DB: db,
	}
	h.initMockData()
	return h
}

func (h *FinanceHandler) initMockData() {
	h.tariffs = []models.TuitionTariff{
		{
			ID:             "ttf-paket-a",
			AcademicYearID: "ay-2026-2027",
			GradeLevelID:   "gl-paket-a-4",
			PackageName:    "Paket A (Setara SD)",
			Amount:         250000.0,
			Description:    "SPP Bulanan Paket A Pendidikan Kesetaraan SD",
			CreatedAt:      time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			ID:             "ttf-paket-b",
			AcademicYearID: "ay-2026-2027",
			GradeLevelID:   "gl-paket-b-7",
			PackageName:    "Paket B (Setara SMP)",
			Amount:         300000.0,
			Description:    "SPP Bulanan Paket B Pendidikan Kesetaraan SMP",
			CreatedAt:      time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			ID:             "ttf-paket-c",
			AcademicYearID: "ay-2026-2027",
			GradeLevelID:   "gl-paket-c-10",
			PackageName:    "Paket C (Setara SMA)",
			Amount:         350000.0,
			Description:    "SPP Bulanan Paket C Pendidikan Kesetaraan SMA",
			CreatedAt:      time.Now().Add(-60 * 24 * time.Hour),
		},
	}

	h.scholarships = []models.ScholarshipGrant{
		{
			ID:                 "sch-001",
			StudentID:          "stu-002",
			StudentName:        "Aisyah Humaira",
			ProgramName:        "Beasiswa Tahfizh Berprestasi Dar El Iman",
			DiscountPercentage: 100,
			DonorName:          "Yayasan Dar El Iman Indonesia",
			IsActive:           true,
			CreatedAt:          time.Now().Add(-90 * 24 * time.Hour),
		},
	}

	h.invoices = []models.StudentInvoice{
		{
			ID:             "inv-2026-10-01",
			StudentID:      "stu-001",
			StudentName:    "Fatih Al-Ayyubi",
			NIS:            "202607001",
			ClassName:      "Paket C - Kelas 10 Abu Bakar Ash-Shiddiq",
			InvoiceNumber:  "INV-2026-10-0045",
			Month:          10,
			Year:           2026,
			Amount:         350000.0,
			DiscountAmount: 0.0,
			FinalAmount:    350000.0,
			DueDate:        "2026-10-15",
			Status:         "PENDING_VERIFICATION",
			CreatedAt:      time.Now().Add(-9 * 24 * time.Hour),
		},
		{
			ID:             "inv-2026-09-01",
			StudentID:      "stu-001",
			StudentName:    "Fatih Al-Ayyubi",
			NIS:            "202607001",
			ClassName:      "Paket C - Kelas 10 Abu Bakar Ash-Shiddiq",
			InvoiceNumber:  "INV-2026-09-0012",
			Month:          9,
			Year:           2026,
			Amount:         350000.0,
			DiscountAmount: 0.0,
			FinalAmount:    350000.0,
			DueDate:        "2026-09-15",
			Status:         "PAID",
			CreatedAt:      time.Now().Add(-40 * 24 * time.Hour),
		},
		{
			ID:             "inv-2026-10-02",
			StudentID:      "stu-002",
			StudentName:    "Aisyah Humaira",
			NIS:            "202607002",
			ClassName:      "Paket C - Kelas 10 Khadijah",
			InvoiceNumber:  "INV-2026-10-0046",
			Month:          10,
			Year:           2026,
			Amount:         350000.0,
			DiscountAmount: 350000.0,
			FinalAmount:    0.0,
			DueDate:        "2026-10-15",
			Status:         "PAID",
			CreatedAt:      time.Now().Add(-9 * 24 * time.Hour),
		},
	}

	h.transactions = []models.PaymentTransaction{
		{
			ID:               "trx-2026-10-01",
			StudentInvoiceID: "inv-2026-10-01",
			InvoiceNumber:    "INV-2026-10-0045",
			StudentName:      "Fatih Al-Ayyubi",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       350000.0,
			ProofFileURL:     "https://cdn.darel-iman.sch.id/spp/bukti_transfer_fatih_okt2026.jpg",
			PaymentDate:      time.Now().Add(-1 * 24 * time.Hour),
			Status:           "PENDING",
		},
	}
}

// GET /api/v1/finance/overview
func (h *FinanceHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	var totalBilled float64
	var totalPaid float64
	var totalPending float64
	var totalUnpaid float64

	for _, inv := range h.invoices {
		totalBilled += inv.FinalAmount
		switch inv.Status {
		case "PAID":
			totalPaid += inv.FinalAmount
		case "PENDING_VERIFICATION":
			totalPending += inv.FinalAmount
		case "UNPAID":
			totalUnpaid += inv.FinalAmount
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_invoices":     len(h.invoices),
		"total_billed":       totalBilled,
		"total_paid":         totalPaid,
		"total_pending":      totalPending,
		"total_unpaid":       totalUnpaid,
		"tariffs":            h.tariffs,
		"active_scholarships": len(h.scholarships),
	})
}

// GET /api/v1/finance/tariffs
func (h *FinanceHandler) ListTariffs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.tariffs)
}

// POST /api/v1/finance/invoices/generate-monthly
func (h *FinanceHandler) GenerateMonthlyInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Month < 1 || payload.Month > 12 {
		payload.Month = int(time.Now().Month())
		payload.Year = time.Now().Year()
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// In this simulation, generate batch invoices for active students
	studentsList := []struct {
		ID        string
		Name      string
		NIS       string
		ClassName string
		Tariff    float64
	}{
		{"stu-001", "Fatih Al-Ayyubi", "202607001", "Paket C - Kelas 10 Abu Bakar", 350000.0},
		{"stu-002", "Aisyah Humaira", "202607002", "Paket C - Kelas 10 Khadijah", 350000.0},
		{"stu-003", "Zaid bin Haritsah", "202607003", "Paket B - Kelas 8", 300000.0},
		{"stu-004", "Hamzah bin Abdul Muthalib", "202607004", "Paket A - Kelas 5", 250000.0},
	}

	createdCount := 0
	for _, stu := range studentsList {
		// Check if invoice already exists for month/year
		exists := false
		for _, inv := range h.invoices {
			if inv.StudentID == stu.ID && inv.Month == payload.Month && inv.Year == payload.Year {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		// Calculate discount if has scholarship
		discount := 0.0
		for _, sch := range h.scholarships {
			if sch.StudentID == stu.ID && sch.IsActive {
				discount = stu.Tariff * float64(sch.DiscountPercentage) / 100.0
				break
			}
		}

		finalAmt := stu.Tariff - discount
		status := "UNPAID"
		if finalAmt <= 0 {
			status = "PAID"
		}

		inv := models.StudentInvoice{
			ID:             fmt.Sprintf("inv-%d-%02d-%s", payload.Year, payload.Month, stu.ID),
			StudentID:      stu.ID,
			StudentName:    stu.Name,
			NIS:            stu.NIS,
			ClassName:      stu.ClassName,
			InvoiceNumber:  fmt.Sprintf("INV-%d-%02d-%04d", payload.Year, payload.Month, len(h.invoices)+1),
			Month:          payload.Month,
			Year:           payload.Year,
			Amount:         stu.Tariff,
			DiscountAmount: discount,
			FinalAmount:    finalAmt,
			DueDate:        fmt.Sprintf("%d-%02d-15", payload.Year, payload.Month),
			Status:         status,
			CreatedAt:      time.Now(),
		}

		h.invoices = append(h.invoices, inv)
		createdCount++
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         fmt.Sprintf("Penerbitan tagihan SPP massal bulan %d/%d berhasil", payload.Month, payload.Year),
		"created_count":   createdCount,
		"month":           payload.Month,
		"year":            payload.Year,
		"total_invoices":  len(h.invoices),
	})
}

// GET /api/v1/finance/invoices/student/:student_id
func (h *FinanceHandler) GetStudentInvoices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := chi.URLParam(r, "student_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var results []models.StudentInvoice
	for _, inv := range h.invoices {
		if studentID == "all" || inv.StudentID == studentID {
			// Attach transactions
			var trxs []models.PaymentTransaction
			for _, t := range h.transactions {
				if t.StudentInvoiceID == inv.ID {
					trxs = append(trxs, t)
				}
			}
			inv.Transactions = trxs
			results = append(results, inv)
		}
	}

	json.NewEncoder(w).Encode(results)
}

// POST /api/v1/finance/invoices/:id/pay
func (h *FinanceHandler) PayInvoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	invID := chi.URLParam(r, "id")

	var req struct {
		PaymentMethod string  `json:"payment_method"`
		AmountPaid    float64 `json:"amount_paid"`
		ProofFileURL  string  `json:"proof_file_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var targetInv *models.StudentInvoice
	for i := range h.invoices {
		if h.invoices[i].ID == invID {
			targetInv = &h.invoices[i]
			h.invoices[i].Status = "PENDING_VERIFICATION"
			break
		}
	}

	if targetInv == nil {
		http.Error(w, `{"error":"Invoice tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	trx := models.PaymentTransaction{
		ID:               "trx-" + uuid.New().String()[:8],
		StudentInvoiceID: invID,
		InvoiceNumber:    targetInv.InvoiceNumber,
		StudentName:      targetInv.StudentName,
		PaymentMethod:    req.PaymentMethod,
		AmountPaid:       req.AmountPaid,
		ProofFileURL:     req.ProofFileURL,
		PaymentDate:      time.Now(),
		Status:           "PENDING",
	}

	h.transactions = append(h.transactions, trx)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Konfirmasi pembayaran SPP berhasil dikirim. Menunggu verifikasi bendahara yayasan.",
		"transaction": trx,
	})
}

// GET /api/v1/finance/payments/pending
func (h *FinanceHandler) GetPendingPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	var pendings []models.PaymentTransaction
	for _, t := range h.transactions {
		if t.Status == "PENDING" {
			pendings = append(pendings, t)
		}
	}

	json.NewEncoder(w).Encode(pendings)
}

// PUT /api/v1/finance/payments/:id/verify
// Maker-checker verifikasi bukti transfer pembayaran
func (h *FinanceHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	trxID := chi.URLParam(r, "id")

	var req struct {
		Action          string `json:"action"` // APPROVED or REJECTED
		RejectionReason string `json:"rejection_reason"`
		VerifiedByName  string `json:"verified_by_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var targetTrx *models.PaymentTransaction
	now := time.Now()
	for i := range h.transactions {
		if h.transactions[i].ID == trxID {
			targetTrx = &h.transactions[i]
			targetTrx.Status = req.Action
			targetTrx.RejectionReason = req.RejectionReason
			targetTrx.VerifiedByName = req.VerifiedByName
			if targetTrx.VerifiedByName == "" {
				targetTrx.VerifiedByName = "Bendahara Yayasan"
			}
			targetTrx.VerifiedAt = &now
			break
		}
	}

	if targetTrx == nil {
		http.Error(w, `{"error":"Transaksi tidak ditemukan"}`, http.StatusNotFound)
		return
	}

	// Update corresponding invoice status
	for i := range h.invoices {
		if h.invoices[i].ID == targetTrx.StudentInvoiceID {
			if req.Action == "APPROVED" {
				h.invoices[i].Status = "PAID"
			} else {
				h.invoices[i].Status = "UNPAID"
			}
			break
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Verifikasi transaksi %s berhasil diproses", req.Action),
		"status":  req.Action,
	})
}

// POST /api/v1/finance/scholarships
func (h *FinanceHandler) GrantScholarship(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.ScholarshipGrant
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	if req.StudentID == "" || req.ProgramName == "" {
		http.Error(w, `{"error":"Student ID dan Program Beasiswa wajib diisi"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = "sch-" + uuid.New().String()[:8]
	req.CreatedAt = time.Now()
	req.IsActive = true
	if req.DiscountPercentage == 0 {
		req.DiscountPercentage = 100
	}

	h.scholarships = append(h.scholarships, req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Pemberian beasiswa santri berhasil disimpan",
		"scholarship": req,
	})
}

// GET /api/v1/finance/scholarships
func (h *FinanceHandler) ListScholarships(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(h.scholarships)
}
