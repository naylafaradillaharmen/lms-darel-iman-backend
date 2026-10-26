-- ==============================================================================
-- BASIS DATA: PKBM DIGITAL DAR EL IMAN
-- TAHAP 2: Kurikulum, Pembelajaran LMS & Halaqah Tahfizh
-- Mesin: MySQL 8.0+ / InnoDB / UTF8MB4
-- ==============================================================================



SET FOREIGN_KEY_CHECKS = 0;

-- 1. TABEL KURIKULUM & MATA PELAJARAN

CREATE TABLE IF NOT EXISTS subjects (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(150) NOT NULL,
    category ENUM('UMUM_KEMENDIKBUD', 'KHAS_DAR_EL_IMAN', 'KETERAMPILAN') NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Pengikatan Mapel ke Jenjang (Paket A, B, atau C)

CREATE TABLE IF NOT EXISTS grade_subjects (
    grade_level_id CHAR(36) NOT NULL,
    subject_id CHAR(36) NOT NULL,
    credit_hours INT DEFAULT 2,
    PRIMARY KEY (grade_level_id, subject_id),
    CONSTRAINT fk_gs_gl FOREIGN KEY (grade_level_id) REFERENCES grade_levels(id) ON DELETE CASCADE,
    CONSTRAINT fk_gs_sub FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Pengikatan Guru Pengampu ke Rombel & Mapel (Course Offering)

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

-- 2. TABEL JADWAL KBM

CREATE TABLE IF NOT EXISTS class_schedules (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) NOT NULL,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 1 AND 7), -- 1=Senin, 7=Ahad
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    meeting_link TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_csch_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. TABEL TOPIK / MODUL PEMBELAJARAN LMS

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

-- 4. TABEL PENUGASAN (ASSIGNMENTS)

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

-- 5. TABEL PRESENSI / KEHADIRAN SESI KBM

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
    status ENUM('PRESENT', 'LATE', 'PERMIT', 'SICK', 'ABSENT') NOT NULL,
    notes TEXT,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_student_session UNIQUE (class_session_id, student_id),
    CONSTRAINT fk_sa_sess FOREIGN KEY (class_session_id) REFERENCES class_sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_sa_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. TABEL HALAQAH TAHFIZH AL-QUR'AN

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

-- Indeks Kinerja
CREATE INDEX idx_course_classes_tutor ON course_classes(tutor_id);
CREATE INDEX idx_course_classes_classroom ON course_classes(classroom_id);
CREATE INDEX idx_tahfizh_logs_student ON tahfizh_logs(student_id, log_date);

SET FOREIGN_KEY_CHECKS = 1;

