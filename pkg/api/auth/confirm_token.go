package auth

import (
	"net/http"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/Brian3647/nakme/pkg/util"
	"github.com/labstack/echo/v4"
)

func ConfirmToken(c echo.Context) error {
	token, err := util.ExpectAuth(c); if err != nil {
		return c.JSON(http.StatusUnauthorized, map [string] string {
			"message": "Failed auth: " + err.Error(),
		})
	}

	err = db.ConfirmIdentity(token); if err != nil {
		return c.JSON(400, map [string] string {
			"message": "Failed confirm identity: " + err.Error(),
		})
	}

	return c.JSON(200, map [string] string {
		"ok": "true",
	})
}
