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

type TahfizhHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	// In-memory fallback
	groups []models.TahfizhGroup
	logs   []models.TahfizhLog
}

func NewTahfizhHandler(db *sql.DB) *TahfizhHandler {
	juz30 := "Juz 30"
	audioURL := "https://cdn.darel-iman.or.id/audio/setoran-fatih-an-naba.mp3"
	note1 := "Hafalan lancar, makhraj ghunnah dan mad thabi'i tepat."
	note2 := "Muraja'ah bagus, perlu diperhatikan waqaf di ayat 15."
	mumtaz := "MUMTAZ"
	jayyidJiddan := "JAYYID_JIDDAN"

	h := &TahfizhHandler{
		DB: db,
		groups: []models.TahfizhGroup{
			{
				ID:               "tg000001-0000-0000-0000-000000000001",
				AcademicYearID:   "22222222-2222-2222-2222-222222222222",
				AcademicYearName: "2026/2027 Ganjil",
				Name:             "Halaqah Imam Nafi' - Ikhwan Kelas X",
				TutorID:          "t0000001-0000-0000-0000-000000000001",
				TutorName:        "Ustadz Ahmad Al-Hafizh, Lc.",
				TargetJuz:        &juz30,
				StudentCount:     1,
				Students: []models.Student{
					{
						ID:       "s0000001-0000-0000-0000-000000000001",
						NIS:      "NIS-2026-001",
						FullName: "Muhammad Al-Fatih",
						Gender:   ptrString("L"),
						Status:   "ACTIVE",
					},
				},
				CreatedAt: time.Now(),
			},
		},
		logs: []models.TahfizhLog{
			{
				ID:             "tl000001-0000-0000-0000-000000000001",
				StudentID:      "s0000001-0000-0000-0000-000000000001",
				StudentName:    "Muhammad Al-Fatih",
				TutorID:        "t0000001-0000-0000-0000-000000000001",
				TutorName:      "Ustadz Ahmad Al-Hafizh, Lc.",
				LogDate:        time.Now().Format("2006-01-02"),
				Type:           "ZIYADAH",
				SurahName:      "An-Naba'",
				StartAyah:      1,
				EndAyah:        40,
				FluencyScore:   &mumtaz,
				TajweedScore:   &mumtaz,
				AudioRecordURL: &audioURL,
				Notes:          &note1,
				CreatedAt:      time.Now(),
			},
			{
				ID:             "tl000002-0000-0000-0000-000000000002",
				StudentID:      "s0000001-0000-0000-0000-000000000001",
				StudentName:    "Muhammad Al-Fatih",
				TutorID:        "t0000001-0000-0000-0000-000000000001",
				TutorName:      "Ustadz Ahmad Al-Hafizh, Lc.",
				LogDate:        time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
				Type:           "MURAJAAH",
				SurahName:      "An-Nazi'at",
				StartAyah:      1,
				EndAyah:        20,
				FluencyScore:   &jayyidJiddan,
				TajweedScore:   &mumtaz,
				Notes:          &note2,
				CreatedAt:      time.Now().AddDate(0, 0, -1),
			},
		},
	}
	return h
}

func ptrString(s string) *string {
	return &s
}

// 1. GET /api/v1/tahfizh/groups
func (h *TahfizhHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	h.mu.RLock()
	defer h.mu.RUnlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    h.groups,
	})
}

// 2. POST /api/v1/tahfizh/groups
type CreateTahfizhGroupRequest struct {
	Name           string  `json:"name"`
	TutorID        string  `json:"tutor_id"`
	AcademicYearID string  `json:"academic_year_id"`
	TargetJuz      *string `json:"target_juz"`
}

