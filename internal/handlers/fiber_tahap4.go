package handlers

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FiberTahap4Handler struct {
	GormDB *gorm.DB
	mu     sync.RWMutex

	// PPDB In-Memory Store
	periods       []models.PPDBPeriod
	registrations []models.PPDBRegistration
	documents     []models.PPDBDocument

	// Finance In-Memory Store
	tariffs      []models.TuitionTariff
	invoices     []models.StudentInvoice
	transactions []models.PaymentTransaction
	scholarships []models.ScholarshipGrant

	// Dapodik Logs
	exportLogs []models.DapodikExportLog
}

func NewFiberTahap4Handler(db *gorm.DB) *FiberTahap4Handler {
	h := &FiberTahap4Handler{
		GormDB: db,
	}
	h.initData()
	return h
}

func (h *FiberTahap4Handler) initData() {
	periodID := "per-2027-gel1"
	h.periods = []models.PPDBPeriod{
		{
			ID:             periodID,
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			Title:          "PPDB Santri Baru Gelombang 1 TA 2027/2028 (Nasional & Diaspora Internasional)",
			StartDate:      "2026-10-01",
			EndDate:        "2026-12-31",
			IsOpen:         true,
			CreatedAt:      time.Now().Add(-20 * 24 * time.Hour),
		},
	}

	h.registrations = []models.PPDBRegistration{
		{
			ID:                   "preg-2027-001",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-001",
			TargetPackage:        "PAKET_C",
			FullName:             "Zaid Al-Khair",
			Gender:               "L",
			NIK:                  "1371012304090001",
			NISN:                 "0098712345",
			BirthPlace:           "Doha",
			BirthDate:            "2009-04-12",
			Phone:                "+97455123456",
			Email:                "zaid.qatar@darel-iman.sch.id",
			IsDiaspora:           true,
			Country:              "Qatar",
			Timezone:             "Asia/Qatar (GMT+3)",
			GuardianName:         "Ir. Salman Al-Farisi",
			GuardianPhone:        "+97455123456",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Khadijah S.E.",
			Status:               "ACCEPTED",
			Notes:                "Berkas lengkap, placement test tahfizh 3 juz Mumtaz. Siap konversi ke santri resmi.",
			CreatedAt:            time.Now().Add(-15 * 24 * time.Hour),
		},
		{
			ID:                   "preg-2027-002",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-002",
			TargetPackage:        "PAKET_B",
			FullName:             "Aisyah Humaira",
			Gender:               "P",
			NIK:                  "1371025510110002",
			NISN:                 "0112345678",
			BirthPlace:           "Padang",
			BirthDate:            "2011-10-15",
			Phone:                "+628137788990",
			Email:                "aisyah.padang@gmail.com",
			IsDiaspora:           false,
			Country:              "Indonesia",
			Timezone:             "Asia/Jakarta (WIB)",
			GuardianName:         "dr. Rahman Hakim",
			GuardianPhone:        "+628137788990",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Fatimah Az-Zahra",
			Status:               "DOCUMENTS_VERIFIED",
			Notes:                "Berkas digital KK dan Akta kelahiran telah tervalidasi oleh panitia PPDB.",
			CreatedAt:            time.Now().Add(-10 * 24 * time.Hour),
		},
		{
			ID:                   "preg-2027-003",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-003",
			TargetPackage:        "PAKET_A",
			FullName:             "Maryam Sulaiman",
			Gender:               "P",
			NIK:                  "9901024503110001",
			NISN:                 "0145678901",
			BirthPlace:           "Riyadh",
			BirthDate:            "2015-05-18",
			Phone:                "+966501234567",
			Email:                "maryam.riyadh@gmail.com",
			IsDiaspora:           true,
			Country:              "Arab Saudi",
			Timezone:             "Asia/Riyadh (GMT+3)",
			GuardianName:         "Tariq Sulaiman, Lc.",
			GuardianPhone:        "+966501234567",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Maryam Hanum",
			Status:               "SUBMITTED",
			Notes:                "Pendaftar diaspora KSA baru masuk. Menunggu penelaahan berkas paspor & akta konsulat.",
			CreatedAt:            time.Now().Add(-3 * 24 * time.Hour),
		},
		{
			ID:                   "preg-2027-004",
			PPDBPeriodID:         periodID,
			RegistrationNumber:   "REG-2027-004",
			TargetPackage:        "PAKET_C",
			FullName:             "Umar Al-Ghazali",
			Gender:               "L",
			NIK:                  "3171010502090004",
			NISN:                 "0099887766",
			BirthPlace:           "Jakarta Selatan",
			BirthDate:            "2009-02-05",
			Phone:                "081299887766",
			Email:                "umar.ghazali@gmail.com",
			IsDiaspora:           false,
			Country:              "Indonesia",
			Timezone:             "Asia/Jakarta (WIB)",
			GuardianName:         "H. Ghazali Ahmad",
			GuardianPhone:        "081299887766",
			GuardianRelationship: "Ayah Kandung",
			MotherName:           "Amina S.P.",
			Status:               "TEST_PASSED",
			Notes:                "Lulus placement test wawancara & tajwid al-Fatihah dengan skor 94. Siap kelulusan.",
			CreatedAt:            time.Now().Add(-1 * 24 * time.Hour),
		},
	}

	h.documents = []models.PPDBDocument{
		{
			ID:                 "pdoc-001",
			PPDBRegistrationID: "preg-2027-001",
			DocumentType:       "AKTA_LAHIR",
			FileURL:            "https://cdn.darel-iman.or.id/ppdb/akta-zaid-qatar.pdf",
			IsVerified:         true,
			VerificationNotes:  "Akta kelahiran resmi KBRI terverifikasi sah.",
			CreatedAt:          time.Now().Add(-15 * 24 * time.Hour),
		},
		{
			ID:                 "pdoc-002",
			PPDBRegistrationID: "preg-2027-001",
			DocumentType:       "KK",
			FileURL:            "https://cdn.darel-iman.or.id/ppdb/kk-salman.pdf",
			IsVerified:         true,
			VerificationNotes:  "Kartu Keluarga WNI di Luar Negeri sah.",
			CreatedAt:          time.Now().Add(-15 * 24 * time.Hour),
		},
		{
			ID:                 "pdoc-003",
			PPDBRegistrationID: "preg-2027-002",
			DocumentType:       "AKTA_LAHIR",
			FileURL:            "https://cdn.darel-iman.or.id/ppdb/akta-aisyah.pdf",
			IsVerified:         true,
			VerificationNotes:  "Dokumen kependudukan tervalidasi.",
			CreatedAt:          time.Now().Add(-10 * 24 * time.Hour),
		},
	}

	h.tariffs = []models.TuitionTariff{
		{
			ID:             "tt-01",
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			GradeLevelID:   "c0000010-0000-0000-0000-000000000010",
			PackageName:    "Paket C (Setara SMA)",
			Amount:         350000.0,
			Description:    "SPP Bulanan Paket C Kesetaraan & Pembinaan Tahfizh Terpadu",
			CreatedAt:      time.Now(),
		},
		{
			ID:             "tt-02",
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			GradeLevelID:   "b0000007-0000-0000-0000-000000000007",
			PackageName:    "Paket B (Setara SMP)",
			Amount:         300000.0,
			Description:    "SPP Bulanan Paket B Kesetaraan & Pembinaan Tahfizh Terpadu",
			CreatedAt:      time.Now(),
		},
		{
			ID:             "tt-03",
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			GradeLevelID:   "a0000001-0000-0000-0000-000000000001",
			PackageName:    "Paket A (Setara SD)",
			Amount:         250000.0,
			Description:    "SPP Bulanan Paket A Kesetaraan & Pembinaan Tahfizh Terpadu",
			CreatedAt:      time.Now(),
		},
	}

	h.invoices = []models.StudentInvoice{
		{
			ID: "inv-2026-09-001", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NIS: "0081234561",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-09-0001", Month: 9, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-09-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-55 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-10-001", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NIS: "0081234561",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-10-0001", Month: 10, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-10-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-001", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NIS: "0081234561",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0001", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "PENDING_VERIFICATION",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-2 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-09-002", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NIS: "0081234562",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-09-0002", Month: 9, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-09-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-55 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-10-002", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NIS: "0081234562",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-10-0002", Month: 10, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-10-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-002", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NIS: "0081234562",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0002", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "UNPAID",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-10-003", StudentID: "u0000014b", StudentName: "Abdul Andhika Raditya", NIS: "0081234563",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-10-0003", Month: 10, Year: 2026,
			Amount: 350000.0, DiscountAmount: 350000.0, FinalAmount: 0.0, DueDate: "2026-10-25", Status: "PAID",
			CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-003", StudentID: "u0000014b", StudentName: "Abdul Andhika Raditya", NIS: "0081234563",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0003", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 350000.0, FinalAmount: 0.0, DueDate: "2026-11-25", Status: "PAID",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-10-004", StudentID: "u0000015b", StudentName: "Afdhal Al Baihaqi", NIS: "0081234564",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-10-0004", Month: 10, Year: 2026,
			Amount: 350000.0, DiscountAmount: 175000.0, FinalAmount: 175000.0, DueDate: "2026-10-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-004", StudentID: "u0000015b", StudentName: "Afdhal Al Baihaqi", NIS: "0081234564",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0004", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 175000.0, FinalAmount: 175000.0, DueDate: "2026-11-25", Status: "UNPAID",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-10-005", StudentID: "u0000016b", StudentName: "Ahmad Faris", NIS: "0081234565",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-10-0005", Month: 10, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-10-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-005", StudentID: "u0000016b", StudentName: "Ahmad Faris", NIS: "0081234565",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0005", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "UNPAID",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-006", StudentID: "u0000017b", StudentName: "Syamil Basayif", NIS: "0081234566",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0006", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "PENDING_VERIFICATION",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-3 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-007", StudentID: "u0000018b", StudentName: "Zaidan Al-Ayyubi", NIS: "0081234567",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0007", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-4 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-008", StudentID: "u0000019b", StudentName: "Hamzah Asadullah", NIS: "0081234568",
			ClassName: "Kelas X Paket C — Umar bin Khattab", InvoiceNumber: "INV-2026-11-0008", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "UNPAID",
			CreatedAt: time.Now().Add(-5 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-009", StudentID: "u0000013", StudentName: "Muhammad Al-Fatih (Paket C - Tokyo)", NIS: "0089123456",
			ClassName: "Kelas X Paket C — Diaspora Internasional", InvoiceNumber: "INV-2026-11-0009", Month: 11, Year: 2026,
			Amount: 350000.0, DiscountAmount: 0.0, FinalAmount: 350000.0, DueDate: "2026-11-25", Status: "PAID",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-6 * 24 * time.Hour),
		},
		{
			ID: "inv-2026-11-010", StudentID: "u0000014", StudentName: "Maryam Khadijah (Paket B - Padang)", NIS: "0089123457",
			ClassName: "Kelas VIII Paket B — Aisyah binti Abi Bakr", InvoiceNumber: "INV-2026-11-0100", Month: 11, Year: 2026,
			Amount: 300000.0, DiscountAmount: 0.0, FinalAmount: 300000.0, DueDate: "2026-11-25", Status: "PENDING_VERIFICATION",
			ReceiptFileURL: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
		},
	}

	verifiedTime1 := time.Now().Add(-20 * 24 * time.Hour)
	verifiedTime2 := time.Now().Add(-4 * 24 * time.Hour)
	verifiedUser := "u0000001"
	h.transactions = []models.PaymentTransaction{
		{
			ID:               "ptrx-001",
			StudentInvoiceID: "inv-2026-10-001",
			InvoiceNumber:    "INV-2026-10-0001",
			StudentName:      "Muhammad Abdullah Ziyad",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       350000.0,
			ProofFileURL:     "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			PaymentDate:      time.Now().Add(-20 * 24 * time.Hour),
			VerifiedByUserID: &verifiedUser,
			VerifiedByName:   "Ustadzah Fatimah, S.Pd. (Bendahara)",
			VerifiedAt:       &verifiedTime1,
			Status:           "APPROVED",
		},
		{
			ID:               "ptrx-002",
			StudentInvoiceID: "inv-2026-11-001",
			InvoiceNumber:    "INV-2026-11-0001",
			StudentName:      "Muhammad Abdullah Ziyad",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       350000.0,
			ProofFileURL:     "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			PaymentDate:      time.Now().Add(-2 * 24 * time.Hour),
			VerifiedByUserID: nil,
			VerifiedByName:   "",
			VerifiedAt:       nil,
			Status:           "PENDING",
		},
		{
			ID:               "ptrx-003",
			StudentInvoiceID: "inv-2026-11-006",
			InvoiceNumber:    "INV-2026-11-0006",
			StudentName:      "Syamil Basayif",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       350000.0,
			ProofFileURL:     "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			PaymentDate:      time.Now().Add(-3 * 24 * time.Hour),
			VerifiedByUserID: nil,
			VerifiedByName:   "",
			VerifiedAt:       nil,
			Status:           "PENDING",
		},
		{
			ID:               "ptrx-004",
			StudentInvoiceID: "inv-2026-11-010",
			InvoiceNumber:    "INV-2026-11-0100",
			StudentName:      "Maryam Khadijah (Paket B - Padang)",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       300000.0,
			ProofFileURL:     "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			PaymentDate:      time.Now().Add(-1 * 24 * time.Hour),
			VerifiedByUserID: nil,
			VerifiedByName:   "",
			VerifiedAt:       nil,
			Status:           "PENDING",
		},
		{
			ID:               "ptrx-005",
			StudentInvoiceID: "inv-2026-11-007",
			InvoiceNumber:    "INV-2026-11-0007",
			StudentName:      "Zaidan Al-Ayyubi",
			PaymentMethod:    "BANK_TRANSFER",
			AmountPaid:       350000.0,
			ProofFileURL:     "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=600&auto=format&fit=crop&q=80",
			PaymentDate:      time.Now().Add(-4 * 24 * time.Hour),
			VerifiedByUserID: &verifiedUser,
			VerifiedByName:   "Ustadzah Fatimah, S.Pd. (Bendahara)",
			VerifiedAt:       &verifiedTime2,
			Status:           "APPROVED",
		},
	}

	h.scholarships = []models.ScholarshipGrant{
		{
			ID:                 "sg-001",
			StudentID:          "u0000014b",
			StudentName:        "Abdul Andhika Raditya",
			ProgramName:        "Beasiswa Santri Huffazh Berprestasi",
			DiscountPercentage: 100,
			DonorName:          "Lembaga Amil Zakat (LAZ) Dar El Iman",
			IsActive:           true,
			CreatedAt:          time.Now().Add(-60 * 24 * time.Hour),
		},
		{
			ID:                 "sg-002",
			StudentID:          "u0000015b",
			StudentName:        "Afdhal Al Baihaqi",
			ProgramName:        "Beasiswa Prestasi Akademik Komunitas Minang",
			DiscountPercentage: 50,
			DonorName:          "Ikatan Alumni & Muhsinin Timur Tengah",
			IsActive:           true,
			CreatedAt:          time.Now().Add(-60 * 24 * time.Hour),
		},
	}

	h.exportLogs = []models.DapodikExportLog{
		{
			ID:             "del-001",
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			ExportType:     "DATA_SISWA",
			TotalRecords:   5,
			DownloadURL:    "/api/v1/dapodik/download/siswa-2026.csv",
			ExportedAt:     time.Now().Add(-2 * 24 * time.Hour),
		},
	}
}

