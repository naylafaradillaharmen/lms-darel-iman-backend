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

type LMSHandler struct {
	DB *sql.DB
	mu sync.RWMutex

	// In-memory fallback
	topics    []models.CourseTopic
	materials []models.TopicMaterial
}

func NewLMSHandler(db *sql.DB) *LMSHandler {
	desc1 := "Membahas niat ikhlas menuntut ilmu dan artikulasi makhraj huruf hijaiyyah"
	desc2 := "Karakteristik bunyi huruf dan hukum Idzhar, Idgham, Ikhfa, Iqlab"

	docURL := "https://cdn.darel-iman.or.id/modul/tajwid-dasar-juz30.pdf"
	videoURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	articleContent := "Rasulullah shallallahu alaihi wa sallam bersabda: Sebaik-baik kalian adalah orang yang belajar Al-Qur'an dan mengajarkannya. (HR. Al-Bukhari). Mari kita luruskan niat dalam menempuh semester ini."

	asgDocURL := "https://cdn.darel-iman.or.id/tugas/petunjuk-tugas-01.pdf"

	h := &LMSHandler{
		DB: db,
		topics: []models.CourseTopic{
			{
				ID:            "top00001-0000-0000-0000-000000000001",
				CourseClassID: "cc000001-0000-0000-0000-000000000001",
				Title:         "Pertemuan 1: Pengantar Adab Penuntut Ilmu & Makharijul Huruf",
				SequenceOrder: 1,
				Description:   &desc1,
				IsPublished:   true,
				Materials: []models.TopicMaterial{
					{
						ID:            "mat00001-0000-0000-0000-000000000001",
						CourseTopicID: "top00001-0000-0000-0000-000000000001",
						Title:         "Buku Panduan Tajwid Asy-Syafi'i Ringkas (PDF)",
						MaterialType:  "DOCUMENT",
						FileURL:       &docURL,
						SequenceOrder: 1,
						CreatedAt:     time.Now(),
					},
					{
						ID:            "mat00002-0000-0000-0000-000000000001",
						CourseTopicID: "top00001-0000-0000-0000-000000000001",
						Title:         "Video Tutorial: Visualisasi 17 Makharijul Huruf",
						MaterialType:  "VIDEO_LINK",
						ExternalLink:  &videoURL,
						SequenceOrder: 2,
						CreatedAt:     time.Now(),
					},
					{
						ID:            "mat00003-0000-0000-0000-000000000001",
						CourseTopicID: "top00001-0000-0000-0000-000000000001",
						Title:         "Nasihat Pembuka Ustadz: Keutamaan Al-Qur'an",
						MaterialType:  "TEXT_ARTICLE",
						ContentBody:   &articleContent,
						SequenceOrder: 3,
						CreatedAt:     time.Now(),
					},
				},
				Assignments: []models.Assignment{
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
				CreatedAt: time.Now(),
			},
			{
				ID:            "top00002-0000-0000-0000-000000000002",
				CourseClassID: "cc000001-0000-0000-0000-000000000001",
				Title:         "Pertemuan 2: Sifatul Huruf & Pengenalan Hukum Nun Mati / Tanwin",
				SequenceOrder: 2,
				Description:   &desc2,
				IsPublished:   true,
				Materials:     []models.TopicMaterial{},
				Assignments:   []models.Assignment{},
				CreatedAt:     time.Now(),
			},
		},
	}
	return h
}

// 1. GET /api/v1/lms/my-courses
func (h *LMSHandler) MyCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	meetLink := "https://meet.google.com/pkbm-thf-x10"
	courses := []models.CourseClass{
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
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    courses,
	})
}

// 2. GET /api/v1/lms/courses/{course_class_id}/topics
func (h *LMSHandler) GetTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ccID := chi.URLParam(r, "course_class_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var matched []models.CourseTopic
	for _, t := range h.topics {
		if t.CourseClassID == ccID || ccID == "" || ccID == "default" {
			matched = append(matched, t)
		}
	}

	if len(matched) == 0 {
		matched = h.topics
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    matched,
	})
}

