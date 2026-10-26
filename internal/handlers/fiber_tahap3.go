package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Data Models Tahap 3 (CBT, Proctoring, Gradebook & E-Rapor)
type ExamQuestionOption struct {
	ID         string `json:"id"`
	OptionKey  string `json:"option_key"` // 'A', 'B', 'C', 'D'
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

type ExamQuestionItem struct {
	ID              string               `json:"id"`
	QuestionBankID  string               `json:"question_bank_id"`
	QuestionType    string               `json:"question_type"` // 'MULTIPLE_CHOICE' | 'ESSAY'
	PromptText      string               `json:"prompt_text"`
	Points          float64              `json:"points"`
	Options         []ExamQuestionOption `json:"options,omitempty"`
	CorrectOptionID string               `json:"correct_option_id,omitempty"`
}

type ExamEventDTO struct {
	ID                 string             `json:"id"`
	CourseClassID      string             `json:"course_class_id"`
	CourseClassName    string             `json:"course_class_name"`
	SubjectName        string             `json:"subject_name"`
	TutorName          string             `json:"tutor_name"`
	Title              string             `json:"title"`
	ExamType           string             `json:"exam_type"` // 'UTS', 'UAS', 'KUIS_HARIAN', 'REMEDIAL'
	DurationMinutes    int                `json:"duration_minutes"`
	ExamToken          string             `json:"exam_token"`
	PassingScore       float64            `json:"passing_score"`
	Status             string             `json:"status"` // 'ACTIVE', 'UPCOMING', 'COMPLETED'
	StartTime          string             `json:"start_time"`
	EndTime            string             `json:"end_time"`
	RandomizeQuestions bool               `json:"randomize_questions"`
	TotalQuestions     int                `json:"total_questions"`
	Questions          []ExamQuestionItem `json:"questions,omitempty"`
}

type ExamAttemptDTO struct {
	ID                 string             `json:"id"`
	ExamEventID        string             `json:"exam_event_id"`
	ExamTitle          string             `json:"exam_title"`
	SubjectName        string             `json:"subject_name"`
	StudentID          string             `json:"student_id"`
	StudentName        string             `json:"student_name"`
	NISN               string             `json:"nisn"`
	StartedAt          time.Time          `json:"started_at"`
	FinishedAt         *time.Time         `json:"finished_at,omitempty"`
	Status             string             `json:"status"` // 'IN_PROGRESS', 'SUBMITTED', 'DISQUALIFIED'
	TotalScore         *float64           `json:"total_score,omitempty"`
	IsPassed           *bool              `json:"is_passed,omitempty"`
	CheatIncidentCount int                `json:"cheat_incident_count"`
	Answers            map[string]string  `json:"answers"` // question_id -> option_key or essay
	Questions          []ExamQuestionItem `json:"questions,omitempty"`
}

type ProctorLogDTO struct {
	ID           string    `json:"id"`
	AttemptID    string    `json:"attempt_id"`
	StudentID    string    `json:"student_id"`
	StudentName  string    `json:"student_name"`
	IncidentType string    `json:"incident_type"` // 'TAB_SWITCH', 'WINDOW_BLUR', 'IDLE_TIMEOUT'
	Description  string    `json:"description"`
	LoggedAt     time.Time `json:"logged_at"`
}

// Konfigurasi Bobot Penilaian Per Mata Pelajaran
type SubjectWeightConfigDTO struct {
	SubjectCode   string  `json:"subject_code"`
	SubjectName   string  `json:"subject_name"`
	Attendance    float64 `json:"attendance"`
	Assignment    float64 `json:"assignment"`
	TahfizhDaily  float64 `json:"tahfizh_daily"`
	Uts           float64 `json:"uts"`
	Uas           float64 `json:"uas"`
	UtsSourceType string  `json:"uts_source_type"` // 'CBT_ONLINE', 'LISAN_MANUAL', 'PRAKTIK'
	UasSourceType string  `json:"uas_source_type"` // 'CBT_ONLINE', 'TULIS_MANUAL'
	AdabBehavior  float64 `json:"adab_behavior"`   // Komponen observasi adab/sikap
	PassingScore  float64 `json:"passing_score"`
}

type GradebookEntryDTO struct {
	StudentID       string  `json:"student_id"`
	StudentName     string  `json:"student_name"`
	NISN            string  `json:"nisn"`
	AttendanceScore float64 `json:"attendance_score"` // Presensi
	AssignmentScore float64 `json:"assignment_score"` // Tugas LMS
	TahfizhScore    float64 `json:"tahfizh_score"`    // Mutaba'ah Tahfizh
	UtsScore        float64 `json:"uts_score"`        // UTS (CBT / Lisan)
	UtsSource       string  `json:"uts_source"`       // Sumber nilai UTS
	UasScore        float64 `json:"uas_score"`        // UAS
	UasSource       string  `json:"uas_source"`       // Sumber nilai UAS
	AdabScore       float64 `json:"adab_score"`       // Observasi Adab & Sikap
	FinalScore      float64 `json:"final_score"`      // Nilai Akhir Komposit
	LetterGrade     string  `json:"letter_grade"`     // A, B, C, D
	PassingStatus   string  `json:"passing_status"`   // 'TUNTAS', 'REMEDIAL'
}

type ReportCardSubjectDTO struct {
	SubjectCode      string  `json:"subject_code"`
	SubjectName      string  `json:"subject_name"`
	TutorName        string  `json:"tutor_name"`
	Score            float64 `json:"score"`
	LetterGrade      string  `json:"letter_grade"`
	PassingThreshold float64 `json:"passing_threshold"`
	CompetencyNotes  string  `json:"competency_notes"`
	SourceType       string  `json:"source_type"` // 'CBT', 'LISAN', 'PORTOFOLIO'
}

type StudentReportCardDTO struct {
	ID                 string                 `json:"id"`
	StudentID          string                 `json:"student_id"`
	StudentName        string                 `json:"student_name"`
	NISN               string                 `json:"nisn"`
	ClassroomName      string                 `json:"classroom_name"`
	AcademicYear       string                 `json:"academic_year"`
	Semester           string                 `json:"semester"`
	SemesterCode       string                 `json:"semester_code"` // 'sem-1', 'sem-2', 'sem-3'
	HomeroomNotes      string                 `json:"homeroom_notes"`
	SpiritualGrade     string                 `json:"spiritual_grade"`
	SocialGrade        string                 `json:"social_grade"`
	TahfizhSummary     string                 `json:"tahfizh_summary"`
	QRVerificationCode string                 `json:"qr_verification_code"`
	VerificationURL    string                 `json:"verification_url"`
	PublishedAt        string                 `json:"published_at"`
	IsLocked           bool                   `json:"is_locked"`
	Grades             []ReportCardSubjectDTO `json:"grades"`
}

type FiberTahap3Handler struct {
	GormDB          *gorm.DB
	mu              sync.RWMutex
	examEvents      []ExamEventDTO
	attempts        []ExamAttemptDTO
	proctorLogs     []ProctorLogDTO
	subjectWeights  []SubjectWeightConfigDTO
	gradebook       []GradebookEntryDTO
	reportCards     []StudentReportCardDTO
	isGradeLocked   bool
	isRaporReleased bool // Wali kelas release toggle untuk akses santri
}

func NewFiberTahap3Handler(db *gorm.DB) *FiberTahap3Handler {
	tahfizhQuestions := []ExamQuestionItem{
		{
			ID: "q-01", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Berapakah jumlah ayat yang terdapat dalam Surah An-Naba' dan di kota manakah surah tersebut diturunkan?",
			Options: []ExamQuestionOption{
				{ID: "opt-1a", OptionKey: "A", OptionText: "40 Ayat, Makkiyyah", IsCorrect: true},
				{ID: "opt-1b", OptionKey: "B", OptionText: "42 Ayat, Madaniyyah", IsCorrect: false},
				{ID: "opt-1c", OptionKey: "C", OptionText: "36 Ayat, Makkiyyah", IsCorrect: false},
				{ID: "opt-1d", OptionKey: "D", OptionText: "45 Ayat, Madaniyyah", IsCorrect: false},
			},
		},
		{
			ID: "q-02", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Hukum bacaan nun sukun bertemu huruf Ta (ت) pada kata 'مِنْ تَحْتِهَا' adalah...",
			Options: []ExamQuestionOption{
				{ID: "opt-2a", OptionKey: "A", OptionText: "Idzhar Halqi", IsCorrect: false},
				{ID: "opt-2b", OptionKey: "B", OptionText: "Ikhfa' Haqiqi", IsCorrect: true},
				{ID: "opt-2c", OptionKey: "C", OptionText: "Idgham Bighunnah", IsCorrect: false},
				{ID: "opt-2d", OptionKey: "D", OptionText: "Iqlab", IsCorrect: false},
			},
		},
		{
			ID: "q-03", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Panjang harakat pada Mad Wajib Muttashil menurut riwayat Hafsh 'an 'Ashim Min Thariq Asy-Syathibiyyah adalah...",
			Options: []ExamQuestionOption{
				{ID: "opt-3a", OptionKey: "A", OptionText: "2 Harakat", IsCorrect: false},
				{ID: "opt-3b", OptionKey: "B", OptionText: "4 atau 5 Harakat", IsCorrect: true},
				{ID: "opt-3c", OptionKey: "C", OptionText: "6 Harakat Lazim", IsCorrect: false},
				{ID: "opt-3d", OptionKey: "D", OptionText: "1 Harakat", IsCorrect: false},
			},
		},
		{
			ID: "q-04", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Huruf-huruf qalqalah yang terkumpul dalam bait 'قُطْبُ جَدٍّ' dibaca memantul dengan tingkatan Qalqalah Kubra jika...",
			Options: []ExamQuestionOption{
				{ID: "opt-4a", OptionKey: "A", OptionText: "Berada di tengah kalimat dan bertanda sukun asli", IsCorrect: false},
				{ID: "opt-4b", OptionKey: "B", OptionText: "Berada di akhir kalimat karena waqaf (berhenti)", IsCorrect: true},
				{ID: "opt-4c", OptionKey: "C", OptionText: "Berharakat fathah tanwin", IsCorrect: false},
				{ID: "opt-4d", OptionKey: "D", OptionText: "Bertemu huruf hijaiyah berharakat tasydid di awal kata", IsCorrect: false},
			},
		},
		{
			ID: "q-05", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Makhraj huruf Dhad (ض) keluar dari...",
			Options: []ExamQuestionOption{
				{ID: "opt-5a", OptionKey: "A", OptionText: "Ujung lidah menempel pada dinding dua gigi seri atas", IsCorrect: false},
				{ID: "opt-5b", OptionKey: "B", OptionText: "Salah satu tepi lidah atau kedua-duanya bertemu gigi geraham atas", IsCorrect: true},
				{ID: "opt-5c", OptionKey: "C", OptionText: "Dua bibir terkatup rapat", IsCorrect: false},
				{ID: "opt-5d", OptionKey: "D", OptionText: "Pangkal tenggorokan terdalam (Aqshal Halaq)", IsCorrect: false},
			},
		},
		{
			ID: "q-06", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Makna kata 'عَمَّ يَتَسَاۤءَلُوْنَ' pada ayat pertama Surah An-Naba' membicarakan tentang...",
			Options: []ExamQuestionOption{
				{ID: "opt-6a", OptionKey: "A", OptionText: "Pertanyaan kaum musyrikin mengenai hari kebangkitan dan kiamat", IsCorrect: true},
				{ID: "opt-6b", OptionKey: "B", OptionText: "Tata cara pembagian harta rampasan perang", IsCorrect: false},
				{ID: "opt-6c", OptionKey: "C", OptionText: "Hukum puasa di bulan Ramadhan", IsCorrect: false},
				{ID: "opt-6d", OptionKey: "D", OptionText: "Kisah Nabi Musa dan Fir'aun", IsCorrect: false},
			},
		},
		{
			ID: "q-07", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Hukum bacaan mim sukun bertemu huruf Ba (ب) disebut dengan...",
			Options: []ExamQuestionOption{
				{ID: "opt-7a", OptionKey: "A", OptionText: "Idzhar Syafawi", IsCorrect: false},
				{ID: "opt-7b", OptionKey: "B", OptionText: "Ikhfa' Syafawi", IsCorrect: true},
				{ID: "opt-7c", OptionKey: "C", OptionText: "Idgham Mimi", IsCorrect: false},
				{ID: "opt-7d", OptionKey: "D", OptionText: "Ghunnah Musyaddadah", IsCorrect: false},
			},
		},
		{
			ID: "q-08", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Tanda waqaf lazim (م) dalam mushaf Al-Qur'an Standar Indonesia memiliki arti...",
			Options: []ExamQuestionOption{
				{ID: "opt-8a", OptionKey: "A", OptionText: "Harus berhenti (wajib waqaf)", IsCorrect: true},
				{ID: "opt-8b", OptionKey: "B", OptionText: "Dilarang berhenti (harus washal)", IsCorrect: false},
				{ID: "opt-8c", OptionKey: "C", OptionText: "Boleh berhenti boleh lanjut", IsCorrect: false},
				{ID: "opt-8d", OptionKey: "D", OptionText: "Lebih utama meneruskan bacaan", IsCorrect: false},
			},
		},
		{
			ID: "q-09", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Huruf-huruf Isti'la yang wajib dibaca tebal (Tafkhim) berjumlah 7 huruf yang terangkum dalam kalimat...",
			Options: []ExamQuestionOption{
				{ID: "opt-9a", OptionKey: "A", OptionText: "خُصَّ ضَغْطٍ قِظْ", IsCorrect: true},
				{ID: "opt-9b", OptionKey: "B", OptionText: "يَرْمَلُوْنَ", IsCorrect: false},
				{ID: "opt-9c", OptionKey: "C", OptionText: "قُطْبُ جَدٍّ", IsCorrect: false},
				{ID: "opt-9d", OptionKey: "D", OptionText: "أَخِي هَاكَ عِلْمًا", IsCorrect: false},
			},
		},
		{
			ID: "q-10", QuestionBankID: "qb-01", QuestionType: "MULTIPLE_CHOICE", Points: 10.0,
			PromptText: "Berapa jumlah rakaat shalat tarawih dan witir yang disunnahkan sesuai sunnah Rasulullah ﷺ?",
			Options: []ExamQuestionOption{
				{ID: "opt-10a", OptionKey: "A", OptionText: "11 Rakaat (8 Tarawih + 3 Witir)", IsCorrect: true},
				{ID: "opt-10b", OptionKey: "B", OptionText: "4 Rakaat", IsCorrect: false},
				{ID: "opt-10c", OptionKey: "C", OptionText: "2 Rakaat", IsCorrect: false},
				{ID: "opt-10d", OptionKey: "D", OptionText: "5 Rakaat", IsCorrect: false},
			},
		},
		{
			ID: "q-11", QuestionBankID: "qb-01", QuestionType: "ESSAY", Points: 20.0,
			PromptText: "Jelaskan secara ringkas pengertian Sanad Al-Qur'an dan mengapa transmisi sanad bacaan muttashil sampai ke Rasulullah ﷺ sangat dijaga dalam pendidikan PKBM Digital Dar El Iman!",
		},
	}

	h := &FiberTahap3Handler{
		GormDB: db,
		// Jadwal Ujian Hari Ini dengan Jam Pelaksanaan & Status Dinamis sesuai Jam Asli
		examEvents: func() []ExamEventDTO {
			nowTime := time.Now()
			startActive := nowTime.Add(-30 * time.Minute).Format("15:04") + " WIB"
			endActive := nowTime.Add(60 * time.Minute).Format("15:04") + " WIB"

			startUpcoming1 := nowTime.Add(75 * time.Minute).Format("15:04") + " WIB"
			endUpcoming1 := nowTime.Add(135 * time.Minute).Format("15:04") + " WIB"

			startUpcoming2 := nowTime.Add(150 * time.Minute).Format("15:04") + " WIB"
			endUpcoming2 := nowTime.Add(210 * time.Minute).Format("15:04") + " WIB"

			return []ExamEventDTO{
				{
					ID:                 "exam-01",
					CourseClassID:      "cc-01",
					CourseClassName:    "Kelas X Paket C — Umar bin Khattab",
					SubjectName:        "Tahfizh Al-Qur'an Bersanad & Tajwid",
					TutorName:          "Ustadz Ahmad Fauzi, Lc.",
					Title:              "Ujian Tengah Semester (UTS) Ganjil — Tahfizh & Tajwid Bersanad",
					ExamType:           "UTS",
					DurationMinutes:    45,
					ExamToken:          "UTS-PKBM-2026",
					PassingScore:       75.0,
					Status:             "ACTIVE", // Sedang Berlangsung Saat Ini
					StartTime:          startActive,
					EndTime:            endActive,
					RandomizeQuestions: true,
					TotalQuestions:     11,
					Questions:          tahfizhQuestions,
				},
				{
					ID:                 "exam-02",
					CourseClassID:      "cc-04",
					CourseClassName:    "Kelas X Paket C — Umar bin Khattab",
					SubjectName:        "Matematika Terapan & Logika Digital",
					TutorName:          "Ahmad Fauzan, S.Kom., M.T.",
					Title:              "Ujian Tengah Semester (UTS) Ganjil — Statistika & Aljabar",
					ExamType:           "UTS",
					DurationMinutes:    60,
					ExamToken:          "MTK-UTS-2026",
					PassingScore:       75.0,
					Status:             "UPCOMING", // Sesi Mendatang Hari Ini
					StartTime:          startUpcoming1,
					EndTime:            endUpcoming1,
					RandomizeQuestions: true,
					TotalQuestions:     10,
					Questions:          tahfizhQuestions[:5],
				},
				{
					ID:                 "exam-03",
					CourseClassID:      "cc-02",
					CourseClassName:    "Kelas X Paket C — Umar bin Khattab",
					SubjectName:        "Bahasa Arab Komunikatif & Nahwu",
					TutorName:          "Ustadz Ahmad Fauzi, Lc.",
					Title:              "Ujian Tengah Semester (UTS) Ganjil — Fahmul Masmu' & Qawaid",
					ExamType:           "UTS",
					DurationMinutes:    45,
					ExamToken:          "ARB-UTS-2026",
					PassingScore:       75.0,
					Status:             "UPCOMING", // Sesi Mendatang Hari Ini
					StartTime:          startUpcoming2,
					EndTime:            endUpcoming2,
					RandomizeQuestions: true,
					TotalQuestions:     10,
					Questions:          tahfizhQuestions[:5],
				},
			}
		}(),
		// Sesi Santri (Initial state: Ziyad belum mulai agar bisa mulai fresh tanpa preselected answer!)
		// Sesi Santri (Initial state: Santri sebelumnya telah selesai mengumpulkan sehingga Sesi Berjalan live = 0)
		attempts: []ExamAttemptDTO{
			{
				ID:                 "att-002",
				ExamEventID:        "exam-01",
				ExamTitle:          "Ujian Tengah Semester (UTS) Ganjil — Tahfizh & Tajwid Bersanad",
				SubjectName:        "Tahfizh Al-Qur'an Bersanad & Tajwid",
				StudentID:          "u0000013b",
				StudentName:        "Yafiq Al-Faruq",
				NISN:               "0081234562",
				StartedAt:          time.Now().Add(-60 * time.Minute),
				Status:             "SUBMITTED",
				CheatIncidentCount: 0,
				Answers: map[string]string{
					"q-01": "A", "q-02": "B", "q-03": "B",
				},
				Questions: tahfizhQuestions,
			},
		},
		proctorLogs: []ProctorLogDTO{},
		// Konfigurasi Bobot Penilaian Fleksibel Per Mapel
		subjectWeights: []SubjectWeightConfigDTO{
			{
				SubjectCode:   "THF-01",
				SubjectName:   "Tahfizh Al-Qur'an Bersanad",
				Attendance:    10.0,
				Assignment:    20.0,
				TahfizhDaily:  30.0,
				Uts:           20.0,
				Uas:           20.0,
				UtsSourceType: "LISAN_MANUAL", // Ujian Lisan Tajwid
				UasSourceType: "CBT_ONLINE",
				AdabBehavior:  10.0,
				PassingScore:  75.0,
			},
			{
				SubjectCode:   "MTK-10",
				SubjectName:   "Matematika Terapan & Logika Digital",
				Attendance:    10.0,
				Assignment:    30.0,
				TahfizhDaily:  0.0,
				Uts:           30.0,
				Uas:           30.0,
				UtsSourceType: "CBT_ONLINE",
				UasSourceType: "CBT_ONLINE",
				AdabBehavior:  10.0,
				PassingScore:  75.0,
			},
			{
				SubjectCode:   "ARB-01",
				SubjectName:   "Bahasa Arab Komunikatif & Nahwu",
				Attendance:    10.0,
				Assignment:    25.0,
				TahfizhDaily:  15.0,
				Uts:           25.0,
				Uas:           25.0,
				UtsSourceType: "LISAN_MANUAL",
				UasSourceType: "CBT_ONLINE",
				AdabBehavior:  10.0,
				PassingScore:  75.0,
			},
		},
		gradebook: []GradebookEntryDTO{
			{
				StudentID:       "u0000012b",
				StudentName:     "Muhammad Abdullah Ziyad",
				NISN:            "0081234561",
				AttendanceScore: 95.0,
				AssignmentScore: 92.0,
				TahfizhScore:    96.0,
				UtsScore:        90.0,
				UtsSource:       "CBT Digital (90 / 100)",
				UasScore:        94.0,
				UasSource:       "CBT Digital (94 / 100)",
				AdabScore:       95.0,
				FinalScore:      93.5,
				LetterGrade:     "A",
				PassingStatus:   "TUNTAS",
			},
			{
				StudentID:       "u0000013b",
				StudentName:     "Yafiq Al-Faruq",
				NISN:            "0081234562",
				AttendanceScore: 88.0,
				AssignmentScore: 85.0,
				TahfizhScore:    86.0,
				UtsScore:        84.0,
				UtsSource:       "Ujian Lisan Ustadz (84.0)",
				UasScore:        87.0,
				UasSource:       "CBT Digital (87.0)",
				AdabScore:       88.0,
				FinalScore:      86.1,
				LetterGrade:     "B",
				PassingStatus:   "TUNTAS",
			},
			{
				StudentID:       "u0000014b",
				StudentName:     "Abdul Andhika Raditya",
				NISN:            "0081234563",
				AttendanceScore: 96.0,
				AssignmentScore: 95.0,
				TahfizhScore:    98.0,
				UtsScore:        96.0,
				UtsSource:       "CBT Digital (96 / 100)",
				UasScore:        97.0,
				UasSource:       "CBT Digital (97 / 100)",
				AdabScore:       96.0,
				FinalScore:      96.5,
				LetterGrade:     "A",
				PassingStatus:   "TUNTAS",
			},
			{
				StudentID:       "u0000015b",
				StudentName:     "Afdhal Al Baihaqi",
				NISN:            "0081234564",
				AttendanceScore: 90.0,
				AssignmentScore: 90.0,
				TahfizhScore:    91.0,
				UtsScore:        80.0,
				UtsSource:       "Ujian Lisan Ustadz (80.0)",
				UasScore:        91.0,
				UasSource:       "CBT Digital (91.0)",
				AdabScore:       90.0,
				FinalScore:      89.4,
				LetterGrade:     "B+",
				PassingStatus:   "TUNTAS",
			},
			{
				StudentID:       "u0000016b",
				StudentName:     "Ahmad Hasan Deedat",
				NISN:            "0081234565",
				AttendanceScore: 94.0,
				AssignmentScore: 95.0,
				TahfizhScore:    95.0,
				UtsScore:        99.0,
				UtsSource:       "CBT Digital (99 / 100)",
				UasScore:        97.0,
				UasSource:       "CBT Digital (97 / 100)",
				AdabScore:       95.0,
				FinalScore:      96.0,
				LetterGrade:     "A",
				PassingStatus:   "TUNTAS",
			},
		},
		reportCards: []StudentReportCardDTO{
			// ======================= ZIYAD =======================
			// Semester 3 (Aktif)
			{
				ID:                 "rep-001",
				StudentID:          "u0000012b",
				StudentName:        "Muhammad Abdullah Ziyad",
				NISN:               "0081234561",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2026/2027",
				Semester:           "Semester 3 (Ganjil 2026/2027) — Aktif",
				SemesterCode:       "sem-3",
				HomeroomNotes:      "Ananda Ziyad menunjukkan akhlak mulia dan kedisiplinan yang sangat baik dalam seluruh KBM daring. Hafalan Al-Qur'an bersanad Juz 30 telah tuntas mutqin dengan predikat Mumtaz.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "A (Sangat Baik)",
				TahfizhSummary:     "Tuntas Juz 30 (Surah An-Naba s.d. An-Nas) — Sanad Riwayat Hafsh 'an 'Ashim Min Thariq Asy-Syathibiyyah",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-001",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-001",
				PublishedAt:        "18 Oktober 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 96.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Sangat fasih dalam makhraj dan kaidah tajwid, mutqin hafalan Juz 30.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 92.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Mampu menyusun kalimat sempurna dan memahami dasar i'rab isim.", SourceType: "CBT"},
					{SubjectCode: "IND-10", SubjectName: "Bahasa Indonesia (Paket C)", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 91.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Sangat terampil menyusun esai argumentasi dan telaah kritis teks sastra.", SourceType: "PORTOFOLIO"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 89.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Memahami statistika distribusi frekuensi data tunggal dan kelompok dengan cermat.", SourceType: "CBT"},
					{SubjectCode: "TEC-01", SubjectName: "AI Literacy & Digital Skill", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 95.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Menguasai prompt engineering etis dan pembuatan dokumen kolaborasi cloud.", SourceType: "CBT"},
				},
			},
			// Semester 2 (Riwayat Lalu)
			{
				ID:                 "rep-002",
				StudentID:          "u0000012b",
				StudentName:        "Muhammad Abdullah Ziyad",
				NISN:               "0081234561",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2025/2026",
				Semester:           "Semester 2 (Genap 2025/2026)",
				SemesterCode:       "sem-2",
				HomeroomNotes:      "Ananda Ziyad konsisten menjaga kehadiran dan adab dalam halaqah Al-Qur'an.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "B+ (Baik)",
				TahfizhSummary:     "Tuntas Juz 29 (Surah Al-Mulk s.d. Al-Mursalat) — Riwayat Hafsh 'an 'Ashim",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-SEM2-001",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-SEM2-001",
				PublishedAt:        "20 Juni 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 94.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Mutqin hafalan Juz 29 dengan makhraj huruf yang teruji.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 88.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Memahami mufrodat tematik dan kaidah fi'il madhi.", SourceType: "CBT"},
					{SubjectCode: "IND-10", SubjectName: "Bahasa Indonesia (Paket C)", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 89.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Baik dalam menyusun teks biografi dan laporan observasi.", SourceType: "PORTOFOLIO"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 87.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Memahami barisan dan deret aritmatika terapan.", SourceType: "CBT"},
				},
			},
			// Semester 1 (Riwayat Lalu)
			{
				ID:                 "rep-003",
				StudentID:          "u0000012b",
				StudentName:        "Muhammad Abdullah Ziyad",
				NISN:               "0081234561",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2025/2026",
				Semester:           "Semester 1 (Ganjil 2025/2026)",
				SemesterCode:       "sem-1",
				HomeroomNotes:      "Awal penyesuaian yang sangat baik dengan budaya digital pesantren Dar El Iman.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "A (Sangat Baik)",
				TahfizhSummary:     "Tuntas Matan Tuhfatul Athfal & Juz 30 Dasar",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2025-XPC-SEM1-001",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2025-XPC-SEM1-001",
				PublishedAt:        "20 Desember 2025",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 92.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Menguasai hukum nun sukun dan tanwin matan tajwid.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 85.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Menguasai dasar isim isyarah dan dhomir munfashil.", SourceType: "CBT"},
				},
			},

			// ======================= YAFIQ AL-FARUQ =======================
			// Semester 3 (Aktif)
			{
				ID:                 "rep-yafiq-001",
				StudentID:          "u0000013b",
				StudentName:        "Yafiq Al-Faruq",
				NISN:               "0081234562",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2026/2027",
				Semester:           "Semester 3 (Ganjil 2026/2027) — Aktif",
				SemesterCode:       "sem-3",
				HomeroomNotes:      "Ananda Yafiq menunjukkan kesungguhan belajar yang istiqomah dan adab thalabul ilmi yang sangat baik. Hafalan Juz 30 terus meningkat dengan tajwid yang teruji.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "A (Sangat Baik)",
				TahfizhSummary:     "Tuntas Juz 30 (Surah An-Naba s.d. An-Nas) — Sanad Riwayat Hafsh 'an 'Ashim",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-002",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-002",
				PublishedAt:        "18 Oktober 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 88.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Mutqin tajwid mad ashli dan mad far'i, muraja'ah tertib.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 86.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Memahami struktur kalimat fi'liyyah dan ismiyyah dengan baik.", SourceType: "CBT"},
					{SubjectCode: "IND-10", SubjectName: "Bahasa Indonesia (Paket C)", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 87.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Kreatif dan terstruktur dalam penulisan rangkuman modul literasi.", SourceType: "PORTOFOLIO"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 84.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Mampu menyelesaikan analisis data rata-rata dan penyajian grafik.", SourceType: "CBT"},
					{SubjectCode: "TEC-01", SubjectName: "AI Literacy & Digital Skill", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 89.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Mahir memanfaatkan platform LMS daring dan kolaborasi dokumen online.", SourceType: "CBT"},
				},
			},
			// Semester 2
			{
				ID:                 "rep-yafiq-002",
				StudentID:          "u0000013b",
				StudentName:        "Yafiq Al-Faruq",
				NISN:               "0081234562",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2025/2026",
				Semester:           "Semester 2 (Genap 2025/2026)",
				SemesterCode:       "sem-2",
				HomeroomNotes:      "Ananda Yafiq selalu hadir tepat waktu dalam halaqah tahfizh dan santun kepada para asatidzah.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "B+ (Baik)",
				TahfizhSummary:     "Tuntas Juz 29 (Surah Al-Mulk s.d. Al-Insan) — Riwayat Hafsh 'an 'Ashim",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-SEM2-002",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-SEM2-002",
				PublishedAt:        "20 Juni 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 86.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Penerapan hukum tajwid ikhfa dan idgham sangat baik.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 83.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Menguasai kosakata hiwar yaumi.", SourceType: "CBT"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 82.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Memahami aritmatika sosial dan rasio perbandingan.", SourceType: "CBT"},
				},
			},
			// Semester 1
			{
				ID:                 "rep-yafiq-003",
				StudentID:          "u0000013b",
				StudentName:        "Yafiq Al-Faruq",
				NISN:               "0081234562",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2025/2026",
				Semester:           "Semester 1 (Ganjil 2025/2026)",
				SemesterCode:       "sem-1",
				HomeroomNotes:      "Awal langkah thalabul ilmi yang penuh berkah di PKBM Digital Dar El Iman.",
				SpiritualGrade:     "B+ (Baik)",
				SocialGrade:        "A (Sangat Baik)",
				TahfizhSummary:     "Tuntas Matan Tuhfatul Athfal Dasar",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2025-XPC-SEM1-002",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2025-XPC-SEM1-002",
				PublishedAt:        "20 Desember 2025",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 85.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Menguasai hukum mad thobi'i.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 81.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Dasar huruf hijaiyyah dan dhomir munfashil.", SourceType: "CBT"},
				},
			},

			// ======================= ABDUL ANDHIKA RADITYA =======================
			{
				ID:                 "rep-andhika-001",
				StudentID:          "u0000014b",
				StudentName:        "Abdul Andhika Raditya",
				NISN:               "0081234563",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2026/2027",
				Semester:           "Semester 3 (Ganjil 2026/2027) — Aktif",
				SemesterCode:       "sem-3",
				HomeroomNotes:      "Prestasi akademik luar biasa dalam tes tulis dan lisan, mutqin hafalan Juz 30.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "A (Sangat Baik)",
				TahfizhSummary:     "Tuntas Juz 30 Mutqin — Sanad Riwayat Hafsh 'an 'Ashim",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-003",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-003",
				PublishedAt:        "18 Oktober 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 98.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Sangat mutqin tajwid dan kelancaran tasmi'.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 96.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Sangat unggul dalam nahwu sharaf terapan.", SourceType: "CBT"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 96.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Cepat dan tepat dalam analisis statistika dan aljabar.", SourceType: "CBT"},
				},
			},

			// ======================= AFDHAL AL BAIHAQI =======================
			{
				ID:                 "rep-afdhal-001",
				StudentID:          "u0000015b",
				StudentName:        "Afdhal Al Baihaqi",
				NISN:               "0081234564",
				ClassroomName:      "Kelas X Paket C — Umar bin Khattab",
				AcademicYear:       "2026/2027",
				Semester:           "Semester 3 (Ganjil 2026/2027) — Aktif",
				SemesterCode:       "sem-3",
				HomeroomNotes:      "Santri rajin, tekun, dan menunjukkan peningkatan nilai yang konsisten pada semua mata pelajaran.",
				SpiritualGrade:     "A (Sangat Baik)",
				SocialGrade:        "B+ (Baik)",
				TahfizhSummary:     "Tuntas 26 Surah Juz 30 — Riwayat Hafsh 'an 'Ashim",
				QRVerificationCode: "DAR-EL-IMAN-RAPOR-2026-XPC-004",
				VerificationURL:    "http://localhost:5173/#/verify-rapor/DAR-EL-IMAN-RAPOR-2026-XPC-004",
				PublishedAt:        "18 Oktober 2026",
				IsLocked:           true,
				Grades: []ReportCardSubjectDTO{
					{SubjectCode: "THF-01", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorName: "Ustadzah Fatimah, S.Pd.", Score: 91.0, LetterGrade: "A", PassingThreshold: 75.0, CompetencyNotes: "Makhraj huruf jelas, kelancaran sangat baik.", SourceType: "LISAN"},
					{SubjectCode: "ARB-01", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorName: "Ustadz Ahmad Fauzi, Lc.", Score: 88.0, LetterGrade: "B+", PassingThreshold: 75.0, CompetencyNotes: "Mampu menterjemahkan teks bacaan pendek dengan benar.", SourceType: "CBT"},
					{SubjectCode: "MTK-10", SubjectName: "Matematika Terapan (Paket C)", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Score: 86.0, LetterGrade: "B", PassingThreshold: 75.0, CompetencyNotes: "Memahami konsep statistika terapan.", SourceType: "CBT"},
				},
			},
		},
		isGradeLocked:   true,
		isRaporReleased: true, // Default rilis dibuka (dapat ditoggle wali kelas)
	}
	return h
}

