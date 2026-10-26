-- ==============================================================================
-- DATA SEEDER AWAL (INITIAL SEED): PKBM DIGITAL DAR EL IMAN
-- Password default akun bawaan adalah: admin123 / tutor123 / santri123
-- ==============================================================================



SET FOREIGN_KEY_CHECKS = 0;

-- 1. LEMBAGA DEFAULT
INSERT INTO institutions (id, name, npsn, address, phone, email, website)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'PKBM Digital Dar El Iman',
    'P9901234',
    'Jl. Gunung Juaro No. 1, Surau Gadang, Nanggalo, Padang, Sumatera Barat',
    '+62-751-123456',
    'info@darel-iman.or.id',
    'https://pkbm.darel-iman.or.id'
) ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 2. TAHUN AJARAN AKTIF (2026/2027 GANJIL)
INSERT INTO academic_years (id, institution_id, name, semester, start_date, end_date, is_active)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    '2026/2027',
    'GANJIL',
    '2026-07-15',
    '2026-12-20',
    TRUE
) ON DUPLICATE KEY UPDATE is_active=TRUE;

-- 3. MASTER JENJANG PENDIDIKAN
INSERT INTO grade_levels (id, package_type, level_number, name) VALUES
('a0000001-0000-0000-0000-000000000001', 'PAKET_A', 1, 'Paket A - Kelas 1'),
('a0000002-0000-0000-0000-000000000002', 'PAKET_A', 2, 'Paket A - Kelas 2'),
('a0000003-0000-0000-0000-000000000003', 'PAKET_A', 3, 'Paket A - Kelas 3'),
('a0000004-0000-0000-0000-000000000004', 'PAKET_A', 4, 'Paket A - Kelas 4'),
('a0000005-0000-0000-0000-000000000005', 'PAKET_A', 5, 'Paket A - Kelas 5'),
('a0000006-0000-0000-0000-000000000006', 'PAKET_A', 6, 'Paket A - Kelas 6'),
('b0000007-0000-0000-0000-000000000007', 'PAKET_B', 7, 'Paket B - Kelas 7'),
('b0000008-0000-0000-0000-000000000008', 'PAKET_B', 8, 'Paket B - Kelas 8'),
('b0000009-0000-0000-0000-000000000009', 'PAKET_B', 9, 'Paket B - Kelas 9'),
('c0000010-0000-0000-0000-000000000010', 'PAKET_C', 10, 'Paket C - Kelas 10'),
('c0000011-0000-0000-0000-000000000011', 'PAKET_C', 11, 'Paket C - Kelas 11'),
('c0000012-0000-0000-0000-000000000012', 'PAKET_C', 12, 'Paket C - Kelas 12')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 4. MASTER 18 ROLES TERPADU PKBM DIGITAL
INSERT INTO roles (id, name, display_name, description) VALUES
('r0000001-0000-0000-0000-000000000001', 'SUPER_ADMIN', 'Super Admin Yayasan', 'Akses penuh seluruh konfigurasi sistem'),
('r0000002-0000-0000-0000-000000000002', 'DEWAN_PEMBINA', 'Dewan Pembina Yayasan', 'Pengawasan eksekutif, BSC, dan tata kelola yayasan'),
('r0000003-0000-0000-0000-000000000003', 'DIREKTUR_PKBM', 'Direktur / Kepala PKBM', 'Otoritas tertinggi akademik, legalisasi rapor & mutu'),
('r0000004-0000-0000-0000-000000000004', 'MANAJER_AKADEMIK', 'Manajer Akademik & Kurikulum', 'Pengelolaan kurikulum, jadwal KBM, dan buku nilai'),
('r0000005-0000-0000-0000-000000000005', 'MANAJER_KESISWAAN', 'Manajer Kesiswaan & Diaspora', 'PPDB nasional & diaspora, penempatan halaqah tahfizh'),
('r0000006-0000-0000-0000-000000000006', 'MANAJER_KEUANGAN', 'Manajer Keuangan & Legal', 'Manajemen kas SPP, tarif, beasiswa donatur'),
('r0000007-0000-0000-0000-000000000007', 'MANAJER_TEKNOLOGI', 'Manajer Teknologi & Sistem', 'Infrastruktur cloud, sinkronisasi Dapodik, CBT & AI'),
('r0000008-0000-0000-0000-000000000008', 'TUTOR', 'Tutor / Pendidik', 'Pengampu ruang kelas LMS, materi, tugas, dan halaqah'),
('r0000009-0000-0000-0000-000000000009', 'WALI_KELAS', 'Wali Kelas', 'Pembimbing rombel santri, evaluasi dan catatan rapor'),
('r0000010-0000-0000-0000-000000000010', 'GURU_BK', 'Guru Bimbingan Konseling', 'Konseling santri, motivasi belajar dan penanganan kendala'),
('r0000011-0000-0000-0000-000000000011', 'BENDAHARA', 'Bendahara Kas & SPP', 'Verifikasi bukti transfer dan rekonsiliasi pembayaran'),
('r0000012-0000-0000-0000-000000000012', 'OPERATOR', 'Staf Operator Dapodik', 'Pengelolaan data induk kesiswaan, NISN, Dapodik Kemendikbud'),
('r0000013-0000-0000-0000-000000000013', 'PENGUJI_CBT', 'Pengawas & Proktor Ujian CBT', 'Monitoring proctoring ujian live & penanganan insiden curang'),
('r0000014-0000-0000-0000-000000000014', 'PENJAMIN_MUTU', 'Auditor & Tim Penjaminan Mutu', 'Pemenuhan 8 Standar BAN-PDM, audit SOP dan register risiko'),
('r0000015-0000-0000-0000-000000000015', 'SISWA', 'Santri / Peserta Didik', 'Peserta KBM, asesmen CBT, hafalan dan presensi mandiri'),
('r0000016-0000-0000-0000-000000000016', 'WALI_SANTRI', 'Orang Tua / Wali Santri', 'Pemantauan rapor anak, jadwal, mutabaah, dan pembayaran'),
('r0000017-0000-0000-0000-000000000017', 'MITRA_LEMBAGA', 'Mitra Pesantren & Diaspora', 'Kolaborasi kurikulum kesetaraan santri pesantren mitra'),
('r0000018-0000-0000-0000-000000000018', 'GUEST', 'Tamu Calon Santri', 'Eksplorasi profil lembaga dan pendaftaran PPDB daring')
ON DUPLICATE KEY UPDATE display_name=VALUES(display_name);

