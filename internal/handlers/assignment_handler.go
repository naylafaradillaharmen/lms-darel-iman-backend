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

type AssignmentHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	// In-memory fallback
	assignments []models.Assignment
	submissions []models.AssignmentSubmission
}

func NewAssignmentHandler(db *sql.DB) *AssignmentHandler {
	asgDocURL := "https://cdn.darel-iman.or.id/tugas/petunjuk-tugas-01.pdf"
	submText := "Bismillah, berikut setoran rekaman audio An-Naba 1-10 dan ringkasan adab saya ustadz: Niat ikhlas semata karena Allah dan berwudhu sebelum membaca Al-Quran."
	submURL := "https://cdn.darel-iman.or.id/submissions/fatih-an-naba-audio.mp3"
	initialScore := 95.00
	initialFeedback := "Masya Allah mumtaz bacaan tajwid sangat baik dan makhorijul huruf jelas. Pertahankan!"
	gradedBy := "t0000001-0000-0000-0000-000000000001"
	gradedByName := "Ustadz Ahmad Al-Hafizh, Lc."
	now := time.Now()

	h := &AssignmentHandler{
		DB: db,
		assignments: []models.Assignment{
			{
				ID:                "asg00001-0000-0000-0000-000000000001",
				CourseTopicID:     "top00001-0000-0000-0000-000000000001",
				Title:             "Tugas 1: Rangkuman Adab & Rekaman Suara Surah An-Naba 1-10",
				Instructions:      "1. Tulis ringkasan 1 paragraf adab Al-Quran.\n2. Rekam bacaan Surah An-Naba 1-10 dan cantumkan link atau lampirkan berkas jawaban.",
				DueDate:           "2026-10-20 23:59:00",
				MaxScore:          100.00,
				FileAttachmentURL: &asgDocURL,
				SubmissionCount:   1,
				CreatedAt:         time.Now(),
			},
		},
		submissions: []models.AssignmentSubmission{
			{
				ID:                "subm0001-0000-0000-0000-000000000001",
				AssignmentID:      "asg00001-0000-0000-0000-000000000001",
				StudentID:         "s0000001-0000-0000-0000-000000000001",
				StudentName:       "Muhammad Al-Fatih",
				SubmissionText:    &submText,
				FileURL:           &submURL,
				SubmittedAt:       time.Now().Add(-2 * time.Hour),
				Score:             &initialScore,
				Feedback:          &initialFeedback,
				GradedByTutorID:   &gradedBy,
				GradedByTutorName: &gradedByName,
				GradedAt:          &now,
			},
		},
	}
	return h
}

// 1. POST /api/v1/assignments/topics/{topic_id}/assignments
type CreateAssignmentRequest struct {
	Title             string   `json:"title"`
	Instructions      string   `json:"instructions"`
	DueDate           string   `json:"due_date"`
	MaxScore          *float64 `json:"max_score"`
	FileAttachmentURL *string  `json:"file_attachment_url"`
}

func (h *AssignmentHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topicID := chi.URLParam(r, "topic_id")

	var req CreateAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Judul dan instruksi tugas wajib diisi"})
		return
	}

	maxScore := 100.0
	if req.MaxScore != nil && *req.MaxScore > 0 {
		maxScore = *req.MaxScore
	}

	newAsg := models.Assignment{
		ID:                uuid.New().String(),
		CourseTopicID:     topicID,
		Title:             req.Title,
		Instructions:      req.Instructions,
		DueDate:           req.DueDate,
		MaxScore:          maxScore,
		FileAttachmentURL: req.FileAttachmentURL,
		SubmissionCount:   0,
		CreatedAt:         time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO assignments (id, course_topic_id, title, instructions, due_date, max_score, file_attachment_url) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			newAsg.ID, newAsg.CourseTopicID, newAsg.Title, newAsg.Instructions, newAsg.DueDate, newAsg.MaxScore, newAsg.FileAttachmentURL)
	}

	h.mu.Lock()
	h.assignments = append(h.assignments, newAsg)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Instruksi tugas berhasil dibuat",
		"data":    newAsg,
	})
}