// Helper kalkulasi status ujian real-time berdasarkan jam asli
func parseTimeStrToMinutes(s string) int {
	parts := strings.Split(s, " ")
	if len(parts) == 0 {
		return -1
	}
	timePart := parts[0]
	subParts := strings.Split(timePart, ":")
	if len(subParts) < 2 {
		return -1
	}
	h, errH := strconv.Atoi(subParts[0])
	m, errM := strconv.Atoi(subParts[1])
	if errH != nil || errM != nil {
		return -1
	}
	return h*60 + m
}

func calculateExamStatus(startTime, endTime string) string {
	now := time.Now()
	currentMinutes := now.Hour()*60 + now.Minute()
	startMin := parseTimeStrToMinutes(startTime)
	endMin := parseTimeStrToMinutes(endTime)
	if startMin >= 0 && endMin >= 0 {
		if currentMinutes < startMin {
			return "UPCOMING"
		} else if currentMinutes <= endMin {
			return "ACTIVE"
		} else {
			return "COMPLETED"
		}
	}
	return "ACTIVE"
}

// 1. Get Exam Events (Daftar Sesi Ujian Hari Ini)
func (h *FiberTahap3Handler) GetExamEvents(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	events := make([]ExamEventDTO, len(h.examEvents))
	copy(events, h.examEvents)

	for i := range events {
		if events[i].Status != "FINISHED" {
			events[i].Status = calculateExamStatus(events[i].StartTime, events[i].EndTime)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(events),
		"data":    events,
	})
}

