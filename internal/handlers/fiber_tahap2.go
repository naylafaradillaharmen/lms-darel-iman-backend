package handlers

import (
	"fmt"
	"sync"
	"time"

	"pkbm-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberTahap2Handler struct {
	GormDB *gorm.DB
	mu     sync.RWMutex

	// Data Tahap 2 (KBM, Modul, Tugas, Tahfizh)
	subjects    []models.Subject
	classes     []models.CourseClass
	schedules   []models.ClassSchedule
	topics      []models.CourseTopic
	assignments []models.Assignment
	submissions []models.AssignmentSubmission
	groups      []models.TahfizhGroup
	tahfizhLogs []models.TahfizhLog
	attendance  []models.AttendanceRecord
}

func NewFiberTahap2Handler(db *gorm.DB) *FiberTahap2Handler {
	desc1 := "Hafalan Al-Qur'an Juz 30-28 dengan kaidah tajwid sanad Hafsh 'an 'Ashim"
	desc2 := "Kaidah bahasa Arab Al-Qur'an, nahwu shorof aplikatif, dan percakapan harian"
	desc3 := "Penanaman aqidah shahihah bermanhaj salaf dan adab penuntut ilmu"
	desc4 := "Pemanfaatan AI etis, logika komputasi, dan literasi teknologi digital santri"
	desc5 := "Kurikulum Merdeka Kesetaraan: Literasi, teks argumentasi, dan komunikasi formal"
	desc6 := "Kurikulum Merdeka Kesetaraan: Aljabar terapan dan pemecahan masalah kontekstual"

	meetLink := "https://meet.google.com/pkbm-live-kbm"
	pdfDoc := "https://cdn.darel-iman.or.id/modul/tajwid-asy-syafii-ringkas.pdf"
	videoTutorial := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	audioSetoran := "https://cdn.darel-iman.or.id/audio/setoran-an-naba-ziyad.mp3"

	score95 := 95.0
	feedbackTutor := "Masya Allah, bacaan tajwid sangat jelas, mad thabi'i konsisten 2 harakat."
	targetJuz30 := "Juz 30 (An-Naba s.d. An-Nas)"
	mumtaz := "MUMTAZ"
	jayyidJiddan := "JAYYID_JIDDAN"
	tutorName := "Ustadz Ahmad Fauzi, Lc."
	tutorID := "t0000009"

	h := &FiberTahap2Handler{
		GormDB: db,
		subjects: []models.Subject{
			{ID: "sub-01", Code: "THF-01", Name: "Tahfizh Al-Qur'an Bersanad", Category: "KHAS_DAR_EL_IMAN", Description: &desc1, CreatedAt: time.Now()},
			{ID: "sub-02", Code: "ARB-01", Name: "Bahasa Arab Komunikatif & Nahwu", Category: "KHAS_DAR_EL_IMAN", Description: &desc2, CreatedAt: time.Now()},
			{ID: "sub-03", Code: "AKD-01", Name: "Akidah & Adab Penuntut Ilmu", Category: "KHAS_DAR_EL_IMAN", Description: &desc3, CreatedAt: time.Now()},
			{ID: "sub-04", Code: "TEC-01", Name: "AI Literacy & Digital Skill", Category: "KHAS_DAR_EL_IMAN", Description: &desc4, CreatedAt: time.Now()},
			{ID: "sub-05", Code: "IND-10", Name: "Bahasa Indonesia (Paket C)", Category: "UMUM_KEMENDIKBUD", Description: &desc5, CreatedAt: time.Now()},
			{ID: "sub-06", Code: "MTK-10", Name: "Matematika Terapan (Paket C)", Category: "UMUM_KEMENDIKBUD", Description: &desc6, CreatedAt: time.Now()},
		},
		classes: []models.CourseClass{
			{
				ID:               "cc-01",
				ClassroomID:      "cr-x-paketc",
				ClassroomName:    "Kelas X Paket C — Umar bin Khattab",
				SubjectID:        "sub-01",
				SubjectCode:      "THF-01",
				SubjectName:      "Tahfizh Al-Qur'an Bersanad",
				TutorID:          "t0000009",
				TutorName:        "Ustadz Ahmad Fauzi, Lc.",
				AcademicYearName: "2026/2027 Ganjil",
			},
			{
				ID:               "cc-02",
				ClassroomID:      "cr-x-paketc",
				ClassroomName:    "Kelas X Paket C — Umar bin Khattab",
				SubjectID:        "sub-02",
				SubjectCode:      "ARB-01",
				SubjectName:      "Bahasa Arab Komunikatif & Nahwu",
				TutorID:          "t0000009",
				TutorName:        "Ustadz Ahmad Fauzi, Lc.",
				AcademicYearName: "2026/2027 Ganjil",
			},
			{
				ID:               "cc-03",
				ClassroomID:      "cr-x-paketc",
				ClassroomName:    "Kelas X Paket C — Umar bin Khattab",
				SubjectID:        "sub-05",
				SubjectCode:      "IND-10",
				SubjectName:      "Bahasa Indonesia (Paket C)",
				TutorID:          "t0000010",
				TutorName:        "Ustadzah Fatimah, S.Pd.",
				AcademicYearName: "2026/2027 Ganjil",
			},
			{
				ID:               "cc-04",
				ClassroomID:      "cr-x-paketc",
				ClassroomName:    "Kelas X Paket C — Umar bin Khattab",
				SubjectID:        "sub-06",
				SubjectCode:      "MTK-10",
				SubjectName:      "Matematika Terapan (Paket C)",
				TutorID:          "t0000007",
				TutorName:        "Ahmad Fauzan, S.Kom., M.T.",
				AcademicYearName: "2026/2027 Ganjil",
			},
		},
		schedules: []models.ClassSchedule{
			{ID: "sch-01", CourseClassID: "cc-01", DayOfWeek: 1, StartTime: "08:00", EndTime: "09:30", MeetingLink: &meetLink, IsActive: true, CreatedAt: time.Now()},
			{ID: "sch-02", CourseClassID: "cc-02", DayOfWeek: 2, StartTime: "08:00", EndTime: "09:30", MeetingLink: &meetLink, IsActive: true, CreatedAt: time.Now()},
			{ID: "sch-03", CourseClassID: "cc-03", DayOfWeek: 3, StartTime: "10:00", EndTime: "11:30", MeetingLink: nil, IsActive: false, CreatedAt: time.Now()},
			{ID: "sch-04", CourseClassID: "cc-04", DayOfWeek: 4, StartTime: "10:00", EndTime: "11:30", MeetingLink: nil, IsActive: false, CreatedAt: time.Now()},
		},
		topics: []models.CourseTopic{
			{
				ID:            "top-01",
				CourseClassID: "cc-01",
				Title:         "Pertemuan 1: Makharijul Huruf & Karakteristik Huruf Halqiyah",
				SequenceOrder: 1,
				Description:   &desc1,
				IsPublished:   true,
				Materials: []models.TopicMaterial{
					{
						ID:            "mat-01",
						CourseTopicID: "top-01",
						Title:         "Buku Saku Tajwid Asy-Syafi'i Praktis (PDF)",
						MaterialType:  "DOCUMENT",
						FileURL:       &pdfDoc,
						SequenceOrder: 1,
						CreatedAt:     time.Now(),
					},
					{
						ID:            "mat-02",
						CourseTopicID: "top-01",
						Title:         "Video Panduan: Simulasi 17 Titik Makhraj Huruf",
						MaterialType:  "VIDEO_LINK",
						ExternalLink:  &videoTutorial,
						SequenceOrder: 2,
						CreatedAt:     time.Now(),
					},
				},
			},
			{
				ID:            "top-02",
				CourseClassID: "cc-01",
				Title:         "Pertemuan 2: Kaidah Mad Thabi'i dan Ahkamul Mad",
				SequenceOrder: 2,
				Description:   &desc1,
				IsPublished:   true,
				Materials: []models.TopicMaterial{
					{
						ID:            "mat-03",
						CourseTopicID: "top-02",
						Title:         "Slide Presentasi Mad Thabi'i & Mad Far'i",
						MaterialType:  "DOCUMENT",
						FileURL:       &pdfDoc,
						SequenceOrder: 1,
						CreatedAt:     time.Now(),
					},
				},
			},
		},
		assignments: []models.Assignment{
			{
				ID:                "asg-01",
				CourseTopicID:     "top-01",
				CourseClassID:     "cc-01",
				SubjectName:       "Tahfizh Al-Qur'an Bersanad",
				TutorID:           "u0000009",
				TutorName:         "Ustadz Ahmad Fauzi, Lc.",
				Title:             "Tugas 1: Rekaman Suara Surah An-Naba' 1-16 & Catatan Tajwid",
				Instructions:      "1. Rekam lantunan Surah An-Naba' ayat 1 s.d. 16 dengan tartil.\n2. Sebutkan minimal 3 hukum mad thabi'i yang Anda temukan.",
				DueDate:           "2026-10-25 23:59:00",
				MaxScore:          100.0,
				FileAttachmentURL: &pdfDoc,
				SubmissionCount:   1,
				CreatedAt:         time.Now(),
			},
			{
				ID:                "asg-02",
				CourseTopicID:     "top-03",
				CourseClassID:     "cc-03",
				SubjectName:       "Bahasa Indonesia (Paket C)",
				TutorID:           "u0000010",
				TutorName:         "Ustadzah Fatimah, S.Pd.",
				Title:             "Tugas 1: Analisis Struktur Teks Argumentasi & Kaidah Kebahasaan",
				Instructions:      "Bacalah teks opini terlampir, tentukan tesis, argumentasi pendukung, dan simpulan dalam format esai.",
				DueDate:           "2026-10-26 23:59:00",
				MaxScore:          100.0,
				FileAttachmentURL: &pdfDoc,
				SubmissionCount:   1,
				CreatedAt:         time.Now(),
			},
			{
				ID:                "asg-03",
				CourseTopicID:     "top-04",
				CourseClassID:     "cc-04",
				SubjectName:       "Matematika Terapan (Paket C)",
				TutorID:           "u0000007",
				TutorName:         "Ahmad Fauzan, S.Kom., M.T.",
				Title:             "Tugas 1: Pengolahan Data Distribusi Frekuensi Nilai & Statistika",
				Instructions:      "Hitung mean, median, dan modus dari tabel data terlampir dan buatkan kesimpulan analisisnya.",
				DueDate:           "2026-10-27 23:59:00",
				MaxScore:          100.0,
				FileAttachmentURL: &pdfDoc,
				SubmissionCount:   1,
				CreatedAt:         time.Now(),
			},
		},
		submissions: []models.AssignmentSubmission{
			{
				ID:                "subm-01",
				AssignmentID:      "asg-01",
				AssignmentTitle:   "Tugas 1: Rekaman Suara Surah An-Naba' 1-16 & Catatan Tajwid",
				SubjectName:       "Tahfizh Al-Qur'an Bersanad",
				TutorID:           "u0000009",
				TutorName:         "Ustadz Ahmad Fauzi, Lc.",
				StudentID:         "u0000012b",
				StudentName:       "Muhammad Abdullah Ziyad",
				SubmissionText:    ptrString("Bismillah ustadz, berikut rekaman hafalan Surah An-Naba 1-16 dan catatan hukum mad."),
				FileURL:           &audioSetoran,
				SubmittedAt:       time.Now().Add(-5 * time.Hour),
				Score:             &score95,
				Feedback:          &feedbackTutor,
				GradedByTutorID:   &tutorID,
				GradedByTutorName: &tutorName,
			},
			{
				ID:                "subm-02",
				AssignmentID:      "asg-02",
				AssignmentTitle:   "Tugas 1: Analisis Struktur Teks Argumentasi & Kaidah Kebahasaan",
				SubjectName:       "Bahasa Indonesia (Paket C)",
				TutorID:           "u0000010",
				TutorName:         "Ustadzah Fatimah, S.Pd.",
				StudentID:         "u0000012b",
				StudentName:       "Muhammad Abdullah Ziyad",
				SubmissionText:    ptrString("Bismillah ustadzah, berikut tugas analisis teks argumentasi."),
				FileURL:           &pdfDoc,
				SubmittedAt:       time.Now().Add(-8 * time.Hour),
				Score:             ptrFloat(88.0),
				Feedback:          ptrString("Analisis struktur esai sangat rapi dan argumentatif. Pertahankan!"),
				GradedByTutorID:   ptrString("u0000010"),
				GradedByTutorName: ptrString("Ustadzah Fatimah, S.Pd."),
			},
			{
				ID:                "subm-03",
				AssignmentID:      "asg-03",
				AssignmentTitle:   "Tugas 1: Pengolahan Data Distribusi Frekuensi Nilai & Statistika",
				SubjectName:       "Matematika Terapan (Paket C)",
				TutorID:           "u0000007",
				TutorName:         "Ahmad Fauzan, S.Kom., M.T.",
				StudentID:         "u0000013b",
				StudentName:       "Yafiq Al-Faruq",
				SubmissionText:    ptrString("Bismillah pak, berikut jawaban perhitungan tabel statistika."),
				FileURL:           &pdfDoc,
				SubmittedAt:       time.Now().Add(-10 * time.Hour),
				Score:             ptrFloat(92.0),
				Feedback:          ptrString("Perhitungan distribusi frekuensi dan mean sangat teliti. Bagus sekali!"),
				GradedByTutorID:   ptrString("u0000007"),
				GradedByTutorName: ptrString("Ahmad Fauzan, S.Kom., M.T."),
			},
		},
		groups: []models.TahfizhGroup{
			{
				ID:               "tg-01",
				AcademicYearName: "2026/2027 Ganjil",
				Name:             "Halaqah Imam Asy-Syathibi (Paket C Ikhwan)",
				TutorID:          "u0000009",
				TutorName:        "Ustadz Ahmad Fauzi, Lc.",
				TargetJuz:        &targetJuz30,
				StudentCount:     15,
			},
		},
		tahfizhLogs: []models.TahfizhLog{
			{
				ID:             "thf-01",
				StudentID:      "u0000012b",
				StudentName:    "Muhammad Abdullah Ziyad",
				TutorID:        "u0000009",
				TutorName:      "Ustadz Ahmad Fauzi, Lc.",
				LogDate:        time.Now().Add(-24 * time.Hour).Format("2006-01-02"),
				Type:           "ZIYADAH",
				SurahName:      "An-Naba'",
				StartAyah:      1,
				EndAyah:        16,
				FluencyScore:   &mumtaz,
				TajweedScore:   &mumtaz,
				Notes:          &feedbackTutor,
				AudioRecordURL: &audioSetoran,
				Status:         "SUDAH_DINILAI",
				CreatedAt:      time.Now().Add(-24 * time.Hour),
			},
			{
				ID:             "thf-02",
				StudentID:      "u0000012b",
				StudentName:    "Muhammad Abdullah Ziyad",
				TutorID:        "u0000009",
				TutorName:      "Ustadz Ahmad Fauzi, Lc.",
				LogDate:        time.Now().Add(-48 * time.Hour).Format("2006-01-02"),
				Type:           "MURAJAAH",
				SurahName:      "An-Nazi'at",
				StartAyah:      1,
				EndAyah:        20,
				FluencyScore:   &jayyidJiddan,
				TajweedScore:   &jayyidJiddan,
				Notes:          ptrString("Hafalan lancar, perhatikan panjang mad jaiz munfashil."),
				AudioRecordURL: &audioSetoran,
				Status:         "SUDAH_DINILAI",
				CreatedAt:      time.Now().Add(-48 * time.Hour),
			},
			{
				ID:             "thf-03",
				StudentID:      "u0000013b",
				StudentName:    "Yafiq Al-Faruq",
				TutorID:        "u0000009",
				TutorName:      "Ustadz Ahmad Fauzi, Lc.",
				LogDate:        time.Now().Format("2006-01-02"),
				Type:           "ZIYADAH",
				SurahName:      "An-Naba'",
				StartAyah:      1,
				EndAyah:        10,
				FluencyScore:   ptrString("MENUNGGU_NILAI"),
				TajweedScore:   ptrString("MENUNGGU_NILAI"),
				Notes:          ptrString("Bismillah ustadz, setoran hafalan Surah An-Naba ayat 1-10."),
				AudioRecordURL: &audioSetoran,
				Status:         "MENUNGGU_NILAI",
				CreatedAt:      time.Now(),
			},
		},
		attendance: []models.AttendanceRecord{
			{ID: "att-01", Date: "2026-10-18", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:55 WIB", Notes: "Tepat waktu di Google Meet"},
			{ID: "att-02", Date: "2026-10-18", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NISN: "0081234562", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:58 WIB", Notes: "Hadir aktif"},
			{ID: "att-03", Date: "2026-10-18", StudentID: "u0000014b", StudentName: "Ahmad Faris", NISN: "0081234563", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:02 WIB", Notes: "Tepat waktu"},
			{ID: "att-04", Date: "2026-10-18", StudentID: "u0000015b", StudentName: "Syamil Basayif", NISN: "0081234564", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "IZIN", ApprovalStatus: "PENDING", Reason: "Keperluan mendesak keluarga di luar kota (Surat terlampir)", TimeRecorded: "07:30 WIB", Notes: "Menunggu persetujuan Ustadz pengampu"},
			{ID: "att-05", Date: "2026-10-18", StudentID: "u0000016b", StudentName: "Zaidan Al-Ayyubi", NISN: "0081234565", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:50 WIB", Notes: "Tepat waktu"},
			{ID: "att-06", Date: "2026-10-18", StudentID: "u0000017b", StudentName: "Hamzah Asadullah", NISN: "0081234566", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:00 WIB", Notes: "Tepat waktu"},
			{ID: "att-07", Date: "2026-10-18", StudentID: "u0000018b", StudentName: "Salman Al-Farisi", NISN: "0081234567", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:54 WIB", Notes: "Tepat waktu"},
			{ID: "att-08", Date: "2026-10-18", StudentID: "u0000019b", StudentName: "Khalid bin Walid", NISN: "0081234568", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:52 WIB", Notes: "Tepat waktu"},
			{ID: "att-09", Date: "2026-10-18", StudentID: "u0000020b", StudentName: "Bilal bin Rabah", NISN: "0081234569", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "SAKIT", ApprovalStatus: "APPROVED", Reason: "Demam tinggi berobat ke klinik", TimeRecorded: "07:15 WIB", Notes: "Disetujui Ustadz (surat dokter terlampir)"},
			{ID: "att-10", Date: "2026-10-18", StudentID: "u0000021b", StudentName: "Thariq bin Ziyad", NISN: "0081234560", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:01 WIB", Notes: "Hadir aktif"},
			{ID: "att-11", Date: "2026-10-18", StudentID: "u0000022b", StudentName: "Usamah bin Zaid", NISN: "0081234571", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:59 WIB", Notes: "Hadir tepat waktu"},
			{ID: "att-12", Date: "2026-10-18", StudentID: "u0000023b", StudentName: "Mus'ab bin Umair", NISN: "0081234572", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:48 WIB", Notes: "Hadir tepat waktu"},
			{ID: "att-13", Date: "2026-10-17", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:00 WIB", Notes: "Hadir aktif"},
			{ID: "att-14", Date: "2026-10-17", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NISN: "0081234562", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Bahasa Arab Komunikatif & Nahwu", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "07:59 WIB", Notes: "Hadir aktif"},
			{ID: "att-15", Date: "2026-10-16", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Matematika Terapan", TutorID: "u0000007", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "10:00 WIB", Notes: "Hadir penuh"},
			{ID: "att-16", Date: "2026-10-15", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Bahasa Indonesia", TutorID: "u0000010", TutorName: "Ustadzah Fatimah, S.Pd.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "10:00 WIB", Notes: "Hadir penuh"},
			{ID: "att-17", Date: "2026-10-10", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:00 WIB", Notes: "Hadir tepat waktu"},
			{ID: "att-18", Date: "2026-10-05", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Tahfizh Al-Qur'an Bersanad", TutorID: "u0000009", TutorName: "Ustadz Ahmad Fauzi, Lc.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "08:00 WIB", Notes: "Hadir tepat waktu"},
			{ID: "att-19", Date: "2026-10-18", StudentID: "u0000013b", StudentName: "Yafiq Al-Faruq", NISN: "0081234562", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Matematika Terapan", TutorID: "u0000007", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Status: "IZIN", ApprovalStatus: "PENDING", Reason: "Acara keluarga mendesak di luar kota", TimeRecorded: "08:15 WIB", Notes: "Menunggu persetujuan Ustadz pengampu"},
			{ID: "att-20", Date: "2026-10-18", StudentID: "u0000012b", StudentName: "Muhammad Abdullah Ziyad", NISN: "0081234561", ClassroomName: "Kelas X Paket C — Umar bin Khattab", SubjectName: "Matematika Terapan", TutorID: "u0000007", TutorName: "Ahmad Fauzan, S.Kom., M.T.", Status: "HADIR", ApprovalStatus: "APPROVED", TimeRecorded: "09:55 WIB", Notes: "Hadir tepat waktu"},
		},
	}
	return h
}

// 1. Get Subjects (Mata Pelajaran)
func (h *FiberTahap2Handler) GetSubjects(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.subjects),
		"data":    h.subjects,
	})
}

// 2. Get Classes & Schedules (Kelas & Jadwal KBM)
func (h *FiberTahap2Handler) GetClasses(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.classes),
		"data":    h.classes,
	})
}

func (h *FiberTahap2Handler) GetSchedules(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.schedules),
		"data":    h.schedules,
	})
}

// 3. Get Course Topics & Materials (Modul Ajar)
func (h *FiberTahap2Handler) GetTopics(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.topics),
		"data":    h.topics,
	})
}

// 4. Assignments & Submissions (Tugas & Pengumpulan)
func (h *FiberTahap2Handler) GetAssignments(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.assignments),
		"data":    h.assignments,
	})
}

func (h *FiberTahap2Handler) GetSubmissions(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.submissions),
		"data":    h.submissions,
	})
}

