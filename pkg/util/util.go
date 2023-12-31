package util

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

// store the value of `Prod()` in a variable so that it only needs to be
// calculated once. (even though getting an env var is a very fast operation)
var (
	once sync.Once
	isProd bool

	senderEmail string
	senderPassword string
	smtpHost string
	smtpPort int

	baseUrl string
)

// InProd returns true if the environment is in production mode. (PROD=true | PROD=1)
func InProd() bool {
	once.Do(func() {
		// convert to lowercase to match `TRUE` and `true`
		prod := strings.ToLower(os.Getenv("PROD"))
		alt := strings.ToLower(os.Getenv("PRODUCTION"))
		isProd = prod == "1" || prod == "true" || alt == "1" || alt == "true"
	})

	return isProd
}

func expectEnvVar(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic("Missing environment variable: " + name)
	}
	return env
}

func SetUp() {
	senderEmail = expectEnvVar("EMAIL_SENDER")
	senderPassword = expectEnvVar("EMAIL_PASSWORD")
	smtpHost = expectEnvVar("EMAIL_HOST")

	portStr := expectEnvVar("EMAIL_PORT")

	var err error
    smtpPort, err = strconv.Atoi(portStr)
    if err != nil {
        panic("Invalid EMAIL_PORT: " + portStr)
    }

	if InProd() {
		baseUrl = "https://nakme.dev"
	} else {
		baseUrl = "http://localhost:3000"
	}
}

// Creates a file of format `logs/YYY-MM-DD.log` and returns a pointer to it.
func GenerateLogFile() *os.File {
	// create the logs directory if it doesn't exist
	os.Mkdir("logs", 0755)

	// create the file
	logFile, err := os.OpenFile("logs/" + GetDate() + ".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	return logFile
}

// GetDate returns the current date in the format `YYYY-MM-DD`.
func GetDate() string {
	return time.Now().Format("2006-01-02")
}

func createEmail(senderEmail, senderPassword, smtpHost string, smtpPort int, recipientEmail, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", senderEmail)
	message.SetHeader("To", recipientEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	err := dialer.DialAndSend(message)
	return err
}

func SendMail(to string, subject string, body string) error {
	return createEmail(senderEmail, senderPassword, smtpHost, smtpPort, to, subject, body)
}

func GetBaseUrl() string {
	return baseUrl
}

func ExpectAuth(c echo.Context) (string, error) {
    token := c.Request().Header.Get("Authorization")
    if token == "" {
        return "", errors.New("missing Authorization header")
    }

    if !strings.HasPrefix(token, "Bearer ") {
        return "", errors.New("invalid Authorization header format (expected Bearer)")
    }

    return strings.TrimSpace(strings.Replace(token, "Bearer ", "", 1)), nil
}
