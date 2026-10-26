package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"pkbm-backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	GormDB *gorm.DB
	SqlDB  *sql.DB
}

func Connect(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Printf("Catatan GORM: MySQL di port %s belum online / database '%s' belum aktif: %v", cfg.DBPort, cfg.DBName, err)
		log.Println("Server tetap aktif dengan mode fallback katalog 18 peran demo.")
		return &Database{GormDB: nil, SqlDB: nil}, nil
	}

	sqlDB, err := gormDB.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	log.Println("Alhamdulillah! Berhasil terhubung ke database MySQL (GORM) PKBM Digital Dar El Iman.")

	// Eksekusi Auto-Migrate skema dan master data saat startup / deploy
	_ = RunAutoMigrate(gormDB)

	return &Database{GormDB: gormDB, SqlDB: sqlDB}, nil
}