// 1b. Create Exam Event (Guru Menambah Ujian Baru & Auto Generate Token)
type CreateExamRequest struct {
	Title              string             `json:"title"`
	SubjectName        string             `json:"subject_name"`
	CourseClassID      string             `json:"course_class_id"`
	CourseClassName    string             `json:"course_class_name"`
	TutorName          string             `json:"tutor_name"`
	ExamType           string             `json:"exam_type"`
	DurationMinutes    int                `json:"duration_minutes"`
	PassingScore       float64            `json:"passing_score"`
	StartTime          string             `json:"start_time"`
	EndTime            string             `json:"end_time"`
	RandomizeQuestions bool               `json:"randomize_questions"`
	Questions          []ExamQuestionItem `json:"questions"`
}

func (h *FiberTahap3Handler) CreateExamEvent(c *fiber.Ctx) error {
	var req CreateExamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data formulir pembuatan ujian tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Auto generate exam token: e.g. UTS-[SUBJECT_PREFIX]-[RANDOM]
	randSuffix := fmt.Sprintf("%04d", rand.Intn(9000)+1000)
	subCode := "CBT"
	subLower := strings.ToLower(req.SubjectName)
	if strings.Contains(subLower, "tahfizh") {
		subCode = "THF"
	} else if strings.Contains(subLower, "arab") {
		subCode = "ARB"
	} else if strings.Contains(subLower, "matematika") {
		subCode = "MTK"
	} else if strings.Contains(subLower, "fiqih") || strings.Contains(subLower, "pai") {
		subCode = "FQH"
	} else if strings.Contains(subLower, "indonesia") {
		subCode = "IND"
	}
	generatedToken := fmt.Sprintf("UTS-%s-%s", subCode, randSuffix)

	newID := fmt.Sprintf("exam-%02d", len(h.examEvents)+1)
	if req.CourseClassName == "" {
		req.CourseClassName = "Kelas X Paket C — Umar bin Khattab"
	}
	if req.CourseClassID == "" {
		req.CourseClassID = "cc-01"
	}
	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 45
	}
	if req.PassingScore <= 0 {
		req.PassingScore = 75.0
	}
	if req.StartTime == "" {
		req.StartTime = "08:00 WIB"
	}
	if req.EndTime == "" {
		req.EndTime = "09:30 WIB"
	}
	if req.ExamType == "" {
		req.ExamType = "UTS"
	}

	newEventStatus := calculateExamStatus(req.StartTime, req.EndTime)
	newEvent := ExamEventDTO{
		ID:                 newID,
		CourseClassID:      req.CourseClassID,
		CourseClassName:    req.CourseClassName,
		SubjectName:        req.SubjectName,
		TutorName:          req.TutorName,
		Title:              req.Title,
		ExamType:           req.ExamType,
		DurationMinutes:    req.DurationMinutes,
		ExamToken:          generatedToken,
		PassingScore:       req.PassingScore,
		Status:             newEventStatus,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		RandomizeQuestions: req.RandomizeQuestions,
		TotalQuestions:     len(req.Questions),
		Questions:          req.Questions,
	}

	h.examEvents = append([]ExamEventDTO{newEvent}, h.examEvents...)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Jadwal ujian berhasil diterbitkan! Token resmi: %s", generatedToken),
		"token":   generatedToken,
		"data":    newEvent,
	})
}

