package auth

import (
	"net/http"

	"github.com/Brian3647/nakme/pkg/db"
	"github.com/labstack/echo/v4"
)

func ConfirmEmail(c echo.Context) error {
	token := c.QueryParam("token")
	email := c.QueryParam("email")

	if token == "" || email == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"message": "Missing token or email in query params",
		})
	}

	err := db.ConfirmIdentity(token)
	if err != nil {
		return c.JSON(400, map[string]string{
			"message": "Failed to confirm token's identity: " + err.Error(),
		})
	}

	err = db.ConfirmSignUp(email, token); if err != nil {
		return c.JSON(400, map[string]string{
			"message": "Failed to confirm email: " + err.Error(),
		})
	}

	return c.Redirect(http.StatusFound, "/")
}
