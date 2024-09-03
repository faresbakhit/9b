package store

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/faresbakhit/9b/config"
)

type User struct {
	Id             int64
	Name           string
	HashedPassword []byte
	CreatedAt      *time.Time
	SessionToken   sql.NullString
}

var ErrInvalidUsername = errors.New("invalid username")

func (s *Store) CreateUser(name string, hashedPassword []byte) (*User, error) {
	if !isUsernameValid(name) {
		return nil, ErrInvalidUsername
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	hashedPassword, err = bcrypt.GenerateFromPassword(hashedPassword, config.PASSWORD_BCRYPT_COST)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:           name,
		HashedPassword: hashedPassword,
		SessionToken:   sql.NullString{String: sessionToken, Valid: true},
	}

	query := `
		INSERT INTO user
		(name, hashed_password, session_token)
		VALUES
		(?, ?, ?)
		RETURNING
		id, created_at`

	row := s.db.QueryRow(query, name, hashedPassword, sessionToken)
	if err := row.Scan(&user.Id, &user.CreatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) UpdateUserPassword(user *User, currentPassword, newPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, currentPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(newPassword, config.PASSWORD_BCRYPT_COST)
	if err != nil {
		return err
	}

	query := `
		UPDATE user
		SET session_token = NULL,
			hashed_password = ?
		WHERE id = ?`
	if _, err := s.db.Exec(query, hashedPassword, user.Id); err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateUserSessionToken(name string, password []byte) (token string, err error) {
	var hashedPassword []byte

	query := `
		SELECT (hashed_password)
		FROM user
		WHERE name = ?`

	row := s.db.QueryRow(query, name)
	if err = row.Scan(&hashedPassword); err != nil {
		return "", err
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return "", err
	}

	query = `
		UPDATE user
		SET session_token = ?
		WHERE name = ?`
	if _, err := s.db.Exec(query, sessionToken, name); err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (s *Store) InvalidateUserSessionToken(sessionToken string) error {
	query := `
		UPDATE user
		SET session_token = NULL
		WHERE session_token = ?`
	if _, err := s.db.Exec(query, sessionToken); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserBySessionToken(sessionToken string) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM user WHERE session_token = ?", sessionToken)
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.SessionToken); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserByName(name string) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM user WHERE name = ?", name)
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.SessionToken); err != nil {
		return nil, err
	}

	return &user, nil

}

func generateSessionToken() (string, error) {
	b := make([]byte, config.SESSION_TOKEN_LENGTH)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// base 16, lower-case, two characters per byte
	return fmt.Sprintf("%x", b), nil
}

// Specification for a correct username:
//   - Length must be in [USERNAME_MINIMUM_LENGTH, USERNAME_MAXIMUM_LENGTH].
//   - Must have one or more character from the sets `[a-z]' and `[0-9]'.
//   - May include zero or more dashes (`-') or dots (`.'), but not at the
//     start, nor the end, nor before or after a dash (`-') or a dot (`.')
//     or an at sign (`@'), nor after a tilde (`~').
//   - May include at most one at sign (`@'), but not at the start, nor the
//     end, nor after a tilde.
//   - May start with at most one tilde (`~').
//
// No, This is not email validation. I just wrote it in hope that any
// username within those constraints will look elegant and easy to type.
func isUsernameValid(username string) bool {
	n := len(username)

	if n < 3 || n > 21 {
		return false
	}

	b1 := username[0]
	b2 := username[1]
	bn := username[n-1]

	if (notLowerAlphaNum(b1) && b1 != '~') ||
		(b1 == '~' && notLowerAlphaNum(b2)) ||
		notLowerAlphaNum(bn) {
		return false
	}

	ateAtSign := false
	for i := 1; i < n-1; i++ {
		bi := username[i]
		bim1 := username[i-1]
		bip1 := username[i+1]
		if bi == '@' {
			if ateAtSign || anyIsDashOrDot(bim1, bip1) {
				return false
			}
			ateAtSign = true
		} else if anyIsDashOrDot(bi) {
			if anyIsDashOrDot(bim1, bip1) {
				return false
			}
		} else if notLowerAlphaNum(bi) {
			return false
		}
	}

	return true
}

func notLowerAlphaNum(b byte) bool {
	return (b < 'a' || b > 'z') && (b < '0' || b > '9')
}

func anyIsDashOrDot(bytes ...byte) bool {
	for _, b := range bytes {
		if b == '-' || b == '.' {
			return true
		}
	}
	return false
}
