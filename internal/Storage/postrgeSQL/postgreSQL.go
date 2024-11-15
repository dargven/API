package postrgeSQL

import (
	storage "API/internal/Storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
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
func (s *Database) Close() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}

func (s *Database) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	// Подготавливаем и выполняем запрос
	query := "INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id"
	var id int64
	err := s.Pool.QueryRow(context.Background(), query, urlToSave, alias).Scan(&id)
	if err != nil {
		// Обрабатываем ошибку уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Код ошибки уникальности в PostgreSQL
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute query: %w", op, err)
	}
	return id, nil
}

func (s *Database) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	// Выполняем запрос для получения URL
	query := "SELECT url FROM url WHERE alias = $1"
	var resURL string
	err := s.Pool.QueryRow(context.Background(), query, alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: URL not found: %w", op, err)
	} else if err != nil {
		return "", fmt.Errorf("%s: query execution failed: %w", op, err)
	}
	return resURL, nil
}