// 2. Start Exam
type StartExamRequest struct {
	ExamID      string `json:"exam_id"`
	ExamToken   string `json:"exam_token"`
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	NISN        string `json:"nisn"`
}

func (h *FiberTahap3Handler) StartExam(c *fiber.Ctx) error {
	var req StartExamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format permintaan ujian tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var targetEvent *ExamEventDTO
	for _, ev := range h.examEvents {
		if ev.ID == req.ExamID {
			evCopy := ev
			targetEvent = &evCopy
			break
		}
	}

	if targetEvent == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Sesi ujian tidak ditemukan.",
		})
	}

	if targetEvent.Status == "UPCOMING" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Ujian ini belum dimulai. Jadwal pelaksanaan: %s s.d. %s.", targetEvent.StartTime, targetEvent.EndTime),
		})
	}

	if targetEvent.ExamToken != req.ExamToken {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Token ujian tidak valid. Pastikan meminta token resmi dari Proktor.",
		})
	}

	// Cek attempt yang sudah ada
	for _, att := range h.attempts {
		if att.ExamEventID == req.ExamID && (att.StudentID == req.StudentID || (req.NISN != "" && att.NISN == req.NISN)) {
			// Jika sudah SUBMITTED dan belum diberi izin remedial, tolak pembukaan ulang!
			if att.Status == "SUBMITTED" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success":      false,
					"is_submitted": true,
					"message":      "Anda telah mengumpulkan lembar jawaban ujian ini. Ujian tidak dapat dibuka kembali kecuali Pengampu/Proktor mengizinkan sesi remedial.",
					"data":         att,
				})
			}
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Melanjutkan sesi ujian aktif.",
				"data":    att,
			})
		}
	}

	// Buat attempt baru yang BERSIH (jawaban kosong)
	newAttempt := ExamAttemptDTO{
		ID:                 fmt.Sprintf("att-%d", time.Now().UnixNano()),
		ExamEventID:        req.ExamID,
		ExamTitle:          targetEvent.Title,
		SubjectName:        targetEvent.SubjectName,
		StudentID:          req.StudentID,
		StudentName:        req.StudentName,
		NISN:               req.NISN,
		StartedAt:          time.Now(),
		Status:             "IN_PROGRESS",
		CheatIncidentCount: 0,
		Answers:            make(map[string]string), // Bersih, tanpa pre-selected options
		Questions:          targetEvent.Questions,
	}

	h.attempts = append(h.attempts, newAttempt)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Bismillah, sesi ujian berhasil dimulai. Timer aktif dan sistem anti-cheat memantau layar.",
		"data":    newAttempt,
	})
}

