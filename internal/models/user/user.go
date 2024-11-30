package user

// здесь же описан сам пользак и функции по его добавлению с проверками и отправкой всего в бд

type User struct {
	ID       int64
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type CreateUserRequest struct {
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
}
type UserResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"` // Возможно придется убрать json
}

// Перенсти в файл userRepository.

//if err != nil {
//	//ошибка при попытке декода
//	logger.Error("failed to decode body", sl.Err(err))
//
//	render.JSON(w, r, resp.Error("failed to decode"))
//
//	return
//}

//if err := validator.New().Struct(user); err != nil {
//	var validateError validator.ValidationErrors
//	errors.As(err, &validateError)
//
//	logger.Error(validateError.Error(), sl.Err(err))
//
//	render.JSON(w, r, resp.Error("invalid request"))
//} Вынести в отдельный обработчик handler

//	if !isEmailValid(user.Email) {
//		logger.Error("email does not match the format")
//
//		render.JSON(w, r, resp.Error("email does not match the format"))
//
//		return
//	}
//
//	if h.isUnique(user.Email) {
//		logger.Error("email already exist")
//
//		render.JSON(w, r, resp.Error("email already exist"))
//
//		return
//	}
//
//	id, err := h.DB.AddUser(user.Email, user.Name, user.Password)
//	if err != nil {
//		logger.Error("failed to create user")
//
//		render.JSON(w, r, resp.Error("failed to create user"))
//
//		return
//	}
//
//	render.JSON(w, r, id)
//
//}
//
//func (h *UserHandler) isUnique(email string) bool {
//	query := `SELECT id FROM users WHERE email = $1`
//
//	var id int64
//
//	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
//
//	err := h.DB.Pool.QueryRow(context.Background(), query, email).Scan(&id)
//	if err != nil {
//		logger.Error("Error checking email uniqueness:", err)
//		return false
//	}
//	if errors.Is(err, sql.ErrNoRows) {
//		return false
//	}
//
//	return true
//}
//
//func isEmailValid(e string) bool {
//	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
//	return emailRegex.MatchString(e)
//}
