package models

import "time"

type QuestionBank struct {
	ID                 string     `json:"id"`
	SubjectID          string     `json:"subject_id"`
	SubjectName        string     `json:"subject_name,omitempty"`
	CreatedByTutorID   *string    `json:"created_by_tutor_id,omitempty"`
	CreatedByTutorName *string    `json:"created_by_tutor_name,omitempty"`
	Title              string     `json:"title"`
	GradeLevelID       *string    `json:"grade_level_id,omitempty"`
	GradeLevelName     *string    `json:"grade_level_name,omitempty"`
	QuestionCount      int        `json:"question_count"`
	CreatedAt          time.Time  `json:"created_at"`
	Questions          []Question `json:"questions,omitempty"`
}

type Question struct {
	ID                  string           `json:"id"`
	QuestionBankID      string           `json:"question_bank_id"`
	QuestionType        string           `json:"question_type"` // MULTIPLE_CHOICE, ESSAY
	PromptText          string           `json:"prompt_text"`
	AttachmentImageURL  *string          `json:"attachment_image_url,omitempty"`
	Points              float64          `json:"points"`
	CreatedAt           time.Time        `json:"created_at"`
	Options             []QuestionOption `json:"options,omitempty"`
}

type QuestionOption struct {
	ID         string `json:"id"`
	QuestionID string `json:"question_id"`
	OptionKey  string `json:"option_key"` // A, B, C, D, E
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct,omitempty"`
}

type ExamEvent struct {
	ID                 string        `json:"id"`
	CourseClassID      string        `json:"course_class_id"`
	CourseClassTitle   string        `json:"course_class_title,omitempty"`
	SubjectName        string        `json:"subject_name,omitempty"`
	QuestionBankID     string        `json:"question_bank_id"`
	Title              string        `json:"title"`
	ExamType           string        `json:"exam_type"` // UTS, UAS, KUIS_HARIAN, REMEDIAL
	StartTime          string        `json:"start_time"`
	EndTime            string        `json:"end_time"`
	DurationMinutes    int           `json:"duration_minutes"`
	ExamToken          string        `json:"exam_token"`
	RandomizeQuestions bool          `json:"randomize_questions"`
	RandomizeOptions   bool          `json:"randomize_options"`
	PassingScore       float64       `json:"passing_score"`
	Status             string        `json:"status"` // SCHEDULED, ACTIVE, COMPLETED
	CreatedAt          time.Time     `json:"created_at"`
	AttemptsCount      int           `json:"attempts_count,omitempty"`
	Questions          []Question    `json:"questions,omitempty"`
}

type ExamAttempt struct {
	ID                 string            `json:"id"`
	ExamEventID        string            `json:"exam_event_id"`
	ExamEventTitle     string            `json:"exam_event_title,omitempty"`
	StudentID          string            `json:"student_id"`
	StudentName        string            `json:"student_name,omitempty"`
	StartedAt          *time.Time        `json:"started_at,omitempty"`
	FinishedAt         *time.Time        `json:"finished_at,omitempty"`
	Status             string            `json:"status"` // IN_PROGRESS, SUBMITTED, DISQUALIFIED
	TotalScore         *float64          `json:"total_score,omitempty"`
	IsPassed           *bool             `json:"is_passed,omitempty"`
	CheatIncidentCount int               `json:"cheat_incident_count"`
	CreatedAt          time.Time         `json:"created_at"`
	Answers            []ExamAnswer      `json:"answers,omitempty"`
	ProctorLogs        []ExamProctorLog  `json:"proctor_logs,omitempty"`
}

type ExamAnswer struct {
	ID               string    `json:"id"`
	ExamAttemptID    string    `json:"exam_attempt_id"`
	QuestionID       string    `json:"question_id"`
	SelectedOptionID *string   `json:"selected_option_id,omitempty"`
	SelectedOptionKey *string  `json:"selected_option_key,omitempty"`
	EssayAnswer      *string   `json:"essay_answer,omitempty"`
	IsAutoGraded     bool      `json:"is_auto_graded"`
	AwardedPoints    *float64  `json:"awarded_points,omitempty"`
	SavedAt          time.Time `json:"saved_at"`
}

type ExamProctorLog struct {
	ID            string    `json:"id"`
	ExamAttemptID string    `json:"exam_attempt_id"`
	StudentName   string    `json:"student_name,omitempty"`
	IncidentType  string    `json:"incident_type"` // TAB_SWITCH, FULLSCREEN_EXIT, IDLE_TIMEOUT
	Description   string    `json:"description"`
	LoggedAt      time.Time `json:"logged_at"`
}

type FinalGradeWeight struct {
	ID               string     `json:"id"`
	CourseClassID    string     `json:"course_class_id"`
	AttendanceWeight int        `json:"attendance_weight"`
	AssignmentWeight int        `json:"assignment_weight"`
	TahfizhWeight    int        `json:"tahfizh_weight"`
	UTSWeight        int        `json:"uts_weight"`
	UASWeight        int        `json:"uas_weight"`
	IsLocked         bool       `json:"is_locked"`
	LockedByUserID   *string    `json:"locked_by_user_id,omitempty"`
	LockedAt         *time.Time `json:"locked_at,omitempty"`
}

type StudentGradebookItem struct {
	StudentID       string  `json:"student_id"`
	StudentName     string  `json:"student_name"`
	NIS             string  `json:"nis"`
	AttendanceScore float64 `json:"attendance_score"`
	AssignmentScore float64 `json:"assignment_score"`
	TahfizhScore    float64 `json:"tahfizh_score"`
	UTSScore        float64 `json:"uts_score"`
	UASScore        float64 `json:"uas_score"`
	FinalScore      float64 `json:"final_score"`
	LetterGrade     string  `json:"letter_grade"` // A, B, C, D
	IsRemedial      bool    `json:"is_remedial"`
}

type StudentReportCard struct {
	ID                     string            `json:"id"`
	StudentID              string            `json:"student_id"`
	StudentName            string            `json:"student_name"`
	NIS                    string            `json:"nis"`
	NISN                   *string           `json:"nisn,omitempty"`
	AcademicYearID         string            `json:"academic_year_id"`
	AcademicYearName       string            `json:"academic_year_name"`
	ClassroomID            string            `json:"classroom_id"`
	ClassroomName          string            `json:"classroom_name"`
	HomeroomNotes          string            `json:"homeroom_notes"`
	SpiritualBehaviorGrade string            `json:"spiritual_behavior_grade"`
	SocialBehaviorGrade    string            `json:"social_behavior_grade"`
	TahfizhSummary         string            `json:"tahfizh_summary"`
	QRVerificationCode     string            `json:"qr_verification_code"`
	PDFFileURL             *string           `json:"pdf_file_url,omitempty"`
	PublishedAt            *time.Time        `json:"published_at,omitempty"`
	Grades                 []ReportCardGrade `json:"grades,omitempty"`
}

type ReportCardGrade struct {
	ID              string  `json:"id"`
	ReportCardID    string  `json:"report_card_id"`
	SubjectID       string  `json:"subject_id"`
	SubjectCode     string  `json:"subject_code"`
	SubjectName     string  `json:"subject_name"`
	NumericalScore  float64 `json:"numerical_score"`
	LetterGrade     string  `json:"letter_grade"`
	CompetencyNotes string  `json:"competency_notes"`
}