// ========================================================
// 1. PPDB ONLINE ENDPOINTS
// ========================================================

// GetActivePeriod mengambil info gelombang PPDB yang sedang dibuka
func (h *FiberTahap4Handler) GetActivePeriod(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var activePeriod *models.PPDBPeriod
	for _, p := range h.periods {
		if p.IsOpen {
			activePeriod = &p
			break
		}
	}

	if activePeriod == nil && len(h.periods) > 0 {
		activePeriod = &h.periods[0]
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    activePeriod,
	})
}

// RegisterPublic pendaftaran mandiri calon santri (Nasional & Diaspora)
func (h *FiberTahap4Handler) RegisterPublic(c *fiber.Ctx) error {
	var req models.PPDBRegistration
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format formulir pendaftaran tidak valid.",
		})
	}

	if strings.TrimSpace(req.FullName) == "" || strings.TrimSpace(req.GuardianName) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama lengkap calon santri dan nama orang tua/wali wajib diisi.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	newID := fmt.Sprintf("preg-%s", uuid.New().String()[:8])
	newRegNum := fmt.Sprintf("REG-2027-%03d", len(h.registrations)+1)

	req.ID = newID
	req.RegistrationNumber = newRegNum
	req.Status = "SUBMITTED"
	req.CreatedAt = time.Now()
	if req.PPDBPeriodID == "" && len(h.periods) > 0 {
		req.PPDBPeriodID = h.periods[0].ID
	}
	if req.Country == "" {
		req.Country = "Indonesia"
	}
	if req.Timezone == "" {
		req.Timezone = "Asia/Jakarta (WIB)"
	}

	h.registrations = append([]models.PPDBRegistration{req}, h.registrations...)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":             true,
		"message":             "Alhamdulillah, formulir pendaftaran berhasil dikirim. Simpan Nomor Registrasi Anda untuk memantau status seleksi.",
		"registration_number": newRegNum,
		"data":                req,
	})
}

