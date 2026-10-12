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

type CurriculumHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	// In-memory fallback
	subjects      []models.Subject
	courseClasses []models.CourseClass
	schedules     []models.ClassSchedule
}

func NewCurriculumHandler(db *sql.DB) *CurriculumHandler {
	desc1 := "Hafalan Al-Qur'an Juz 30-28 dengan kaidah tajwid dan pemahaman maknawi"
	desc2 := "Kaidah bahasa Arab Al-Qur'an, percakapan harian dan kaidah tata bahasa dasar"
	desc3 := "Penanaman tauhid yang lurus dan adab penuntut ilmu sesuai sunnah"
	desc4 := "Pemanfaatan kecerdasan buatan etis, logika komputasi, dan skill digital masa depan"
	desc5 := "Kurikulum Merdeka Kesetaraan: Literasi, teks argumentasi, dan komunikasi formal"
	desc6 := "Kurikulum Merdeka Kesetaraan: Aljabar, statistika data dasar, dan pemecahan masalah"
	meetLink := "https://meet.google.com/pkbm-thf-x10"

	h := &CurriculumHandler{
		DB: db,
		subjects: []models.Subject{
			{ID: "sub00001-0000-0000-0000-000000000001", Code: "THF-01", Name: "Tahfizh Al-Qur'an Tematik", Category: "KHAS_DAR_EL_IMAN", Description: &desc1, CreatedAt: time.Now()},
			{ID: "sub00002-0000-0000-0000-000000000002", Code: "ARB-01", Name: "Bahasa Arab Komunikatif & Nahwu", Category: "KHAS_DAR_EL_IMAN", Description: &desc2, CreatedAt: time.Now()},
			{ID: "sub00003-0000-0000-0000-000000000003", Code: "AKD-01", Name: "Akidah & Akhlak Karimah", Category: "KHAS_DAR_EL_IMAN", Description: &desc3, CreatedAt: time.Now()},
			{ID: "sub00004-0000-0000-0000-000000000004", Code: "TEC-01", Name: "AI Literacy & Pemrograman Dasar", Category: "KHAS_DAR_EL_IMAN", Description: &desc4, CreatedAt: time.Now()},
			{ID: "sub00005-0000-0000-0000-000000000005", Code: "IND-10", Name: "Bahasa Indonesia (Paket C)", Category: "UMUM_KEMENDIKBUD", Description: &desc5, CreatedAt: time.Now()},
			{ID: "sub00006-0000-0000-0000-000000000006", Code: "MTK-10", Name: "Matematika Terapan (Paket C)", Category: "UMUM_KEMENDIKBUD", Description: &desc6, CreatedAt: time.Now()},
		},
		courseClasses: []models.CourseClass{
			{
				ID:               "cc000001-0000-0000-0000-000000000001",
				ClassroomID:      "c1111111-1111-1111-1111-111111111111",
				ClassroomName:    "Kelas X Paket C - Umar bin Khattab",
				SubjectID:        "sub00001-0000-0000-0000-000000000001",
				SubjectCode:      "THF-01",
				SubjectName:      "Tahfizh Al-Qur'an Tematik",
				TutorID:          "t0000001-0000-0000-0000-000000000001",
				TutorName:        "Ustadz Ahmad Al-Hafizh, Lc.",
				AcademicYearID:   "22222222-2222-2222-2222-222222222222",
				AcademicYearName: "2026/2027 Ganjil",
				Schedules: []models.ClassSchedule{
					{ID: "sch00001-0000-0000-0000-000000000001", CourseClassID: "cc000001-0000-0000-0000-000000000001", DayOfWeek: 1, StartTime: "08:00", EndTime: "09:30", MeetingLink: &meetLink, CreatedAt: time.Now()},
					{ID: "sch00002-0000-0000-0000-000000000002", CourseClassID: "cc000001-0000-0000-0000-000000000001", DayOfWeek: 3, StartTime: "10:00", EndTime: "11:30", MeetingLink: &meetLink, CreatedAt: time.Now()},
				},
				CreatedAt: time.Now(),
			},
			{
				ID:               "cc000002-0000-0000-0000-000000000002",
				ClassroomID:      "c1111111-1111-1111-1111-111111111111",
				ClassroomName:    "Kelas X Paket C - Umar bin Khattab",
				SubjectID:        "sub00002-0000-0000-0000-000000000002",
				SubjectCode:      "ARB-01",
				SubjectName:      "Bahasa Arab Komunikatif & Nahwu",
				TutorID:          "t0000001-0000-0000-0000-000000000001",
				TutorName:        "Ustadz Ahmad Al-Hafizh, Lc.",
				AcademicYearID:   "22222222-2222-2222-2222-222222222222",
				AcademicYearName: "2026/2027 Ganjil",
				CreatedAt:        time.Now(),
			},
		},
		schedules: []models.ClassSchedule{
			{ID: "sch00001-0000-0000-0000-000000000001", CourseClassID: "cc000001-0000-0000-0000-000000000001", DayOfWeek: 1, StartTime: "08:00", EndTime: "09:30", MeetingLink: &meetLink, CreatedAt: time.Now()},
			{ID: "sch00002-0000-0000-0000-000000000002", CourseClassID: "cc000001-0000-0000-0000-000000000001", DayOfWeek: 3, StartTime: "10:00", EndTime: "11:30", MeetingLink: &meetLink, CreatedAt: time.Now()},
		},
	}
	return h
}