-- 5. AKUN PENGGUNA RESMI (USERS)
-- Bcrypt Hash password default masing-masing akun (contoh: admin123, tutor123, santri123):
INSERT INTO users (id, username, email, password_hash, full_name, phone, is_active) VALUES
('u0000001-0000-0000-0000-000000000001', 'admin', 'admin@darel-iman.or.id', '$2a$10$8eyOlq04ovnJjqqLRES1oOBv6Gbi3qKvylveMJo5JewQmsCv6EbQm', 'Administrator Utama Yayasan', '+6281234567890', TRUE),
('u0000002-0000-0000-0000-000000000002', 'pembina', 'pembina@darel-iman.or.id', '$2a$10$P5su97z.5mXwZgZPX1EHxuESnXLdbg3xG7MkP0QNXYC9A9Pir1LwS', 'Prof. Dr. H. Faisal, M.A. (Dewan Pembina)', '+6281234567802', TRUE),
('u0000003-0000-0000-0000-000000000003', 'direktur', 'direktur@darel-iman.or.id', '$2a$10$FxHxbVt0xDZ.zivx0EIwZeTvvSiHyZQNrLAAnZQLIpgJgl26hVCce', 'Ustadz H. Ridwan Abu Ziyad, Lc., M.Pd.', '+6281234567803', TRUE),
('u0000004-0000-0000-0000-000000000004', 'mgr.akademik', 'akademik@darel-iman.or.id', '$2a$10$r5b2jMiTc/qa9lD9zr5wcuYXONR3D8lTwulvZaz0Belg2QbtQ8O/.', 'Ustadz Syamsul Hadi, S.Pd.I.', '+6281234567804', TRUE),
('u0000005-0000-0000-0000-000000000005', 'mgr.kesiswaan', 'kesiswaan@darel-iman.or.id', '$2a$10$TneEOVZ2/dgQSmVBOzyqGOFUfxFlBa6fXqU7Ko2Q2imtOgq/IjNUi', 'Ustadzah Nurul Hidayah, S.Sos.', '+6281234567805', TRUE),
('u0000006-0000-0000-0000-000000000006', 'mgr.keuangan', 'keuangan@darel-iman.or.id', '$2a$10$3OcOHVJv.lMcW2JXh28WEeqIzARwtoP.ZyRgQgD2K5jGVmbrq5/Bi', 'Hj. Ratna Juwita, S.E., Ak.', '+6281234567806', TRUE),
('u0000007-0000-0000-0000-000000000007', 'mgr.it', 'it@darel-iman.or.id', '$2a$10$3BjQ.g/uWmcyHXiitjk2ru9Bhu3va7bRKDaw5zuOQeustW3M1Wh4S', 'Fajar Rahmad, S.Kom., M.T.', '+6281234567807', TRUE),
('u0000008-0000-0000-0000-000000000008', 'tutor.ahmad', 'ahmad@darel-iman.or.id', '$2a$10$U/HNe0ktUbb6Yu0oEi76GeGdwYphPb0crHjSguMT7KuhXRLR2Ywqy', 'Ustadz Ahmad Al-Hafizh, Lc.', '+6281234567892', TRUE),
('u0000009-0000-0000-0000-000000000009', 'walikelas', 'walikelas@darel-iman.or.id', '$2a$10$el0823neJsE6tv8g.xKdx.GDQDvOk25wsa4E/1JT6AkvbruOu4z0e', 'Ustadz M. Ilham, S.Pd.', '+6281234567893', TRUE),
('u0000010-0000-0000-0000-000000000010', 'bk', 'bk@darel-iman.or.id', '$2a$10$Dz2oDdyCGuhrrIk/Tdct1esQIRZLEWgjGSa2/WZj3ZyHs4vaa0m7a', 'Ustadzah Fatimah, S.Psi.', '+6281234567810', TRUE),
('u0000011-0000-0000-0000-000000000011', 'bendahara', 'bendahara@darel-iman.or.id', '$2a$10$ssltiLJ4QDRW.n92fho/G.qq8bMNJrwqFhwlcW3cQuIpxP.Knz.Py', 'Siti Rahmah, S.E.', '+6281234567811', TRUE),
('u0000012-0000-0000-0000-000000000012', 'operator', 'operator@darel-iman.or.id', '$2a$10$pnwSAVivVSledUw2TG2mkub0udrLon9KjG7Ijh3.NYItZuQVoRa6i', 'Staf Operator Dapodik', '+6281234567891', TRUE),
('u0000013-0000-0000-0000-000000000013', 'proktor', 'proktor@darel-iman.or.id', '$2a$10$aBeb/E2/wh6qkFMRloOVJuXSr/Sr5DPylhXQSgyfQgq/0Vvc/lfUO', 'Ustadz Fauzi, S.Kom.', '+6281234567813', TRUE),
('u0000014-0000-0000-0000-000000000014', 'mutu', 'mutu@darel-iman.or.id', '$2a$10$9KTlvzJl8EMF05Jj0PJxFeq7njY2LilDbWlOaHJeQBswXa6wbIQly', 'Ir. M. Hidayat, M.M.', '+6281234567814', TRUE),
('u0000015-0000-0000-0000-000000000015', 'santri.fatih', 'fatih@santri.darel-iman.or.id', '$2a$10$GLQCKlc3lEaE4HXwGzx9We9nzc8E1tesigJHqFDlsGxX36mI9Iv9W', 'Muhammad Al-Fatih', '+6281299887766', TRUE),
('u0000016-0000-0000-0000-000000000016', 'santri.zaki', 'zaki@santri.darel-iman.or.id', '$2a$10$s.PKnC6SqzvxesRbTElxnOjklqJYB/E4aaB8WIsiN2UClCSrnL3D.', 'Ahmad Zaki', '+6281299887755', TRUE),
('u0000017-0000-0000-0000-000000000017', 'wali.abdullah', 'abdullah@wali.darel-iman.or.id', '$2a$10$kCzk2IIkXiNKC9NbaP8nLejlkCWKqvdP3PNpDVXD5Dd0z94KEOxA2', 'Abdullah (Wali Al-Fatih)', '+6281122334455', TRUE),
('u0000018-0000-0000-0000-000000000018', 'mitra.bukittinggi', 'pesantren.dq@gmail.com', '$2a$10$aqN4dgILo0tBHlCyCrg7TebfXIRxUoDhAb1q/vUB0Two5nXJGzW.a', 'Ust. Lukman Hakim, Lc.', '+628126789001', TRUE),
('u0000019-0000-0000-0000-000000000019', 'santri.ziyad', 'santri.ziyad@student.darel-iman.sch.id', '$2a$10$GLQCKlc3lEaE4HXwGzx9We9nzc8E1tesigJHqFDlsGxX36mI9Iv9W', 'Muhammad Abdullah Ziyad', '+6281276543210', TRUE),
('u0000020-0000-0000-0000-000000000020', 'santri.yafiq', 'yafiq@student.darel-iman.sch.id', '$2a$10$s.PKnC6SqzvxesRbTElxnOjklqJYB/E4aaB8WIsiN2UClCSrnL3D.', 'Yafiq Al-Faruq', '+6281276543211', TRUE)
ON DUPLICATE KEY UPDATE password_hash=VALUES(password_hash), full_name=VALUES(full_name);

