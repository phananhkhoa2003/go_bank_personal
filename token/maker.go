package token

import "time"

// Maker is an interface for creating tokens.
type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid and returns the payload if it is.
	VerifyToken(token string) (*Payload, error)
}
