package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AITutorHandler struct {
	DB       *sql.DB
	mu       sync.RWMutex
	sessions []models.AIChatSession
	messages []models.AIChatMessage
}

func NewAITutorHandler(db *sql.DB) *AITutorHandler {
	h := &AITutorHandler{
		DB: db,
	}
	h.initMockData()
	return h
}

func (h *AITutorHandler) initMockData() {
	defSessID := "sess-ai-001"
	h.sessions = []models.AIChatSession{
		{
			ID:          defSessID,
			StudentID:   "stu-001",
			SubjectName: "Tahfizh & Tajwid Al-Qur'an",
			StartedAt:   time.Now().Add(-2 * time.Hour),
		},
	}

	h.messages = []models.AIChatMessage{
		{
			ID:         "msg-01",
			SessionID:  defSessID,
			SenderType: "ASSISTANT",
			Content:    "Assalamu'alaikum warahmatullah wa barakatuh! Ahlan wa sahlan ananda di Tutor AI PKBM Digital Dar El Iman. Ada materi pelajaran, kaidah tajwid, bahasa Arab, atau fikih yang ingin kita pelajari bersama hari ini?",
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
	}
}

// POST /api/v1/ai-tutor/sessions
func (h *AITutorHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		StudentID   string `json:"student_id"`
		SubjectName string `json:"subject_name"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	h.mu.Lock()
	defer h.mu.Unlock()

	sessID := "sess-ai-" + uuid.New().String()[:8]
	if req.StudentID == "" {
		req.StudentID = "stu-001"
	}
	if req.SubjectName == "" {
		req.SubjectName = "Umum & Kepesantrenan"
	}

	session := models.AIChatSession{
		ID:          sessID,
		StudentID:   req.StudentID,
		SubjectName: req.SubjectName,
		StartedAt:   time.Now(),
	}

	welcomeMsg := models.AIChatMessage{
		ID:         "msg-" + uuid.New().String()[:8],
		SessionID:  sessID,
		SenderType: "ASSISTANT",
		Content:    fmt.Sprintf("Bismillah, sesi tanya jawab mata pelajaran %s telah dibuka. Silakan ajukan pertanyaan seputar materi Anda!", req.SubjectName),
		CreatedAt:  time.Now(),
	}

	h.sessions = append(h.sessions, session)
	h.messages = append(h.messages, welcomeMsg)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id":      sessID,
		"welcome_message": welcomeMsg.Content,
		"subject_name":    req.SubjectName,
	})
}

// POST /api/v1/ai-tutor/sessions/:id/ask
func (h *AITutorHandler) AskAI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessionID := chi.URLParam(r, "id")

	var req models.AIChatAskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Question) == "" {
		http.Error(w, `{"error":"Pertanyaan wajib diisi"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 1. Simpan pesan santri
	userMsg := models.AIChatMessage{
		ID:         "msg-" + uuid.New().String()[:8],
		SessionID:  sessionID,
		SenderType: "USER",
		Content:    req.Question,
		CreatedAt:  time.Now(),
	}
	h.messages = append(h.messages, userMsg)

	// 2. Mesin Pengetahuan Konseptual Silabus PKBM Digital Dar El Iman
	answer, source := generateDarElImanAIResponse(req.Question)

	// 3. Simpan jawaban Tutor AI
	aiMsg := models.AIChatMessage{
		ID:         "msg-" + uuid.New().String()[:8],
		SessionID:  sessionID,
		SenderType: "ASSISTANT",
		Content:    answer,
		CreatedAt:  time.Now(),
	}
	h.messages = append(h.messages, aiMsg)

	response := models.AIChatAskResponse{
		SessionID: sessionID,
		Answer:    answer,
		Source:    source,
		Timestamp: time.Now().Format("15:04 WIB"),
	}

	json.NewEncoder(w).Encode(response)
}

// GET /api/v1/ai-tutor/sessions/:id/history
func (h *AITutorHandler) GetSessionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessionID := chi.URLParam(r, "id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var history []models.AIChatMessage
	for _, m := range h.messages {
		if m.SessionID == sessionID {
			history = append(history, m)
		}
	}

	json.NewEncoder(w).Encode(history)
}

// generateDarElImanAIResponse menghasilkan respon terkurasi sesuai silabus kurikulum Dar El Iman
func generateDarElImanAIResponse(q string) (string, string) {
	qLower := strings.ToLower(q)

	if strings.Contains(qLower, "tajwid") || strings.Contains(qLower, "ikhfa") || strings.Contains(qLower, "idgham") || strings.Contains(qLower, "nun mati") {
		return "Hukum nun sukun (نْ) dan tanwin terbagi menjadi 4 hukum utama:\n" +
			"1. **Idzhar Halqi** (Jelas): Apabila bertemu huruf halq (ء هـ ع ح غ خ).\n" +
			"2. **Idgham**: Terbagi dua, yaitu *Bighunnah* (dengan dengung: ي ن م و) dan *Bilaghunnah* (tanpa dengung: ل ر).\n" +
			"3. **Iqlab**: Apabila bertemu huruf ب, bunyi nun/tanwin berubah menjadi mim disertai ghunnah.\n" +
			"4. **Ikhfa Haqiqi** (Samar): Bertemu dengan 15 huruf hijaiyah lainnya.\n\n" +
			"*Nasihat:* Latihlah bacaan talaqqi langsung di halaqah bersama Ustadz agar makhraj dan kadar panjang mad tepat sesuai kaidah Asy-Syathibiyyah.", "Modul Tajwid Praktis Kelas X PKBM Dar El Iman"
	}

	if strings.Contains(qLower, "adab") || strings.Contains(qLower, "ilmu") || strings.Contains(qLower, "niat") {
		return "Imam Malik rahimahullah berpesan: *'Pelajarilah adab sebelum engkau mempelajari ilmu.'*\n" +
			"Adab utama seorang penuntut ilmu:\n" +
			"1. Mengikhlaskan niat semata-mata karena mengharap ridha Allah Subhanahu wa Ta'ala.\n" +
			"2. Bersikap tawadhu' dan menghormati para asatidzah/tutor.\n" +
			"3. Mengikat ilmu dengan tulisan dan mengamalkannya dalam kehidupan sehari-hari.\n" +
			"4. Senantiasa berwudhu dan menjaga kesucian saat berinteraksi dengan Al-Qur'an.", "Kitab Tadzkirotus Sami' wal Mutakallim & Silabus Adab PKBM"
	}

	if strings.Contains(qLower, "ai") || strings.Contains(qLower, "kecerdasan buatan") || strings.Contains(qLower, "teknologi") || strings.Contains(qLower, "koding") {
		return "Dalam mata pelajaran Khas Dar El Iman: *AI Literacy & Pemrograman Dasar*, kita memandang kecerdasan buatan (AI) sebagai instrumen wasilah untuk kemaslahatan umat:\n" +
			"1. **Etika Islami:** AI harus digunakan untuk tujuan kebaikan, kejujuran ilmiah (dilarang plagiasi tanpa atribusi), dan menjaga batasan syar'i.\n" +
			"2. **Logika Komputasi:** Berpikir sistematis melalui algoritma, perulangan, dan pemecahan masalah kompleks.\n" +
			"3. Gunakan teknologi modern ini untuk menyebarkan dakwah, memudahkan pembelajaran Al-Qur'an, dan memajukan peradaban Islam.", "Silabus Khas Dar El Iman: AI Literacy & Digital Skills"
	}

	if strings.Contains(qLower, "diaspora") || strings.Contains(qLower, "jadwal") || strings.Contains(qLower, "waktu") {
		return "Bagi santri binaan diaspora luar negeri (Arab Saudi, Uni Emirat Arab, Jepang, Malaysia, Eropa), seluruh materi rekaman sesi KBM dan forum diskusi dapat diakses secara asinkron 24/7 di ruang LMS.\n" +
			"Untuk halaqah talaqqi langsung, jadwal disesuaikan dengan zona waktu masing-masing negara melalui kesepakatan kelompok halaqah bersama Ustadz.", "Panduan Santri Diaspora PKBM Digital Dar El Iman"
	}

	return fmt.Sprintf("Alhamdulillah, pertanyaan ananda mengenai *\"%s\"* adalah topik pembelajaran yang sangat berbobot.\n\n"+
		"Berdasarkan silabus Kurikulum PKBM Digital Dar El Iman, pastikan ananda menelaah materi pada modul topik terkait di Ruang LMS, mencatat poin-poin penting, dan jika membutuhkan bimbingan langsung, ananda dapat berdiskusi pada sesi KBM daring berikutnya atau menekan tombol *'Konsultasi Bersama Ustadz Pembimbing'*. Semoga Allah memberkahi proses tholabul 'ilmi ananda!", q), "Sistem Pembelajaran Mandiri Terpadu PKBM Dar El Iman"
}
