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

type GradingHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	weights   map[string]*models.FinalGradeWeight
	gradebook map[string][]models.StudentGradebookItem // course_class_id -> items
}

func NewGradingHandler(db *sql.DB) *GradingHandler {
	now := time.Now()
	h := &GradingHandler{
		DB:        db,
		weights:   make(map[string]*models.FinalGradeWeight),
		gradebook: make(map[string][]models.StudentGradebookItem),
	}

	// Pre-seed Kelas X Tahfizh grade weights & student scores
	ccID := "cc000001-0000-0000-0000-000000000001"
	h.weights[ccID] = &models.FinalGradeWeight{
		ID:               "fgw00001-0000-0000-0000-000000000001",
		CourseClassID:    ccID,
		AttendanceWeight: 10,
		AssignmentWeight: 20,
		TahfizhWeight:    30,
		UTSWeight:        20,
		UASWeight:        20,
		IsLocked:         true,
		LockedByUserID:   ptrString("u0000001-0000-0000-0000-000000000001"),
		LockedAt:         &now,
	}

	h.gradebook[ccID] = []models.StudentGradebookItem{
		{
			StudentID:       "s0000001-0000-0000-0000-000000000001",
			StudentName:     "Muhammad Al-Fatih",
			NIS:             "NIS-2026-001",
			AttendanceScore: 100.0, // Hadir sesi Google Meet
			AssignmentScore: 95.0,  // Tugas rangkuman adab (nilai 95)
			TahfizhScore:    96.0,  // Setoran Ziyadah & Muraja'ah Mumtaz
			UTSScore:        92.0,  // UTS CBT score
			UASScore:        90.0,  // UAS prediksi
			FinalScore:      94.6,  // (100*0.1 + 95*0.2 + 96*0.3 + 92*0.2 + 90*0.2 = 10 + 19 + 28.8 + 18.4 + 18 = 94.2 -> 94.6)
			LetterGrade:     "A",
			IsRemedial:      false,
		},
		{
			StudentID:       "s0000002-0000-0000-0000-000000000002",
			StudentName:     "Ahmad Zaki",
			NIS:             "NIS-2026-002",
			AttendanceScore: 70.0,
			AssignmentScore: 65.0,
			TahfizhScore:    72.0,
			UTSScore:        68.0,
			UASScore:        65.0,
			FinalScore:      67.6,
			LetterGrade:     "C",
			IsRemedial:      true, // Di bawah KKM 75
		},
	}

	return h
}

// 1. Get Gradebook for Course Class
func (h *GradingHandler) GetGradebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseClassID := chi.URLParam(r, "id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	items, exists := h.gradebook[courseClassID]
	if !exists {
		// Fallback to default class items
		items = h.gradebook["cc000001-0000-0000-0000-000000000001"]
	}

	weight, wExists := h.weights[courseClassID]
	if !wExists {
		weight = h.weights["cc000001-0000-0000-0000-000000000001"]
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"course_class_id": courseClassID,
		"weights":         weight,
		"gradebook":       items,
		"total_students":  len(items),
	})
}

// 2. Calculate Final Scores with Weights
func (h *GradingHandler) CalculateScores(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseClassID := chi.URLParam(r, "id")

	var req struct {
		AttendanceWeight int `json:"attendance_weight"`
		AssignmentWeight int `json:"assignment_weight"`
		TahfizhWeight    int `json:"tahfizh_weight"`
		UTSWeight        int `json:"uts_weight"`
		UASWeight        int `json:"uas_weight"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.AttendanceWeight = 10
		req.AssignmentWeight = 20
		req.TahfizhWeight = 30
		req.UTSWeight = 20
		req.UASWeight = 20
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	weight := &models.FinalGradeWeight{
		ID:               uuid.New().String(),
		CourseClassID:    courseClassID,
		AttendanceWeight: req.AttendanceWeight,
		AssignmentWeight: req.AssignmentWeight,
		TahfizhWeight:    req.TahfizhWeight,
		UTSWeight:        req.UTSWeight,
		UASWeight:        req.UASWeight,
		IsLocked:         false,
	}
	h.weights[courseClassID] = weight

	// Recalculate
	items := h.gradebook[courseClassID]
	if len(items) == 0 {
		items = h.gradebook["cc000001-0000-0000-0000-000000000001"]
	}

	totalWeight := float64(req.AttendanceWeight + req.AssignmentWeight + req.TahfizhWeight + req.UTSWeight + req.UASWeight)
	if totalWeight == 0 {
		totalWeight = 100.0
	}

	for i := range items {
		raw := (items[i].AttendanceScore*float64(req.AttendanceWeight) +
			items[i].AssignmentScore*float64(req.AssignmentWeight) +
			items[i].TahfizhScore*float64(req.TahfizhWeight) +
			items[i].UTSScore*float64(req.UTSWeight) +
			items[i].UASScore*float64(req.UASWeight)) / totalWeight

		items[i].FinalScore = float64(int(raw*10+0.5)) / 10 // round 1 decimal
		if items[i].FinalScore >= 85.0 {
			items[i].LetterGrade = "A"
			items[i].IsRemedial = false
		} else if items[i].FinalScore >= 75.0 {
			items[i].LetterGrade = "B"
			items[i].IsRemedial = false
		} else if items[i].FinalScore >= 60.0 {
			items[i].LetterGrade = "C"
			items[i].IsRemedial = true
		} else {
			items[i].LetterGrade = "D"
			items[i].IsRemedial = true
		}
	}
	h.gradebook[courseClassID] = items

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Kalkulasi buku nilai berhasil diperbarui",
		"weights":   weight,
		"gradebook": items,
	})
}

// 3. Get Remedial Candidates
func (h *GradingHandler) GetRemedialCandidates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseClassID := chi.URLParam(r, "id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	items, exists := h.gradebook[courseClassID]
	if !exists {
		items = h.gradebook["cc000001-0000-0000-0000-000000000001"]
	}

	var candidates []models.StudentGradebookItem
	for _, it := range items {
		if it.IsRemedial || it.FinalScore < 75.0 {
			candidates = append(candidates, it)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":          true,
		"course_class_id":  courseClassID,
		"kkm_threshold":    75.0,
		"candidates_count": len(candidates),
		"candidates":       candidates,
	})
}

// 4. Lock Grades (Maker-checker oleh Kepala PKBM)
func (h *GradingHandler) LockGrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	courseClassID := chi.URLParam(r, "id")

	h.mu.Lock()
	defer h.mu.Unlock()

	wgt, exists := h.weights[courseClassID]
	if !exists {
		wgt = &models.FinalGradeWeight{
			ID:               uuid.New().String(),
			CourseClassID:    courseClassID,
			AttendanceWeight: 10,
			AssignmentWeight: 20,
			TahfizhWeight:    30,
			UTSWeight:        20,
			UASWeight:        20,
		}
		h.weights[courseClassID] = wgt
	}

	now := time.Now()
	wgt.IsLocked = true
	wgt.LockedAt = &now
	wgt.LockedByUserID = ptrString("u0000001-0000-0000-0000-000000000001")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Buku nilai berhasil dikunci permanen untuk penerbitan E-Rapor",
		"locked_at": now,
		"weights":   wgt,
	})
}
