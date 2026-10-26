-- ==============================================================================
-- BASIS DATA: PKBM DIGITAL DAR EL IMAN
-- TAHAP 1: Fondasi Sistem, Autentikasi & Master Data
-- Mesin: MySQL 8.0+ / InnoDB / UTF8MB4
-- ==============================================================================





SET FOREIGN_KEY_CHECKS = 0;

-- 1. TABEL LEMBAGA

CREATE TABLE IF NOT EXISTS institutions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL DEFAULT 'PKBM Digital Dar El Iman',
    npsn VARCHAR(20) UNIQUE,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(100),
    website VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. TABEL TAHUN AJARAN & SEMESTER

CREATE TABLE IF NOT EXISTS academic_years (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    institution_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL, -- Contoh: '2026/2027'
    semester ENUM('GANJIL', 'GENAP') NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_academic_year_semester UNIQUE (institution_id, name, semester),
    CONSTRAINT fk_ay_institution FOREIGN KEY (institution_id) REFERENCES institutions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. TABEL JENJANG PENDIDIKAN (Paket A, B, C)

CREATE TABLE IF NOT EXISTS grade_levels (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    package_type ENUM('PAKET_A', 'PAKET_B', 'PAKET_C') NOT NULL,
    level_number INT NOT NULL, -- 1-6 (Paket A), 7-9 (Paket B), 10-12 (Paket C)
    name VARCHAR(50) NOT NULL, -- Contoh: 'Paket C - Kelas 10'
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. TABEL PENGGUNA & RBAC

CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    avatar_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS roles (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(50) UNIQUE NOT NULL, -- 'SUPER_ADMIN', 'OPERATOR', 'TUTOR', 'SISWA', 'WALI_SANTRI'
    display_name VARCHAR(100) NOT NULL,
    description TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS user_roles (
    user_id CHAR(36) NOT NULL,
    role_id CHAR(36) NOT NULL,
    institution_id CHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_ur_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ur_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_ur_inst FOREIGN KEY (institution_id) REFERENCES institutions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS refresh_tokens (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    revoked_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_rt_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. TABEL MASTER DATA TUTOR (PENDIDIK / USTADZ)

CREATE TABLE IF NOT EXISTS tutors (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) UNIQUE NOT NULL,
    nip VARCHAR(50) UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    gender ENUM('L', 'P'),
    specialization VARCHAR(100), -- Contoh: 'Tahfizh Al-Qur'an', 'Bahasa Arab', 'Matematika'
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_tutor_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. TABEL MASTER DATA SANTRI (SISWA) & WALI

CREATE TABLE IF NOT EXISTS students (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) UNIQUE NOT NULL,
    nis VARCHAR(50) UNIQUE NOT NULL,
    nisn VARCHAR(20) UNIQUE,
    nik VARCHAR(20),
    full_name VARCHAR(255) NOT NULL,
    gender ENUM('L', 'P'),
    birth_place VARCHAR(100),
    birth_date DATE,
    address TEXT,
    
    -- Atribut Khusus Diaspora Internasional
    is_diaspora BOOLEAN DEFAULT FALSE,
    country VARCHAR(100) DEFAULT 'Indonesia',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    
    -- Data Orang Tua / Wali
    guardian_name VARCHAR(255),
    guardian_phone VARCHAR(50),
    guardian_relationship VARCHAR(50),
    
    status ENUM('ACTIVE', 'ON_LEAVE', 'GRADUATED', 'DROPPED_OUT') NOT NULL DEFAULT 'ACTIVE',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_student_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. TABEL ROMBONGAN BELAJAR (ROMBEL/KELAS)

CREATE TABLE IF NOT EXISTS classrooms (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    academic_year_id CHAR(36) NOT NULL,
    grade_level_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL, -- Contoh: 'Kelas X Paket C - Al-Fatih'
    homeroom_tutor_id CHAR(36),
    capacity INT DEFAULT 30,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_classroom_name_per_year UNIQUE (academic_year_id, name),
    CONSTRAINT fk_cr_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cr_gl FOREIGN KEY (grade_level_id) REFERENCES grade_levels(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cr_tutor FOREIGN KEY (homeroom_tutor_id) REFERENCES tutors(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS classroom_students (
    classroom_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (classroom_id, student_id),
    CONSTRAINT fk_cs_cr FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_cs_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. TABEL AUDIT LOG DASAR

CREATE TABLE IF NOT EXISTS audit_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36),
    action VARCHAR(100) NOT NULL,
    target_table VARCHAR(100) NOT NULL,
    target_id CHAR(36),
    details JSON,
    ip_address VARCHAR(45),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_audit_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ==============================================================================
-- TAHAP 2: Kurikulum, Pembelajaran LMS & Halaqah Tahfizh
-- ==============================================================================

-- 9. TABEL KURIKULUM & MATA PELAJARAN

CREATE TABLE IF NOT EXISTS subjects (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    category ENUM('UMUM_KEMENDIKBUD', 'KHAS_DAR_EL_IMAN', 'KETERAMPILAN') NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS grade_subjects (
    grade_level_id CHAR(36) NOT NULL,
    subject_id CHAR(36) NOT NULL,
    credit_hours INT DEFAULT 2,
    PRIMARY KEY (grade_level_id, subject_id),
    CONSTRAINT fk_gs_gl FOREIGN KEY (grade_level_id) REFERENCES grade_levels(id) ON DELETE CASCADE,
    CONSTRAINT fk_gs_sub FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS course_classes (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    classroom_id CHAR(36) NOT NULL,
    subject_id CHAR(36) NOT NULL,
    tutor_id CHAR(36) NOT NULL,
    academic_year_id CHAR(36) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_class_subject UNIQUE (classroom_id, subject_id),
    CONSTRAINT fk_cc_cr FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_cc_sub FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cc_tutor FOREIGN KEY (tutor_id) REFERENCES tutors(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cc_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 10. TABEL JADWAL KBM

CREATE TABLE IF NOT EXISTS class_schedules (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) NOT NULL,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    meeting_link TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_csch_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 11. TABEL TOPIK / MODUL PEMBELAJARAN LMS

CREATE TABLE IF NOT EXISTS course_topics (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) NOT NULL,
    title VARCHAR(200) NOT NULL,
    sequence_order INT NOT NULL DEFAULT 1,
    description TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ct_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS topic_materials (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_topic_id CHAR(36) NOT NULL,
    title VARCHAR(200) NOT NULL,
    material_type ENUM('DOCUMENT', 'VIDEO_LINK', 'TEXT_ARTICLE') NOT NULL,
    content_body TEXT,
    file_url TEXT,
    external_link TEXT,
    sequence_order INT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tm_ct FOREIGN KEY (course_topic_id) REFERENCES course_topics(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 12. TABEL PENUGASAN (ASSIGNMENTS)

CREATE TABLE IF NOT EXISTS assignments (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_topic_id CHAR(36) NOT NULL,
    title VARCHAR(200) NOT NULL,
    instructions TEXT NOT NULL,
    due_date DATETIME NOT NULL,
    max_score DECIMAL(5,2) DEFAULT 100.00,
    file_attachment_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_asg_ct FOREIGN KEY (course_topic_id) REFERENCES course_topics(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS assignment_submissions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    assignment_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    submission_text TEXT,
    file_url TEXT,
    submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    score DECIMAL(5,2),
    feedback TEXT,
    graded_by_tutor_id CHAR(36),
    graded_at DATETIME,
    CONSTRAINT uq_student_assignment UNIQUE (assignment_id, student_id),
    CONSTRAINT fk_asub_asg FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE,
    CONSTRAINT fk_asub_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_asub_tutor FOREIGN KEY (graded_by_tutor_id) REFERENCES tutors(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 13. TABEL PRESENSI / KEHADIRAN SESI KBM

CREATE TABLE IF NOT EXISTS class_sessions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) NOT NULL,
    session_title VARCHAR(150) NOT NULL,
    session_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_csess_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS session_attendances (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    class_session_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    status ENUM('PRESENT', 'PERMIT', 'SICK', 'ABSENT') NOT NULL,
    notes TEXT,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_student_session UNIQUE (class_session_id, student_id),
    CONSTRAINT fk_sa_sess FOREIGN KEY (class_session_id) REFERENCES class_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_sa_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 14. TABEL HALAQAH TAHFIZH AL-QUR'AN

CREATE TABLE IF NOT EXISTS tahfizh_groups (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    academic_year_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    tutor_id CHAR(36) NOT NULL,
    target_juz VARCHAR(50),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tg_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_tg_tutor FOREIGN KEY (tutor_id) REFERENCES tutors(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS tahfizh_group_students (
    tahfizh_group_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    PRIMARY KEY (tahfizh_group_id, student_id),
    CONSTRAINT fk_tgs_tg FOREIGN KEY (tahfizh_group_id) REFERENCES tahfizh_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_tgs_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS tahfizh_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    student_id CHAR(36) NOT NULL,
    tutor_id CHAR(36) NOT NULL,
    log_date DATE NOT NULL,
    type ENUM('ZIYADAH', 'MURAJAAH') NOT NULL,
    surah_name VARCHAR(100) NOT NULL,
    start_ayah INT NOT NULL,
    end_ayah INT NOT NULL,
    fluency_score ENUM('MUMTAZ', 'JAYYID_JIDDAN', 'JAYYID', 'MAQBUL'),
    tajweed_score ENUM('MUMTAZ', 'JAYYID_JIDDAN', 'JAYYID', 'MAQBUL'),
    audio_record_url TEXT,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tl_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_tl_tutor FOREIGN KEY (tutor_id) REFERENCES tutors(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

-- Indeks
CREATE INDEX idx_users_username_email ON users(username, email);
CREATE INDEX idx_students_nis_nisn ON students(nis, nisn);
CREATE INDEX idx_students_status ON students(status);
CREATE INDEX idx_classroom_students_cid ON classroom_students(classroom_id);
CREATE INDEX idx_course_classes_tutor ON course_classes(tutor_id);
CREATE INDEX idx_course_classes_classroom ON course_classes(classroom_id);
CREATE INDEX idx_tahfizh_logs_student ON tahfizh_logs(student_id, log_date);

