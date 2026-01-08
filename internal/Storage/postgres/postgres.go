package postgres

import (
	storage "API/internal/Storage"
	"API/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage представляет PostgreSQL хранилище
type Storage struct {
	pool *pgxpool.Pool
}

// New создает новое подключение к PostgreSQL
func New(host string, port int, user, password, dbname, sslmode string) (*Storage, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: parse config: %w", op, err)
	}

	// Настройка пула соединений
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	// Проверка соединения
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

// Close закрывает пул соединений
func (s *Storage) Close() {
	s.pool.Close()
}

// SaveURL сохраняет URL с алиасом
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64
	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id`,
		urlToSave, alias,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetURL возвращает URL по алиасу
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var resURL string
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT url FROM url WHERE alias = $1`,
		alias,
	).Scan(&resURL)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resURL, nil
}

// CreateUser создает нового пользователя
func (s *Storage) CreateUser(email, name, passwordHash string) (*models.User, error) {
	const op = "storage.postgres.CreateUser"

	var user models.User
	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO users(email, name, password_hash, created_at, updated_at) 
		 VALUES($1, $2, $3, $4, $4) 
		 RETURNING id, email, name, password_hash, created_at, updated_at`,
		email, name, passwordHash, time.Now(),
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, storage.ErrUserExists
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// GetUserByEmail возвращает пользователя по email
func (s *Storage) GetUserByEmail(email string) (*models.User, error) {
	const op = "storage.postgres.GetUserByEmail"

	var user models.User
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT id, email, name, password_hash, created_at, updated_at 
		 FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// GetUserByID возвращает пользователя по ID
func (s *Storage) GetUserByID(id int64) (*models.User, error) {
	const op = "storage.postgres.GetUserByID"

	var user models.User
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT id, email, name, password_hash, created_at, updated_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// EmailExists проверяет существование пользователя с данным email
func (s *Storage) EmailExists(email string) (bool, error) {
	const op = "storage.postgres.EmailExists"

	var exists bool
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

// CreateEvent создает новое мероприятие
func (s *Storage) CreateEvent(event *models.Event) (*models.Event, error) {
	const op = "storage.postgres.CreateEvent"

	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO events(title, description, location, start_time, end_time, creator_id, max_slots, created_at, updated_at) 
		 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $8) 
		 RETURNING id, title, description, location, start_time, end_time, creator_id, max_slots, created_at, updated_at`,
		event.Title, event.Description, event.Location, event.StartTime, event.EndTime, event.CreatorID, event.MaxSlots, time.Now(),
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.CreatorID,
		&event.MaxSlots,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return event, nil
}

// GetEventByID возвращает мероприятие по ID
func (s *Storage) GetEventByID(id int64) (*models.Event, error) {
	const op = "storage.postgres.GetEventByID"

	var event models.Event
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT id, title, description, location, start_time, end_time, creator_id, max_slots, created_at, updated_at 
		 FROM events WHERE id = $1`,
		id,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.CreatorID,
		&event.MaxSlots,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &event, nil
}

// GetAllEvents возвращает все мероприятия с пагинацией
func (s *Storage) GetAllEvents(limit, offset int) ([]*models.Event, error) {
	const op = "storage.postgres.GetAllEvents"

	rows, err := s.pool.Query(
		context.Background(),
		`SELECT id, title, description, location, start_time, end_time, creator_id, max_slots, created_at, updated_at 
		 FROM events 
		 ORDER BY start_time ASC 
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.CreatorID,
			&event.MaxSlots,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return events, nil
}

// GetEventsByCreator возвращает мероприятия созданные пользователем
func (s *Storage) GetEventsByCreator(creatorID int64) ([]*models.Event, error) {
	const op = "storage.postgres.GetEventsByCreator"

	rows, err := s.pool.Query(
		context.Background(),
		`SELECT id, title, description, location, start_time, end_time, creator_id, max_slots, created_at, updated_at 
		 FROM events 
		 WHERE creator_id = $1
		 ORDER BY start_time ASC`,
		creatorID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.CreatorID,
			&event.MaxSlots,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return events, nil
}