type FiberSubmitAssignmentReq struct {
	AssignmentID   string `json:"assignment_id"`
	StudentID      string `json:"student_id"`
	StudentName    string `json:"student_name"`
	SubmissionText string `json:"submission_text"`
	FileURL        string `json:"file_url"`
}

func (h *FiberTahap2Handler) SubmitAssignment(c *fiber.Ctx) error {
	var req FiberSubmitAssignmentReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tugas tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	newID := fmt.Sprintf("subm-%d", time.Now().Unix())
	newSubm := models.AssignmentSubmission{
		ID:             newID,
		AssignmentID:   req.AssignmentID,
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		SubmissionText: &req.SubmissionText,
		FileURL:        &req.FileURL,
		SubmittedAt:    time.Now(),
	}

	h.submissions = append([]models.AssignmentSubmission{newSubm}, h.submissions...)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alhamdulillah, tugas berhasil dikirim ke tutor pengampu!",
		"data":    newSubm,
	})
}

type FiberGradeSubmissionReq struct {
	SubmissionID string  `json:"submission_id"`
	Score        float64 `json:"score"`
	Feedback     string  `json:"feedback"`
	TutorID      string  `json:"tutor_id"`
	TutorName    string  `json:"tutor_name"`
}

func (h *FiberTahap2Handler) GradeSubmission(c *fiber.Ctx) error {
	var req FiberGradeSubmissionReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format penilaian tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, s := range h.submissions {
		if s.ID == req.SubmissionID {
			h.submissions[i].Score = &req.Score
			h.submissions[i].Feedback = &req.Feedback
			h.submissions[i].GradedByTutorID = &req.TutorID
			h.submissions[i].GradedByTutorName = &req.TutorName
			now := time.Now()
			h.submissions[i].GradedAt = &now

			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Alhamdulillah, nilai %.1f dan masukan berhasil disimpan!", req.Score),
				"data":    h.submissions[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Data pengumpulan tugas tidak ditemukan.",
	})
}

// 5. Tahfizh Groups & Logs (Halaqah Tahfizh Bersanad)
func (h *FiberTahap2Handler) GetTahfizhGroups(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.groups),
		"data":    h.groups,
	})
}