// 3. POST /api/v1/lms/courses/{course_class_id}/topics
type CreateTopicRequest struct {
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	SequenceOrder int     `json:"sequence_order"`
	IsPublished   bool    `json:"is_published"`
}

func (h *LMSHandler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ccID := chi.URLParam(r, "course_class_id")

	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Judul topik pertemuan wajib diisi"})
		return
	}

	h.mu.Lock()
	seq := req.SequenceOrder
	if seq == 0 {
		seq = len(h.topics) + 1
	}

	newTopic := models.CourseTopic{
		ID:            uuid.New().String(),
		CourseClassID: ccID,
		Title:         req.Title,
		SequenceOrder: seq,
		Description:   req.Description,
		IsPublished:   req.IsPublished,
		Materials:     []models.TopicMaterial{},
		Assignments:   []models.Assignment{},
		CreatedAt:     time.Now(),
	}

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO course_topics (id, course_class_id, title, sequence_order, description, is_published) VALUES (?, ?, ?, ?, ?, ?)`,
			newTopic.ID, newTopic.CourseClassID, newTopic.Title, newTopic.SequenceOrder, newTopic.Description, newTopic.IsPublished)
	}

	h.topics = append(h.topics, newTopic)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Topik pertemuan baru berhasil dibuat",
		"data":    newTopic,
	})
}

// 4. POST /api/v1/lms/topics/{topic_id}/materials
type AddMaterialRequest struct {
	Title        string  `json:"title"`
	MaterialType string  `json:"material_type"` // DOCUMENT, VIDEO_LINK, TEXT_ARTICLE
	ContentBody  *string `json:"content_body"`
	FileURL      *string `json:"file_url"`
	ExternalLink *string `json:"external_link"`
}

func (h *LMSHandler) AddMaterial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topicID := chi.URLParam(r, "topic_id")

	var req AddMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Judul materi wajib diisi"})
		return
	}

	newMat := models.TopicMaterial{
		ID:            uuid.New().String(),
		CourseTopicID: topicID,
		Title:         req.Title,
		MaterialType:  req.MaterialType,
		ContentBody:   req.ContentBody,
		FileURL:       req.FileURL,
		ExternalLink:  req.ExternalLink,
		CreatedAt:     time.Now(),
	}

	h.mu.Lock()
	for i, t := range h.topics {
		if t.ID == topicID {
			newMat.SequenceOrder = len(t.Materials) + 1
			h.topics[i].Materials = append(h.topics[i].Materials, newMat)
			break
		}
	}
	h.materials = append(h.materials, newMat)
	h.mu.Unlock()

	if h.DB != nil {
		_, _ = h.DB.Exec(`INSERT INTO topic_materials (id, course_topic_id, title, material_type, content_body, file_url, external_link, sequence_order) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			newMat.ID, newMat.CourseTopicID, newMat.Title, newMat.MaterialType, newMat.ContentBody, newMat.FileURL, newMat.ExternalLink, newMat.SequenceOrder)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Materi ajar berhasil ditambahkan ke modul",
		"data":    newMat,
	})
}

// 5. PUT /api/v1/lms/topics/{topic_id}/publish
func (h *LMSHandler) TogglePublishTopic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	topicID := chi.URLParam(r, "topic_id")

	h.mu.Lock()
	var updated *models.CourseTopic
	for i, t := range h.topics {
		if t.ID == topicID {
			h.topics[i].IsPublished = !h.topics[i].IsPublished
			updated = &h.topics[i]
			break
		}
	}
	h.mu.Unlock()

	if updated != nil {
		if h.DB != nil {
			_, _ = h.DB.Exec(`UPDATE course_topics SET is_published = ? WHERE id = ?`, updated.IsPublished, topicID)
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Status publikasi materi berhasil diperbarui",
			"data":    updated,
		})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Topik tidak ditemukan"})
}
