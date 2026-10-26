-- ==============================================================================
-- BASIS DATA: PKBM DIGITAL DAR EL IMAN
-- TAHAP 4: PPDB Online, Keuangan SPP & Ekspor Dapodik Kemendikbudristek
-- Mesin: MySQL 8.0+ / InnoDB / UTF8MB4
-- ==============================================================================



SET FOREIGN_KEY_CHECKS = 0;

-- 1. TABEL PENERIMAAN PESERTA DIDIK BARU (PPDB)




CREATE TABLE IF NOT EXISTS ppdb_periods (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    academic_year_id CHAR(36) NOT NULL,
    title VARCHAR(150) NOT NULL, -- Contoh: 'PPDB Gelombang 1 TA 2027/2028'
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_open BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ppdb_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ppdb_registrations (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    ppdb_period_id CHAR(36) NOT NULL,
    registration_number VARCHAR(50) UNIQUE NOT NULL, -- Contoh: 'REG-2027-001'
    target_package ENUM('PAKET_A', 'PAKET_B', 'PAKET_C') NOT NULL,
    
    -- Biodata Calon Santri
    full_name VARCHAR(255) NOT NULL,
    gender ENUM('L', 'P') NOT NULL,
    nik VARCHAR(20),
    nisn VARCHAR(20),
    birth_place VARCHAR(100),
    birth_date DATE,
    phone VARCHAR(50),
    email VARCHAR(150),
    
    -- Atribut Diaspora
    is_diaspora BOOLEAN DEFAULT FALSE,
    country VARCHAR(100) DEFAULT 'Indonesia',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    
    -- Biodata Orang Tua
    guardian_name VARCHAR(255) NOT NULL,
    guardian_phone VARCHAR(50) NOT NULL,
    guardian_relationship VARCHAR(50) NOT NULL,
    mother_name VARCHAR(255),
    
    status ENUM('SUBMITTED', 'DOCUMENTS_VERIFIED', 'TEST_PASSED', 'ACCEPTED', 'REJECTED') DEFAULT 'SUBMITTED',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_preg_period FOREIGN KEY (ppdb_period_id) REFERENCES ppdb_periods(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS ppdb_documents (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    ppdb_registration_id CHAR(36) NOT NULL,
    document_type ENUM('KK', 'AKTA_LAHIR', 'IJAZAH_SEBELUMNYA', 'KTP_WALI', 'SURAT_KONSULAT') NOT NULL,
    file_url TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pdoc_reg FOREIGN KEY (ppdb_registration_id) REFERENCES ppdb_registrations(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. TABEL KEUANGAN, SPP & PEMBAYARAN






CREATE TABLE IF NOT EXISTS tuition_tariffs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    academic_year_id CHAR(36) NOT NULL,
    grade_level_id CHAR(36) NOT NULL,
    amount DECIMAL(12,2) NOT NULL, -- Contoh: Rp 350.000 / bulan
    description VARCHAR(200),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tt_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT,
    CONSTRAINT fk_tt_gl FOREIGN KEY (grade_level_id) REFERENCES grade_levels(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS student_invoices (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    student_id CHAR(36) NOT NULL,
    invoice_number VARCHAR(50) UNIQUE NOT NULL, -- Contoh: 'INV-2026-10-0045'
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year INT NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    discount_amount DECIMAL(12,2) DEFAULT 0.00,
    final_amount DECIMAL(12,2) NOT NULL,
    due_date DATE NOT NULL,
    status ENUM('UNPAID', 'PENDING_VERIFICATION', 'PAID', 'CANCELLED') DEFAULT 'UNPAID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sinv_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS payment_transactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    student_invoice_id CHAR(36) NOT NULL,
    payment_method ENUM('BANK_TRANSFER', 'PAYMENT_GATEWAY', 'CASH', 'BEASISWA') NOT NULL,
    amount_paid DECIMAL(12,2) NOT NULL,
    proof_file_url TEXT, -- Bukti struk transfer
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    verified_by_user_id CHAR(36),
    verified_at DATETIME,
    status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING',
    rejection_reason TEXT,
    CONSTRAINT fk_ptrx_inv FOREIGN KEY (student_invoice_id) REFERENCES student_invoices(id) ON DELETE RESTRICT,
    CONSTRAINT fk_ptrx_user FOREIGN KEY (verified_by_user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS scholarship_grants (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    student_id CHAR(36) NOT NULL,
    program_name VARCHAR(150) NOT NULL, -- Contoh: 'Beasiswa Tahfizh Berprestasi Dar El Iman'
    discount_percentage INT DEFAULT 100, -- 100% (bebas penuh) atau 50%
    donor_name VARCHAR(150),
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sg_stu FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. TABEL LOG EKSPOR DAPODIK
CREATE TABLE IF NOT EXISTS dapodik_export_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    exported_by_user_id CHAR(36),
    academic_year_id CHAR(36) NOT NULL,
    export_type ENUM('DATA_SISWA', 'DATA_TUTOR', 'DATA_ROMBEL') NOT NULL,
    total_records INT NOT NULL,
    download_url TEXT NOT NULL,
    exported_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_del_user FOREIGN KEY (exported_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_del_ay FOREIGN KEY (academic_year_id) REFERENCES academic_years(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indeks Kinerja
CREATE INDEX idx_ppdb_reg_number ON ppdb_registrations(registration_number);
CREATE INDEX idx_student_invoices_status ON student_invoices(student_id, status);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);

SET FOREIGN_KEY_CHECKS = 1;