func (h *FiberTahap2Handler) GetTahfizhLogs(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(h.tahfizhLogs),
		"data":    h.tahfizhLogs,
	})
}

type FiberCreateTahfizhLogReq struct {
	StudentID   string `json:"student_id"`
	StudentName string `json:"student_name"`
	LogType     string `json:"log_type"`
	SurahName   string `json:"surah_name"`
	AyatStart   int    `json:"ayat_start"`
	AyatEnd     int    `json:"ayat_end"`
	Grade       string `json:"grade"`
	Notes       string `json:"notes"`
	AudioURL    string `json:"audio_url"`
	TutorName   string `json:"tutor_name"`
	TutorID     string `json:"tutor_id"`
}

func (h *FiberTahap2Handler) CreateTahfizhLog(c *fiber.Ctx) error {
	var req FiberCreateTahfizhLogReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format setoran tahfizh tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	newID := fmt.Sprintf("thf-%d", time.Now().Unix())

	grade := req.Grade
	status := "SUDAH_DINILAI"
	if grade == "" || grade == "MENUNGGU_NILAI" || grade == "BELUM_DINILAI" {
		grade = "MENUNGGU_NILAI"
		status = "MENUNGGU_NILAI"
	}

	newLog := models.TahfizhLog{
		ID:             newID,
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		TutorID:        req.TutorID,
		TutorName:      req.TutorName,
		LogDate:        time.Now().Format("2006-01-02"),
		Type:           req.LogType,
		SurahName:      req.SurahName,
		StartAyah:      req.AyatStart,
		EndAyah:        req.AyatEnd,
		FluencyScore:   &grade,
		TajweedScore:   &grade,
		Notes:          &req.Notes,
		AudioRecordURL: &req.AudioURL,
		Status:         status,
		CreatedAt:      time.Now(),
	}

	h.tahfizhLogs = append([]models.TahfizhLog{newLog}, h.tahfizhLogs...)

	msg := fmt.Sprintf("Alhamdulillah, setoran hafalan Surah %s ayat %d-%d berhasil dikirim dan menunggu penilaian Musyrif!", req.SurahName, req.AyatStart, req.AyatEnd)
	if status == "SUDAH_DINILAI" {
		msg = fmt.Sprintf("Alhamdulillah, setoran hafalan Surah %s ayat %d-%d berhasil dicatat dengan predikat %s!", req.SurahName, req.AyatStart, req.AyatEnd, grade)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    newLog,
	})
}

