package entities

type User struct {
	ID           int    `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"password"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

type UserCreateRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	ID int `json:"id" db:"id"`
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	ID int `json:"id" db:"id"`
}

type UserGet struct {
	ID        int    `json:"id" db:"id"`
	Login     string `json:"login" db:"login"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type UserUpdateLoginRequest struct {
	ID    int    `json:"id" db:"id"`
	Login string `json:"login" db:"login"`
}

type UserUpdateLoginResponse struct {
	ID int `json:"id" db:"id"`
}

type UserUpdatePasswordRequest struct {
	ID       int    `json:"id" db:"id"`
	Password string `json:"password" db:"password"`
}

type UserUpdatePasswordResponse struct {
	ID int `json:"id" db:"id"`
}
