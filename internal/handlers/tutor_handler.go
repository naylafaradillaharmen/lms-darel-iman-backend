package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pkbm-backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TutorHandler struct {
	DB *sql.DB
}

func (h *TutorHandler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT t.id, t.user_id, t.nip, t.full_name, t.gender, t.specialization, t.phone,
		       u.email, u.username, t.is_active, t.created_at
		FROM tutors t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`
	rows, err := h.DB.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": []models.Tutor{}})
		return
	}
	defer rows.Close()

	var tutors []models.Tutor
	for rows.Next() {
		var t models.Tutor
		_ = rows.Scan(&t.ID, &t.UserID, &t.NIP, &t.FullName, &t.Gender, &t.Specialization, &t.Phone,
			&t.Email, &t.Username, &t.IsActive, &t.CreatedAt)
		tutors = append(tutors, t)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "data": tutors})
}

type CreateTutorRequest struct {
	FullName       string  `json:"full_name"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	NIP            *string `json:"nip"`
	Gender         *string `json:"gender"`
	Specialization *string `json:"specialization"`
	Phone          *string `json:"phone"`
}

func (h *TutorHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateTutorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Permintaan tidak valid"})
		return
	}

	userID := uuid.New().String()
	tutorID := uuid.New().String()

	pass := req.Password
	if pass == "" {
		pass = "tutor123"
	}
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

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
	`, userID, req.Username, req.Email, string(hashedPass), req.FullName, req.Phone)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Username atau email sudah digunakan"})
		return
	}

	// 2. Insert Role TUTOR
	_, _ = tx.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		SELECT ?, id FROM roles WHERE name = 'TUTOR' LIMIT 1
	`, userID)

	// 3. Insert Tutor
	_, err = tx.Exec(`
		INSERT INTO tutors (id, user_id, nip, full_name, gender, specialization, phone)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, tutorID, userID, req.NIP, req.FullName, req.Gender, req.Specialization, req.Phone)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Gagal menyimpan data profil tutor"})
		return
	}

	_ = tx.Commit()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Data Tutor dan Akun Pengguna berhasil dibuat",
		"data": map[string]string{
			"id":       tutorID,
			"user_id":  userID,
			"username": req.Username,
		},
	})
}
