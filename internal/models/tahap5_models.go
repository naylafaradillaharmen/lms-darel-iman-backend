package models

import "time"

// Partner mewakili institusi mitra (pesantren, majelis taklim, komunitas diaspora)
type Partner struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	PartnerType   string    `json:"partner_type"` // PESANTREN, MAJELIS_TAKLIM, KOMUNITAS_DIASPORA, LEMBAGA_BEASISWA
	Country       string    `json:"country"`
	City          string    `json:"city"`
	ContactPerson string    `json:"contact_person"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	IsActive      bool      `json:"is_active"`
	StudentCount  int       `json:"student_count,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// PartnerStudent mewakili relasi santri di bawah naungan mitra
type PartnerStudent struct {
	PartnerID    string    `json:"partner_id"`
	StudentID    string    `json:"student_id"`
	StudentName  string    `json:"student_name,omitempty"`
	NIS          string    `json:"nis,omitempty"`
	ClassName    string    `json:"class_name,omitempty"`
	TahfizhJuz   string    `json:"tahfizh_juz,omitempty"`
	RegisteredAt time.Time `json:"registered_at"`
}

// AccreditationStandard mewakili 8 standar akreditasi BAN-PDM
type AccreditationStandard struct {
	ID            string              `json:"id"`
	Code          string              `json:"code"`
	StandardName  string              `json:"standard_name"`
	ComponentName string              `json:"component_name"`
	MaxScore      float64             `json:"max_score"`
	Items         []AccreditationItem `json:"items,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

// AccreditationItem mewakili butir instrumen eviden akreditasi
type AccreditationItem struct {
	ID                   string    `json:"id"`
	StandardID           string    `json:"standard_id"`
	StandardCode         string    `json:"standard_code,omitempty"`
	ItemNumber           string    `json:"item_number"`
	IndicatorDescription string    `json:"indicator_description"`
	FulfillmentStatus    string    `json:"fulfillment_status"` // NOT_STARTED, IN_PROGRESS, READY_AUDIT, FULFILLED
	EvidenceFileURL      string    `json:"evidence_file_url"`
	Notes                string    `json:"notes"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// InstitutionSOP mewakili katalog Standar Operasional Prosedur lembaga
type InstitutionSOP struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Title         string    `json:"title"`
	Department    string    `json:"department"`
	Version       string    `json:"version"`
	FileURL       string    `json:"file_url"`
	EffectiveDate string    `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}

// RiskRegister mewakili pemetaan risiko mutu dan mitigasi
type RiskRegister struct {
	ID             string    `json:"id"`
	RiskTitle      string    `json:"risk_title"`
	Category       string    `json:"category"` // OPERASIONAL, TEKNOLOGI, REPUTASI, KEUANGAN
	SeverityLevel  string    `json:"severity_level"` // LOW, MEDIUM, HIGH, CRITICAL
	MitigationPlan string    `json:"mitigation_plan"`
	OwnerUserID    *string   `json:"owner_user_id,omitempty"`
	OwnerName      string    `json:"owner_name,omitempty"`
	Status         string    `json:"status"` // OPEN, MITIGATED, CLOSED
	CreatedAt      time.Time `json:"created_at"`
}

// BSCIndicator mewakili indikator Balanced Scorecard
type BSCIndicator struct {
	ID            string    `json:"id"`
	Perspective   string    `json:"perspective"` // FINANCIAL, CUSTOMER_STAKEHOLDER, INTERNAL_PROCESS, LEARNING_GROWTH
	IndicatorName string    `json:"indicator_name"`
	TargetValue   float64   `json:"target_value"`
	CurrentValue  float64   `json:"current_value"`
	Unit          string    `json:"unit"`
	PeriodYear    int       `json:"period_year"`
	AchievementPct float64  `json:"achievement_pct"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BSCDashboardSummary mewakili ringkasan eksekutif 4 perspektif BSC
type BSCDashboardSummary struct {
	FinancialScore       float64        `json:"financial_score"`
	CustomerScore        float64        `json:"customer_score"`
	InternalProcessScore float64        `json:"internal_process_score"`
	LearningGrowthScore  float64        `json:"learning_growth_score"`
	OverallIndex         float64        `json:"overall_index"`
	Indicators           []BSCIndicator `json:"indicators"`
}

// DiasporaCountryRecord mewakili titik sebaran santri per negara/zona waktu
type DiasporaCountryRecord struct {
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	StudentCount int    `json:"student_count"`
	Timezone     string `json:"timezone"`
	FlagEmoji    string `json:"flag_emoji"`
}

// DiasporaMapSummary mewakili agregasi persebaran santri global
type DiasporaMapSummary struct {
	TotalStudents   int                     `json:"total_students"`
	NationalCount   int                     `json:"national_count"`
	DiasporaCount   int                     `json:"diaspora_count"`
	CountriesCount  int                     `json:"countries_count"`
	Countries       []DiasporaCountryRecord `json:"countries"`
}

// AIChatSession mewakili sesi tanya jawab santri dengan Tutor AI
type AIChatSession struct {
	ID          string          `json:"id"`
	StudentID   string          `json:"student_id"`
	SubjectID   *string         `json:"subject_id,omitempty"`
	SubjectName string          `json:"subject_name,omitempty"`
	StartedAt   time.Time       `json:"started_at"`
	Messages    []AIChatMessage `json:"messages,omitempty"`
}

// AIChatMessage mewakili pesan interaksi chat AI
type AIChatMessage struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	SenderType string    `json:"sender_type"` // USER, ASSISTANT
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// AIChatAskRequest permintaan tanya materi belajar ke Tutor AI
type AIChatAskRequest struct {
	Question  string `json:"question"`
	SubjectID string `json:"subject_id,omitempty"`
}

// AIChatAskResponse respon jawaban pintar dari Tutor AI
type AIChatAskResponse struct {
	SessionID string `json:"session_id"`
	Answer    string `json:"answer"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}
