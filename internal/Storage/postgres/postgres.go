package postgres

import (
	storage "API/internal/Storage"
	"API/internal/models"
	"context"
	"crypto/rand"
	"encoding/hex"
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

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: connect: %w", op, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: ping: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

// Close закрывает пул соединений
func (s *Storage) Close() {
	s.pool.Close()
}

// ==================== URL Methods ====================

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
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
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

// ==================== User Methods ====================

// CreateUser создает нового пользователя
func (s *Storage) CreateUser(email, name, passwordHash string) (*models.User, error) {
	const op = "storage.postgres.CreateUser"

	var user models.User
	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO users(email, name, password_hash, balance, created_at, updated_at) 
		 VALUES($1, $2, $3, 0, $4, $4) 
		 RETURNING id, email, name, password_hash, phone, avatar_url, bio, balance, created_at, updated_at`,
		email, name, passwordHash, time.Now(),
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Phone, &user.AvatarURL, &user.Bio, &user.Balance,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
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
		`SELECT id, email, name, password_hash, phone, avatar_url, bio, balance, created_at, updated_at 
		 FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Phone, &user.AvatarURL, &user.Bio, &user.Balance,
		&user.CreatedAt, &user.UpdatedAt,
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
		`SELECT id, email, name, password_hash, phone, avatar_url, bio, balance, created_at, updated_at 
		 FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Phone, &user.AvatarURL, &user.Bio, &user.Balance,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (s *Storage) UpdateUserProfile(userID int64, name string, phone, avatarURL, bio *string) (*models.User, error) {
	const op = "storage.postgres.UpdateUserProfile"

	var user models.User
	err := s.pool.QueryRow(
		context.Background(),
		`UPDATE users SET name = $1, phone = $2, avatar_url = $3, bio = $4, updated_at = $5
		 WHERE id = $6
		 RETURNING id, email, name, password_hash, phone, avatar_url, bio, balance, created_at, updated_at`,
		name, phone, avatarURL, bio, time.Now(), userID,
	).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash,
		&user.Phone, &user.AvatarURL, &user.Bio, &user.Balance,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// UpdateUserBalance обновляет баланс пользователя
func (s *Storage) UpdateUserBalance(userID int64, amount float64) error {
	const op = "storage.postgres.UpdateUserBalance"

	result, err := s.pool.Exec(
		context.Background(),
		`UPDATE users SET balance = balance + $1, updated_at = $2 WHERE id = $3`,
		amount, time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}

// EmailExists проверяет существование email
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

// ==================== Event Methods ====================

// CreateEvent создает новое мероприятие
func (s *Storage) CreateEvent(event *models.Event) (*models.Event, error) {
	const op = "storage.postgres.CreateEvent"

	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO events(title, description, category, image_url, venue, address, price, capacity, available_tickets, start_time, end_time, creator_id, created_at, updated_at) 
		 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $8, $9, $10, $11, $12, $12) 
		 RETURNING id, title, description, category, image_url, venue, address, price, capacity, available_tickets, start_time, end_time, creator_id, created_at, updated_at`,
		event.Title, event.Description, event.Category, event.ImageURL, event.Venue, event.Address,
		event.Price, event.Capacity, event.StartTime, event.EndTime, event.CreatorID, time.Now(),
	).Scan(
		&event.ID, &event.Title, &event.Description, &event.Category, &event.ImageURL,
		&event.Venue, &event.Address, &event.Price, &event.Capacity, &event.AvailableTickets,
		&event.StartTime, &event.EndTime, &event.CreatorID, &event.CreatedAt, &event.UpdatedAt,
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
		`SELECT id, title, description, category, image_url, venue, address, price, capacity, available_tickets, start_time, end_time, creator_id, created_at, updated_at 
		 FROM events WHERE id = $1`,
		id,
	).Scan(
		&event.ID, &event.Title, &event.Description, &event.Category, &event.ImageURL,
		&event.Venue, &event.Address, &event.Price, &event.Capacity, &event.AvailableTickets,
		&event.StartTime, &event.EndTime, &event.CreatorID, &event.CreatedAt, &event.UpdatedAt,
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
		`SELECT id, title, description, category, image_url, venue, address, price, capacity, available_tickets, start_time, end_time, creator_id, created_at, updated_at 
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
			&event.ID, &event.Title, &event.Description, &event.Category, &event.ImageURL,
			&event.Venue, &event.Address, &event.Price, &event.Capacity, &event.AvailableTickets,
			&event.StartTime, &event.EndTime, &event.CreatorID, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

// SearchEvents выполняет полнотекстовый поиск мероприятий
func (s *Storage) SearchEvents(query string, category string, dateFrom, dateTo *time.Time, priceMin, priceMax *float64, limit, offset int) ([]*models.Event, int, error) {
	const op = "storage.postgres.SearchEvents"

	// Базовый запрос с условиями
	baseQuery := `
		FROM events 
		WHERE 1=1
	`
	args := []interface{}{}
	argNum := 1

	// Полнотекстовый поиск
	if query != "" {
		baseQuery += fmt.Sprintf(` AND (
			to_tsvector('russian', title || ' ' || COALESCE(description, '') || ' ' || venue || ' ' || COALESCE(address, '')) 
			@@ plainto_tsquery('russian', $%d)
			OR title ILIKE $%d
			OR venue ILIKE $%d
		)`, argNum, argNum+1, argNum+2)
		args = append(args, query, "%"+query+"%", "%"+query+"%")
		argNum += 3
	}

	// Фильтр по категории
	if category != "" {
		baseQuery += fmt.Sprintf(` AND category = $%d`, argNum)
		args = append(args, category)
		argNum++
	}

	// Фильтр по датам
	if dateFrom != nil {
		baseQuery += fmt.Sprintf(` AND start_time >= $%d`, argNum)
		args = append(args, *dateFrom)
		argNum++
	}
	if dateTo != nil {
		baseQuery += fmt.Sprintf(` AND start_time <= $%d`, argNum)
		args = append(args, *dateTo)
		argNum++
	}

	// Фильтр по цене
	if priceMin != nil {
		baseQuery += fmt.Sprintf(` AND price >= $%d`, argNum)
		args = append(args, *priceMin)
		argNum++
	}
	if priceMax != nil {
		baseQuery += fmt.Sprintf(` AND price <= $%d`, argNum)
		args = append(args, *priceMax)
		argNum++
	}

	// Получаем общее количество
	var total int
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := s.pool.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: count: %w", op, err)
	}

	// Получаем записи с пагинацией
	selectQuery := `SELECT id, title, description, category, image_url, venue, address, price, capacity, available_tickets, start_time, end_time, creator_id, created_at, updated_at ` + baseQuery
	selectQuery += fmt.Sprintf(` ORDER BY start_time ASC LIMIT $%d OFFSET $%d`, argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := s.pool.Query(context.Background(), selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.Category, &event.ImageURL,
			&event.Venue, &event.Address, &event.Price, &event.Capacity, &event.AvailableTickets,
			&event.StartTime, &event.EndTime, &event.CreatorID, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: scan: %w", op, err)
		}
		events = append(events, &event)
	}

	return events, total, rows.Err()
}

// ==================== Booking Methods ====================

// generateBookingCode генерирует уникальный код бронирования
func generateBookingCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "BK-" + hex.EncodeToString(bytes)
}

// CreateBooking создает бронирование с транзакцией
func (s *Storage) CreateBooking(userID, eventID int64, quantity int) (*models.Booking, error) {
	const op = "storage.postgres.CreateBooking"

	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Проверяем доступность билетов и получаем цену
	var availableTickets int
	var price float64
	err = tx.QueryRow(ctx,
		`SELECT available_tickets, price FROM events WHERE id = $1 FOR UPDATE`,
		eventID,
	).Scan(&availableTickets, &price)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: get event: %w", op, err)
	}

	if availableTickets < quantity {
		return nil, storage.ErrNoTickets
	}

	totalPrice := price * float64(quantity)

	// Проверяем баланс пользователя
	var balance float64
	err = tx.QueryRow(ctx, `SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: get balance: %w", op, err)
	}

	if balance < totalPrice {
		return nil, storage.ErrInsufficientBalance
	}

	// Списываем с баланса
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE id = $2`, totalPrice, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: deduct balance: %w", op, err)
	}

	// Уменьшаем количество доступных билетов
	_, err = tx.Exec(ctx, `UPDATE events SET available_tickets = available_tickets - $1 WHERE id = $2`, quantity, eventID)
	if err != nil {
		return nil, fmt.Errorf("%s: update tickets: %w", op, err)
	}

	// Создаем бронирование
	bookingCode := generateBookingCode()
	var booking models.Booking
	err = tx.QueryRow(ctx,
		`INSERT INTO bookings(user_id, event_id, quantity, total_price, status, booking_code, created_at) 
		 VALUES($1, $2, $3, $4, $5, $6, $7) 
		 RETURNING id, user_id, event_id, quantity, total_price, status, booking_code, created_at`,
		userID, eventID, quantity, totalPrice, models.BookingStatusConfirmed, bookingCode, time.Now(),
	).Scan(
		&booking.ID, &booking.UserID, &booking.EventID, &booking.Quantity,
		&booking.TotalPrice, &booking.Status, &booking.BookingCode, &booking.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, storage.ErrBookingExists
		}
		return nil, fmt.Errorf("%s: insert booking: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: commit: %w", op, err)
	}

	return &booking, nil
}

// GetBookingsByUserID возвращает все бронирования пользователя
func (s *Storage) GetBookingsByUserID(userID int64) ([]*models.BookingWithEvent, error) {
	const op = "storage.postgres.GetBookingsByUserID"

	rows, err := s.pool.Query(
		context.Background(),
		`SELECT b.id, b.user_id, b.event_id, b.quantity, b.total_price, b.status, b.booking_code, b.created_at,
		        e.title, e.start_time, e.venue
		 FROM bookings b
		 JOIN events e ON b.event_id = e.id
		 WHERE b.user_id = $1
		 ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bookings []*models.BookingWithEvent
	for rows.Next() {
		var b models.BookingWithEvent
		err := rows.Scan(
			&b.ID, &b.UserID, &b.EventID, &b.Quantity, &b.TotalPrice, &b.Status, &b.BookingCode, &b.CreatedAt,
			&b.EventTitle, &b.EventDate, &b.Venue,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		bookings = append(bookings, &b)
	}

	return bookings, rows.Err()
}

// GetBookingByID возвращает бронирование по ID
func (s *Storage) GetBookingByID(bookingID, userID int64) (*models.BookingWithEvent, error) {
	const op = "storage.postgres.GetBookingByID"

	var b models.BookingWithEvent
	err := s.pool.QueryRow(
		context.Background(),
		`SELECT b.id, b.user_id, b.event_id, b.quantity, b.total_price, b.status, b.booking_code, b.created_at,
		        e.title, e.start_time, e.venue
		 FROM bookings b
		 JOIN events e ON b.event_id = e.id
		 WHERE b.id = $1 AND b.user_id = $2`,
		bookingID, userID,
	).Scan(
		&b.ID, &b.UserID, &b.EventID, &b.Quantity, &b.TotalPrice, &b.Status, &b.BookingCode, &b.CreatedAt,
		&b.EventTitle, &b.EventDate, &b.Venue,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrBookingNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &b, nil
}

// CancelBooking отменяет бронирование и возвращает деньги
func (s *Storage) CancelBooking(bookingID, userID int64) error {
	const op = "storage.postgres.CancelBooking"

	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Получаем бронирование
	var eventID int64
	var quantity int
	var totalPrice float64
	var status models.BookingStatus
	err = tx.QueryRow(ctx,
		`SELECT event_id, quantity, total_price, status FROM bookings WHERE id = $1 AND user_id = $2 FOR UPDATE`,
		bookingID, userID,
	).Scan(&eventID, &quantity, &totalPrice, &status)

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrBookingNotFound
	}
	if err != nil {
		return fmt.Errorf("%s: get booking: %w", op, err)
	}

	if status == models.BookingStatusCancelled {
		return errors.New("booking already cancelled")
	}

	// Обновляем статус бронирования
	_, err = tx.Exec(ctx, `UPDATE bookings SET status = $1 WHERE id = $2`, models.BookingStatusCancelled, bookingID)
	if err != nil {
		return fmt.Errorf("%s: update status: %w", op, err)
	}

	// Возвращаем билеты
	_, err = tx.Exec(ctx, `UPDATE events SET available_tickets = available_tickets + $1 WHERE id = $2`, quantity, eventID)
	if err != nil {
		return fmt.Errorf("%s: return tickets: %w", op, err)
	}

	// Возвращаем деньги
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, totalPrice, userID)
	if err != nil {
		return fmt.Errorf("%s: refund: %w", op, err)
	}

	return tx.Commit(ctx)
}