-- 6. PEMETAAN USER ROLES
INSERT INTO user_roles (user_id, role_id, institution_id) VALUES
('u0000001-0000-0000-0000-000000000001', 'r0000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111'),
('u0000002-0000-0000-0000-000000000002', 'r0000002-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111'),
('u0000003-0000-0000-0000-000000000003', 'r0000003-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111'),
('u0000004-0000-0000-0000-000000000004', 'r0000004-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111'),
('u0000005-0000-0000-0000-000000000005', 'r0000005-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111'),
('u0000006-0000-0000-0000-000000000006', 'r0000006-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111'),
('u0000007-0000-0000-0000-000000000007', 'r0000007-0000-0000-0000-000000000007', '11111111-1111-1111-1111-111111111111'),
('u0000008-0000-0000-0000-000000000008', 'r0000008-0000-0000-0000-000000000008', '11111111-1111-1111-1111-111111111111'),
('u0000009-0000-0000-0000-000000000009', 'r0000009-0000-0000-0000-000000000009', '11111111-1111-1111-1111-111111111111'),
('u0000010-0000-0000-0000-000000000010', 'r0000010-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111111'),
('u0000011-0000-0000-0000-000000000011', 'r0000011-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111111'),
('u0000012-0000-0000-0000-000000000012', 'r0000012-0000-0000-0000-000000000012', '11111111-1111-1111-1111-111111111111'),
('u0000013-0000-0000-0000-000000000013', 'r0000013-0000-0000-0000-000000000013', '11111111-1111-1111-1111-111111111111'),
('u0000014-0000-0000-0000-000000000014', 'r0000014-0000-0000-0000-000000000014', '11111111-1111-1111-1111-111111111111'),
('u0000015-0000-0000-0000-000000000015', 'r0000015-0000-0000-0000-000000000015', '11111111-1111-1111-1111-111111111111'),
('u0000016-0000-0000-0000-000000000016', 'r0000015-0000-0000-0000-000000000015', '11111111-1111-1111-1111-111111111111'),
('u0000017-0000-0000-0000-000000000017', 'r0000016-0000-0000-0000-000000000016', '11111111-1111-1111-1111-111111111111'),
('u0000018-0000-0000-0000-000000000018', 'r0000017-0000-0000-0000-000000000017', '11111111-1111-1111-1111-111111111111'),
('u0000019-0000-0000-0000-000000000019', 'r0000015-0000-0000-0000-000000000015', '11111111-1111-1111-1111-111111111111'),
('u0000020-0000-0000-0000-000000000020', 'r0000015-0000-0000-0000-000000000015', '11111111-1111-1111-1111-111111111111')
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);

