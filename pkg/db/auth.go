package db

import (
	"fmt"
	"time"

	"github.com/Brian3647/nakme/pkg/util"
	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type TempUser struct {
	Username string
	Token string
	HashedPassword []byte
}

var (
	pendingAccounts map [string] TempUser
)

func GenerateToken(email string, username string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  email, // subject
		"name": username,
		"iat": time.Now().Unix(), // issued at
	})

	token, err := jwtToken.SignedString(jwtSecret)
	return token, err
}

// signs up an user and returns its JWT
func SignUp(username string, password string, email string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
        return err
    }

    if _, ok := pendingAccounts[email]; ok {
		return fmt.Errorf("email already in use")
    }

    token, err := GenerateToken(email, username); if err != nil {
	   return err
    }

	stmt, err := db.Prepare("SELECT id FROM users WHERE email = $1"); if err != nil {
		return err
	}

	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(email).Scan(&id); if err == nil {
		return fmt.Errorf("email already in use")
	}

	pendingAccounts[email] = TempUser{
		Username: username,
		Token: token,
		HashedPassword: hashedPassword,
	}

	time.AfterFunc(time.Hour, func() {
		delete(pendingAccounts, email)
	})

	return util.SendMail(email, "nakme: Confirm your account",
		fmt.Sprintf(
			"Please confirm your account by clicking <a href=\"%s/api/auth/confirm_email?token=%s&email=%s\">here</a>",
			util.GetBaseUrl(), token, email,
		) + "<br/>" + "If you didn't sign up, please ignore this email.",
	)
}

func ConfirmSignUp(email string, token string) error {
	account, ok := pendingAccounts[email]
    if !ok || account.Token != token {
        return fmt.Errorf("invalid or expired token")
    }

	delete(pendingAccounts, email)

	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id"

	stmt, err := db.Prepare(query); if err != nil {
		return err
	}

	defer stmt.Close()

	var id uint
	err = stmt.QueryRow(account.Username, account.HashedPassword, email).Scan(&id); if err != nil {
		return err
	}

	return nil
}

// logs in an user and returns its JWT & usernameemail
func LogIn(email string, password string) (jwtoken string, uname string, e error) {
	query := "SELECT password, id, username FROM users WHERE email = $1"
	stmt, err := db.Prepare(query); if err != nil {
		return "", "", err
	}

	defer stmt.Close()

	var hashedPassword string
	var id uint
	var username string
	err = stmt.QueryRow(email).Scan(&hashedPassword, &id, &username); if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); if err != nil {
		return "", "", err
	}

	jwt, err := GenerateToken(email, username)

	return jwt, username, err
}

func ConfirmIdentity(token string) error {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })

    if err != nil {
        return err
    }

    if claims, ok := parsed.Claims.(jwt.MapClaims); ok && claims.Valid() == nil {
        return nil
    } else {
        return fmt.Errorf("invalid token")
    }
}

func SetUpAuth() {
	pendingAccounts = make(map[string] TempUser)
}
