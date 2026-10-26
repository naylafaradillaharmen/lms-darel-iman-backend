-- ========================================================
-- SEED DATA TAHAP 5: KEMITRAAN, MUTU BAN-PDM & TUTOR AI
-- PKBM DIGITAL DAR EL IMAN
-- ========================================================

-- 1. MITRA LEMBAGA & KOMUNITAS
INSERT INTO partners (id, name, partner_type, country, city, contact_person, phone, email, is_active) VALUES
('part-001', 'Pondok Pesantren Tahfizh Rabbani Padang', 'PESANTREN', 'Indonesia', 'Padang', 'Ustadz Ilham Siddiq, M.Ag.', '+6281234567891', 'kerjasama@pesantren-rabbani.sch.id', TRUE),
('part-002', 'Komunitas Diaspora Muslim Indonesia di Qatar', 'KOMUNITAS_DIASPORA', 'Qatar', 'Doha', 'Ir. Salman Al-Farisi', '+97455123456', 'diaspora.qatar@darel-iman.sch.id', TRUE),
('part-003', 'Majelis Taklim WNI Riyadh & Wilayah Tengah KSA', 'KOMUNITAS_DIASPORA', 'Arab Saudi', 'Riyadh', 'Ustadz Abu Rayyan Al-Atsary', '+966501234567', 'taklim.riyadh@gmail.com', TRUE),
('part-004', 'Yayasan Wakaf & Beasiswa Dhuafa Nusantara', 'LEMBAGA_BEASISWA', 'Indonesia', 'Jakarta Selatan', 'H. Abdullah Mas''ud', '+628119876543', 'beasiswa@dhuafanusantara.or.id', TRUE)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 2. PARTNER STUDENTS MAPPING
INSERT INTO partner_students (partner_id, student_id, registered_at) VALUES
('part-001', '11111111-1111-1111-1111-111111111111', '2026-07-15 08:00:00'),
('part-002', '22222222-2222-2222-2222-222222222222', '2026-07-16 09:30:00'),
('part-003', '33333333-3333-3333-3333-333333333333', '2026-07-17 10:15:00')
ON DUPLICATE KEY UPDATE registered_at=VALUES(registered_at);

-- 3. STANDAR AKREDITASI BAN-PDM (8 STANDAR NASIONAL PENDIDIKAN)
INSERT INTO accreditation_standards (id, code, standard_name, component_name, max_score) VALUES
('std-01', 'BAN-PDM-01', 'Standar Kompetensi Lulusan (SKL)', 'Capaian Akademik & Hafalan Al-Qur''an Terstandar', 4.00),
('std-02', 'BAN-PDM-02', 'Standar Isi Kurikulum', 'Struktur Kurikulum Merdeka Terpadu & Modul Digital', 4.00),
('std-03', 'BAN-PDM-03', 'Standar Proses Pembelajaran', 'KBM Daring Fleksibel, Presensi Digital & Halaqah Virtual', 4.00),
('std-04', 'BAN-PDM-04', 'Standar Pendidik & Tenaga Kependidikan', 'Kualifikasi S1/S2 Tutor & Sanad Tahfizh Hafsh', 4.00),
('std-05', 'BAN-PDM-05', 'Standar Sarana & Prasarana Digital', 'Server LMS, Portal CBT Anti-Cheat & Studio Rekaman', 4.00),
('std-06', 'BAN-PDM-06', 'Standar Pengelolaan Lembaga', 'Tata Kelola Berbasis SOP, Manajemen Risiko & Transparansi', 4.00),
('std-07', 'BAN-PDM-07', 'Standar Pembiayaan & Transparansi', 'Akuntabilitas SPP, Beasiswa Donatur & Laporan Kas', 4.00),
('std-08', 'BAN-PDM-08', 'Standar Penilaian Pendidikan', 'Bank Soal CBT, Validasi E-Rapor Ber-QR Code & Dapodik', 4.00)
ON DUPLICATE KEY UPDATE standard_name=VALUES(standard_name);

