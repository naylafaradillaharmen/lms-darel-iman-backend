package main

import (
	"fmt"
	"log"
	"os"

	"pkbm-backend/internal/config"
	"pkbm-backend/internal/database"
	"pkbm-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	// Inisialisasi Database GORM (dengan graceful fallback katalog demo jika MySQL offline)
	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Catatan GORM: Menggunakan mode in-memory: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "PKBM Digital Dar El Iman API v2.0 (Fiber + GORM)",
		ServerHeader: "Fiber",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://127.0.0.1:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: true,
	}))

	// Fiber Auth Handler
	var gormDB = db.GormDB
	authH := handlers.NewFiberAuthHandler(gormDB, cfg.JWTSecret)

	// API Routes v1
	v1 := app.Group("/api/v1")

	// System Health & Info
	v1.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "connected"
		if db.GormDB == nil {
			dbStatus = "fallback_mock (MySQL belum aktif di port " + cfg.DBPort + ")"
		}
		return c.JSON(fiber.Map{
			"success":   true,
			"service":   "PKBM Digital Dar El Iman - Backend Core",
			"framework": "Go Fiber v2 + GORM",
			"database":  cfg.DBName,
			"db_status": dbStatus,
			"timestamp": fiber.Map{
				"timezone": "WIB (UTC+7)",
			},
		})
	})

	// Auth Endpoints
	auth := v1.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/register", authH.Register)
	auth.Post("/ppdb-guest", authH.PPDBGuestLogin)
	auth.Get("/demo-accounts", authH.GetDemoAccounts)
	auth.Post("/change-password", authH.ChangePassword)

	// Tahap 2 Handler (KBM, Kurikulum, LMS, Tahfizh)
	tahap2H := handlers.NewFiberTahap2Handler(gormDB)

	// Curriculum & KBM Routes
	curr := v1.Group("/curriculum")
	curr.Get("/subjects", tahap2H.GetSubjects)
	curr.Get("/classes", tahap2H.GetClasses)
	curr.Get("/schedules", tahap2H.GetSchedules)
	curr.Post("/schedules/meet", tahap2H.UpdateMeetingLink)
	curr.Get("/attendance", tahap2H.GetAttendance)
	curr.Post("/attendance", tahap2H.SaveAttendance)
	curr.Post("/attendance/check-in", tahap2H.StudentCheckIn)
	curr.Post("/attendance/approve", tahap2H.ApproveAttendance)

	// LMS Topics & Materials
	lms := v1.Group("/lms")
	lms.Get("/topics", tahap2H.GetTopics)

	// Assignments & Submissions
	asg := v1.Group("/assignments")
	asg.Get("/", tahap2H.GetAssignments)
	asg.Get("/submissions", tahap2H.GetSubmissions)
	asg.Post("/submit", tahap2H.SubmitAssignment)
	asg.Post("/grade", tahap2H.GradeSubmission)

	// Tahfizh Halaqah & Mutaba'ah
	thf := v1.Group("/tahfizh")
	thf.Get("/groups", tahap2H.GetTahfizhGroups)
	thf.Get("/logs", tahap2H.GetTahfizhLogs)
	thf.Post("/logs", tahap2H.CreateTahfizhLog)
	thf.Post("/evaluate", tahap2H.EvaluateTahfizhLog)

	// Tahap 3 Handler (CBT Online, Proctoring Grid, Gradebook & E-Rapor Ber-QR)
	tahap3H := handlers.NewFiberTahap3Handler(gormDB)

	// 1. CBT Exam Engine
	exam := v1.Group("/exams")
	exam.Get("/", tahap3H.GetExamEvents)
	exam.Post("/create", tahap3H.CreateExamEvent)
	exam.Post("/start", tahap3H.StartExam)
	exam.Post("/save-answer", tahap3H.SaveAnswer)
	exam.Post("/log-incident", tahap3H.LogIncident)
	exam.Post("/finish", tahap3H.FinishExam)

	// 2. Proctoring & Anti-Cheat Monitoring
	proctor := v1.Group("/proctor")
	proctor.Get("/monitoring", tahap3H.GetProctorMonitoring)
	proctor.Post("/reset-session", tahap3H.ResetStudentSession)
	proctor.Post("/grant-remedial", tahap3H.GrantRemedial)

	// 3. Gradebook (Buku Nilai Terpadu & Bobot Fleksibel)
	grading := v1.Group("/grading")
	grading.Get("/gradebook", tahap3H.GetGradebook)
	grading.Post("/manual-grade", tahap3H.UpdateManualGrade)
	grading.Post("/batch-ledger", tahap3H.BatchSaveLedger)
	grading.Get("/weights", tahap3H.GetSubjectWeights)
	grading.Post("/weights", tahap3H.UpdateSubjectWeights)
	grading.Post("/toggle-lock", tahap3H.ToggleLockGrades)

	// 4. E-Rapor Digital & Public Verification
	rapor := v1.Group("/report-cards")
	rapor.Get("/", tahap3H.GetReportCards)
	rapor.Post("/toggle-release", tahap3H.ToggleReleaseRapor)
	rapor.Get("/verify/:qr_code", tahap3H.VerifyReportCardQR)
	v1.Get("/verify-rapor/:qr_code", tahap3H.VerifyReportCardQR)

	// Tahap 4 Handler (PPDB Online, Keuangan SPP, & Ekspor Dapodik)
	tahap4H := handlers.NewFiberTahap4Handler(gormDB)

	// 1. Modul PPDB Online (Nasional & Diaspora)
	ppdb := v1.Group("/ppdb")
	ppdb.Get("/periods/active", tahap4H.GetActivePeriod)
	ppdb.Post("/register", tahap4H.RegisterPublic)
	ppdb.Post("/:registration_id/documents", tahap4H.UploadDocument)
	ppdb.Get("/my-status", tahap4H.GetMyStatus)
	ppdb.Get("/admin/registrations", tahap4H.GetAdminRegistrations)
	ppdb.Get("/registrations", tahap4H.GetAdminRegistrations)
	ppdb.Put("/admin/registrations/:id/verify", tahap4H.VerifyRegistration)
	ppdb.Put("/registrations/:id/verify", tahap4H.VerifyRegistration)
	ppdb.Post("/admin/registrations/:id/accept-and-convert", tahap4H.AcceptAndConvertStudent)
	ppdb.Post("/registrations/:id/accept-and-convert", tahap4H.AcceptAndConvertStudent)

	// 2. Modul Keuangan SPP & Transaksi
	finance := v1.Group("/finance")
	finance.Get("/tariffs", tahap4H.GetTariffs)
	finance.Post("/invoices/generate-monthly", tahap4H.GenerateMonthlyInvoices)
	finance.Get("/invoices", tahap4H.GetAllInvoices)
	finance.Get("/invoices/student/:student_id", tahap4H.GetStudentInvoices)
	finance.Post("/invoices/:id/pay", tahap4H.ConfirmPayment)
	finance.Get("/payments/pending", tahap4H.GetPendingPayments)
	finance.Put("/payments/:id/verify", tahap4H.VerifyPayment)
	finance.Get("/scholarships", tahap4H.GetScholarships)
	finance.Post("/scholarships", tahap4H.CreateScholarship)
	finance.Get("/summary", tahap4H.GetFinanceSummary)

	// 3. Modul Integrasi & Ekspor Dapodik
	dapodik := v1.Group("/dapodik")
	dapodik.Get("/validation-check", tahap4H.GetValidationAudit)
	dapodik.Get("/validation-audit", tahap4H.GetValidationAudit)
	dapodik.Get("/export-students", tahap4H.ExportStudents)
	dapodik.Post("/export-students", tahap4H.ExportStudents)
	dapodik.Get("/export-tutors", tahap4H.ExportTutors)
	dapodik.Post("/export-tutors", tahap4H.ExportTutors)
	dapodik.Get("/export-classrooms", tahap4H.ExportClassrooms)
	dapodik.Post("/export-classrooms", tahap4H.ExportClassrooms)
	dapodik.Get("/export-logs", tahap4H.GetExportLogs)

	// Tahap 5 Handler (Kemitraan, Penjaminan Mutu BAN-PDM, BSC, dan Tutor AI)
	tahap5H := handlers.NewFiberTahap5Handler(gormDB)

	// 1. Modul Kemitraan & Diaspora
	partners := v1.Group("/partners")
	partners.Get("/", tahap5H.GetPartners)
	partners.Post("/", tahap5H.CreatePartner)
	partners.Post("/:id/bulk-enroll", tahap5H.BulkEnrollStudents)
	partners.Get("/analytics/diaspora-map", tahap5H.GetDiasporaMapAnalytics)

	// 2. Modul Penjaminan Mutu, SOP, Risiko & BSC
	quality := v1.Group("/quality")
	quality.Get("/accreditation/checklist", tahap5H.GetAccreditationChecklist)
	quality.Put("/accreditation/items/:id", tahap5H.UpdateAccreditationItem)
	quality.Get("/sops", tahap5H.GetSOPs)
	quality.Post("/risks", tahap5H.CreateRiskRegister)
	quality.Get("/bsc/dashboard", tahap5H.GetBSCDashboard)

	// 3. Modul Asisten Cerdas Tutor AI
	aiTutor := v1.Group("/ai-tutor")
	aiTutor.Post("/sessions", tahap5H.CreateAISession)
	aiTutor.Post("/sessions/:id/ask", tahap5H.AskAITutor)
	aiTutor.Get("/sessions/:id/history", tahap5H.GetAISessionHistory)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 PKBM Digital Dar El Iman Go Fiber server berjalan di http://localhost%s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Gagal menjalankan server Fiber: %v", err)
	}
}