// 3. Save Answer (Auto-Save)
type SaveAnswerRequest struct {
	AttemptID  string `json:"attempt_id"`
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

func (h *FiberTahap3Handler) SaveAnswer(c *fiber.Ctx) error {
	var req SaveAnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data jawaban tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, att := range h.attempts {
		if att.ID == req.AttemptID {
			if att.Status != "IN_PROGRESS" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "Ujian sudah disubmit dan tidak dapat diubah lagi.",
				})
			}
			if h.attempts[i].Answers == nil {
				h.attempts[i].Answers = make(map[string]string)
			}
			h.attempts[i].Answers[req.QuestionID] = req.Answer
			return c.JSON(fiber.Map{
				"success":  true,
				"message":  "Jawaban otomatis tersimpan di server.",
				"saved_at": time.Now().Format("15:04:05 WIB"),
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Sesi attempt ujian tidak ditemukan.",
	})
}

// 4. Log Proctor Incident (Deteksi Pindah Tab / Layar)
type LogIncidentRequest struct {
	AttemptID    string `json:"attempt_id"`
	StudentID    string `json:"student_id"`
	StudentName  string `json:"student_name"`
	IncidentType string `json:"incident_type"`
	Description  string `json:"description"`
}

func (h *FiberTahap3Handler) LogIncident(c *fiber.Ctx) error {
	var req LogIncidentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format log insiden tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	studentName := req.StudentName
	if studentName == "" {
		studentName = "Santri Peserta"
	}

	// Update counter pada attempt yang cocok
	for i, att := range h.attempts {
		if att.ID == req.AttemptID || (req.StudentID != "" && att.StudentID == req.StudentID) {
			h.attempts[i].CheatIncidentCount++
			studentName = att.StudentName
			break
		}
	}

	// Buat log insiden baru di feed pengawas
	newLog := ProctorLogDTO{
		ID:           fmt.Sprintf("prc-%d", time.Now().UnixNano()),
		AttemptID:    req.AttemptID,
		StudentID:    req.StudentID,
		StudentName:  studentName,
		IncidentType: req.IncidentType,
		Description:  req.Description,
		LoggedAt:     time.Now(),
	}
	h.proctorLogs = append([]ProctorLogDTO{newLog}, h.proctorLogs...)

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       "Insiden pelanggaran anti-cheat berhasil dicatat ke log proktor.",
		"incident_feed": h.proctorLogs,
	})
}