// 2b. Update Schedule Meeting Link (Google Meet Control by Tutor)
type FiberUpdateMeetingLinkReq struct {
	ScheduleID  string `json:"schedule_id"`
	MeetingLink string `json:"meeting_link"`
	IsActive    bool   `json:"is_active"`
}

func (h *FiberTahap2Handler) UpdateMeetingLink(c *fiber.Ctx) error {
	var req FiberUpdateMeetingLinkReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format data meeting link tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, sch := range h.schedules {
		if sch.ID == req.ScheduleID {
			if req.MeetingLink != "" {
				h.schedules[i].MeetingLink = &req.MeetingLink
			} else if !req.IsActive {
				h.schedules[i].MeetingLink = nil
			}
			h.schedules[i].IsActive = req.IsActive

			statusText := "diaktifkan"
			if !req.IsActive {
				statusText = "dinonaktifkan/ditutup"
			}

			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Alhamdulillah, ruang Google Meet jadwal %s berhasil %s!", sch.ID, statusText),
				"data":    h.schedules[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Jadwal KBM tidak ditemukan.",
	})
}

// 2c. Presensi & Absensi KBM (Rekap Bulanan, Pekanan, Harian & PDF)
func (h *FiberTahap2Handler) GetAttendance(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	period := c.Query("period")       // "day", "week", "month", "all"
	status := c.Query("status")       // "HADIR", "IZIN", "SAKIT", "ALPA", "all"
	studentID := c.Query("student_id")

	var filtered []models.AttendanceRecord
	for _, att := range h.attendance {
		if studentID != "" && att.StudentID != studentID {
			continue
		}
		if status != "" && status != "ALL" && att.Status != status {
			continue
		}
		if period == "day" && att.Date != "2026-10-18" {
			continue
		}
		filtered = append(filtered, att)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(filtered),
		"data":    filtered,
	})
}