-- 4. BUTIR INSTRUMEN & EVIDEN AKREDITASI
INSERT INTO accreditation_items (id, standard_id, item_number, indicator_description, fulfillment_status, evidence_file_url, notes) VALUES
('item-01', 'std-01', '1.1', 'Persentase kelulusan santri kesetaraan Paket A, B, C dengan rata-rata nilai di atas KKM lembaga (75.00).', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/skl-rekapitulasi-nilai-2026.pdf', 'Audit internal menyatakan 100% peserta lulus dengan predikat Baik/Amat Baik.'),
('item-02', 'std-01', '1.2', 'Bukti sertifikasi tasmi'' dan sanad hafalan Al-Qur''an santri minimal 3 Juz mutqin.', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/syahadah-tahfizh-bersanad.pdf', 'Sertifikat tasmi'' 30 juz dan sanad hafsh ditandatangani Ustadz Ahmad Fauzi Lc.'),
('item-03', 'std-02', '2.1', 'Ketersediaan dokumen Kurikulum Operasional Satuan Pendidikan (KOSP) terpadu sunnah.', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/kosp-pkbm-darel-iman-2026.pdf', 'Disahkan oleh Dinas Pendidikan Kota Padang.'),
('item-04', 'std-03', '3.1', 'Keterlaksanaan proses pembelajaran daring sinkron (Google Meet) dan asinkron (LMS modul interaktif).', 'READY_AUDIT', 'https://cdn.darel-iman.sch.id/eviden/log-kbm-presensi-digital.pdf', 'Log presensi geo-tagging dan rekaman video KBM terdokumentasi rapi di LMS.'),
('item-05', 'std-04', '4.1', 'Rasio tutor berijazah sarjana linier dan bersertifikat pendidik profesional.', 'READY_AUDIT', 'https://cdn.darel-iman.sch.id/eviden/ptk-sk-pengangkatan-tutor.pdf', 'Semua tutor berpendidikan S1/S2 dengan NUPTK resmi.'),
('item-06', 'std-05', '5.1', 'Infrastruktur cloud server LMS, sistem ujian proctoring, dan jaringan bandwith memadai.', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/arsitektur-cloud-server-pkbm.pdf', 'Backend Go Fiber microsecond latency & sistem proctoring anti-cheat aktif.'),
('item-07', 'std-06', '6.1', 'Struktur organisasi dan manual SOP operasional yang tersosialisasikan kepada civitas lembaga.', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/katalog-sop-kelembagaan.pdf', 'Tersedia 12 SOP operasional digital.'),
('item-08', 'std-07', '7.1', 'Laporan audit keuangan berkala dan alokasi dana beasiswa santri berprestasi/dhuafa.', 'READY_AUDIT', 'https://cdn.darel-iman.sch.id/eviden/laporan-keuangan-beasiswa-2026.pdf', 'Laporan penerimaan kas SPP dan penyaluran beasiswa transparan.'),
('item-09', 'std-08', '8.1', 'Sistem penilaian terpadu e-rapor digital ber-QR code verifikasi publik dan sinkronisasi Dapodik.', 'FULFILLED', 'https://cdn.darel-iman.sch.id/eviden/sampel-e-rapor-qr-code.pdf', 'Tervalidasi 100% pada portal verifikasi publik.')
ON DUPLICATE KEY UPDATE indicator_description=VALUES(indicator_description);

-- 5. STANDAR OPERASIONAL PROSEDUR (SOP)
INSERT INTO institution_sops (id, code, title, department, version, file_url, effective_date) VALUES
('sop-001', 'SOP-AKD-001', 'Standar Prosedur Operasional KBM Daring & Tatap Maya', 'Akademik', '2.0', 'https://cdn.darel-iman.sch.id/sop/sop-kbm-daring.pdf', '2026-07-01'),
('sop-002', 'SOP-TFZ-002', 'Standar Prosedur Setoran & Ujian Tasmi'' Tahfizh Bersanad', 'Tahfizh & Al-Qur''an', '2.1', 'https://cdn.darel-iman.sch.id/sop/sop-tahfizh-tasmi.pdf', '2026-07-01'),
('sop-003', 'SOP-CBT-003', 'Standar Prosedur Integritas Proctoring & Ujian Berbasis Komputer', 'Teknologi & CBT', '2.0', 'https://cdn.darel-iman.sch.id/sop/sop-cbt-proctoring.pdf', '2026-08-15'),
('sop-004', 'SOP-KEU-004', 'Standar Prosedur Pembayaran SPP & Verifikasi Maker-Checker', 'Keuangan & Kas', '1.5', 'https://cdn.darel-iman.sch.id/sop/sop-spp-maker-checker.pdf', '2026-09-01'),
('sop-005', 'SOP-DPD-005', 'Standar Prosedur Validasi Data Pokok & Ekspor Terstandar Dapodik', 'Tata Usaha & Dapodik', '1.2', 'https://cdn.darel-iman.sch.id/sop/sop-validasi-dapodik.pdf', '2026-09-15')
ON DUPLICATE KEY UPDATE title=VALUES(title);

-- 6. REGISTER RISIKO OPERASIONAL
INSERT INTO risk_registers (id, risk_title, category, severity_level, mitigation_plan, status) VALUES
('risk-001', 'Gangguan jaringan internet santri diaspora saat ujian sumatif CBT terjadwal', 'TEKNOLOGI', 'MEDIUM', 'Pemberian token remedial otomatis dan mode cache lokal penyimpanan jawaban sementara pada peramban santri.', 'MITIGATED'),
('risk-002', 'Ketidaklengkapan dokumen kependudukan (NIK/NISN) calon santri baru untuk Dapodik', 'OPERASIONAL', 'HIGH', 'Validasi prasyarat pra-submit pada wizard pendaftaran PPDB dan koordinasi dengan KBRI/Konsulat Jenderal bagi diaspora.', 'MITIGATED'),
('risk-003', 'Keterlambatan pembayaran SPP santri mandiri', 'KEUANGAN', 'MEDIUM', 'Pengiriman notifikasi otomatis via WhatsApp H-3 jatuh tempo dan alokasi program beasiswa donatur bagi yang terkendala.', 'OPEN'),
('risk-004', 'Perbedaan zona waktu KBM interaktif santri diaspora Timur Tengah (selisih 4 jam)', 'OPERASIONAL', 'LOW', 'Penyediaan jadwal halaqah tahfizh shift sore WIB (ba''da Ashar / Maghrib) yang bertepatan dengan siang hari di Timur Tengah.', 'CLOSED')
ON DUPLICATE KEY UPDATE risk_title=VALUES(risk_title);

-- 7. INDIKATOR BALANCED SCORECARD (BSC)
INSERT INTO bsc_indicators (id, perspective, indicator_name, target_value, current_value, unit, period_year) VALUES
('bsc-01', 'FINANCIAL', 'Tingkat Kolektibilitas Kas SPP Bulanan Tepat Waktu', 90.00, 88.50, '%', 2026),
('bsc-02', 'FINANCIAL', 'Kemandirian Dana Operasional Non-SPP (Infaq & Wakaf Digital)', 25.00, 27.20, '%', 2026),
('bsc-03', 'CUSTOMER_STAKEHOLDER', 'Indeks Kepuasan Wali Santri Diaspora & Domestik (NPS)', 90.00, 94.60, '%', 2026),
('bsc-04', 'CUSTOMER_STAKEHOLDER', 'Rasio Retensi Santri Melanjutkan ke Jenjang Berikutnya', 95.00, 98.00, '%', 2026),
('bsc-05', 'INTERNAL_PROCESS', 'Kepatuhan Pemenuhan Instrumen Akreditasi BAN-PDM', 100.00, 95.00, '%', 2026),
('bsc-06', 'INTERNAL_PROCESS', 'Tingkat Kehadiran & Ketepatan KBM Guru/Tutor', 95.00, 97.40, '%', 2026),
('bsc-07', 'LEARNING_GROWTH', 'Persentase Tutor Tersertifikasi Pendidik / Hafizh Sanad', 85.00, 90.00, '%', 2026),
('bsc-08', 'LEARNING_GROWTH', 'Adopsi Teknologi AI Literacy & Proctoring LMS', 80.00, 85.00, '%', 2026)
ON DUPLICATE KEY UPDATE target_value=VALUES(target_value);

-- 8. INITIAL AI TUTOR CHAT SESSION & SEED MESSAGES
INSERT INTO ai_chat_sessions (id, student_id, started_at) VALUES
('sess-demo-001', '11111111-1111-1111-1111-111111111111', '2026-10-19 14:00:00')
ON DUPLICATE KEY UPDATE started_at=VALUES(started_at);

INSERT INTO ai_chat_messages (id, session_id, sender_type, content, created_at) VALUES
('msg-001', 'sess-demo-001', 'USER', 'Assalamu''alaikum Ustadz AI, bisa tolong jelaskan kembali rukun-rukun wudhu dan hukum membaca bismillah sebelum wudhu sesuai modul Fiqh?', '2026-10-19 14:00:05'),
('msg-002', 'sess-demo-001', 'ASSISTANT', 'Wa''alaikumussalam warahmatullah wabarakatuh ya akhi Muhammad Abdullah Ziyad. Barakallahu feek.\n\nBerdasarkan silabus Fiqh Ibadah PKBM Digital Dar El Iman (rujukan Kitab Al-Wajiz fi Fiqhis Sunnah):\n\n1. **Rukun Wudhu (Fardhu Wudhu)** sebagaimana disebutkan dalam Surah Al-Ma''idah ayat 6:\n   - Niat di dalam hati.\n   - Membasuh seluruh wajah (termasuk berkumur dan istinsyaq menurut sebagian ulama).\n   - Membasuh kedua tangan sampai ke siku.\n   - Mengusap seluruh kepala (termasuk kedua telinga).\n   - Membasuh kedua kaki sampai ke mata kaki.\n   - Tertib (berurutan) dan Muwalah (bersambung tanpa jeda lama).\n\n2. **Hukum Membaca Bismillah**:\n   Menurut jumhur ulama hukumnya adalah Sunnah Muakkadah yang sangat dianjurkan, berdasarkan hadits Rasulullah ﷺ: *"Tidak sempurna wudhu orang yang tidak menyebut nama Allah atasnya."* (HR. Abu Dawud & Ibnu Majah, dinilai shahih oleh Syaikh Al-Albani).\n\nApakah ada bagian tata cara membasuh yang ingin kamu diskusikan lebih lanjut?', '2026-10-19 14:00:08')
ON DUPLICATE KEY UPDATE content=VALUES(content);

