package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pkbm-backend/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ClassroomHandler struct {
	DB *sql.DB
}

func (h *ClassroomHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT c.id, c.academic_year_id, ay.name as ay_name, c.grade_level_id, gl.name as gl_name,
		       c.name, c.homeroom_tutor_id, t.full_name as tutor_name, c.capacity,
		       (SELECT COUNT(*) FROM classroom_students cs WHERE cs.classroom_id = c.id) as student_count
		FROM classrooms c
		JOIN academic_years ay ON c.academic_year_id = ay.id
		JOIN grade_levels gl ON c.grade_level_id = gl.id
		LEFT JOIN tutors t ON c.homeroom_tutor_id = t.id
		ORDER BY gl.level_number ASC, c.name ASC
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": []models.Classroom{}})
		return
	}
	defer rows.Close()

	var classrooms []models.Classroom
	for rows.Next() {
		var cr models.Classroom
		_ = rows.Scan(&cr.ID, &cr.AcademicYearID, &cr.AcademicYearName, &cr.GradeLevelID, &cr.GradeLevelName,
			&cr.Name, &cr.HomeroomTutorID, &cr.HomeroomTutorName, &cr.Capacity, &cr.StudentCount)
		classrooms = append(classrooms, cr)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": classrooms})
}

type CreateClassroomRequest struct {
	AcademicYearID  string  `json:"academic_year_id"`
	GradeLevelID    string  `json:"grade_level_id"`
	Name            string  `json:"name"`
	HomeroomTutorID *string `json:"homeroom_tutor_id"`
	Capacity        int     `json:"capacity"`
}

func (h *ClassroomHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateClassroomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Permintaan tidak valid"})
		return
	}

	cap := req.Capacity
	if cap <= 0 {
		cap = 30
	}

	crID := uuid.New().String()
	_, err := h.DB.Exec(`
		INSERT INTO classrooms (id, academic_year_id, grade_level_id, name, homeroom_tutor_id, capacity)
		VALUES (?, ?, ?, ?, ?, ?)
	`, crID, req.AcademicYearID, req.GradeLevelID, req.Name, req.HomeroomTutorID, cap)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Gagal membuat rombel: nama kelas mungkin sudah ada"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Rombongan Belajar berhasil dibentuk",
		"data":    map[string]string{"id": crID},
	})
}

func (h *ClassroomHandler) AssignStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	classroomID := chi.URLParam(r, "id")

	var req struct {
		StudentID string `json:"student_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	_, err := h.DB.Exec(`
		INSERT INTO classroom_students (classroom_id, student_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE joined_at = NOW()
	`, classroomID, req.StudentID)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Gagal memasukkan santri ke rombel"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Santri berhasil dimasukkan ke rombel"})
}
