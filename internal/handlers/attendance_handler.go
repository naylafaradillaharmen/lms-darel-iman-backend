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

type AttendanceHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	// In-memory fallback
	sessions    []models.ClassSession
	attendances []models.SessionAttendance
}

func NewAttendanceHandler(db *sql.DB) *AttendanceHandler {
	note := "Santri hadir tepat waktu via Google Meet"
	h := &AttendanceHandler{
		DB: db,
		sessions: []models.ClassSession{
			{
				ID:            "sess0001-0000-0000-0000-000000000001",
				CourseClassID: "cc000001-0000-0000-0000-000000000001",
				SessionTitle:  "Sesi 1: Pembukaan & Talaqqi Surah An-Naba",
				SessionDate:   time.Now().Format("2006-01-02"),
				StartTime:     "08:00",
				EndTime:       "09:30",
				CreatedAt:     time.Now(),
			},
		},
		attendances: []models.SessionAttendance{
			{
				ID:             "att00001-0000-0000-0000-000000000001",
				ClassSessionID: "sess0001-0000-0000-0000-000000000001",
				StudentID:      "s0000001-0000-0000-0000-000000000001",
				StudentName:    "Muhammad Al-Fatih",
				Status:         "PRESENT",
				Notes:          &note,
				RecordedAt:     time.Now(),
			},
		},
	}
	return h
}

// 1. POST /api/v1/attendances/courses/{course_class_id}/sessions
type CreateSessionRequest struct {
	SessionTitle string `json:"session_title"`
	SessionDate  string `json:"session_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

func (h *AttendanceHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ccID := chi.URLParam(r, "course_class_id")

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SessionTitle == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Judul sesi kelas wajib diisi"})
		return
	}

	dateStr := req.SessionDate
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	newSess := models.ClassSession{
		ID:            uuid.New().String(),
		CourseClassID: ccID,
		SessionTitle:  req.SessionTitle,
		SessionDate:   dateStr,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		CreatedAt:     time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO class_sessions (id, course_class_id, session_title, session_date, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?)`,
			newSess.ID, newSess.CourseClassID, newSess.SessionTitle, newSess.SessionDate, newSess.StartTime, newSess.EndTime)
	}

	h.mu.Lock()
	h.sessions = append(h.sessions, newSess)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Sesi presensi KBM berhasil dibuka",
		"data":    newSess,
	})
}

// 2. GET /api/v1/attendances/sessions/{session_id}
func (h *AttendanceHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessID := chi.URLParam(r, "session_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, s := range h.sessions {
		if s.ID == sessID || sessID == "latest" {
			var atts []models.SessionAttendance
			for _, a := range h.attendances {
				if a.ClassSessionID == s.ID {
					atts = append(atts, a)
				}
			}
			s.Attendances = atts
			json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": s})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Sesi pertemuan tidak ditemukan"})
}

// 3. POST /api/v1/attendances/sessions/{session_id}/record
type RecordAttendanceItem struct {
	StudentID   string  `json:"student_id"`
	StudentName string  `json:"student_name"`
	Status      string  `json:"status"` // PRESENT, PERMIT, SICK, ABSENT
	Notes       *string `json:"notes"`
}

type RecordAttendanceRequest struct {
	Attendances []RecordAttendanceItem `json:"attendances"`
}

func (h *AttendanceHandler) RecordAttendance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessID := chi.URLParam(r, "session_id")

	var req RecordAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data absensi tidak valid"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, item := range req.Attendances {
		found := false
		for i, a := range h.attendances {
			if a.ClassSessionID == sessID && a.StudentID == item.StudentID {
				h.attendances[i].Status = item.Status
				h.attendances[i].Notes = item.Notes
				h.attendances[i].RecordedAt = time.Now()
				found = true
				break
			}
		}
		if !found {
			newAtt := models.SessionAttendance{
				ID:             uuid.New().String(),
				ClassSessionID: sessID,
				StudentID:      item.StudentID,
				StudentName:    item.StudentName,
				Status:         item.Status,
				Notes:          item.Notes,
				RecordedAt:     time.Now(),
			}
			h.attendances = append(h.attendances, newAtt)
			if h.DB != nil {
				_, _ = h.DB.Exec(`INSERT INTO session_attendances (id, class_session_id, student_id, status, notes) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE status=VALUES(status), notes=VALUES(notes)`,
					newAtt.ID, newAtt.ClassSessionID, newAtt.StudentID, newAtt.Status, newAtt.Notes)
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Data presensi santri berhasil dicatat",
	})
}

// 4. POST /api/v1/attendances/sessions/{session_id}/check-in
type CheckInRequest struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
}

func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessID := chi.URLParam(r, "session_id")

	var req CheckInRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.StudentID == "" {
		req.StudentID = "s0000001-0000-0000-0000-000000000001"
	}
	if req.StudentName == "" {
		req.StudentName = "Muhammad Al-Fatih"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	note := "Check-in mandiri oleh santri via web portal"
	for i, a := range h.attendances {
		if a.ClassSessionID == sessID && a.StudentID == req.StudentID {
			h.attendances[i].Status = "PRESENT"
			h.attendances[i].Notes = &note
			h.attendances[i].RecordedAt = time.Now()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Kehadiran mandiri berhasil dikonfirmasi (HADIR)",
				"data":    h.attendances[i],
			})
			return
		}
	}

	newAtt := models.SessionAttendance{
		ID:             uuid.New().String(),
		ClassSessionID: sessID,
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		Status:         "PRESENT",
		Notes:          &note,
		RecordedAt:     time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO session_attendances (id, class_session_id, student_id, status, notes) VALUES (?, ?, ?, ?, ?)`,
			newAtt.ID, newAtt.ClassSessionID, newAtt.StudentID, newAtt.Status, newAtt.Notes)
	}

	h.attendances = append(h.attendances, newAtt)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Kehadiran mandiri berhasil dikonfirmasi (HADIR)",
		"data":    newAtt,
	})
}
