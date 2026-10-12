package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pkbm-backend/internal/models"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type StudentHandler struct {
	DB *sql.DB
}

func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT s.id, s.user_id, s.nis, s.nisn, s.nik, s.full_name, s.gender, s.birth_place, s.birth_date,
		       s.address, s.is_diaspora, s.country, s.timezone, s.guardian_name, s.guardian_phone,
		       s.guardian_relationship, s.status, u.username, u.email,
		       c.name as classroom_name, s.created_at
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classroom_students cs ON s.id = cs.student_id
		LEFT JOIN classrooms c ON cs.classroom_id = c.id
		ORDER BY s.created_at DESC
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": []models.Student{}})
		return
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		_ = rows.Scan(&s.ID, &s.UserID, &s.NIS, &s.NISN, &s.NIK, &s.FullName, &s.Gender, &s.BirthPlace, &s.BirthDate,
			&s.Address, &s.IsDiaspora, &s.Country, &s.Timezone, &s.GuardianName, &s.GuardianPhone,
			&s.GuardianRelationship, &s.Status, &s.Username, &s.Email, &s.ClassroomName, &s.CreatedAt)
		students = append(students, s)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": students})
}

type CreateStudentRequest struct {
	FullName             string  `json:"full_name"`
	NIS                  string  `json:"nis"`
	NISN                 *string `json:"nisn"`
	NIK                  *string `json:"nik"`
	Gender               *string `json:"gender"`
	BirthPlace           *string `json:"birth_place"`
	BirthDate            *string `json:"birth_date"`
	Address              *string `json:"address"`
	IsDiaspora           bool    `json:"is_diaspora"`
	Country              string  `json:"country"`
	Timezone             string  `json:"timezone"`
	GuardianName         *string `json:"guardian_name"`
	GuardianPhone        *string `json:"guardian_phone"`
	GuardianRelationship *string `json:"guardian_relationship"`
	Email                string  `json:"email"`
}

func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format data tidak valid"})
		return
	}

	userID := uuid.New().String()
	studentID := uuid.New().String()

	username := fmt.Sprintf("santri.%s", strings.ToLower(strings.ReplaceAll(req.NIS, "-", "")))
	defaultPassword := "santri123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)

	email := req.Email
	if email == "" {
		email = fmt.Sprintf("%s@santri.darel-iman.or.id", username)
	}

	country := req.Country
	if country == "" {
		country = "Indonesia"
	}
	tz := req.Timezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, `{"success":false,"message":"Gagal memulai transaksi"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// 1. Insert User
	_, err = tx.Exec(`
		INSERT INTO users (id, username, email, password_hash, full_name, phone, is_active)
		VALUES (?, ?, ?, ?, ?, ?, TRUE)
	`, userID, username, email, string(hashedPass), req.FullName, req.GuardianPhone)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "NIS atau username sudah digunakan"})
		return
	}

	// 2. Assign Role SISWA
	_, _ = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		SELECT ?, id FROM roles WHERE name = 'SISWA' LIMIT 1
	`, userID)

	// 3. Insert Student
	_, err = tx.Exec(`
		INSERT INTO students (id, user_id, nis, nisn, nik, full_name, gender, birth_place, birth_date,
		                     address, is_diaspora, country, timezone, guardian_name, guardian_phone,
		                     guardian_relationship, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'ACTIVE')
	`, studentID, userID, req.NIS, req.NISN, req.NIK, req.FullName, req.Gender, req.BirthPlace, req.BirthDate,
		req.Address, req.IsDiaspora, country, tz, req.GuardianName, req.GuardianPhone, req.GuardianRelationship)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Gagal menyimpan data santri"})
		return
	}

	_ = tx.Commit()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Data Santri berhasil disimpan dan akun login otomatis aktif",
		"data": map[string]string{
			"id":       studentID,
			"username": username,
		},
	})
}

// Bulk Import CSV
func (h *StudentHandler) BulkImportCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Berkas CSV tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip Header
	_, err = reader.Read()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Format berkas kosong"})
		return
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}

		// Format CSV: NIS, NamaLengkap, JenisKelamin(L/P), NISN, Wali, TelpWali
		nis := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		gender := strings.TrimSpace(record[2])

		if nis == "" || name == "" {
			continue
		}

		userID := uuid.New().String()
		studentID := uuid.New().String()
		username := fmt.Sprintf("santri.%s", strings.ToLower(strings.ReplaceAll(nis, "-", "")))
		email := fmt.Sprintf("%s@santri.darel-iman.or.id", username)
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("santri123"), bcrypt.DefaultCost)

		_, _ = h.DB.Exec(`
			INSERT IGNORE INTO users (id, username, email, password_hash, full_name, is_active)
			VALUES (?, ?, ?, ?, ?, TRUE)
		`, userID, username, email, string(hashedPass), name)

		_, _ = h.DB.Exec(`
			INSERT IGNORE INTO user_roles (user_id, role_id)
			SELECT ?, id FROM roles WHERE name = 'SISWA' LIMIT 1
		`, userID)

		_, err = h.DB.Exec(`
			INSERT IGNORE INTO students (id, user_id, nis, full_name, gender, status)
			VALUES (?, ?, ?, ?, ?, 'ACTIVE')
		`, studentID, userID, nis, name, gender)

		if err == nil {
			count++
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Berhasil mengimpor %d santri ke dalam sistem", count),
		"count":   count,
	})
}
