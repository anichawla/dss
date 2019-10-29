package auth

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	errMissingOrEmptySubject = errors.New("Missing or empty subject")
	errTokenExpireTooFar     = errors.New("Token expiration time is too far in the furture, Max token duration is 1 Hour")
	errMissingIssuer         = errors.New("Missing Issuer URI")
	// Now allows test to override with specific time values
	Now = time.Now
)

// scopes models a set of scopes.
type scopes map[string]struct{}

func (s *scopes) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*s = map[string]struct{}{}

	for _, scope := range strings.Split(str, " ") {
		(*s)[scope] = struct{}{}
	}

	return nil
}

type claims struct {
	jwt.StandardClaims
	Scopes scopes `json:"scope"`
}

func (c *claims) Valid() error {
	if c.Subject == "" {
		return errMissingOrEmptySubject
	}
	now := Now()

	c.VerifyExpiresAt(now.Unix(), true)

	if c.ExpiresAt > now.Add(time.Hour).Unix() {
		return errTokenExpireTooFar
	}

	if c.Issuer == "" {
		return errMissingIssuer
	}

	return c.StandardClaims.Valid()
}
