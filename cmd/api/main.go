package main

import (
	"log"
	"net/http"

	"pkbm-backend/internal/config"
	"pkbm-backend/internal/database"
	"pkbm-backend/internal/handlers"
	"pkbm-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Catatan: Server tetap berjalan dengan mode fallback jika MySQL belum online: %v", err)
	}

	r := chi.NewRouter()

	// Middlewares
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Handlers
	authH := &handlers.AuthHandler{DB: db, JWTSecret: cfg.JWTSecret}
	instH := &handlers.InstitutionHandler{DB: db}
	tutorH := &handlers.TutorHandler{DB: db}
	studentH := &handlers.StudentHandler{DB: db}
	classH := &handlers.ClassroomHandler{DB: db}

	// Tahap 2 Handlers
	currH := handlers.NewCurriculumHandler(db)
	lmsH := handlers.NewLMSHandler(db)
	asgH := handlers.NewAssignmentHandler(db)
	attH := handlers.NewAttendanceHandler(db)
	thfH := handlers.NewTahfizhHandler(db)

	// Tahap 3 Handlers (CBT, Proctoring, Gradebook, and E-Rapor)
	examH := handlers.NewExamHandler(db)
	gradingH := handlers.NewGradingHandler(db)
	reportH := handlers.NewReportCardHandler(db)

	// Tahap 4 Handlers (PPDB, Keuangan SPP, and Ekspor Dapodik)
	ppdbH := handlers.NewPPDBHandler(db)
	financeH := handlers.NewFinanceHandler(db)
	dapodikH := handlers.NewDapodikHandler(db)

	// Tahap 5 Handlers (Kemitraan, Penjaminan Mutu, BSC & Tutor AI)
	partnerH := handlers.NewPartnerHandler(db)
	qualityH := handlers.NewQualityHandler(db)
	aiTutorH := handlers.NewAITutorHandler(db)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth
		r.Post("/auth/login", authH.Login)

		// 1. Kurikulum & Mapel
		r.Get("/curriculum/subjects", currH.ListSubjects)
		r.Post("/curriculum/subjects", currH.CreateSubject)
		r.Get("/curriculum/course-classes", currH.ListCourseClasses)
		r.Post("/curriculum/course-classes", currH.CreateCourseClass)
		r.Get("/curriculum/schedules/classroom/{classroom_id}", currH.GetClassroomSchedules)

		// 2. Ruang Kelas LMS & Materi
		r.Get("/lms/my-courses", lmsH.MyCourses)
		r.Get("/lms/courses/{course_class_id}/topics", lmsH.GetTopics)
		r.Post("/lms/courses/{course_class_id}/topics", lmsH.CreateTopic)
		r.Post("/lms/topics/{topic_id}/materials", lmsH.AddMaterial)
		r.Put("/lms/topics/{topic_id}/publish", lmsH.TogglePublishTopic)

		// 3. Tugas & Submission
		r.Post("/assignments/topics/{topic_id}/assignments", asgH.CreateAssignment)
		r.Get("/assignments/{assignment_id}", asgH.GetAssignment)
		r.Post("/assignments/{assignment_id}/submit", asgH.SubmitAssignment)
		r.Get("/assignments/{assignment_id}/submissions", asgH.GetSubmissions)
		r.Put("/assignments/submissions/{submission_id}/grade", asgH.GradeSubmission)

		// 4. Presensi KBM
		r.Post("/attendances/courses/{course_class_id}/sessions", attH.CreateSession)
		r.Get("/attendances/sessions/{session_id}", attH.GetSession)
		r.Post("/attendances/sessions/{session_id}/record", attH.RecordAttendance)
		r.Post("/attendances/sessions/{session_id}/check-in", attH.CheckIn)

		// 5. Halaqah Tahfizh Al-Qur'an
		r.Get("/tahfizh/groups", thfH.ListGroups)
		r.Post("/tahfizh/groups", thfH.CreateGroup)
		r.Post("/tahfizh/groups/{id}/students", thfH.AddStudentToGroup)
		r.Get("/tahfizh/students/{student_id}/logs", thfH.GetStudentLogs)
		r.Post("/tahfizh/logs", thfH.RecordLog)

		// 6. Bank Soal, Engine Ujian CBT & Proctoring (Tahap 3)
		r.Get("/exams/banks", examH.ListQuestionBanks)
		r.Post("/exams/banks", examH.CreateQuestionBank)
		r.Post("/exams/banks/{id}/questions", examH.AddQuestion)
		r.Get("/exams/events", examH.ListExamEvents)
		r.Post("/exams/events", examH.CreateExamEvent)
		r.Post("/exams/events/{id}/start", examH.StartExam)
		r.Post("/exams/attempts/{attempt_id}/answer", examH.SaveAnswer)
		r.Post("/exams/attempts/{attempt_id}/log-incident", examH.LogProctorIncident)
		r.Post("/exams/attempts/{attempt_id}/finish", examH.FinishExam)
		r.Get("/exams/events/{id}/monitoring", examH.ProctorMonitoring)

		// 7. Buku Nilai & Remedial (Tahap 3)
		r.Get("/grading/courses/{id}/gradebook", gradingH.GetGradebook)
		r.Post("/grading/courses/{id}/calculate-scores", gradingH.CalculateScores)
		r.Get("/grading/courses/{id}/remedial-candidates", gradingH.GetRemedialCandidates)
		r.Post("/grading/courses/{id}/lock-grades", gradingH.LockGrades)

		// 8. E-Rapor Digital & Verifikasi QR Publik (Tahap 3)
		r.Post("/report-cards/generate-batch", reportH.GenerateBatch)
		r.Get("/report-cards/student/{student_id}/semester/{year_id}", reportH.GetStudentReportCard)
		r.Get("/report-cards/verify/{qr_code}", reportH.VerifyQRCode)

		// 9. PPDB Online Publik & Panitia (Tahap 4)
		r.Get("/ppdb/periods/active", ppdbH.GetActivePeriod)
		r.Post("/ppdb/register", ppdbH.Register)
		r.Post("/ppdb/{registration_id}/documents", ppdbH.UploadDocument)
		r.Get("/ppdb/my-status", ppdbH.GetMyStatus)
		r.Get("/ppdb/admin/registrations", ppdbH.ListRegistrations)
		r.Put("/ppdb/admin/registrations/{id}/verify", ppdbH.VerifyRegistration)
		r.Post("/ppdb/admin/registrations/{id}/accept-and-convert", ppdbH.AcceptAndConvert)

		// 10. Manajemen Keuangan & SPP (Tahap 4)
		r.Get("/finance/overview", financeH.GetOverview)
		r.Get("/finance/tariffs", financeH.ListTariffs)
		r.Post("/finance/invoices/generate-monthly", financeH.GenerateMonthlyInvoices)
		r.Get("/finance/invoices/student/{student_id}", financeH.GetStudentInvoices)
		r.Post("/finance/invoices/{id}/pay", financeH.PayInvoice)
		r.Get("/finance/payments/pending", financeH.GetPendingPayments)
		r.Put("/finance/payments/{id}/verify", financeH.VerifyPayment)
		r.Get("/finance/scholarships", financeH.ListScholarships)
		r.Post("/finance/scholarships", financeH.GrantScholarship)

		// 11. Integrasi & Ekspor Dapodik Kemendikbudristek (Tahap 4)
		r.Get("/dapodik/validation-check", dapodikH.ValidationCheck)
		r.Get("/dapodik/export-students", dapodikH.ExportStudents)
		r.Post("/dapodik/export-students", dapodikH.ExportStudents)
		r.Get("/dapodik/export-tutors", dapodikH.ExportTutors)
		r.Post("/dapodik/export-tutors", dapodikH.ExportTutors)
		r.Get("/dapodik/export-classrooms", dapodikH.ExportClassrooms)
		r.Post("/dapodik/export-classrooms", dapodikH.ExportClassrooms)
		r.Get("/dapodik/logs", dapodikH.GetLogs)

		// 12. Kemitraan Pesantren & Komunitas Diaspora (Tahap 5)
		r.Get("/partners", partnerH.ListPartners)
		r.Post("/partners", partnerH.CreatePartner)
		r.Get("/partners/{id}/students", partnerH.GetPartnerStudents)
		r.Post("/partners/{id}/bulk-enroll", partnerH.BulkEnrollPartnerStudents)
		r.Get("/partners/analytics/diaspora-map", partnerH.GetDiasporaMap)

		// 13. Penjaminan Mutu, Akreditasi BAN-PDM & BSC (Tahap 5)
		r.Get("/quality/accreditation/checklist", qualityH.GetAccreditationChecklist)
		r.Put("/quality/accreditation/items/{id}", qualityH.UpdateAccreditationItem)
		r.Get("/quality/sops", qualityH.ListSOPs)
		r.Post("/quality/sops", qualityH.CreateSOP)
		r.Get("/quality/risks", qualityH.ListRisks)
		r.Post("/quality/risks", qualityH.CreateRisk)
		r.Get("/quality/bsc/dashboard", qualityH.GetBSCDashboard)

		// 14. Asisten Cerdas Tutor AI (Tahap 5)
		r.Post("/ai-tutor/sessions", aiTutorH.CreateSession)
		r.Post("/ai-tutor/sessions/{id}/ask", aiTutorH.AskAI)
		r.Get("/ai-tutor/sessions/{id}/history", aiTutorH.GetSessionHistory)

		// Protected Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			r.Get("/auth/me", authH.Me)

			// Master & Institution
			r.Get("/master/overview", instH.GetOverview)
			r.Get("/master/academic-years", instH.GetAcademicYears)

			// Tutors
			r.Get("/tutors", tutorH.List)
			r.Post("/tutors", tutorH.Create)

			// Students
			r.Get("/students", studentH.List)
			r.Post("/students", studentH.Create)
			r.Post("/students/bulk-import", studentH.BulkImportCSV)

			// Classrooms
			r.Get("/classrooms", classH.List)
			r.Post("/classrooms", classH.Create)
			r.Post("/classrooms/{id}/students", classH.AssignStudent)
		})
	})

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"online","service":"PKBM Digital Dar El Iman API (Go + MySQL)","stage":5}`))
	})

	log.Printf("Backend PKBM Digital Dar El Iman (Go) aktif di port :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
