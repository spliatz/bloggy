package structs

type AuthInput struct {
    UserId   *int    `json:"user_id" db:"user_id"`
    Password *string `json:"password" db:"password"`
}

type AuthResponse struct {
    Refresh
    Access
}

type Refresh struct {
    Token     string
    ExpiresAt string
}

type Access struct {
    Token     string
    ExpiresAt string
}