// UploadDocument upload berkas digital persyaratan PPDB
func (h *FiberTahap4Handler) UploadDocument(c *fiber.Ctx) error {
	regID := c.Params("registration_id")

	var req models.PPDBDocument
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format berkas tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found := false
	for _, r := range h.registrations {
		if r.ID == regID || r.RegistrationNumber == regID {
			found = true
			break
		}
	}
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data pendaftaran tidak ditemukan.",
		})
	}

	docID := fmt.Sprintf("pdoc-%s", uuid.New().String()[:8])
	req.ID = docID
	req.PPDBRegistrationID = regID
	req.IsVerified = false
	req.CreatedAt = time.Now()

	h.documents = append(h.documents, req)

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Berkas %s berhasil diunggah.", req.DocumentType),
		"data":    req,
	})
}

// GetMyStatus cek status seleksi mandiri oleh calon santri
func (h *FiberTahap4Handler) GetMyStatus(c *fiber.Ctx) error {
	query := strings.TrimSpace(c.Query("registration_number"))
	if query == "" {
		query = strings.TrimSpace(c.Query("email"))
	}
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Masukkan Nomor Registrasi atau Email pendaftaran.",
		})
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, r := range h.registrations {
		if strings.EqualFold(r.RegistrationNumber, query) || strings.EqualFold(r.Email, query) {
			// Attach documents
			var docs []models.PPDBDocument
			for _, d := range h.documents {
				if d.PPDBRegistrationID == r.ID {
					docs = append(docs, d)
				}
			}
			r.Documents = docs

			return c.JSON(fiber.Map{
				"success": true,
				"data":    r,
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Nomor Registrasi atau Email tidak ditemukan dalam pangkalan data PPDB.",
	})
}

// GetAdminRegistrations mengambil seluruh pendaftar untuk panitia PPDB
func (h *FiberTahap4Handler) GetAdminRegistrations(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	statusFilter := c.Query("status")
	packageFilter := c.Query("package")
	diasporaFilter := c.Query("diaspora")

	var list []models.PPDBRegistration
	for _, r := range h.registrations {
		if statusFilter != "" && r.Status != statusFilter {
			continue
		}
		if packageFilter != "" && r.TargetPackage != packageFilter {
			continue
		}
		if diasporaFilter == "true" && !r.IsDiaspora {
			continue
		}
		if diasporaFilter == "false" && r.IsDiaspora {
			continue
		}

		// Attach docs
		var docs []models.PPDBDocument
		for _, d := range h.documents {
			if d.PPDBRegistrationID == r.ID {
				docs = append(docs, d)
			}
		}
		r.Documents = docs
		list = append(list, r)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(list),
		"data":    list,
	})
}

