# PKBM Digital Dar El Iman - Backend Core API

Sistem Backend Core berbasis **Go (Golang)** menggunakan framework **Fiber v2** dan **GORM** untuk platform pembelajaran digital PKBM Dar El Iman Padang (Kurikulum Merdeka Kesetaraan Paket A, B, C & Halaqah Tahfizh).

---

## 🚀 Fitur Utama
- **Go Fiber v2 & GORM**: Kinerja tinggi, arsitektur modular, latency rendah.
- **Auto-Migrate on Deploy**: Skema database (Tahap 1 s.d. Tahap 5) dan seed master data (18 Roles, akun demo default, tahun ajaran, kurikulum) otomatis dimigrasikan ke database saat server startup.
- **RBAC 18 Peran**: Hak akses granular untuk Santri, Tutor/Ustadz, Wali Kelas, Wali Santri, dan Manajemen Yayasan.
- **Autentikasi JWT**: Manajemen token login, register PPDB, dan ganti password.
- **Graceful Fallback**: Tetap berjalan dengan katalog akun mock jika koneksi database MySQL belum aktif di lokal.

---

## 🛠️ Konfigurasi Environment (`.env`)
Salin file `.env.example` menjadi `.env`:
```env
PORT=8080
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=pkbm_digital_darel_iman
JWT_SECRET=dar-el-iman-secret-key-2027-2032
FRONTEND_URL=http://localhost:5173
AUTO_MIGRATE=true
```

---

## 💻 Menjalankan Secara Lokal

```bash
# Unduh dependensi
go mod download

# Jalankan server
go run cmd/server/main.go
```
Server akan aktif di `http://localhost:8080`.

---

## 🐳 Deployment (Docker / Easypanel)

Repository ini telah dilengkapi dengan `Dockerfile` siap pakai:
```bash
docker build -t pkbm-backend:latest .
docker run -p 8080:8080 --env-file .env pkbm-backend:latest
```
Ketika container dijalankan, sistem secara otomatis mengeksekusi migrasi skema database (`AUTO_MIGRATE=true`).
