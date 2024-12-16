package userRepository

import (
	"API/internal/Storage/postrgeSQL"
	"API/internal/models/user"
	"context"
	"errors"
	"fmt"
)

// UserRepository управляет взаимодействием с таблицей пользователей
type UserRepository struct {
	db *postrgeSQL.Database
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *postrgeSQL.Database) *UserRepository {
	return &UserRepository{db: db}
}

// NewUser добавляет нового пользователя в базу данных
func (r *UserRepository) NewUser(ctx context.Context, req user.CreateUserRequest) (*user.User, error) {
	const query = `
		INSERT INTO users (email, name, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, email, name
	`

	newUser := &user.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	err := r.db.Pool.QueryRow(ctx, query, req.Email, req.Name, req.Password).
		Scan(&newUser.ID, &newUser.Email, &newUser.Name)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil, errors.New("user with this email already exists")
		}
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return newUser, nil
}

// IsEmailUnique проверяет, уникален ли email
func (r *UserRepository) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	const query = "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email uniqueness: %w", err)
	}

	return count == 0, nil
}

// GetUserByID возвращает пользователя по его ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*user.User, error) {
	const query = "SELECT id, email, name FROM users WHERE id = $1"

	var u user.User
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&u.ID, &u.Email, &u.Name)
	if err != nil {
		if isNotFoundError(err) {
			return nil, errors.New("u not found")
		}
		return nil, fmt.Errorf("failed to get u by ID: %w", err)
	}

	return &u, nil
}

// DeleteUser удаляет пользователя по его ID
func (r *UserRepository) DeleteUser(ctx context.Context, userID uint) error {
	const query = "DELETE FROM users WHERE id = $1"

	result, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	var u user.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

// isUniqueViolationError проверяет, является ли ошибка нарушением уникальности
func isUniqueViolationError(err error) bool {
	return err != nil && err.Error() == "unique_violation"
}

// isNotFoundError проверяет, является ли ошибка отсутствием данных
func isNotFoundError(err error) bool {
	return err != nil && err.Error() == "no rows in result set"
}