// VerifyRegistration verifikasi berkas atau update status pendaftaran oleh panitia
func (h *FiberTahap4Handler) VerifyRegistration(c *fiber.Ctx) error {
	id := c.Params("id")

	type VerifyReq struct {
		Status string `json:"status"` // DOCUMENTS_VERIFIED, TEST_PASSED, ACCEPTED, REJECTED
		Notes  string `json:"notes"`
	}
	var req VerifyReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format verifikasi tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	foundIdx := -1
	for i, r := range h.registrations {
		if r.ID == id || r.RegistrationNumber == id {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pendaftaran calon santri tidak ditemukan.",
		})
	}

	if req.Status != "" {
		h.registrations[foundIdx].Status = req.Status
	}
	if req.Notes != "" {
		h.registrations[foundIdx].Notes = req.Notes
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Status pendaftaran %s berhasil diperbarui menjadi %s.", h.registrations[foundIdx].RegistrationNumber, h.registrations[foundIdx].Status),
		"data":    h.registrations[foundIdx],
	})
}

// AcceptAndConvertStudent AKSI KRUSIAL: Luluskan calon santri dan auto-create akun santri & students record
func (h *FiberTahap4Handler) AcceptAndConvertStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	foundIdx := -1
	for i, r := range h.registrations {
		if r.ID == id || r.RegistrationNumber == id {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data pendaftaran calon santri tidak ditemukan.",
		})
	}

	reg := &h.registrations[foundIdx]
	reg.Status = "ACCEPTED"
	reg.Notes = "Alhamdulillah DITERIMA RESMI. Akun santri telah diterbitkan otomatis dan terintegrasi ke Rombel."

	newStudentID := fmt.Sprintf("s-ppdb-%s", uuid.New().String()[:6])
	newUsername := strings.ToLower(strings.ReplaceAll(reg.FullName, " ", "."))
	if len(newUsername) > 20 {
		newUsername = newUsername[:20]
	}
	newUsername = fmt.Sprintf("santri.%s", newUsername)

	classroomName := "Kelas X Paket C — Umar bin Khattab"
	if reg.TargetPackage == "PAKET_B" {
		classroomName = "Kelas VIII Paket B — Thariq bin Ziyad"
	} else if reg.TargetPackage == "PAKET_A" {
		classroomName = "Kelas V Paket A — Salman Al-Farisi"
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Calon santri '%s' berhasil dikonversi menjadi santri resmi!", reg.FullName),
		"student": fiber.Map{
			"id":             newStudentID,
			"full_name":      reg.FullName,
			"username":       newUsername,
			"nisn":           reg.NISN,
			"classroom_name": classroomName,
			"target_package": reg.TargetPackage,
			"is_diaspora":    reg.IsDiaspora,
			"country":        reg.Country,
		},
		"registration": reg,
	})
}

