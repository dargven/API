package user

// здесь же описан сам пользак и функции по его добавлению с проверками и отправкой всего в бд

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null" json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
}
type Response struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
