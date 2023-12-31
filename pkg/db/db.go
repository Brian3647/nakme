package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

var (
	once sync.Once
	jwtSecret []byte
	dbSecret string
	db *sql.DB
)

func SetupDB() (*sql.DB, error) {
    dataSourceName := fmt.Sprintf("user=postgres password=%s dbname=postgres sslmode=disable", dbSecret)
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func setUpSecret(name string) (string, error) {
	env := os.Getenv(name); if env == "" {
		return "", fmt.Errorf("environment variable %s is not set", name)
	}

	return env, nil
}

func SetUpSecrets() error {
    var err error
    jwtSecretString, err := setUpSecret("JWT_SECRET")
    if err != nil {
        return err
    }

	jwtSecret = []byte(jwtSecretString)

    dbSecret, err = setUpSecret("POSTGRES_PASSWORD")
    if err != nil {
        return err
    }

    return nil
}

func SetUp() error {
	var err error
    once.Do(func() {
        err = SetUpSecrets()
        if err != nil {
            return
        }

        db, err = SetupDB()
		SetUpAuth()
    })

    return err
}