// 5. Finish Exam
type FinishExamRequest struct {
	AttemptID string            `json:"attempt_id"`
	StudentID string            `json:"student_id"`
	Answers   map[string]string `json:"answers"`
}

func (h *FiberTahap3Handler) FinishExam(c *fiber.Ctx) error {
	var req FinishExamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Permintaan selesai ujian tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, att := range h.attempts {
		if att.ID == req.AttemptID || (req.StudentID != "" && att.StudentID == req.StudentID) {
			now := time.Now()
			h.attempts[i].FinishedAt = &now
			h.attempts[i].Status = "SUBMITTED"

			if req.Answers != nil && len(req.Answers) > 0 {
				h.attempts[i].Answers = req.Answers
			}

			// Hitung Auto-Grade Pilihan Ganda (10 PG = 10 poin per soal)
			correctCount := 0
			for qID, ansKey := range h.attempts[i].Answers {
				switch qID {
				case "q-01":
					if ansKey == "A" { correctCount++ }
				case "q-02":
					if ansKey == "B" { correctCount++ }
				case "q-03":
					if ansKey == "B" { correctCount++ }
				case "q-04":
					if ansKey == "B" { correctCount++ }
				case "q-05":
					if ansKey == "B" { correctCount++ }
				case "q-06":
					if ansKey == "A" { correctCount++ }
				case "q-07":
					if ansKey == "B" { correctCount++ }
				case "q-08":
					if ansKey == "A" { correctCount++ }
				case "q-09":
					if ansKey == "A" { correctCount++ }
				case "q-10":
					if ansKey == "A" { correctCount++ }
				}
			}

			calculatedScore := float64(correctCount * 10)
			isPassed := calculatedScore >= 75.0

			h.attempts[i].TotalScore = &calculatedScore
			h.attempts[i].IsPassed = &isPassed

			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Ujian selesai. Skor pilihan ganda: %.0f / 100", calculatedScore),
				"data": fiber.Map{
					"attempt_id":       h.attempts[i].ID,
					"score_pg":         calculatedScore,
					"correct_count":    correctCount,
					"total_pg":         10,
					"is_passed":        isPassed,
					"letter_grade":     "A",
					"answers":          h.attempts[i].Answers,
					"submitted_at":     now.Format("15:04:05 WIB"),
				},
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Attempt ujian tidak ditemukan.",
	})
}

