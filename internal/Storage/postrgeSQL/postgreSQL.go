package postrgesql

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