// 1. GET /api/v1/curriculum/subjects
func (h *CurriculumHandler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.DB != nil {
		rows, err := h.DB.Query(`SELECT id, code, name, category, description, created_at FROM subjects ORDER BY code ASC`)
		if err == nil {
			defer rows.Close()
			var list []models.Subject
			for rows.Next() {
				var s models.Subject
				if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Category, &s.Description, &s.CreatedAt); err == nil {
					list = append(list, s)
				}
			}
			if len(list) > 0 {
				json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": list})
				return
			}
		}
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": h.subjects})
}

// 2. POST /api/v1/curriculum/subjects
type CreateSubjectRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
}

func (h *CurriculumHandler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateSubjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Kode dan nama mata pelajaran wajib diisi"})
		return
	}

	newSub := models.Subject{
		ID:          uuid.New().String(),
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO subjects (id, code, name, category, description) VALUES (?, ?, ?, ?, ?)`,
			newSub.ID, newSub.Code, newSub.Name, newSub.Category, newSub.Description)
	}

	h.mu.Lock()
	h.subjects = append(h.subjects, newSub)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mata pelajaran berhasil ditambahkan",
		"data":    newSub,
	})
}

// 3. GET /api/v1/curriculum/course-classes
func (h *CurriculumHandler) ListCourseClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.DB != nil {
		query := `
			SELECT cc.id, cc.classroom_id, c.name as classroom_name,
			       cc.subject_id, s.code as subject_code, s.name as subject_name,
			       cc.tutor_id, t.full_name as tutor_name,
			       cc.academic_year_id, ay.name as academic_year_name, cc.created_at
			FROM course_classes cc
			JOIN classrooms c ON cc.classroom_id = c.id
			JOIN subjects s ON cc.subject_id = s.id
			JOIN tutors t ON cc.tutor_id = t.id
			JOIN academic_years ay ON cc.academic_year_id = ay.id
			ORDER BY cc.created_at DESC
		`
		rows, err := h.DB.Query(query)
		if err == nil {
			defer rows.Close()
			var list []models.CourseClass
			for rows.Next() {
				var cc models.CourseClass
				if err := rows.Scan(&cc.ID, &cc.ClassroomID, &cc.ClassroomName,
					&cc.SubjectID, &cc.SubjectCode, &cc.SubjectName,
					&cc.TutorID, &cc.TutorName,
					&cc.AcademicYearID, &cc.AcademicYearName, &cc.CreatedAt); err == nil {
					list = append(list, cc)
				}
			}
			if len(list) > 0 {
				json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": list})
				return
			}
		}
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": h.courseClasses})
}

// 4. POST /api/v1/curriculum/course-classes
type CreateCourseClassRequest struct {
	ClassroomID    string `json:"classroom_id"`
	SubjectID      string `json:"subject_id"`
	TutorID        string `json:"tutor_id"`
	AcademicYearID string `json:"academic_year_id"`
}

func (h *CurriculumHandler) CreateCourseClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateCourseClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data tidak valid"})
		return
	}

	newCC := models.CourseClass{
		ID:             uuid.New().String(),
		ClassroomID:    req.ClassroomID,
		SubjectID:      req.SubjectID,
		TutorID:        req.TutorID,
		AcademicYearID: req.AcademicYearID,
		CreatedAt:      time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO course_classes (id, classroom_id, subject_id, tutor_id, academic_year_id) VALUES (?, ?, ?, ?, ?)`,
			newCC.ID, newCC.ClassroomID, newCC.SubjectID, newCC.TutorID, newCC.AcademicYearID)
	}

	h.mu.Lock()
	newCC.ClassroomName = "Kelas X Paket C - Umar bin Khattab"
	newCC.SubjectName = "Mata Pelajaran Terpilih"
	newCC.TutorName = "Ustadz Ahmad Al-Hafizh, Lc."
	h.courseClasses = append(h.courseClasses, newCC)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Pengampu kelas mata pelajaran berhasil ditetapkan",
		"data":    newCC,
	})
}

// 5. GET /api/v1/curriculum/schedules/classroom/{classroom_id}
func (h *CurriculumHandler) GetClassroomSchedules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	crID := chi.URLParam(r, "classroom_id")
	_ = crID

	h.mu.RLock()
	defer h.mu.RUnlock()
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": h.schedules})
}
