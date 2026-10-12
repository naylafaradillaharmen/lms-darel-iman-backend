package models

import "time"

type Subject struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Category    string    `json:"category"` // UMUM_KEMENDIKBUD, KHAS_DAR_EL_IMAN, KETERAMPILAN
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type GradeSubject struct {
	GradeLevelID string `json:"grade_level_id"`
	SubjectID    string `json:"subject_id"`
	CreditHours  int    `json:"credit_hours"`
}

type CourseClass struct {
	ID               string          `json:"id"`
	ClassroomID      string          `json:"classroom_id"`
	ClassroomName    string          `json:"classroom_name,omitempty"`
	SubjectID        string          `json:"subject_id"`
	SubjectCode      string          `json:"subject_code,omitempty"`
	SubjectName      string          `json:"subject_name,omitempty"`
	TutorID          string          `json:"tutor_id"`
	TutorName        string          `json:"tutor_name,omitempty"`
	AcademicYearID   string          `json:"academic_year_id"`
	AcademicYearName string          `json:"academic_year_name,omitempty"`
	Schedules        []ClassSchedule `json:"schedules,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

type ClassSchedule struct {
	ID            string    `json:"id"`
	CourseClassID string    `json:"course_class_id"`
	DayOfWeek     int       `json:"day_of_week"` // 1=Senin, 7=Ahad
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	MeetingLink   *string   `json:"meeting_link,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type CourseTopic struct {
	ID            string          `json:"id"`
	CourseClassID string          `json:"course_class_id"`
	Title         string          `json:"title"`
	SequenceOrder int             `json:"sequence_order"`
	Description   *string         `json:"description,omitempty"`
	IsPublished   bool            `json:"is_published"`
	Materials     []TopicMaterial `json:"materials,omitempty"`
	Assignments   []Assignment    `json:"assignments,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

type TopicMaterial struct {
	ID            string    `json:"id"`
	CourseTopicID string    `json:"course_topic_id"`
	Title         string    `json:"title"`
	MaterialType  string    `json:"material_type"` // DOCUMENT, VIDEO_LINK, TEXT_ARTICLE
	ContentBody   *string   `json:"content_body,omitempty"`
	FileURL       *string   `json:"file_url,omitempty"`
	ExternalLink  *string   `json:"external_link,omitempty"`
	SequenceOrder int       `json:"sequence_order"`
	CreatedAt     time.Time `json:"created_at"`
}

type Assignment struct {
	ID                string                 `json:"id"`
	CourseTopicID     string                 `json:"course_topic_id"`
	Title             string                 `json:"title"`
	Instructions      string                 `json:"instructions"`
	DueDate           string                 `json:"due_date"`
	MaxScore          float64                `json:"max_score"`
	FileAttachmentURL *string                `json:"file_attachment_url,omitempty"`
	SubmissionCount   int                    `json:"submission_count,omitempty"`
	Submissions       []AssignmentSubmission `json:"submissions,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
}

type AssignmentSubmission struct {
	ID                 string     `json:"id"`
	AssignmentID       string     `json:"assignment_id"`
	StudentID          string     `json:"student_id"`
	StudentName        string     `json:"student_name,omitempty"`
	SubmissionText     *string    `json:"submission_text,omitempty"`
	FileURL            *string    `json:"file_url,omitempty"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	Score              *float64   `json:"score,omitempty"`
	Feedback           *string    `json:"feedback,omitempty"`
	GradedByTutorID    *string    `json:"graded_by_tutor_id,omitempty"`
	GradedByTutorName  *string    `json:"graded_by_tutor_name,omitempty"`
	GradedAt           *time.Time `json:"graded_at,omitempty"`
}

type ClassSession struct {
	ID            string              `json:"id"`
	CourseClassID string              `json:"course_class_id"`
	SessionTitle  string              `json:"session_title"`
	SessionDate   string              `json:"session_date"`
	StartTime     string              `json:"start_time"`
	EndTime       string              `json:"end_time"`
	Attendances   []SessionAttendance `json:"attendances,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
}

type SessionAttendance struct {
	ID             string    `json:"id"`
	ClassSessionID string    `json:"class_session_id"`
	StudentID      string    `json:"student_id"`
	StudentName    string    `json:"student_name,omitempty"`
	Status         string    `json:"status"` // PRESENT, PERMIT, SICK, ABSENT
	Notes          *string   `json:"notes,omitempty"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type TahfizhGroup struct {
	ID               string    `json:"id"`
	AcademicYearID   string    `json:"academic_year_id"`
	AcademicYearName string    `json:"academic_year_name,omitempty"`
	Name             string    `json:"name"`
	TutorID          string    `json:"tutor_id"`
	TutorName        string    `json:"tutor_name,omitempty"`
	TargetJuz        *string   `json:"target_juz,omitempty"`
	StudentCount     int       `json:"student_count,omitempty"`
	Students         []Student `json:"students,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type TahfizhLog struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	StudentName    string    `json:"student_name,omitempty"`
	TutorID        string    `json:"tutor_id"`
	TutorName      string    `json:"tutor_name,omitempty"`
	LogDate        string    `json:"log_date"`
	Type           string    `json:"type"` // ZIYADAH, MURAJAAH
	SurahName      string    `json:"surah_name"`
	StartAyah      int       `json:"start_ayah"`
	EndAyah        int       `json:"end_ayah"`
	FluencyScore   *string   `json:"fluency_score,omitempty"` // MUMTAZ, JAYYID_JIDDAN, JAYYID, MAQBUL
	TajweedScore   *string   `json:"tajweed_score,omitempty"` // MUMTAZ, JAYYID_JIDDAN, JAYYID, MAQBUL
	AudioRecordURL *string   `json:"audio_record_url,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