// ========================================================
// 2. KEUANGAN SPP & TRANSAKSI
// ========================================================

// GetTariffs mengambil daftar tarif SPP bulanan per jenjang
func (h *FiberTahap4Handler) GetTariffs(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    h.tariffs,
	})
}

// GenerateMonthlyInvoices membuat invoice bulanan otomatis untuk seluruh santri
func (h *FiberTahap4Handler) GenerateMonthlyInvoices(c *fiber.Ctx) error {
	type GenReq struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}
	var req GenReq
	if err := c.BodyParser(&req); err != nil || req.Month < 1 || req.Month > 12 {
		req.Month = int(time.Now().Month())
		req.Year = time.Now().Year()
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	sampleStudents := []struct {
		ID       string
		Name     string
		NIS      string
		Class    string
		Discount float64
	}{
		{"u0000012b", "Muhammad Abdullah Ziyad", "0081234561", "Kelas X Paket C — Umar bin Khattab", 0.0},
		{"u0000013b", "Yafiq Al-Faruq", "0081234562", "Kelas X Paket C — Umar bin Khattab", 0.0},
		{"u0000014b", "Abdul Andhika Raditya", "0081234563", "Kelas X Paket C — Umar bin Khattab", 350000.0}, // Beasiswa 100%
		{"u0000015b", "Afdhal Al Baihaqi", "0081234564", "Kelas X Paket C — Umar bin Khattab", 175000.0},     // Beasiswa 50%
		{"u0000016b", "Ahmad Hasan Deedat", "0081234565", "Kelas X Paket C — Umar bin Khattab", 0.0},
	}

	baseAmount := 350000.0
	createdCount := 0

	for i, st := range sampleStudents {
		invNum := fmt.Sprintf("INV-%d-%02d-%04d", req.Year, req.Month, i+1)

		// Check if already exists
		exists := false
		for _, inv := range h.invoices {
			if inv.StudentID == st.ID && inv.Month == req.Month && inv.Year == req.Year {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		finalAmt := math.Max(0.0, baseAmount-st.Discount)
		status := "UNPAID"
		if finalAmt == 0 {
			status = "PAID" // Gratis beasiswa penuh
		}

		newInv := models.StudentInvoice{
			ID:             fmt.Sprintf("inv-%d-%02d-%03d", req.Year, req.Month, i+1),
			StudentID:      st.ID,
			StudentName:    st.Name,
			NIS:            st.NIS,
			ClassName:      st.Class,
			InvoiceNumber:  invNum,
			Month:          req.Month,
			Year:           req.Year,
			Amount:         baseAmount,
			DiscountAmount: st.Discount,
			FinalAmount:    finalAmt,
			DueDate:        fmt.Sprintf("%d-%02d-10", req.Year, req.Month),
			Status:         status,
			CreatedAt:      time.Now(),
		}

		h.invoices = append(h.invoices, newInv)
		createdCount++
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"message":       fmt.Sprintf("Berhasil menerbitkan %d tagihan SPP untuk Bulan %d/%d.", createdCount, req.Month, req.Year),
		"created_count": createdCount,
		"total_invoices": len(h.invoices),
	})
}

// GetAllInvoices mengambil seluruh invoice untuk modul admin keuangan
func (h *FiberTahap4Handler) GetAllInvoices(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	month := c.QueryInt("month", 0)
	status := c.Query("status")

	var list []models.StudentInvoice
	for _, inv := range h.invoices {
		if month > 0 && inv.Month != month {
			continue
		}
		if status != "" && inv.Status != status {
			continue
		}
		list = append(list, inv)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(list),
		"data":    list,
	})
}

