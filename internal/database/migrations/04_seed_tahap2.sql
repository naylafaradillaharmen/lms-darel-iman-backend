-- ==============================================================================
-- SEED DATA: PKBM DIGITAL DAR EL IMAN - TAHAP 2
-- Kurikulum LMS, Jadwal KBM, Penugasan & Halaqah Tahfizh Terpadu
-- ==============================================================================



SET FOREIGN_KEY_CHECKS = 0;

-- 1. MATA PELAJARAN
INSERT INTO subjects (id, code, name, category, description) VALUES
('sub00001-0000-0000-0000-000000000001', 'THF-01', 'Tahfizh Al-Qur\'an Tematik', 'KHAS_DAR_EL_IMAN', 'Hafalan Al-Qur\'an Juz 30-28 dengan kaidah tajwid dan pemahaman maknawi'),
('sub00002-0000-0000-0000-000000000002', 'ARB-01', 'Bahasa Arab Komunikatif & Nahwu Shorof Dasar', 'KHAS_DAR_EL_IMAN', 'Kaidah bahasa Arab Al-Qur\'an, percakapan harian dan kaidah tata bahasa dasar'),
('sub00003-0000-0000-0000-000000000003', 'AKD-01', 'Akidah & Akhlak Karimah', 'KHAS_DAR_EL_IMAN', 'Penanaman tauhid yang lurus dan adab penuntut ilmu sesuai sunnah'),
('sub00004-0000-0000-0000-000000000004', 'TEC-01', 'AI Literacy & Pemrograman Dasar', 'KHAS_DAR_EL_IMAN', 'Pemanfaatan kecerdasan buatan etis, logika komputasi, dan skill digital masa depan'),
('sub00005-0000-0000-0000-000000000005', 'IND-10', 'Bahasa Indonesia (Paket C)', 'UMUM_KEMENDIKBUD', 'Kurikulum Merdeka Kesetaraan: Literasi, teks argumentasi, dan komunikasi formal'),
('sub00006-0000-0000-0000-000000000006', 'MTK-10', 'Matematika Terapan (Paket C)', 'UMUM_KEMENDIKBUD', 'Kurikulum Merdeka Kesetaraan: Aljabar, statistika data dasar, dan pemecahan masalah'),
('sub00007-0000-0000-0000-000000000007', 'ENG-10', 'Bahasa Inggris Global (Paket C)', 'UMUM_KEMENDIKBUD', 'Kurikulum Merdeka Kesetaraan: English for academic and global communication')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 2. PENGIKATAN MAPEL KE KELAS 10 PAKET C
INSERT INTO grade_subjects (grade_level_id, subject_id, credit_hours) VALUES
('c0000010-0000-0000-0000-000000000010', 'sub00001-0000-0000-0000-000000000001', 4),
('c0000010-0000-0000-0000-000000000010', 'sub00002-0000-0000-0000-000000000002', 3),
('c0000010-0000-0000-0000-000000000010', 'sub00003-0000-0000-0000-000000000003', 2),
('c0000010-0000-0000-0000-000000000010', 'sub00004-0000-0000-0000-000000000004', 2),
('c0000010-0000-0000-0000-000000000010', 'sub00005-0000-0000-0000-000000000005', 3),
('c0000010-0000-0000-0000-000000000010', 'sub00006-0000-0000-0000-000000000006', 3)
ON DUPLICATE KEY UPDATE credit_hours=VALUES(credit_hours);

-- 3. COURSE OFFERING (RUANG KELAS AJAR LMS)
INSERT INTO course_classes (id, classroom_id, subject_id, tutor_id, academic_year_id) VALUES
('cc000001-0000-0000-0000-000000000001', 'c1111111-1111-1111-1111-111111111111', 'sub00001-0000-0000-0000-000000000001', 't0000001-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222'),
('cc000002-0000-0000-0000-000000000002', 'c1111111-1111-1111-1111-111111111111', 'sub00002-0000-0000-0000-000000000002', 't0000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222')
ON DUPLICATE KEY UPDATE tutor_id=VALUES(tutor_id);

