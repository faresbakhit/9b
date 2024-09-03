// Runtime configuration files are hard to manage, hard to get right, and
// needless complexity for any kind of HTTP service.
//
// Instead, using centralized compile-time data for configuration is the
// way to go, following https://suckless.org.

package config

import "golang.org/x/crypto/bcrypt"

const (
	HTTP_SERVER_ADDR     = ":8000"
	SQLITE_SOURCE_NAME   = "nineb.db"
	PAGES_GLOB_PATTERN   = "pages/*"
	TLS_CERTIFICATE_FILE = ""
	TLS_PRIVATE_KEY_FILE = ""
	PASSWORD_BCRYPT_COST = bcrypt.DefaultCost
	SESSION_TOKEN_LENGTH = 16
)
