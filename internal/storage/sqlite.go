package storage

import (
	"database/sql"
	"fmt"
	"homelab-inventory/pkg/model"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage initializes the DB and ensures tables exist
func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	dataSource := fmt.Sprintf("file:%s?cache=shared&mode=rwc", path)
	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB: %w", err)
	}

	s := &SQLiteStorage{db: db}
	if err := s.ensureSchema(); err != nil {
		return nil, err
	}

	return s, nil
}

// ensureSchema creates the necessary tables if they do not exist
func (s *SQLiteStorage) ensureSchema() error {
	const schema = `
    CREATE TABLE IF NOT EXISTS system_info (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hostname TEXT,
        os TEXT,
        platform TEXT,
        arch TEXT,
        cpu_model TEXT,
        cpu_cores INTEGER,
        memory_gb REAL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS system_disk (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        system_info_id INTEGER,
        mountpoint TEXT,
        total_gb REAL,
        used_gb REAL,
        free_gb REAL,
        FOREIGN KEY(system_info_id) REFERENCES system_info(id) ON DELETE CASCADE
    );
    `
	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

// SaveSystemInfo stores the system info and its related disk info in a transaction
func (s *SQLiteStorage) SaveSystemInfo(info *model.SystemInfo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	systemID, err := insertSystemInfo(tx, info)
	if err != nil {
		return err
	}

	if err := insertDisks(tx, systemID, info.Disks); err != nil {
		return err
	}

	return tx.Commit()
}

// insertSystemInfo inserts into system_info and returns the new ID
func insertSystemInfo(tx *sql.Tx, info *model.SystemInfo) (int64, error) {
	result, err := tx.Exec(`
        INSERT INTO system_info (hostname, os, platform, arch, cpu_model, cpu_cores, memory_gb)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		info.Hostname, info.OS, info.Platform, info.Arch,
		info.CPUModel, info.CPUCores, info.MemoryGB,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert system_info: %w", err)
	}

	return result.LastInsertId()
}

// insertDisks inserts related disks
func insertDisks(tx *sql.Tx, systemID int64, disks []model.DiskInfo) error {
	for _, d := range disks {
		_, err := tx.Exec(`
            INSERT INTO system_disk (system_info_id, mountpoint, total_gb, used_gb, free_gb)
            VALUES (?, ?, ?, ?, ?)`,
			systemID, d.Mountpoint, d.TotalGB, d.UsedGB, d.FreeGB,
		)
		if err != nil {
			return fmt.Errorf("failed to insert disk info: %w", err)
		}
	}
	return nil
}
