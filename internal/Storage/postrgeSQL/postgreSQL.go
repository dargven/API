package postrgeSQL

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"API/internal/config"

	"github.com/jackc/pgx/v5/pgxpool" // pgx для подключения к PostgreSQL
)

type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase создает и возвращает подключение к базе данных
func NewDatabase(cfg *config.DataBase) (*Database, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	// Настройка пула соединений
	dbCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Настройка тайм-аутов
	dbCfg.ConnConfig.ConnectTimeout = 5 * time.Second

	// Создаем пул соединений
	pool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Проверяем соединение
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("Successfully connected to the database")
	db := &Database{Pool: pool}

	if err := db.checkUserMigration(); err != nil {
		log.Fatalf("Failed to run user migrations: %v", err)
	}

	return db, nil
}

// Close закрывает все соединения пула
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *Database) AddUser(email, name, password string) (int64, error) {
	query := `INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := db.Pool.QueryRow(context.Background(), query, email, name, password).Scan(&id)
	if err != nil {
		logger.Error("failed to create user")

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (db *Database) checkUserMigration() error { // Миграции не так пишутся
	//logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	//query :=
	//	`CREATE TABLE IF NOT EXISTS
	//	(
	//		id SERIAL PRIMARY KEY
	//		email TEXT NOT NULL UNIQUE,
	//		name TEXT NOT NULL,
	//		password TEXT NOT NULL
	//	)`
	//
	//_, err := db.Pool.Exec(context.Background(), query)
	//if err != nil {
	//	logger.Error("Failed to create migration for table users", sl.Err(err))
	//	return fmt.Errorf("failed to create migration for table users: %w", err)
	//}
	//logger.Info("Migration completed successfully (table created or already exists)") //как то разделить
	return nil
}