type FiberSaveAttendanceReq struct {
	StudentID     string `json:"student_id"`
	StudentName   string `json:"student_name"`
	NISN          string `json:"nisn"`
	ClassroomName string `json:"classroom_name"`
	SubjectName   string `json:"subject_name"`
	Status        string `json:"status"` // HADIR, IZIN, SAKIT, ALPA
	Notes         string `json:"notes"`
	Date          string `json:"date"`
}

func (h *FiberTahap2Handler) SaveAttendance(c *fiber.Ctx) error {
	var req FiberSaveAttendanceReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format data presensi tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	dateStr := req.Date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	for i, att := range h.attendance {
		if att.StudentID == req.StudentID && att.Date == dateStr {
			h.attendance[i].Status = req.Status
			h.attendance[i].Notes = req.Notes
			h.attendance[i].TimeRecorded = time.Now().Format("15:04") + " WIB"
			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Presensi %s diperbarui menjadi %s.", req.StudentName, req.Status),
				"data":    h.attendance[i],
			})
		}
	}

	newAtt := models.AttendanceRecord{
		ID:            fmt.Sprintf("att-%d", time.Now().Unix()),
		Date:          dateStr,
		StudentID:     req.StudentID,
		StudentName:   req.StudentName,
		NISN:          req.NISN,
		ClassroomName: req.ClassroomName,
		SubjectName:   req.SubjectName,
		Status:        req.Status,
		TimeRecorded:  time.Now().Format("15:04") + " WIB",
		Notes:         req.Notes,
	}

	h.attendance = append([]models.AttendanceRecord{newAtt}, h.attendance...)
	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Presensi %s berhasil dicatat (%s).", req.StudentName, req.Status),
		"data":    newAtt,
	})
}