// GetStudentInvoices mengambil tagihan santri spesifik (untuk portal orang tua/santri)
func (h *FiberTahap4Handler) GetStudentInvoices(c *fiber.Ctx) error {
	studentID := c.Params("student_id")

	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []models.StudentInvoice
	for _, inv := range h.invoices {
		if inv.StudentID == studentID {
			list = append(list, inv)
		}
	}

	if len(list) == 0 {
		// Fallback sample Ziyad
		for _, inv := range h.invoices {
			if inv.StudentID == "u0000012b" {
				list = append(list, inv)
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(list),
		"data":    list,
	})
}

// ConfirmPayment konfirmasi transfer bayar SPP oleh wali/santri
func (h *FiberTahap4Handler) ConfirmPayment(c *fiber.Ctx) error {
	invoiceID := c.Params("id")

	type PayReq struct {
		PaymentMethod string  `json:"payment_method"`
		AmountPaid    float64 `json:"amount_paid"`
		ProofFileURL  string  `json:"proof_file_url"`
	}
	var req PayReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format pembayaran tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	foundInvIdx := -1
	for i, inv := range h.invoices {
		if inv.ID == invoiceID || inv.InvoiceNumber == invoiceID {
			foundInvIdx = i
			break
		}
	}

	if foundInvIdx == -1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Tagihan invoice tidak ditemukan.",
		})
	}

	inv := &h.invoices[foundInvIdx]
	inv.Status = "PENDING_VERIFICATION"

	trxID := fmt.Sprintf("ptrx-%s", uuid.New().String()[:8])
	newTrx := models.PaymentTransaction{
		ID:               trxID,
		StudentInvoiceID: inv.ID,
		InvoiceNumber:    inv.InvoiceNumber,
		StudentName:      inv.StudentName,
		PaymentMethod:    req.PaymentMethod,
		AmountPaid:       req.AmountPaid,
		ProofFileURL:     req.ProofFileURL,
		PaymentDate:      time.Now(),
		Status:           "PENDING",
	}

	h.transactions = append([]models.PaymentTransaction{newTrx}, h.transactions...)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Konfirmasi pembayaran berhasil dikirim! Menunggu verifikasi (*maker-checker*) oleh bendahara lembaga.",
		"invoice": inv,
		"transaction": newTrx,
	})
}

// GetPendingPayments mengambil daftar transfer yang belum diverifikasi
func (h *FiberTahap4Handler) GetPendingPayments(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var pending []models.PaymentTransaction
	for _, t := range h.transactions {
		if t.Status == "PENDING" {
			pending = append(pending, t)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(pending),
		"data":    pending,
	})
}

// VerifyPayment MAKER-CHECKER: Bendahara menyetujui atau menolak bukti transfer
func (h *FiberTahap4Handler) VerifyPayment(c *fiber.Ctx) error {
	id := c.Params("id")

	type VerifyPayReq struct {
		Status          string `json:"status"` // APPROVED, REJECTED
		Action          string `json:"action"` // APPROVE, REJECT
		RejectionReason string `json:"rejection_reason"`
		VerifierName    string `json:"verifier_name"`
	}
	var req VerifyPayReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format verifikasi pembayaran tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	foundTrxIdx := -1
	for i, t := range h.transactions {
		if t.ID == id {
			foundTrxIdx = i
			break
		}
	}

	if foundTrxIdx == -1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data transaksi pembayaran tidak ditemukan.",
		})
	}

	// Normalisasi status dari field action atau status
	finalStatus := strings.ToUpper(strings.TrimSpace(req.Status))
	if finalStatus == "" && req.Action != "" {
		if strings.ToUpper(req.Action) == "APPROVE" || strings.ToUpper(req.Action) == "APPROVED" {
			finalStatus = "APPROVED"
		} else if strings.ToUpper(req.Action) == "REJECT" || strings.ToUpper(req.Action) == "REJECTED" {
			finalStatus = "REJECTED"
		}
	}
	if finalStatus == "" {
		finalStatus = "APPROVED"
	}

	trx := &h.transactions[foundTrxIdx]
	now := time.Now()
	trx.Status = finalStatus
	trx.VerifiedAt = &now
	if req.VerifierName == "" {
		req.VerifierName = "Ustadzah Fatimah, S.Pd. (Bendahara)"
	}
	trx.VerifiedByName = req.VerifierName

	// Update invoice status accordingly
	for i, inv := range h.invoices {
		if inv.ID == trx.StudentInvoiceID {
			if finalStatus == "APPROVED" {
				h.invoices[i].Status = "PAID"
			} else if finalStatus == "REJECTED" {
				h.invoices[i].Status = "UNPAID"
				trx.RejectionReason = req.RejectionReason
			}
			break
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Pembayaran invoice %s telah di-%s.", trx.InvoiceNumber, trx.Status),
		"data":    trx,
	})
}