-- 4. JADWAL KBM SINKRON
INSERT INTO class_schedules (id, course_class_id, day_of_week, start_time, end_time, meeting_link) VALUES
('sch00001-0000-0000-0000-000000000001', 'cc000001-0000-0000-0000-000000000001', 1, '08:00:00', '09:30:00', 'https://meet.google.com/pkbm-thf-x10'),
('sch00002-0000-0000-0000-000000000002', 'cc000002-0000-0000-0000-000000000002', 3, '10:00:00', '11:30:00', 'https://meet.google.com/pkbm-arb-x10')
ON DUPLICATE KEY UPDATE meeting_link=VALUES(meeting_link);

-- 5. TOPIK SILABUS MINGGUAN LMS
INSERT INTO course_topics (id, course_class_id, title, sequence_order, description, is_published) VALUES
('top00001-0000-0000-0000-000000000001', 'cc000001-0000-0000-0000-000000000001', 'Pertemuan 1: Pengantar Adab Penuntut Ilmu & Makharijul Huruf', 1, 'Membahas niat ikhlas menuntut ilmu dan artikulasi makhraj huruf hijaiyyah', TRUE),
('top00002-0000-0000-0000-000000000002', 'cc000001-0000-0000-0000-000000000001', 'Pertemuan 2: Sifatul Huruf & Pengenalan Hukum Nun Mati / Tanwin', 2, 'Karakteristik bunyi huruf dan hukum Idzhar, Idgham, Ikhfa, Iqlab', TRUE)
ON DUPLICATE KEY UPDATE title=VALUES(title);

-- 6. MATERI AJAR LMS
INSERT INTO topic_materials (id, course_topic_id, title, material_type, content_body, file_url, external_link, sequence_order) VALUES
('mat00001-0000-0000-0000-000000000001', 'top00001-0000-0000-0000-000000000001', 'Buku Panduan Tajwid Asy-Syafi\'i Ringkas (PDF)', 'DOCUMENT', NULL, 'https://cdn.darel-iman.or.id/modul/tajwid-dasar-juz30.pdf', NULL, 1),
('mat00002-0000-0000-0000-000000000001', 'top00001-0000-0000-0000-000000000001', 'Video Tutorial: Visualisasi 17 Makharijul Huruf', 'VIDEO_LINK', NULL, NULL, 'https://www.youtube.com/watch?v=dQw4w9WgXcQ', 2),
('mat00003-0000-0000-0000-000000000001', 'top00001-0000-0000-0000-000000000001', 'Catatan Pembuka Ustadz: Keutamaan Al-Qur\'an', 'TEXT_ARTICLE', 'Rasulullah shallallahu alaihi wa sallam bersabda: Sebaik-baik kalian adalah orang yang belajar Al-Qur\'an dan mengajarkannya. (HR. Al-Bukhari). Mari kita luruskan niat dalam menempuh semester ini.', NULL, NULL, 3)
ON DUPLICATE KEY UPDATE title=VALUES(title);

-- 7. PENUGASAN (ASSIGNMENT)
INSERT INTO assignments (id, course_topic_id, title, instructions, due_date, max_score, file_attachment_url) VALUES
('asg00001-0000-0000-0000-000000000001', 'top00001-0000-0000-0000-000000000001', 'Tugas 1: Rangkuman Adab & Rekaman Suara Surah An-Naba 1-10', '1. Tulis ringkasan 1 paragraf tentang adab Al-Qur\'an.\n2. Rekam audio bacaan Surah An-Naba\' ayat 1 sampai 10 dan cantumkan link atau lampirkan berkas jawaban.', '2026-10-20 23:59:00', 100.00, 'https://cdn.darel-iman.or.id/tugas/petunjuk-tugas-01.pdf')
ON DUPLICATE KEY UPDATE title=VALUES(title);

-- SUBMISSION SANTRI (KEDUA SANTRI SINKRON DENGAN BUKU NILAI)
INSERT INTO assignment_submissions (id, assignment_id, student_id, submission_text, file_url, score, feedback, graded_by_tutor_id, graded_at) VALUES
('subm0001-0000-0000-0000-000000000001', 'asg00001-0000-0000-0000-000000000001', 's0000001-0000-0000-0000-000000000001', 'Bismillah, berikut setoran rekaman audio An-Naba 1-10 dan ringkasan adab saya ustadz: Niat ikhlas semata karena Allah dan berwudhu sebelum membaca Al-Quran.', 'https://cdn.darel-iman.or.id/submissions/fatih-an-naba-audio.mp3', 95.00, 'Masya Allah mumtaz bacaan tajwid sangat baik dan makhorijul huruf jelas. Pertahankan!', 't0000001-0000-0000-0000-000000000001', '2026-10-09 10:00:00'),
('subm0002-0000-0000-0000-000000000002', 'asg00001-0000-0000-0000-000000000001', 's0000002-0000-0000-0000-000000000002', 'Bismillah, tugas rekaman surah An-Naba 1-10 dan adab penuntut ilmu dari Ahmad Zaki.', 'https://cdn.darel-iman.or.id/submissions/zaki-an-naba-audio.mp3', 65.00, 'Perhatikan kaidah makhraj dan mad thabi\'i, latihan lebih intensif ya.', 't0000001-0000-0000-0000-000000000001', '2026-10-09 11:00:00')
ON DUPLICATE KEY UPDATE score=VALUES(score);

