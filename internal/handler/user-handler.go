package handler

import (
    "github.com/Intellect-Bloggy/bloggy-backend/internal/services"
)

type UserHandler struct {
    userService services.User
}

func newUserHandler(userService services.User) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}