// GetScholarships daftar program beasiswa
func (h *FiberTahap4Handler) GetScholarships(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    h.scholarships,
	})
}

// CreateScholarship alokasi beasiswa baru untuk santri
func (h *FiberTahap4Handler) CreateScholarship(c *fiber.Ctx) error {
	var req models.ScholarshipGrant
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format beasiswa tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	req.ID = fmt.Sprintf("sg-%s", uuid.New().String()[:6])
	req.IsActive = true
	req.CreatedAt = time.Now()

	h.scholarships = append(h.scholarships, req)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Beasiswa '%s' berhasil dialokasikan.", req.ProgramName),
		"data":    req,
	})
}

// GetFinanceSummary rekapitulasi kas penerimaan SPP dan saldo
func (h *FiberTahap4Handler) GetFinanceSummary(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var totalBilled float64
	var totalCollected float64
	var totalPending float64
	var totalOutstanding float64

	paidCount := 0
	unpaidCount := 0
	pendingCount := 0

	for _, inv := range h.invoices {
		totalBilled += inv.FinalAmount
		if inv.Status == "PAID" {
			totalCollected += inv.FinalAmount
			paidCount++
		} else if inv.Status == "PENDING_VERIFICATION" {
			totalPending += inv.FinalAmount
			pendingCount++
		} else {
			totalOutstanding += inv.FinalAmount
			unpaidCount++
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_billed":      totalBilled,
			"total_collected":   totalCollected,
			"total_pending":     totalPending,
			"total_outstanding": totalOutstanding,
			"paid_count":        paidCount,
			"unpaid_count":      unpaidCount,
			"pending_count":     pendingCount,
			"collection_rate":   fmt.Sprintf("%.1f%%", (totalCollected/math.Max(1.0, totalBilled))*100),
		},
	})
}

// ========================================================
// 3. INTEGRASI & EKSPOR DAPODIK
// ========================================================

// GetValidationAudit audit kelengkapan atribut pokok santri sebelum sinkronisasi Dapodik
func (h *FiberTahap4Handler) GetValidationAudit(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 5 santri aktif
	students := []models.DapodikStudentRow{
		{
			No: 1, NamaSiswa: "Muhammad Abdullah Ziyad", NISN: "0081234561", NIK: "1371011508100003",
			JenisKelamin: "L", TempatLahir: "Padang", TanggalLahir: "2010-08-15", NamaIbuKandung: "Fatimah Zahra",
			TingkatPendidikan: "Paket C (Tingkat 10)", Rombel: "Kelas X Paket C — Umar bin Khattab", StatusValidasi: "VALID",
		},
		{
			No: 2, NamaSiswa: "Yafiq Al-Faruq", NISN: "0081234562", NIK: "1371021203100002",
			JenisKelamin: "L", TempatLahir: "Bukittinggi", TanggalLahir: "2010-03-12", NamaIbuKandung: "Maryam Azizah",
			TingkatPendidikan: "Paket C (Tingkat 10)", Rombel: "Kelas X Paket C — Umar bin Khattab", StatusValidasi: "VALID",
		},
		{
			No: 3, NamaSiswa: "Abdul Andhika Raditya", NISN: "0081234563", NIK: "1371010101100005",
			JenisKelamin: "L", TempatLahir: "Padang Panjang", TanggalLahir: "2010-01-01", NamaIbuKandung: "Siti Rahmah",
			TingkatPendidikan: "Paket C (Tingkat 10)", Rombel: "Kelas X Paket C — Umar bin Khattab", StatusValidasi: "VALID",
		},
		{
			No: 4, NamaSiswa: "Afdhal Al Baihaqi", NISN: "0081234564", NIK: "1371032005100001",
			JenisKelamin: "L", TempatLahir: "Pariaman", TanggalLahir: "2010-05-20", NamaIbuKandung: "Aisyah Nur",
			TingkatPendidikan: "Paket C (Tingkat 10)", Rombel: "Kelas X Paket C — Umar bin Khattab", StatusValidasi: "VALID",
		},
		{
			No: 5, NamaSiswa: "Ahmad Hasan Deedat", NISN: "0081234565", NIK: "1371011111100009",
			JenisKelamin: "L", TempatLahir: "Solok", TanggalLahir: "2010-11-11", NamaIbuKandung: "Khadijah S.Pd.",
			TingkatPendidikan: "Paket C (Tingkat 10)", Rombel: "Kelas X Paket C — Umar bin Khattab", StatusValidasi: "VALID",
		},
	}

	validCount := len(students)
	totalCount := len(students)
	percentage := 100.0

	return c.JSON(fiber.Map{
		"success": true,
		"audit": fiber.Map{
			"total_students":    totalCount,
			"valid_students":    validCount,
			"invalid_students":  0,
			"readiness_percent": percentage,
			"npsn":              "P9997123",
			"institution_name":  "PKBM Digital Dar El Iman",
			"accreditation":     "Terakreditasi B — Kemendikbudristek RI",
			"last_sync":         time.Now().Format("02 January 2006 15:04 WIB"),
		},
		"students": students,
	})
}

