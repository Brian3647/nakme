package auth

import (
	"net/http"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/Brian3647/nakme/pkg/util"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type DeleteUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func DeleteUser(c echo.Context) error {
	token, err := util.ExpectAuth(c); if err != nil {
		return c.JSON(http.StatusUnauthorized, map [string] string {
			"message": "Invalid auth: " + err.Error(),
		})
	}

	var req DeleteUserRequest

    if err := c.Bind(&req); err != nil {
        return c.JSON(400, map [string] string {
            "message": "Error parsing request body: " + err.Error(),
        })
    }

    validate := validator.New()
    err = validate.Struct(req)
    if err != nil {
        return c.JSON(400, map [string] string {
            "message": "Invalid request body: " + err.Error(),
        })
    }

	err = db.DeleteAccount(req.Email, token); if err != nil {
		return c.JSON(400, map [string] string {
			"message": "Failed to delete account: " + err.Error(),
		})
	}

	return c.JSON(200, map [string] string {
		"ok": "true",
	})
}
