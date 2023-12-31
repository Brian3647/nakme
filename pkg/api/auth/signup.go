package auth

import (
	"github.com/go-playground/validator"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/labstack/echo/v4"
)

type SignUpRequest struct {
	Username string `json:"username" validate:"min=1,max=100,required"`
	Password string `json:"password" validate:"min=8,max=80,required"`
	Email string `json:"email" validate:"email,required,min=1,max=100"`
}

func SignUp(c echo.Context) error {
	var req SignUpRequest

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

	err = db.SignUp(req.Username, req.Password, req.Email); if err != nil {
		return c.JSON(500, map [string] string {
			"message": "Failed to create user: " + err.Error(),
		})
	}

	return c.JSON(200, map [string] string {
		"ok": "true",
	})
}