// 2. GET /api/v1/assignments/{assignment_id}
func (h *AssignmentHandler) GetAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	asgID := chi.URLParam(r, "assignment_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, a := range h.assignments {
		if a.ID == asgID {
			var subs []models.AssignmentSubmission
			for _, s := range h.submissions {
				if s.AssignmentID == asgID {
					subs = append(subs, s)
				}
			}
			a.Submissions = subs
			a.SubmissionCount = len(subs)
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": a})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Tugas tidak ditemukan"})
}

// 3. POST /api/v1/assignments/{assignment_id}/submit
type SubmitAssignmentRequest struct {
	StudentID      string  `json:"student_id"`
	StudentName    string  `json:"student_name"`
	SubmissionText *string `json:"submission_text"`
	FileURL        *string `json:"file_url"`
}

func (h *AssignmentHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	asgID := chi.URLParam(r, "assignment_id")

	var req SubmitAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data submit tidak valid"})
		return
	}

	if req.StudentID == "" {
		req.StudentID = "s0000001-0000-0000-0000-000000000001"
	}
	if req.StudentName == "" {
		req.StudentName = "Muhammad Al-Fatih"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Cek apakah sudah pernah submit (update jika ada)
	for i, s := range h.submissions {
		if s.AssignmentID == asgID && s.StudentID == req.StudentID {
			h.submissions[i].SubmissionText = req.SubmissionText
			h.submissions[i].FileURL = req.FileURL
			h.submissions[i].SubmittedAt = time.Now()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Jawaban tugas berhasil diperbarui",
				"data":    h.submissions[i],
			})
			return
		}
	}

	newSub := models.AssignmentSubmission{
		ID:             uuid.New().String(),
		AssignmentID:   asgID,
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		SubmissionText: req.SubmissionText,
		FileURL:        req.FileURL,
		SubmittedAt:    time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO assignment_submissions (id, assignment_id, student_id, submission_text, file_url) VALUES (?, ?, ?, ?, ?)`,
			newSub.ID, newSub.AssignmentID, newSub.StudentID, newSub.SubmissionText, newSub.FileURL)
	}

	h.submissions = append(h.submissions, newSub)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tugas berhasil dikumpulkan",
		"data":    newSub,
	})
}

// 4. GET /api/v1/assignments/{assignment_id}/submissions
func (h *AssignmentHandler) GetSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	asgID := chi.URLParam(r, "assignment_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []models.AssignmentSubmission
	for _, s := range h.submissions {
		if s.AssignmentID == asgID || asgID == "" || asgID == "all" {
			result = append(result, s)
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// 5. PUT /api/v1/assignments/submissions/{submission_id}/grade
type GradeSubmissionRequest struct {
	Score             float64 `json:"score"`
	Feedback          string  `json:"feedback"`
	GradedByTutorID   *string `json:"graded_by_tutor_id"`
	GradedByTutorName *string `json:"graded_by_tutor_name"`
}

func (h *AssignmentHandler) GradeSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	submID := chi.URLParam(r, "submission_id")

	var req GradeSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format nilai tidak valid"})
		return
	}

	tutorName := "Ustadz Ahmad Al-Hafizh, Lc."
	if req.GradedByTutorName != nil && *req.GradedByTutorName != "" {
		tutorName = *req.GradedByTutorName
	}
	tutorID := "t0000001-0000-0000-0000-000000000001"
	if req.GradedByTutorID != nil && *req.GradedByTutorID != "" {
		tutorID = *req.GradedByTutorID
	}
	now := time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, s := range h.submissions {
		if s.ID == submID {
			h.submissions[i].Score = &req.Score
			h.submissions[i].Feedback = &req.Feedback
			h.submissions[i].GradedByTutorID = &tutorID
			h.submissions[i].GradedByTutorName = &tutorName
			h.submissions[i].GradedAt = &now

			if h.DB != nil {
				_, _ = h.DB.Exec(`UPDATE assignment_submissions SET score = ?, feedback = ?, graded_by_tutor_id = ?, graded_at = NOW() WHERE id = ?`,
					req.Score, req.Feedback, tutorID, submID)
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Nilai dan evaluasi berhasil disimpan",
				"data":    h.submissions[i],
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Pengumpulan tugas tidak ditemukan"})
}
