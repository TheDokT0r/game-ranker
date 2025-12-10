package database

import (
	"context"
	"game-ranker/users-manager/internal"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type User = internal.User

func InitDbTable() {
	dbUrl, present := os.LookupEnv("DATABASE_URL")
	if !present {
		log.Fatal("Invalid DATABASE_URL env variable")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Unable to enable pgcrypto extension: %v", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
		username TEXT NOT NULL,
		pass TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		role TEXT NOT NULL
	);
	`

	if _, err = conn.Exec(ctx, sqlStmt); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Users table initialized")
}

func Connect() *pgx.Conn {
	dbUrl, present := os.LookupEnv("DATABASE_URL")
	if !present {
		log.Fatal("Invalid DATABASE_URL env variable")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	return conn
}

func Close(conn *pgx.Conn) {
	err := conn.Close(context.Background())
	if err != nil {
		log.Fatal("DB close failed:", err)
	}
}

func AddUser(user User) error {
	conn := Connect()
	defer Close(conn)

	ctx := context.Background()

	sqlStmt := `
		INSERT INTO users (username, pass, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING public_id;
	`

	return conn.QueryRow(ctx, sqlStmt,
		user.Username,
		user.HashedPass,
		user.Email,
		user.Role,
	).Scan(&user.ID)
}

func GetUser(email string) (*User, error) {
	conn := Connect()
	defer Close(conn)

	ctx := context.Background()

	sqlStmt := `
		SELECT public_id, username, pass, email, role
		FROM users
		WHERE email = $1;
	`

	var user User

	err := conn.QueryRow(ctx, sqlStmt, email).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPass,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