// 2d. Student Check-In (Absen Mandiri & Pengajuan Izin/Sakit Santri)
type FiberStudentCheckInReq struct {
	StudentID     string `json:"student_id"`
	StudentName   string `json:"student_name"`
	NISN          string `json:"nisn"`
	ClassroomName string `json:"classroom_name"`
	SubjectName   string `json:"subject_name"`
	Status        string `json:"status"` // HADIR, IZIN, SAKIT
	Reason        string `json:"reason"`
	AttachmentURL string `json:"attachment_url"`
	Date          string `json:"date"`
}

func (h *FiberTahap2Handler) StudentCheckIn(c *fiber.Ctx) error {
	var req FiberStudentCheckInReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format presensi santri tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	dateStr := req.Date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	approvalStatus := "APPROVED"
	notes := "Santri hadir tepat waktu di kelas daring"
	successMsg := fmt.Sprintf("Alhamdulillah, presensi HADIR %s berhasil dikonfirmasi!", req.StudentName)

	if req.Status == "IZIN" || req.Status == "SAKIT" {
		approvalStatus = "PENDING"
		notes = fmt.Sprintf("Menunggu konfirmasi Ustadz: %s", req.Reason)
		successMsg = fmt.Sprintf("Alhamdulillah, pengajuan %s Anda berhasil dikirim dan menunggu verifikasi Ustadz pengampu.", req.Status)
	}

	for i, att := range h.attendance {
		if att.StudentID == req.StudentID && att.Date == dateStr {
			h.attendance[i].Status = req.Status
			h.attendance[i].ApprovalStatus = approvalStatus
			h.attendance[i].Reason = req.Reason
			h.attendance[i].AttachmentURL = req.AttachmentURL
			h.attendance[i].RejectionReason = ""
			h.attendance[i].Notes = notes
			h.attendance[i].TimeRecorded = time.Now().Format("15:04") + " WIB"

			return c.JSON(fiber.Map{
				"success": true,
				"message": successMsg,
				"data":    h.attendance[i],
			})
		}
	}

	newAtt := models.AttendanceRecord{
		ID:             fmt.Sprintf("att-%d", time.Now().Unix()),
		Date:           dateStr,
		StudentID:      req.StudentID,
		StudentName:    req.StudentName,
		NISN:           req.NISN,
		ClassroomName:  req.ClassroomName,
		SubjectName:    req.SubjectName,
		Status:         req.Status,
		ApprovalStatus: approvalStatus,
		Reason:         req.Reason,
		AttachmentURL:  req.AttachmentURL,
		TimeRecorded:   time.Now().Format("15:04") + " WIB",
		Notes:          notes,
		CreatedAt:      time.Now(),
	}

	h.attendance = append([]models.AttendanceRecord{newAtt}, h.attendance...)
	return c.JSON(fiber.Map{
		"success": true,
		"message": successMsg,
		"data":    newAtt,
	})
}