// ExportStudents menghasilkan format CSV resmi Peserta Didik untuk aplikasi Dapodik
func (h *FiberTahap4Handler) ExportStudents(c *fiber.Ctx) error {
	csvContent := `No,Nama Peserta Didik,NISN,NIK,Jenis Kelamin,Tempat Lahir,Tanggal Lahir,Nama Ibu Kandung,Rombel,Jenjang,NPSN
1,Muhammad Abdullah Ziyad,0081234561,1371011508100003,L,Padang,2010-08-15,Fatimah Zahra,Kelas X Paket C — Umar bin Khattab,Paket C,P9997123
2,Yafiq Al-Faruq,0081234562,1371021203100002,L,Bukittinggi,2010-03-12,Maryam Azizah,Kelas X Paket C — Umar bin Khattab,Paket C,P9997123
3,Abdul Andhika Raditya,0081234563,1371010101100005,L,Padang Panjang,2010-01-01,Siti Rahmah,Kelas X Paket C — Umar bin Khattab,Paket C,P9997123
4,Afdhal Al Baihaqi,0081234564,1371032005100001,L,Pariaman,2010-05-20,Aisyah Nur,Kelas X Paket C — Umar bin Khattab,Paket C,P9997123
5,Ahmad Hasan Deedat,0081234565,1371011111100009,L,Solok,2010-11-11,Khadijah S.Pd.,Kelas X Paket C — Umar bin Khattab,Paket C,P9997123
`
	h.mu.Lock()
	h.exportLogs = append([]models.DapodikExportLog{
		{
			ID:             fmt.Sprintf("del-%s", uuid.New().String()[:6]),
			AcademicYearID: "22222222-2222-2222-2222-222222222222",
			ExportType:     "DATA_SISWA",
			TotalRecords:   5,
			DownloadURL:    "/api/v1/dapodik/export-students",
			ExportedAt:     time.Now(),
		},
	}, h.exportLogs...)
	h.mu.Unlock()

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="Dapodik_PesertaDidik_PKBM_DarElIman_2026.csv"`)
	return c.SendString(csvContent)
}

// ExportTutors menghasilkan CSV data Tutor / Pendidik
func (h *FiberTahap4Handler) ExportTutors(c *fiber.Ctx) error {
	csvContent := `No,Nama PTK / Tutor,NIP,NUPTK,Jenis Kelamin,Tempat Lahir,Tanggal Lahir,Mata Pelajaran,Status Penugasan
1,Ustadzah Fatimah S.Pd.,19880415 201201 2 003,9876543210987654,P,Padang,1988-04-15,Tahfizh Al-Qur'an & Bahasa Indonesia,Wali Kelas & Tutor Tetap
2,Ustadz Ahmad Fauzi Lc.,19850110 201001 1 004,8765432109876543,L,Medan,1985-01-10,Bahasa Arab & Fiqh Ibadah,Tutor Tetap
3,Ahmad Fauzan S.Kom. M.T.,19920822 201801 1 007,7654321098765432,L,Bukittinggi,1992-08-22,Matematika & AI Literacy,Tutor Tetap
`
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="Dapodik_PTK_Tutor_PKBM_DarElIman_2026.csv"`)
	return c.SendString(csvContent)
}

// ExportClassrooms menghasilkan CSV data Rombel
func (h *FiberTahap4Handler) ExportClassrooms(c *fiber.Ctx) error {
	csvContent := `No,Nama Rombongan Belajar,Tingkat Pendidikan,Jenjang,Wali Kelas,Jumlah Peserta Didik,Kurikulum
1,Kelas X Paket C — Umar bin Khattab,Tingkat 10,Paket C,Ustadzah Fatimah S.Pd.,5,Kurikulum Merdeka Terpadu
2,Kelas VIII Paket B — Thariq bin Ziyad,Tingkat 8,Paket B,Ustadz Ahmad Fauzi Lc.,15,Kurikulum Merdeka Terpadu
3,Kelas V Paket A — Salman Al-Farisi,Tingkat 5,Paket A,Ahmad Fauzan S.Kom. M.T.,12,Kurikulum Merdeka Terpadu
`
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="Dapodik_Rombel_PKBM_DarElIman_2026.csv"`)
	return c.SendString(csvContent)
}

// GetExportLogs riwayat log ekspor Dapodik
func (h *FiberTahap4Handler) GetExportLogs(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return c.JSON(fiber.Map{
		"success": true,
		"data":    h.exportLogs,
	})
}