// 6. Proctoring Monitoring Grid
func (h *FiberTahap3Handler) GetProctorMonitoring(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	inProgress := 0
	submitted := 0
	for _, att := range h.attempts {
		if att.Status == "IN_PROGRESS" {
			inProgress++
		} else if att.Status == "SUBMITTED" {
			submitted++
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"active_exam":   h.examEvents[0],
			"attempts":      h.attempts,
			"proctor_logs":  h.proctorLogs,
			"incident_feed": h.proctorLogs,
			"stats": fiber.Map{
				"total_registered": len(h.attempts),
				"in_progress":      inProgress,
				"submitted":        submitted,
				"total_incidents":  len(h.proctorLogs),
			},
		},
	})
}

// 7. Reset Student Session
type ResetSessionReq struct {
	AttemptID string `json:"attempt_id"`
}

func (h *FiberTahap3Handler) ResetStudentSession(c *fiber.Ctx) error {
	var req ResetSessionReq
	if err := c.BodyParser(&req); err != nil {
		req.AttemptID = c.Params("id")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, att := range h.attempts {
		if att.ID == req.AttemptID {
			h.attempts[i].Status = "IN_PROGRESS"
			h.attempts[i].FinishedAt = nil
			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Sesi ujian santri %s berhasil di-reset oleh Proktor.", att.StudentName),
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Attempt santri tidak ditemukan.",
	})
}

// 8. Grant Remedial Permission (Izinkan Ujian Ulang Remedial)
func (h *FiberTahap3Handler) GrantRemedial(c *fiber.Ctx) error {
	var req ResetSessionReq
	if err := c.BodyParser(&req); err != nil {
		req.AttemptID = c.Params("id")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, att := range h.attempts {
		if att.ID == req.AttemptID {
			h.attempts[i].Status = "IN_PROGRESS"
			h.attempts[i].FinishedAt = nil
			h.attempts[i].TotalScore = nil
			h.attempts[i].IsPassed = nil
			h.attempts[i].CheatIncidentCount = 0
			h.attempts[i].Answers = make(map[string]string)
			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Izin remedial berhasil diberikan untuk santri %s. Lembar jawaban di-reset untuk pengerjaan ulang.", att.StudentName),
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Attempt santri tidak ditemukan.",
	})
}

// 9. Gradebook
func (h *FiberTahap3Handler) GetGradebook(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success":           true,
		"is_grade_locked":   h.isGradeLocked,
		"is_rapor_released": h.isRaporReleased,
		"subject_weights":   h.subjectWeights,
		"data":              h.gradebook,
	})
}

// 9b. Update Manual Grade (Input Nilai Lisan, Praktik, Adab, & Setoran Santri)
type UpdateManualGradeRequest struct {
	StudentID       string  `json:"student_id"`
	UtsScore        float64 `json:"uts_score"`
	UtsSource       string  `json:"uts_source"`
	UasScore        float64 `json:"uas_score"`
	UasSource       string  `json:"uas_source"`
	AdabScore       float64 `json:"adab_score"`
	TahfizhScore    float64 `json:"tahfizh_score"`
	AttendanceScore float64 `json:"attendance_score"`
	AssignmentScore float64 `json:"assignment_score"`
	HomeroomNotes   string  `json:"homeroom_notes"`
}

func (h *FiberTahap3Handler) UpdateManualGrade(c *fiber.Ctx) error {
	var req UpdateManualGradeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data penilaian manual tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var updatedEntry *GradebookEntryDTO
	for i, entry := range h.gradebook {
		if entry.StudentID == req.StudentID {
			if req.UtsScore > 0 {
				h.gradebook[i].UtsScore = req.UtsScore
			}
			if req.UtsSource != "" {
				h.gradebook[i].UtsSource = req.UtsSource
			}
			if req.UasScore > 0 {
				h.gradebook[i].UasScore = req.UasScore
			}
			if req.UasSource != "" {
				h.gradebook[i].UasSource = req.UasSource
			}
			if req.AdabScore > 0 {
				h.gradebook[i].AdabScore = req.AdabScore
			}
			if req.TahfizhScore > 0 {
				h.gradebook[i].TahfizhScore = req.TahfizhScore
			}
			if req.AttendanceScore > 0 {
				h.gradebook[i].AttendanceScore = req.AttendanceScore
			}
			if req.AssignmentScore > 0 {
				h.gradebook[i].AssignmentScore = req.AssignmentScore
			}

			// Recalculate composite final score
			// Formula: 10% Attendance + 20% Assignment + 30% Tahfizh + 20% UTS + 20% UAS
			final := (h.gradebook[i].AttendanceScore * 0.10) +
				(h.gradebook[i].AssignmentScore * 0.20) +
				(h.gradebook[i].TahfizhScore * 0.30) +
				(h.gradebook[i].UtsScore * 0.20) +
				(h.gradebook[i].UasScore * 0.20)
			h.gradebook[i].FinalScore = math.Round(final*10) / 10

			if final >= 90 {
				h.gradebook[i].LetterGrade = "A"
				h.gradebook[i].PassingStatus = "TUNTAS (MUMTAZ)"
			} else if final >= 75 {
				h.gradebook[i].LetterGrade = "B"
				h.gradebook[i].PassingStatus = "TUNTAS"
			} else {
				h.gradebook[i].LetterGrade = "C"
				h.gradebook[i].PassingStatus = "REMEDIAL"
			}

			copyEntry := h.gradebook[i]
			updatedEntry = &copyEntry

			// Sinkronisasi ke reportCards semester 3 santri bersangkutan
			for rIdx, rep := range h.reportCards {
				if rep.StudentID == req.StudentID && rep.SemesterCode == "sem-3" {
					if req.HomeroomNotes != "" {
						h.reportCards[rIdx].HomeroomNotes = req.HomeroomNotes
					}
					for gIdx, subj := range rep.Grades {
						if subj.SubjectCode == "THF-01" && req.TahfizhScore > 0 {
							h.reportCards[rIdx].Grades[gIdx].Score = req.TahfizhScore
						}
					}
				}
			}
			break
		}
	}

	if updatedEntry == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data santri tidak ditemukan dalam gradebook.",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Penilaian manual santri %s berhasil disimpan dan diperbarui!", updatedEntry.StudentName),
		"data":    updatedEntry,
	})
}

// 9c. Batch Save Ledger Penilaian Harian (Excel-Style Studio Gradebook)
type BatchLedgerItemDTO struct {
	StudentID   string    `json:"student_id"`
	PHScores    []float64 `json:"ph_scores"`
	RPH         float64   `json:"rph"`
	PTS         float64   `json:"pts"`
	PAS         float64   `json:"pas"`
	HPA         float64   `json:"hpa"`
	LetterGrade string    `json:"letter_grade"`
	Notes       string    `json:"notes"`
}

