-- ==============================================================================
-- BASIS DATA: PKBM DIGITAL DAR EL IMAN
-- TAHAP 3: Asesmen, Evaluasi CBT & E-Rapor Digital
-- Mesin: MySQL 8.0+ / InnoDB / UTF8MB4
-- ==============================================================================



SET FOREIGN_KEY_CHECKS = 0;

-- 1. TABEL BANK SOAL




CREATE TABLE IF NOT EXISTS question_banks (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    subject_id CHAR(36) NOT NULL,
    created_by_tutor_id CHAR(36),
    title VARCHAR(200) NOT NULL,
    grade_level_id CHAR(36),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_qb_sub FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_qb_tutor FOREIGN KEY (created_by_tutor_id) REFERENCES tutors(id) ON DELETE SET NULL,
    CONSTRAINT fk_qb_gl FOREIGN KEY (grade_level_id) REFERENCES grade_levels(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS questions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    question_bank_id CHAR(36) NOT NULL,
    question_type ENUM('MULTIPLE_CHOICE', 'ESSAY') NOT NULL,
    prompt_text TEXT NOT NULL,
    attachment_image_url TEXT,
    points DECIMAL(5,2) DEFAULT 10.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_q_qb FOREIGN KEY (question_bank_id) REFERENCES question_banks(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS question_options (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    question_id CHAR(36) NOT NULL,
    option_key VARCHAR(5) NOT NULL, -- 'A', 'B', 'C', 'D', 'E'
    option_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_qo_q FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. TABEL SESI UJIAN CBT & PARTICIPATION





CREATE TABLE IF NOT EXISTS exam_events (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) NOT NULL,
    question_bank_id CHAR(36) NOT NULL,
    title VARCHAR(200) NOT NULL,
    exam_type ENUM('UTS', 'UAS', 'KUIS_HARIAN', 'REMEDIAL') NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 90,
    exam_token VARCHAR(20),
    randomize_questions BOOLEAN DEFAULT TRUE,
    randomize_options BOOLEAN DEFAULT TRUE,
    passing_score DECIMAL(5,2) DEFAULT 75.00,
    status ENUM('SCHEDULED', 'ACTIVE', 'COMPLETED') DEFAULT 'SCHEDULED',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ee_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_ee_qb FOREIGN KEY (question_bank_id) REFERENCES question_banks(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS exam_attempts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    exam_event_id CHAR(36) NOT NULL,
    student_id CHAR(36) NOT NULL,
    started_at DATETIME,
    finished_at DATETIME,
    status ENUM('IN_PROGRESS', 'SUBMITTED', 'DISQUALIFIED') DEFAULT 'IN_PROGRESS',
    total_score DECIMAL(5,2),
    is_passed BOOLEAN,
    cheat_incident_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_student_exam_attempt UNIQUE (exam_event_id, student_id),
    CONSTRAINT fk_ea_ee FOREIGN KEY (exam_event_id) REFERENCES exam_events(id) ON DELETE CASCADE,
    CONSTRAINT fk_ea_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS exam_answers (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    exam_attempt_id CHAR(36) NOT NULL,
    question_id CHAR(36) NOT NULL,
    selected_option_id CHAR(36),
    essay_answer TEXT,
    is_auto_graded BOOLEAN DEFAULT FALSE,
    awarded_points DECIMAL(5,2),
    saved_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT uq_attempt_question UNIQUE (exam_attempt_id, question_id),
    CONSTRAINT fk_eans_att FOREIGN KEY (exam_attempt_id) REFERENCES exam_attempts(id) ON DELETE CASCADE,
    CONSTRAINT fk_eans_q FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    CONSTRAINT fk_eans_opt FOREIGN KEY (selected_option_id) REFERENCES question_options(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS exam_proctor_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    exam_attempt_id CHAR(36) NOT NULL,
    incident_type VARCHAR(50) NOT NULL, -- 'TAB_SWITCH', 'FULLSCREEN_EXIT', 'IDLE_TIMEOUT'
    description TEXT,
    logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_epl_att FOREIGN KEY (exam_attempt_id) REFERENCES exam_attempts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. TABEL BUKU NILAI & E-RAPOR




CREATE TABLE IF NOT EXISTS final_grade_weights (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    course_class_id CHAR(36) UNIQUE NOT NULL,
    attendance_weight INT DEFAULT 10,
    assignment_weight INT DEFAULT 20,
    tahfizh_weight INT DEFAULT 30,
    uts_weight INT DEFAULT 20,
    uas_weight INT DEFAULT 20,
    is_locked BOOLEAN DEFAULT FALSE,
    locked_by_user_id CHAR(36),
    locked_at DATETIME,
    CONSTRAINT fk_fgw_cc FOREIGN KEY (course_class_id) REFERENCES course_classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_fgw_user FOREIGN KEY (locked_by_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS student_report_cards (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    student_id CHAR(36) NOT NULL,
    academic_year_id CHAR(36) NOT NULL,
    classroom_id CHAR(36) NOT NULL,
    homeroom_notes TEXT,
    spiritual_behavior_grade VARCHAR(5) DEFAULT 'A',
    social_behavior_grade VARCHAR(5) DEFAULT 'A',
    tahfizh_summary TEXT,
    qr_verification_code VARCHAR(100) UNIQUE NOT NULL,
    pdf_file_url TEXT,
    published_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_student_report UNIQUE (student_id, academic_year_id),
    CONSTRAINT fk_src_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_src_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_src_cr FOREIGN KEY (classroom_id) REFERENCES classrooms(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS report_card_grades (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    report_card_id CHAR(36) NOT NULL,
    subject_id CHAR(36) NOT NULL,
    numerical_score DECIMAL(5,2) NOT NULL,
    letter_grade VARCHAR(5) NOT NULL,
    competency_notes TEXT,
    CONSTRAINT uq_rc_grade UNIQUE (report_card_id, subject_id),
    CONSTRAINT fk_rcg_rc FOREIGN KEY (report_card_id) REFERENCES student_report_cards(id) ON DELETE CASCADE,
    CONSTRAINT fk_rcg_sub FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indeks Kinerja
CREATE INDEX idx_exam_answers_attempt ON exam_answers(exam_attempt_id);
CREATE INDEX idx_student_report_qr ON student_report_cards(qr_verification_code);
CREATE INDEX idx_exam_events_class ON exam_events(course_class_id);
CREATE INDEX idx_questions_qb ON questions(question_bank_id);

SET FOREIGN_KEY_CHECKS = 1;

