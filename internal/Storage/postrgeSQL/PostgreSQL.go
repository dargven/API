package postrgeSQL

import (
	"context"
	"fmt"
	"log"
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
	return &Database{Pool: pool}, nil
}

// Close закрывает все соединения пула
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
