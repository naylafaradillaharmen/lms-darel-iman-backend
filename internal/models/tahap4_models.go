package models

import "time"

// PPDBPeriod mewakili gelombang pendaftaran PPDB
type PPDBPeriod struct {
	ID             string    `json:"id"`
	AcademicYearID string    `json:"academic_year_id"`
	Title          string    `json:"title"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	IsOpen         bool      `json:"is_open"`
	CreatedAt      time.Time `json:"created_at"`
}

// PPDBRegistration mewakili data pendaftar calon santri baru
type PPDBRegistration struct {
	ID                   string         `json:"id"`
	PPDBPeriodID         string         `json:"ppdb_period_id"`
	RegistrationNumber   string         `json:"registration_number"`
	TargetPackage        string         `json:"target_package"` // PAKET_A, PAKET_B, PAKET_C
	FullName             string         `json:"full_name"`
	Gender               string         `json:"gender"` // L, P
	NIK                  string         `json:"nik"`
	NISN                 string         `json:"nisn"`
	BirthPlace           string         `json:"birth_place"`
	BirthDate            string         `json:"birth_date"`
	Phone                string         `json:"phone"`
	Email                string         `json:"email"`
	IsDiaspora           bool           `json:"is_diaspora"`
	Country              string         `json:"country"`
	Timezone             string         `json:"timezone"`
	GuardianName         string         `json:"guardian_name"`
	GuardianPhone        string         `json:"guardian_phone"`
	GuardianRelationship string         `json:"guardian_relationship"`
	MotherName           string         `json:"mother_name"`
	Status               string         `json:"status"` // SUBMITTED, DOCUMENTS_VERIFIED, TEST_PASSED, ACCEPTED, REJECTED
	Notes                string         `json:"notes"`
	CreatedAt            time.Time      `json:"created_at"`
	Documents            []PPDBDocument `json:"documents,omitempty"`
}

// PPDBDocument mewakili dokumen berkas persyaratan calon santri
type PPDBDocument struct {
	ID                 string    `json:"id"`
	PPDBRegistrationID string    `json:"ppdb_registration_id"`
	DocumentType       string    `json:"document_type"` // KK, AKTA_LAHIR, IJAZAH_SEBELUMNYA, KTP_WALI, SURAT_KONSULAT
	FileURL            string    `json:"file_url"`
	IsVerified         bool      `json:"is_verified"`
	VerificationNotes  string    `json:"verification_notes"`
	CreatedAt          time.Time `json:"created_at"`
}

// TuitionTariff mewakili tarif bulanan SPP per jenjang
type TuitionTariff struct {
	ID             string    `json:"id"`
	AcademicYearID string    `json:"academic_year_id"`
	GradeLevelID   string    `json:"grade_level_id"`
	PackageName    string    `json:"package_name"`
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

// StudentInvoice mewakili tagihan SPP bulanan santri
type StudentInvoice struct {
	ID             string                `json:"id"`
	StudentID      string                `json:"student_id"`
	StudentName    string                `json:"student_name,omitempty"`
	NIS            string                `json:"nis,omitempty"`
	ClassName      string                `json:"class_name,omitempty"`
	InvoiceNumber  string                `json:"invoice_number"`
	Month          int                   `json:"month"`
	Year           int                   `json:"year"`
	Amount         float64               `json:"amount"`
	DiscountAmount float64               `json:"discount_amount"`
	FinalAmount    float64               `json:"final_amount"`
	DueDate        string                `json:"due_date"`
	Status         string                `json:"status"` // UNPAID, PENDING_VERIFICATION, PAID, CANCELLED
	CreatedAt      time.Time             `json:"created_at"`
	Transactions   []PaymentTransaction  `json:"transactions,omitempty"`
}

// PaymentTransaction mewakili transaksi setoran/pembayaran SPP
type PaymentTransaction struct {
	ID                string    `json:"id"`
	StudentInvoiceID  string    `json:"student_invoice_id"`
	InvoiceNumber     string    `json:"invoice_number,omitempty"`
	StudentName       string    `json:"student_name,omitempty"`
	PaymentMethod     string    `json:"payment_method"` // BANK_TRANSFER, PAYMENT_GATEWAY, CASH, BEASISWA
	AmountPaid        float64   `json:"amount_paid"`
	ProofFileURL      string    `json:"proof_file_url"`
	PaymentDate       time.Time `json:"payment_date"`
	VerifiedByUserID  *string   `json:"verified_by_user_id,omitempty"`
	VerifiedByName    string    `json:"verified_by_name,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	Status            string    `json:"status"` // PENDING, APPROVED, REJECTED
	RejectionReason   string    `json:"rejection_reason,omitempty"`
}

// ScholarshipGrant mewakili beasiswa potongan biaya santri
type ScholarshipGrant struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	StudentName        string    `json:"student_name,omitempty"`
	ProgramName        string    `json:"program_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	DonorName          string    `json:"donor_name"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
}

// DapodikExportLog mewakili riwayat ekspor data Dapodik
type DapodikExportLog struct {
	ID               string    `json:"id"`
	ExportedByUserID string    `json:"exported_by_user_id"`
	AcademicYearID   string    `json:"academic_year_id"`
	ExportType       string    `json:"export_type"` // DATA_SISWA, DATA_TUTOR, DATA_ROMBEL
	TotalRecords     int       `json:"total_records"`
	DownloadURL      string    `json:"download_url"`
	ExportedAt       time.Time `json:"exported_at"`
}

// DapodikValidationAudit mewakili hasil audit kelengkapan data Dapodik
type DapodikValidationAudit struct {
	TotalStudents       int                  `json:"total_students"`
	CompleteStudents    int                  `json:"complete_students"`
	IncompleteStudents  int                  `json:"incomplete_students"`
	ReadinessPercentage float64              `json:"readiness_percentage"`
	Issues              []DapodikIssueRecord `json:"issues"`
}

// DapodikIssueRecord mencatat rincian santri dengan atribut Dapodik belum lengkap
type DapodikIssueRecord struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	NIS         string `json:"nis"`
	Missing     string `json:"missing"`
}
