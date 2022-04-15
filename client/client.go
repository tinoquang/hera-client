package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sethvargo/go-retry"
)

type client struct {
	*sql.DB
}

func newClient(db *sql.DB) *client {
	return &client{db}
}

// Tx execs fn inside a transaction.
func (c *client) tx(ctx context.Context, fn func(TxConn) error) (err error) {
	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx err: %v, rb err: %v\n", err, rbErr)
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Tx executes fn inside a transaction, with maximum 3 retries.
func (c *client) Tx(ctx context.Context, conn DBConn, fn func(TxConn) error) error {
	return c.TxWithRetry(ctx, conn, 3, fn)
}

// TxWithRetry executes fn inside a transaction
// with fibonacci backoff and retry.
// Retry is done with maximum retries.
// Retry strategy: 1s -> 1s -> 2s -> 3s -> 5s -> ...
func (c *client) TxWithRetry(ctx context.Context, conn DBConn, maxAttempts int, fn func(TxConn) error) error {
	if maxAttempts <= 0 {
		return fmt.Errorf("maxAttempts must be greater than 0")
	}

	b := retry.NewExponential(1 * time.Second)

	// Set max retries
	b = retry.WithMaxRetries(uint64(maxAttempts), b)
	return retry.Do(ctx, b, func(ctx context.Context) error {
		return retry.RetryableError(c.tx(ctx, fn))
	})
}