// 2e. Tutor Approve / Reject Izin Santri
type FiberApproveAttendanceReq struct {
	AttendanceID    string `json:"attendance_id"`
	Action          string `json:"action"` // APPROVE, REJECT
	RejectionReason string `json:"rejection_reason"`
	TutorName       string `json:"tutor_name"`
}

func (h *FiberTahap2Handler) ApproveAttendance(c *fiber.Ctx) error {
	var req FiberApproveAttendanceReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format persetujuan izin tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, att := range h.attendance {
		if att.ID == req.AttendanceID {
			if req.Action == "APPROVE" {
				h.attendance[i].ApprovalStatus = "APPROVED"
				h.attendance[i].Notes = fmt.Sprintf("Izin disetujui oleh %s", req.TutorName)
				return c.JSON(fiber.Map{
					"success": true,
					"message": fmt.Sprintf("Alhamdulillah, pengajuan izin santri %s berhasil disetujui!", att.StudentName),
					"data":    h.attendance[i],
				})
			} else if req.Action == "REJECT" {
				h.attendance[i].ApprovalStatus = "REJECTED"
				h.attendance[i].RejectionReason = req.RejectionReason
				h.attendance[i].Notes = fmt.Sprintf("Izin ditolak: %s. Santri dapat mengisi ulang kehadiran.", req.RejectionReason)
				return c.JSON(fiber.Map{
					"success": true,
					"message": fmt.Sprintf("Pengajuan izin %s telah ditolak. Santri dapat melakukan absen ulang.", att.StudentName),
					"data":    h.attendance[i],
				})
			}
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Data permohonan presensi tidak ditemukan.",
	})
}