-- 8. SESI KELAS & PRESENSI (KEDUA SANTRI SINKRON)
INSERT INTO class_sessions (id, course_class_id, session_title, session_date, start_time, end_time) VALUES
('sess0001-0000-0000-0000-000000000001', 'cc000001-0000-0000-0000-000000000001', 'Sesi 1: Pembukaan & Talaqqi Surah An-Naba', '2026-10-08', '08:00:00', '09:30:00')
ON DUPLICATE KEY UPDATE session_title=VALUES(session_title);

INSERT INTO session_attendances (id, class_session_id, student_id, status, notes) VALUES
('att00001-0000-0000-0000-000000000001', 'sess0001-0000-0000-0000-000000000001', 's0000001-0000-0000-0000-000000000001', 'PRESENT', 'Santri hadir tepat waktu via Google Meet (100%)'),
('att00002-0000-0000-0000-000000000002', 'sess0001-0000-0000-0000-000000000001', 's0000002-0000-0000-0000-000000000002', 'LATE', 'Santri hadir terlambat 15 menit (Presensi 70%)')
ON DUPLICATE KEY UPDATE status=VALUES(status);

-- 9. HALAQAH TAHFIZH AL-QUR'AN (2 SANTRI AKTIF)
INSERT INTO tahfizh_groups (id, academic_year_id, name, tutor_id, target_juz) VALUES
('tg000001-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'Halaqah Imam Nafi\' - Ikhwan Kelas X', 't0000001-0000-0000-0000-000000000001', 'Juz 30')
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- Anggota Halaqah: Fatih & Ahmad Zaki
INSERT INTO tahfizh_group_students (tahfizh_group_id, student_id) VALUES
('tg000001-0000-0000-0000-000000000001', 's0000001-0000-0000-0000-000000000001'),
('tg000001-0000-0000-0000-000000000001', 's0000002-0000-0000-0000-000000000002')
ON DUPLICATE KEY UPDATE student_id=VALUES(student_id);

-- Log Setoran Tahfizh
INSERT INTO tahfizh_logs (id, student_id, tutor_id, log_date, type, surah_name, start_ayah, end_ayah, fluency_score, tajweed_score, audio_record_url, notes) VALUES
('tl000001-0000-0000-0000-000000000001', 's0000001-0000-0000-0000-000000000001', 't0000001-0000-0000-0000-000000000001', '2026-10-08', 'ZIYADAH', 'An-Naba\'', 1, 40, 'MUMTAZ', 'MUMTAZ', 'https://cdn.darel-iman.or.id/audio/setoran-fatih-an-naba.mp3', 'Hafalan lancar, makhraj ghunnah dan mad thabi\'i tepat.'),
('tl000002-0000-0000-0000-000000000001', 's0000001-0000-0000-0000-000000000001', 't0000001-0000-0000-0000-000000000001', '2026-10-07', 'MURAJAAH', 'An-Nazi\'at', 1, 20, 'JAYYID_JIDDAN', 'MUMTAZ', NULL, 'Muraja\'ah bagus, perlu diperhatikan waqaf di ayat 15.'),
('tl000003-0000-0000-0000-000000000003', 's0000002-0000-0000-0000-000000000002', 't0000001-0000-0000-0000-000000000001', '2026-10-07', 'ZIYADAH', 'An-Naba\'', 1, 20, 'JAYYID', 'JAYYID', NULL, 'Hafalan lancar, makhraj huruf istila perlu dipertebal (Nilai 72.0).')
ON DUPLICATE KEY UPDATE surah_name=VALUES(surah_name);

SET FOREIGN_KEY_CHECKS = 1;