func (h *TahfizhHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateTahfizhGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Nama halaqah wajib diisi"})
		return
	}

	if req.TutorID == "" {
		req.TutorID = "t0000001-0000-0000-0000-000000000001"
	}
	if req.AcademicYearID == "" {
		req.AcademicYearID = "22222222-2222-2222-2222-222222222222"
	}

	newGroup := models.TahfizhGroup{
		ID:               uuid.New().String(),
		AcademicYearID:   req.AcademicYearID,
		AcademicYearName: "2026/2027 Ganjil",
		Name:             req.Name,
		TutorID:          req.TutorID,
		TutorName:        "Ustadz Ahmad Al-Hafizh, Lc.",
		TargetJuz:        req.TargetJuz,
		StudentCount:     0,
		Students:         []models.Student{},
		CreatedAt:        time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO tahfizh_groups (id, academic_year_id, name, tutor_id, target_juz) VALUES (?, ?, ?, ?, ?)`,
			newGroup.ID, newGroup.AcademicYearID, newGroup.Name, newGroup.TutorID, newGroup.TargetJuz)
	}

	h.mu.Lock()
	h.groups = append(h.groups, newGroup)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Kelompok halaqah tahfizh berhasil dibentuk",
		"data":    newGroup,
	})
}

// 3. POST /api/v1/tahfizh/groups/{id}/students
type AddStudentToGroupRequest struct {
	StudentID string `json:"student_id"`
}

func (h *TahfizhHandler) AddStudentToGroup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groupID := chi.URLParam(r, "id")

	var req AddStudentToGroupRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.StudentID == "" {
		req.StudentID = "s0000001-0000-0000-0000-000000000001"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, g := range h.groups {
		if g.ID == groupID {
			h.groups[i].StudentCount++
			h.groups[i].Students = append(h.groups[i].Students, models.Student{
				ID:       req.StudentID,
				FullName: "Muhammad Al-Fatih",
				NIS:      "NIS-2026-001",
				Gender:   ptrString("L"),
			})
			break
		}
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO tahfizh_group_students (tahfizh_group_id, student_id) VALUES (?, ?)`, groupID, req.StudentID)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Santri berhasil dimasukkan ke halaqah",
	})
}

// 4. GET /api/v1/tahfizh/students/{student_id}/logs
func (h *TahfizhHandler) GetStudentLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stuID := chi.URLParam(r, "student_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []models.TahfizhLog
	for _, l := range h.logs {
		if l.StudentID == stuID || stuID == "" || stuID == "all" {
			result = append(result, l)
		}
	}

	if len(result) == 0 {
		result = h.logs
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// 5. POST /api/v1/tahfizh/logs
type RecordTahfizhLogRequest struct {
	StudentID      string  `json:"student_id"`
	StudentName    string  `json:"student_name"`
	TutorID        string  `json:"tutor_id"`
	LogDate        string  `json:"log_date"`
	Type           string  `json:"type"` // ZIYADAH, MURAJAAH
	SurahName      string  `json:"surah_name"`
	StartAyah      int     `json:"start_ayah"`
	EndAyah        int     `json:"end_ayah"`
	FluencyScore   *string `json:"fluency_score"`
	TajweedScore   *string `json:"tajweed_score"`
	AudioRecordURL *string `json:"audio_record_url"`
	Notes          *string `json:"notes"`
}

func (h *TahfizhHandler) RecordLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RecordTahfizhLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.SurahName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Nama surah dan rentang ayat wajib diisi"})
		return
	}

	if req.StudentID == "" {
		req.StudentID = "s0000001-0000-0000-0000-000000000001"
	}
	if req.StudentName == "" {
		req.StudentName = "Muhammad Al-Fatih"
	}
	if req.TutorID == "" {
		req.TutorID = "t0000001-0000-0000-0000-000000000001"
	}
	if req.LogDate == "" {
		req.LogDate = time.Now().Format("2006-01-02")
	}
	if req.Type == "" {
		req.Type = "ZIYADAH"
	}
	if req.FluencyScore == nil {
		f := "MUMTAZ"
		req.FluencyScore = &f
	}
	if req.TajweedScore == nil {
		t := "MUMTAZ"
		req.TajweedScore = &t
	}

	newLog := models.TahfizhLog{
		ID:             uuid.New().String(),
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		TutorID:        req.TutorID,
		TutorName:      "Ustadz Ahmad Al-Hafizh, Lc.",
		LogDate:        req.LogDate,
		Type:           req.Type,
		SurahName:      req.SurahName,
		StartAyah:      req.StartAyah,
		EndAyah:        req.EndAyah,
		FluencyScore:   req.FluencyScore,
		TajweedScore:   req.TajweedScore,
		AudioRecordURL: req.AudioRecordURL,
		Notes:          req.Notes,
		CreatedAt:      time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO tahfizh_logs (id, student_id, tutor_id, log_date, type, surah_name, start_ayah, end_ayah, fluency_score, tajweed_score, audio_record_url, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			newLog.ID, newLog.StudentID, newLog.TutorID, newLog.LogDate, newLog.Type, newLog.SurahName, newLog.StartAyah, newLog.EndAyah, newLog.FluencyScore, newLog.TajweedScore, newLog.AudioRecordURL, newLog.Notes)
	}

	h.mu.Lock()
	h.logs = append([]models.TahfizhLog{newLog}, h.logs...)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Setoran hafalan tahfizh berhasil dicatat",
		"data":    newLog,
	})
}