// 5b. Evaluate Tahfizh Log by Musyrif (Beri Nilai & Evaluasi Tajwid)
type FiberEvaluateTahfizhLogReq struct {
	LogID        string `json:"log_id"`
	Grade        string `json:"grade"` // MUMTAZ, JAYYID_JIDDAN, JAYYID, MAQBUL
	TajweedScore string `json:"tajweed_score"`
	MakhrajScore string `json:"makhraj_score"`
	Notes        string `json:"notes"`
	TutorID      string `json:"tutor_id"`
	TutorName    string `json:"tutor_name"`
}

func (h *FiberTahap2Handler) EvaluateTahfizhLog(c *fiber.Ctx) error {
	var req FiberEvaluateTahfizhLogReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format evaluasi tahfizh tidak valid.",
		})
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, log := range h.tahfizhLogs {
		if log.ID == req.LogID {
			h.tahfizhLogs[i].FluencyScore = &req.Grade
			h.tahfizhLogs[i].TajweedScore = &req.TajweedScore
			h.tahfizhLogs[i].MakhrajScore = &req.MakhrajScore
			h.tahfizhLogs[i].Notes = &req.Notes
			h.tahfizhLogs[i].Status = "SUDAH_DINILAI"
			if req.TutorID != "" {
				h.tahfizhLogs[i].TutorID = req.TutorID
			}
			if req.TutorName != "" {
				h.tahfizhLogs[i].TutorName = req.TutorName
			}
			now := time.Now()
			h.tahfizhLogs[i].EvaluatedAt = &now

			return c.JSON(fiber.Map{
				"success": true,
				"message": fmt.Sprintf("Alhamdulillah, penilaian setoran santri %s berhasil disimpan dengan predikat %s!", log.StudentName, req.Grade),
				"data":    h.tahfizhLogs[i],
			})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"message": "Catatan setoran tahfizh tidak ditemukan.",
	})
}

func ptrFloat(f float64) *float64 {
	return &f
}
