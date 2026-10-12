package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"pkbm-backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("Peringatan: Belum dapat terhubung ke MySQL (%s). Pastikan MySQL aktif.", err)
		return db, nil
	}

	log.Println("Berhasil terhubung ke database MySQL PKBM Digital Dar El Iman.")
	return db, nil
}
