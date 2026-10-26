-- ========================================================
-- TAHAP 5: KEMITRAAN, PENJAMINAN MUTU BAN-PDM & TUTOR AI
-- PKBM DIGITAL DAR EL IMAN
-- ========================================================

-- 1. TABEL KEMITRAAN & DIASPORA
CREATE TABLE IF NOT EXISTS partners (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    partner_type VARCHAR(50) NOT NULL, -- PESANTREN, MAJELIS_TAKLIM, KOMUNITAS_DIASPORA, LEMBAGA_BEASISWA
    country VARCHAR(100) DEFAULT 'Indonesia',
    city VARCHAR(100),
    contact_person VARCHAR(150),
    phone VARCHAR(50),
    email VARCHAR(150),
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS partner_students (
    partner_id VARCHAR(50) NOT NULL,
    student_id VARCHAR(50) NOT NULL,
    registered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (partner_id, student_id)
);

-- 2. TABEL PENJAMINAN MUTU, SOP & AKREDITASI BAN-PDM
CREATE TABLE IF NOT EXISTS accreditation_standards (
    id VARCHAR(50) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    standard_name VARCHAR(200) NOT NULL,
    component_name VARCHAR(200) NOT NULL,
    max_score DECIMAL(5,2) DEFAULT 4.00,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS accreditation_items (
    id VARCHAR(50) PRIMARY KEY,
    standard_id VARCHAR(50) NOT NULL,
    item_number VARCHAR(20) NOT NULL,
    indicator_description TEXT NOT NULL,
    fulfillment_status VARCHAR(30) DEFAULT 'NOT_STARTED', -- NOT_STARTED, IN_PROGRESS, READY_AUDIT, FULFILLED
    evidence_file_url TEXT,
    notes TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS institution_sops (
    id VARCHAR(50) PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    department VARCHAR(100) NOT NULL,
    version VARCHAR(20) DEFAULT '1.0',
    file_url TEXT NOT NULL,
    effective_date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS risk_registers (
    id VARCHAR(50) PRIMARY KEY,
    risk_title VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL, -- OPERASIONAL, TEKNOLOGI, REPUTASI, KEUANGAN
    severity_level VARCHAR(20) NOT NULL, -- LOW, MEDIUM, HIGH, CRITICAL
    mitigation_plan TEXT NOT NULL,
    owner_user_id VARCHAR(50),
    status VARCHAR(20) DEFAULT 'OPEN', -- OPEN, MITIGATED, CLOSED
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. TABEL BALANCED SCORECARD (BSC)
CREATE TABLE IF NOT EXISTS bsc_indicators (
    id VARCHAR(50) PRIMARY KEY,
    perspective VARCHAR(50) NOT NULL, -- FINANCIAL, CUSTOMER_STAKEHOLDER, INTERNAL_PROCESS, LEARNING_GROWTH
    indicator_name VARCHAR(200) NOT NULL,
    target_value DECIMAL(10,2) NOT NULL,
    current_value DECIMAL(10,2) DEFAULT 0.00,
    unit VARCHAR(50) NOT NULL, -- %, Orang, Rp, Skor
    period_year INT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 4. TABEL ASISTEN CERDAS TUTOR AI
CREATE TABLE IF NOT EXISTS ai_chat_sessions (
    id VARCHAR(50) PRIMARY KEY,
    student_id VARCHAR(50) NOT NULL,
    subject_id VARCHAR(50),
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ai_chat_messages (
    id VARCHAR(50) PRIMARY KEY,
    session_id VARCHAR(50) NOT NULL,
    sender_type VARCHAR(10) NOT NULL, -- USER, ASSISTANT
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indeks
CREATE INDEX IF NOT EXISTS idx_partner_students_pid ON partner_students(partner_id);
CREATE INDEX IF NOT EXISTS idx_accreditation_items_sid ON accreditation_items(standard_id);
CREATE INDEX IF NOT EXISTS idx_ai_messages_session ON ai_chat_messages(session_id);

