package models

import "time"

// Partner mewakili lembaga mitra (pesantren, majelis taklim, komunitas diaspora luar negeri)
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
	CreatedAt     time.Time `json:"created_at"`
	StudentCount  int       `json:"student_count,omitempty"`
}

// PartnerStudent relasi antara lembaga mitra dan santri binaan
type PartnerStudent struct {
	PartnerID    string    `json:"partner_id"`
	StudentID    string    `json:"student_id"`
	StudentName  string    `json:"student_name,omitempty"`
	NIS          string    `json:"nis,omitempty"`
	ClassName    string    `json:"class_name,omitempty"`
	TahfizhJuz   string    `json:"tahfizh_juz,omitempty"`
	RegisteredAt time.Time `json:"registered_at"`
}

// AccreditationStandard mewakili 8 Standar Nasional Pendidikan (BAN-PDM)
type AccreditationStandard struct {
	ID            string              `json:"id"`
	Code          string              `json:"code"` // BAN-PDM-01 s/d BAN-PDM-08
	StandardName  string              `json:"standard_name"`
	ComponentName string              `json:"component_name"`
	MaxScore      float64             `json:"max_score"`
	CreatedAt     time.Time           `json:"created_at"`
	Items         []AccreditationItem `json:"items,omitempty"`
}

// AccreditationItem mewakili butir instrumen akreditasi dan bukti fisik eviden
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

// InstitutionSOP mewakili Standar Operasional Prosedur kelembagaan
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

// RiskRegister mewakili register risiko mutu & operasional
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

// BSCIndicator mewakili indikator Balanced Scorecard 4 perspektif
type BSCIndicator struct {
	ID             string    `json:"id"`
	Perspective    string    `json:"perspective"` // FINANCIAL, CUSTOMER_STAKEHOLDER, INTERNAL_PROCESS, LEARNING_GROWTH
	IndicatorName  string    `json:"indicator_name"`
	TargetValue    float64   `json:"target_value"`
	CurrentValue   float64   `json:"current_value"`
	Unit           string    `json:"unit"` // %, Orang, Rp, Skor
	PeriodYear     int       `json:"period_year"`
	AchievementPct float64   `json:"achievement_pct,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BSCDashboardSummary ringkasan eksekutif Balanced Scorecard
type BSCDashboardSummary struct {
	FinancialScore       float64        `json:"financial_score"`
	CustomerScore        float64        `json:"customer_score"`
	InternalProcessScore float64        `json:"internal_process_score"`
	LearningGrowthScore  float64        `json:"learning_growth_score"`
	OverallIndex         float64        `json:"overall_index"`
	Indicators           []BSCIndicator `json:"indicators"`
}

// AIChatSession mewakili sesi tanya jawab belajar santri dengan Tutor AI
type AIChatSession struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	SubjectID   *string   `json:"subject_id,omitempty"`
	SubjectName string    `json:"subject_name,omitempty"`
	StartedAt   time.Time `json:"started_at"`
}

// AIChatMessage mewakili pesan tanya jawab chat AI
type AIChatMessage struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	SenderType string    `json:"sender_type"` // USER, ASSISTANT
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// AIChatAskRequest permintaan tanya jawab AI
type AIChatAskRequest struct {
	Question string `json:"question,omitempty"`
	Prompt   string `json:"prompt,omitempty"`
}

// AIChatAskResponse respon tanya jawab AI
type AIChatAskResponse struct {
	SessionID string        `json:"session_id"`
	Answer    string        `json:"answer,omitempty"`
	Source    string        `json:"source,omitempty"`
	Timestamp string        `json:"timestamp,omitempty"`
	UserMsg   AIChatMessage `json:"user_msg,omitempty"`
	ReplyMsg  AIChatMessage `json:"reply_msg,omitempty"`
}

// DiasporaMapData struktur untuk peta visual sebaran santri dunia & nasional
type DiasporaMapData struct {
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Timezone     string  `json:"timezone"`
	TotalStudent int     `json:"total_student"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	FlagEmoji    string  `json:"flag_emoji"`
}

// DiasporaCountryRecord rekor negara sebaran santri
type DiasporaCountryRecord struct {
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	StudentCount int    `json:"student_count"`
	Timezone     string `json:"timezone"`
	FlagEmoji    string `json:"flag_emoji"`
}

// DiasporaMapSummary ringkasan agregasi diaspora
type DiasporaMapSummary struct {
	TotalStudents  int                     `json:"total_students"`
	NationalCount  int                     `json:"national_count"`
	DiasporaCount  int                     `json:"diaspora_count"`
	CountriesCount int                     `json:"countries_count"`
	Countries      []DiasporaCountryRecord `json:"countries"`
}
