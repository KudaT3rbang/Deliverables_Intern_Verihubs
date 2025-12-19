package entity

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type RegisterUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
