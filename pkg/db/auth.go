package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
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
	PendingAccounts map [string] TempUser
	PendingAccountsMutex sync.RWMutex
)

var (
    selectUserByEmailStmt *sql.Stmt
    insertUserStmt        *sql.Stmt
    deleteUserByEmailStmt *sql.Stmt
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

    if _, ok := PendingAccounts[email]; ok {
		return fmt.Errorf("email already in use")
    }

    token, err := GenerateToken(email, username); if err != nil {
	   return err
    }

	var id uint
	err = selectUserByEmailStmt.QueryRow(email).Scan(&id); if err == nil {
		return fmt.Errorf("email already in use")
	}

	PendingAccountsMutex.Lock()
	PendingAccounts[email] = TempUser{
		Username: username,
		Token: token,
		HashedPassword: hashedPassword,
	}
	PendingAccountsMutex.Unlock()

	time.AfterFunc(time.Hour, func() {
		delete(PendingAccounts, email)
	})

	return util.SendMail(email, "nakme: Confirm your account",
		fmt.Sprintf(
			"Please confirm your account by clicking <a href=\"%s/api/auth/confirm_email?token=%s&email=%s\">here</a>",
			util.GetBaseUrl(), token, email,
		) + "<br/>" + "If you didn't sign up, please ignore this email.",
	)
}

func ConfirmSignUp(email string, token string) error {
	PendingAccountsMutex.RLock()
    account, ok := PendingAccounts[email]
    PendingAccountsMutex.RUnlock()

    if !ok || account.Token != token {
        return fmt.Errorf("invalid or expired token")
    }

	delete(PendingAccounts, email)

	var id string
	err := insertUserStmt.QueryRow(account.Username, account.HashedPassword, email).Scan(&id); if err != nil {
		return err
	}

	return nil
}

// logs in an user and returns its JWT & username
func LogIn(email string, password string) (string, string, error) {
	var hashedPassword string
	var id string
	var username string
	err := selectUserByEmailStmt.QueryRow(email).Scan(&hashedPassword, &username, &id); if err != nil {
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

func DeleteAccount(email string, token string) error {
	err := ConfirmIdentity(token); if err != nil {
		return err
	}

	_, err = deleteUserByEmailStmt.Exec(email); if err != nil {
		return err
	}

	return nil
}

func SetUpAuth() {
    var err error

    PendingAccounts = make(map[string]TempUser)

    selectUserByEmailStmt, err = db.Prepare("SELECT password, username, id FROM users WHERE email = $1")
    if err != nil {
        log.Fatal(err)
    }

    insertUserStmt, err = db.Prepare("INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id")
    if err != nil {
        log.Fatal(err)
    }

    deleteUserByEmailStmt, err = db.Prepare("DELETE FROM users WHERE email = $1")
    if err != nil {
        log.Fatal(err)
    }
}
