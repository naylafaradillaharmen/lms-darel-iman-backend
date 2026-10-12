package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pkbm-backend/internal/models"
)

type InstitutionHandler struct {
	DB *sql.DB
}

func (h *InstitutionHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var totalStudents, totalTutors, totalClassrooms int

	_ = h.DB.QueryRow("SELECT COUNT(*) FROM students WHERE status = 'ACTIVE'").Scan(&totalStudents)
	_ = h.DB.QueryRow("SELECT COUNT(*) FROM tutors WHERE is_active = TRUE").Scan(&totalTutors)
	_ = h.DB.QueryRow("SELECT COUNT(*) FROM classrooms").Scan(&totalClassrooms)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"institution_name": "PKBM Digital Dar El Iman",
			"academic_year":    "2026/2027 Ganjil",
			"total_students":   totalStudents,
			"total_tutors":     totalTutors,
			"total_classrooms": totalClassrooms,
		},
	})
}

func (h *InstitutionHandler) GetAcademicYears(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.DB.Query("SELECT id, institution_id, name, semester, start_date, end_date, is_active FROM academic_years ORDER BY created_at DESC")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": []models.AcademicYear{}})
		return
	}
	defer rows.Close()

	var list []models.AcademicYear
	for rows.Next() {
		var a models.AcademicYear
		_ = rows.Scan(&a.ID, &a.InstitutionID, &a.Name, &a.Semester, &a.StartDate, &a.EndDate, &a.IsActive)
		list = append(list, a)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": list})
}
