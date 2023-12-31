package auth

import (
	"github.com/go-playground/validator"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/labstack/echo/v4"
)

type ConfirmTokenRequest struct {
	Token string `json:"token"`
}

func ConfirmToken(c echo.Context) error {
	var req ConfirmTokenRequest

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

	err = db.ConfirmIdentity(req.Token); if err != nil {
		return c.JSON(400, map [string] string {
			"message": "Failed confirm identity: " + err.Error(),
		})
	}

	return c.JSON(200, map [string] string {
		"ok": "true",
	})
}
