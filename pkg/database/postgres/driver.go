package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adaggerboy/utasksd/models/config"
	"github.com/adaggerboy/utasksd/pkg/database"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	database.IDatabase
	db *sql.DB
	// timeout time.Duration
	tx *sql.Tx

	// statements map[string]*sql.Stmt
}

func init() {
	database.GetDatabaseFabric().RegisterDatabaseFabric("postgres", getDI)
}

func createConnect(endpoint config.DatabaseEndpointConfig, user, password string) (conn *sql.DB, err error) {
	conn, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		password,
		endpoint.Host,
		endpoint.Port,
		endpoint.Database))
	return
}

func getDI(ctx context.Context, endpoint config.DatabaseEndpointConfig, user, password string) (database.IDatabase, error) {
	conn, err := createConnect(endpoint, user, password)
	if err != nil {
		return nil, fmt.Errorf("postgres database connect error: %s", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres database ping error: %s", err)
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("postgres start transaction error: %s", err)
	}
	_, err = tx.Exec("set constraints all deferred")
	if err != nil {
		return nil, fmt.Errorf("postgres start transaction error: %s", err)
	}
	// conn.Close()

	reader := &PostgresConnection{
		db: conn,
		// timeout: time.Duration(endpoint.Timeout),
		tx: tx,
	}
	return reader, err
}

func (db *PostgresConnection) CreateDBUser(ctx context.Context, username string, password string) (err error) {
	quotedUsername := pq.QuoteIdentifier(username)
	quotedPassword := pq.QuoteLiteral(password)
	_, err = db.db.ExecContext(ctx, fmt.Sprintf("create user %s with encrypted password %s; grant generic to %s", quotedUsername, quotedPassword, quotedUsername))
	return
}
func (db *PostgresConnection) DeleteDBUser(ctx context.Context, username string) (err error) {
	quotedUsername := pq.QuoteIdentifier(username)
	_, err = db.db.ExecContext(ctx, fmt.Sprintf("drop user %s ", quotedUsername))
	return
}
func (db *PostgresConnection) ChangeDBUserPassword(ctx context.Context, username string, password string) (err error) {
	quotedUsername := pq.QuoteIdentifier(username)
	quotedPassword := pq.QuoteLiteral(password)
	_, err = db.db.ExecContext(ctx, fmt.Sprintf("alter user %s with encrypted password %s", quotedUsername, quotedPassword))
	return
}

func (db *PostgresConnection) SetRole(ctx context.Context, username, role string, set bool) (err error) {
	action := "grant"
	slog := "to"
	if !set {
		action = "revoke"
		slog = "from"
	}
	quotedUsername := pq.QuoteIdentifier(username)
	_, err = db.db.ExecContext(ctx, fmt.Sprintf("%s %s %s %s", action, role, slog, quotedUsername))
	return
}
