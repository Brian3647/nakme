// Go "script" to set-up the postgres & redis databases.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func exec(db *sql.DB, query string) {
    _, err := db.Exec(query)
    if err != nil {
        if strings.Contains(err.Error(), "already exists") {
            return
        }
        log.Fatal(err)
    }
}

func main() {
	fmt.Println("Warning: this WILL drop tables if they exist.")
	fmt.Println("Sleeping for 10 seconds to give you time to cancel this.")
	time.Sleep(10 * time.Second)

    adminPwd := os.Getenv("POSTGRES_PASSWORD")

    conn := fmt.Sprintf("user=postgres password=%s dbname=postgres sslmode=disable", adminPwd)

    db, err := sql.Open("postgres", conn)
    if err != nil {
        log.Fatal(err)
    }

    exec(db, `CREATE DATABASE nakme;`)

    db.Close()

    conn = fmt.Sprintf("user=postgres password=%s dbname=nakme sslmode=disable", adminPwd)

    db, err = sql.Open("postgres", conn)
    if err != nil {
        log.Fatal(err)
    }

	db.Exec("DROP TABLE users;")

    exec(db, `CREATE TABLE users (
        id serial PRIMARY KEY,
        username VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
        password bytea NOT NULL, -- hashed
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`)

    defer db.Close()
}
