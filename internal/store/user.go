package store

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/faresbakhit/9b/internal/config"
)

type User struct {
	Id             int
	Username       string
	HashedPassword []byte
	CreatedAt      *time.Time
	SessionToken   sql.NullString
}

var UserErrUsername = errors.New("invalid username")

func (s *Store) UserNew(username string, hashedPassword []byte) (*User, error) {
	// There is no password "requirements" validation, as I think, in my humble
	// opinion, that this is not a software's job. The user should be the one
	// to blame for weak passwords and rejecting them will only make the user,
	// who previous chose a weak password to choose the next weak password
	// within the arbitrary constraints a software has put.

	if !isUsernameValid(username) {
		return nil, UserErrUsername
	}

	sessionToken, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	hashedPassword, err = generateHashedPassword(hashedPassword)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:       username,
		HashedPassword: hashedPassword,
		SessionToken:   sql.NullString{String: sessionToken, Valid: true},
	}

	query := `
		INSERT INTO user
		(username, hashed_password, session_token)
		VALUES (?, ?, ?)
		RETURNING id, created_at`

	row := s.db.QueryRow(query, username, hashedPassword, sessionToken)
	if err := row.Scan(&user.Id, &user.CreatedAt); err != nil {
		log.Print(err)
		return nil, err
	}

	return &user, nil
}

func (s *Store) UserUpdatePassword(user *User, currentPassword, newPassword []byte) error {
	if err := compareHashAndPassword(user.HashedPassword, currentPassword); err != nil {
		return err
	}

	hashedPassword, err := generateHashedPassword(newPassword)
	if err != nil {
		return err
	}

	query := `
		UPDATE user
		SET session_token = NULL,
			hashed_password = ?
		WHERE id = ?`
	_, err = s.db.Exec(query, hashedPassword, user.Id)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s *Store) UserUpdateSessionToken(username string, password []byte) (sessionToken string, err error) {
	var hashedPassword []byte

	query := `
		SELECT (hashed_password)
		FROM user
		WHERE username = ?`
	row := s.db.QueryRow(query, username)
	if err = row.Scan(&hashedPassword); err != nil {
		log.Print(err)
		return
	}

	sessionToken, err = generateSessionToken()
	if err != nil {
		log.Print(err)
		return
	}

	if err = compareHashAndPassword(hashedPassword, password); err != nil {
		log.Print(err)
		return
	}

	query = `
		UPDATE user
		SET session_token = ?
		WHERE username = ?`
	if _, err = s.db.Exec(query, sessionToken, username); err != nil {
		log.Print(err)
	}

	return
}

func (s *Store) UserDeleteSessionToken(sessionToken string) error {
	query := `
		UPDATE user
		SET session_token = NULL
		WHERE session_token = ?`
	_, err := s.db.Exec(query, sessionToken)

	if err != nil {
		log.Print(err)
	}

	return err
}

func (s *Store) UserFromSessionToken(sessionToken string) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM user WHERE session_token = ?", sessionToken)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.HashedPassword,
		&user.SessionToken,
		&user.CreatedAt)

	if err != nil {
		log.Print(err)
	}

	return &user, err
}

func (s *Store) UserFromUsername(username string) (*User, error) {
	user := User{}

	row := s.db.QueryRow("SELECT * FROM user WHERE username = ?", username)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.HashedPassword,
		&user.SessionToken,
		&user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func compareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func generateHashedPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, config.PASSWORD_BCRYPT_COST)
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
//   - Length must be in [3, 21].
//   - Must have one or more character from the sets `[a-z]' and `[0-9]'.
//   - May include zero or more dashes (`-') or dots (`.'), but not at the
//     start, nor the end, nor before or after a dash (`-') or a dot (`.')
//     or an at sign (`@'), nor after a tilde (`~').
//   - May include at most one at sign (`@'), but not at the start, nor the
//     end, nor after a tilde.
//   - May start with at most one tilde (`~').
//
// No, This is not email vali_ation. I just wrote it in hope that any
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
		(notLowerAlphaNum(b2) && b1 == '~') ||
		notLowerAlphaNum(bn) {
		return false
	}

	ateAtSign := false
	for i := 1; i < n-1; i++ {
		bi := username[i]
		bin := username[i-1]
		bip := username[i+1]
		if bi == '@' {
			if ateAtSign || anyIsDashOrDot(bin, bip) {
				return false
			}
			ateAtSign = true
		} else if anyIsDashOrDot(bi) {
			if anyIsDashOrDot(bin, bip) {
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
