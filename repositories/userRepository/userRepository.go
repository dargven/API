package userRepository

import (
	"API/internal/Storage/postrgeSQL"
	resp "API/internal/lib/api/response"
	"API/internal/models/user"
	"context"
	"errors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type UserRepository struct {
	db *postrgeSQL.Database
}

func NewUserRepository(db *postrgeSQL.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (h *UserRepository) NewUser(w http.ResponseWriter, r *http.Request) {
	var userRequest user.CreateUserRequest
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	err := render.DecodeJSON(r.Body, &userRequest)
	if errors.Is(err, io.EOF) {
		// Если тело запроса пусто
		logger.Error("body is empty")
		render.JSON(w, r, resp.Error("body is empty"))
		return
	}
	if err != nil {
		logger.Error("failed to decode JSON", "error", err)
		render.JSON(w, r, resp.Error("invalid request body"))
		return
	}

	const query = "INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id"
	ctx := context.Background()

	var userID int64
	err = h.db.Pool.QueryRow(ctx, query, userRequest.Email, userRequest.Name, userRequest.Password).Scan(&userID)
	if err != nil {
		logger.Error("failed to create user", "error", err)
		render.JSON(w, r, resp.Error("failed to create user"))
		return
	}

	logger.Info("user created successfully", "userID", userID)
	render.JSON(w, r, resp.Success(map[string]interface{}{
		"user_id": userID,
		"message": "User created successfully",
	}))
}
func (h *UserRepository) IsEmailUnique(email string) (bool, error) {
	if email == "" {
		log.Println("[DEBUG] Email is empty")
		return false, errors.New("email cannot be empty")
	}

	const query = "SELECT COUNT(*) FROM users WHERE email = $1"
	log.Printf("[DEBUG] Executing query: %s with email: %s", query, email)

	ctx := context.Background()
	rows, err := h.db.Pool.Query(ctx, query, email)
	if err != nil {
		log.Printf("[ERROR] Failed to execute query: %v", err)
		return false, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Printf("[ERROR] Failed to scan row: %v", err)
			return false, err
		}
	}

	log.Printf("[DEBUG] Query result - count: %d", count)

	isUnique := count == 0
	log.Printf("[DEBUG] Email is unique: %v", isUnique)

	return isUnique, nil
}