-- 7. DATA TUTOR / USTADZ
INSERT INTO tutors (id, user_id, nip, full_name, gender, specialization, phone) VALUES
('t0000001-0000-0000-0000-000000000001', 'u0000008-0000-0000-0000-000000000008', 'NIP-19850101-001', 'Ustadz Ahmad Al-Hafizh, Lc.', 'L', 'Tahfizh Al-Quran & Adab', '+6281234567892'),
('t0000002-0000-0000-0000-000000000002', 'u0000009-0000-0000-0000-000000000009', 'NIP-19880415-002', 'Ustadz M. Ilham, S.Pd.', 'L', 'Bahasa Arab & Kurikulum Kesetaraan', '+6281234567893')
ON DUPLICATE KEY UPDATE full_name=VALUES(full_name);

-- 8. DATA ROMBEL (KELAS X PAKET C)
INSERT INTO classrooms (id, academic_year_id, grade_level_id, name, homeroom_tutor_id, capacity)
VALUES (
    'c1111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    'c0000010-0000-0000-0000-000000000010',
    'Kelas X Paket C - Umar bin Khattab',
    't0000001-0000-0000-0000-000000000001',
    30
) ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 9. DATA SANTRI RESMI
INSERT INTO students (
    id, user_id, nis, nisn, nik, full_name, gender, birth_place, birth_date,
    address, is_diaspora, country, timezone, guardian_name, guardian_phone, guardian_relationship, status
) VALUES
(
    's0000001-0000-0000-0000-000000000001',
    'u0000015-0000-0000-0000-000000000015',
    'NIS-2026-001',
    '0089123456',
    '1371012345670001',
    'Muhammad Al-Fatih',
    'L',
    'Padang',
    '2010-05-12',
    'Jl. Belimbing No. 12, Padang',
    FALSE,
    'Indonesia',
    'Asia/Jakarta',
    'Abdullah',
    '+6281122334455',
    'Ayah',
    'ACTIVE'
),
(
    's0000002-0000-0000-0000-000000000002',
    'u0000016-0000-0000-0000-000000000016',
    'NIS-2026-002',
    '0089123457',
    '1371012345670002',
    'Ahmad Zaki',
    'L',
    'Padang',
    '2010-08-20',
    'Jl. Gajah Mada No. 45, Padang',
    FALSE,
    'Indonesia',
    'Asia/Jakarta',
    'Zulkifli',
    '+6281122334466',
    'Ayah',
    'ACTIVE'
),
(
    's0000003-0000-0000-0000-000000000003',
    'u0000019-0000-0000-0000-000000000019',
    'NIS-2026-003',
    '0081234561',
    '1371052809070002',
    'Muhammad Abdullah Ziyad',
    'L',
    'Padang',
    '2007-09-28',
    'Jl. Tuanku Tambusai No. 12, Kel. Padang Barat, Kota Padang',
    FALSE,
    'Indonesia',
    'Asia/Jakarta',
    'H. Abu Ziyad',
    '+6281276543210',
    'Ayah',
    'ACTIVE'
),
(
    's0000004-0000-0000-0000-000000000004',
    'u0000020-0000-0000-0000-000000000020',
    'NIS-2026-004',
    '0081234562',
    '1371052809070003',
    'Yafiq Al-Faruq',
    'L',
    'Padang',
    '2007-11-14',
    'Jl. Andalas No. 5, Padang Timur, Kota Padang',
    FALSE,
    'Indonesia',
    'Asia/Jakarta',
    'Faruq Basalamah',
    '+6281276543211',
    'Ayah',
    'ACTIVE'
)
ON DUPLICATE KEY UPDATE full_name=VALUES(full_name);

-- 10. PEMETAAN SANTRI KE ROMBEL KELAS X
INSERT INTO classroom_students (classroom_id, student_id) VALUES
('c1111111-1111-1111-1111-111111111111', 's0000001-0000-0000-0000-000000000001'),
('c1111111-1111-1111-1111-111111111111', 's0000002-0000-0000-0000-000000000002'),
('c1111111-1111-1111-1111-111111111111', 's0000003-0000-0000-0000-000000000003'),
('c1111111-1111-1111-1111-111111111111', 's0000004-0000-0000-0000-000000000004')
ON DUPLICATE KEY UPDATE joined_at=CURRENT_TIMESTAMP;

SET FOREIGN_KEY_CHECKS = 1;