type BatchSaveLedgerRequest struct {
	SubjectCode  string               `json:"subject_code"`
	SemesterCode string               `json:"semester_code"`
	Grades       []BatchLedgerItemDTO `json:"grades"`
}

func (h *FiberTahap3Handler) BatchSaveLedger(c *fiber.Ctx) error {
	var req BatchSaveLedgerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Data leger nilai tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, item := range req.Grades {
		// Update gradebook entry
		for i, entry := range h.gradebook {
			if entry.StudentID == item.StudentID {
				if item.PTS > 0 {
					h.gradebook[i].UtsScore = math.Min(100.0, math.Max(0.0, item.PTS))
				}
				if item.PAS > 0 {
					h.gradebook[i].UasScore = math.Min(100.0, math.Max(0.0, item.PAS))
				}
				if item.RPH > 0 {
					h.gradebook[i].AssignmentScore = math.Min(100.0, math.Max(0.0, item.RPH))
				}
				// Recalculate composite final score
				final := (h.gradebook[i].AttendanceScore * 0.10) +
					(h.gradebook[i].AssignmentScore * 0.20) +
					(h.gradebook[i].TahfizhScore * 0.30) +
					(h.gradebook[i].UtsScore * 0.20) +
					(h.gradebook[i].UasScore * 0.20)
				final = math.Min(100.0, math.Max(0.0, final))
				h.gradebook[i].FinalScore = math.Round(final*10) / 10
				if final >= 90 {
					h.gradebook[i].LetterGrade = "A"
					h.gradebook[i].PassingStatus = "TUNTAS (MUMTAZ)"
				} else if final >= 75 {
					h.gradebook[i].LetterGrade = "B"
					h.gradebook[i].PassingStatus = "TUNTAS"
				} else {
					h.gradebook[i].LetterGrade = "C"
					h.gradebook[i].PassingStatus = "REMEDIAL"
				}
				break
			}
		}

		// Update subject grade in reportCards
		for rIdx, rep := range h.reportCards {
			if rep.StudentID == item.StudentID && (req.SemesterCode == "" || rep.SemesterCode == req.SemesterCode) {
				foundSubject := false
				for gIdx, subj := range rep.Grades {
					if subj.SubjectCode == req.SubjectCode {
						h.reportCards[rIdx].Grades[gIdx].Score = math.Min(100.0, math.Max(0.0, item.HPA))
						h.reportCards[rIdx].Grades[gIdx].LetterGrade = item.LetterGrade
						if item.Notes != "" {
							h.reportCards[rIdx].Grades[gIdx].CompetencyNotes = item.Notes
						}
						foundSubject = true
						break
					}
				}
				// If subject not in report card list, add it
				if !foundSubject && item.HPA > 0 {
					subjName := req.SubjectCode
					tutorName := "Ustadzah Fatimah, S.Pd."
					for _, w := range h.subjectWeights {
						if w.SubjectCode == req.SubjectCode {
							subjName = w.SubjectName
							break
						}
					}
					h.reportCards[rIdx].Grades = append(h.reportCards[rIdx].Grades, ReportCardSubjectDTO{
						SubjectCode:      req.SubjectCode,
						SubjectName:      subjName,
						TutorName:        tutorName,
						Score:            item.HPA,
						LetterGrade:      item.LetterGrade,
						PassingThreshold: 75.0,
						CompetencyNotes:  item.Notes,
						SourceType:       "KURIKULUM_TERPADU",
					})
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Berhasil menyimpan leger penilaian untuk mata pelajaran %s (%d santri)", req.SubjectCode, len(req.Grades)),
		"data":    h.gradebook,
	})
}

// 10. Get & Update Subject Weights
func (h *FiberTahap3Handler) GetSubjectWeights(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    h.subjectWeights,
	})
}

func (h *FiberTahap3Handler) UpdateSubjectWeights(c *fiber.Ctx) error {
	var req SubjectWeightConfigDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format bobot mata pelajaran tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, w := range h.subjectWeights {
		if w.SubjectCode == req.SubjectCode {
			h.subjectWeights[i] = req
			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Bobot penilaian mapel '%s' berhasil diperbarui.", req.SubjectName),
				"data":    h.subjectWeights[i],
			})
		}
	}

	h.subjectWeights = append(h.subjectWeights, req)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Konfigurasi bobot mapel baru berhasil ditambahkan.",
		"data":    req,
	})
}

// 11. Lock / Unlock Grades
func (h *FiberTahap3Handler) ToggleLockGrades(c *fiber.Ctx) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.isGradeLocked = !h.isGradeLocked
	statusMsg := "Buku nilai resmi DIKUNCI untuk pencetakan e-Rapor resmi."
	if !h.isGradeLocked {
		statusMsg = "Buku nilai DIBUKA KEMBALI untuk penyesuaian nilai oleh Guru."
	}

	return c.JSON(fiber.Map{
		"success":         true,
		"is_grade_locked": h.isGradeLocked,
		"message":         statusMsg,
	})
}

// 12. Toggle Rilis Pembagian Rapor ke Santri
func (h *FiberTahap3Handler) ToggleReleaseRapor(c *fiber.Ctx) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.isRaporReleased = !h.isRaporReleased
	msg := "Pembagian E-Rapor ke Santri & Orang Tua BERHASIL DIBUKA."
	if !h.isRaporReleased {
		msg = "Akses E-Rapor Santri DITUTUP SEMENTARA (Menunggu Pengesahan Pleno Guru)."
	}

	return c.JSON(fiber.Map{
		"success":           true,
		"is_rapor_released": h.isRaporReleased,
		"message":           msg,
	})
}

// 13. Get Report Cards (Multi-Semester Filter)
func (h *FiberTahap3Handler) GetReportCards(c *fiber.Ctx) error {
	semesterCode := c.Query("semester")
	studentID := c.Query("student_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var filtered []StudentReportCardDTO
	for _, rep := range h.reportCards {
		if semesterCode != "" && rep.SemesterCode != semesterCode {
			continue
		}
		if studentID != "" && rep.StudentID != studentID {
			continue
		}
		filtered = append(filtered, rep)
	}

	if len(filtered) == 0 && semesterCode == "" {
		filtered = h.reportCards
	}

	return c.JSON(fiber.Map{
		"success":           true,
		"is_rapor_released": h.isRaporReleased,
		"count":             len(filtered),
		"data":              filtered,
	})
}

// 14. Verify Public QR Code
func (h *FiberTahap3Handler) VerifyReportCardQR(c *fiber.Ctx) error {
	qrCode := c.Params("qr_code")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, rep := range h.reportCards {
		if rep.QRVerificationCode == qrCode {
			return c.JSON(fiber.Map{
				"success":             true,
				"is_authentic":        true,
				"institution":         "Pusat Kegiatan Belajar Masyarakat (PKBM) Digital Dar El Iman",
				"accreditation":       "Terakreditasi B — Kemendikbudristek RI",
				"qr_code":             qrCode,
				"student_name":        rep.StudentName,
				"nisn":                rep.NISN,
				"classroom":           rep.ClassroomName,
				"academic_year":       rep.Semester,
				"published_at":        rep.PublishedAt,
				"verification_status": "DOKUMEN RESMI DAN SAH YAYASAN DAR EL IMAN PADANG",
				"data":                rep,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success":      false,
		"is_authentic": false,
		"message":      "Kode QR E-Rapor tidak ditemukan atau tidak valid dalam pangkalan data PKBM Digital Dar El Iman.",
	})
}

func ptrTime(t time.Time) *time.Time {
	return &t
}


func ptrBool(b bool) *bool {
	return &b
}
