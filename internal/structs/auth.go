package structs

type SignUpRequest struct {
    Username string `json:"username" db:"username" binding:"required"`
    Password string `json:"password" db:"password" binding:"required"`

    Name     *string `json:"name"`
    Birthday *string `json:"birthday"` // Use "2000-12-31" format
    Email    *string `json:"email"`
    Phone    *string `json:"phone"`
}

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
