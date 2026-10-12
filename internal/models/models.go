package models

import "time"

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Phone        *string    `json:"phone,omitempty"`
	AvatarURL    *string    `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	RoleName     string     `json:"role_name,omitempty"`
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type Institution struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	NPSN      *string   `json:"npsn,omitempty"`
	Address   *string   `json:"address,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Website   *string   `json:"website,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type AcademicYear struct {
	ID            string    `json:"id"`
	InstitutionID string    `json:"institution_id"`
	Name          string    `json:"name"`
	Semester      string    `json:"semester"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type GradeLevel struct {
	ID          string `json:"id"`
	PackageType string `json:"package_type"`
	LevelNumber int    `json:"level_number"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
}

type Tutor struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	NIP            *string   `json:"nip,omitempty"`
	FullName       string    `json:"full_name"`
	Gender         *string   `json:"gender,omitempty"`
	Specialization *string   `json:"specialization,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type Student struct {
	ID                   string     `json:"id"`
	UserID               string     `json:"user_id"`
	NIS                  string     `json:"nis"`
	NISN                 *string    `json:"nisn,omitempty"`
	NIK                  *string    `json:"nik,omitempty"`
	FullName             string     `json:"full_name"`
	Gender               *string    `json:"gender,omitempty"`
	BirthPlace           *string    `json:"birth_place,omitempty"`
	BirthDate            *string    `json:"birth_date,omitempty"`
	Address              *string    `json:"address,omitempty"`
	IsDiaspora           bool       `json:"is_diaspora"`
	Country              string     `json:"country"`
	Timezone             string     `json:"timezone"`
	GuardianName         *string    `json:"guardian_name,omitempty"`
	GuardianPhone        *string    `json:"guardian_phone,omitempty"`
	GuardianRelationship *string    `json:"guardian_relationship,omitempty"`
	Status               string     `json:"status"`
	Username             string     `json:"username"`
	Email                string     `json:"email"`
	ClassroomName        *string    `json:"classroom_name,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
}

type Classroom struct {
	ID               string `json:"id"`
	AcademicYearID   string `json:"academic_year_id"`
	AcademicYearName string `json:"academic_year_name,omitempty"`
	GradeLevelID     string `json:"grade_level_id"`
	GradeLevelName   string `json:"grade_level_name,omitempty"`
	Name             string `json:"name"`
	HomeroomTutorID  *string `json:"homeroom_tutor_id,omitempty"`
	HomeroomTutorName *string `json:"homeroom_tutor_name,omitempty"`
	Capacity         int    `json:"capacity"`
	StudentCount     int    `json:"student_count"`
}
