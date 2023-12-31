package auth

import (
	"github.com/go-playground/validator"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/labstack/echo/v4"
)

type LogInRequest struct {
	Email string `json:"email" validate:"email,required,min=1,max=100"`
	Password string `json:"password" validate:"min=8,max=80,required"`
}

func LogIn(c echo.Context) error {
	var req LogInRequest

    if err := c.Bind(&req); err != nil {
        return c.JSON(400, map [string] string {
            "message": "Error parsing request body: " + err.Error(),
        })
    }

    validate := validator.New()
    err := validate.Struct(req)
    if err != nil {
        return c.JSON(400, map [string] string {
            "message": "Invalid request body: " + err.Error(),
        })
    }

	token, username, err := db.LogIn(req.Email, req.Password); if err != nil {
		return c.JSON(500, map [string] string {
			"message": "Failed to log-in: " + err.Error(),
		})
	}

	return c.JSON(200, map [string] string {
		"username": username,
		"token": token,
	})
}
