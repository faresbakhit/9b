// Runtime configuration files cost us:
// - Extra input at runtime. One or more files!
// - Needless arguments about the 'right' langugage.
// - Dealing with languages defined by sheer luck (or not at all).
//
// Using centralized compile-time data for configuration is the only correct
// way to manage software configuration, following https://suckless.org.

package config

import "golang.org/x/crypto/bcrypt"

const (
	HTTP_SERVER_ADDR     = ":8000"
	PAGES_GLOB_PATTERN   = "templates/*.html"
	PAGES_RELOAD         = false
	PASSWORD_BCRYPT_COST = bcrypt.DefaultCost
	SESSION_TOKEN_LENGTH = 16
	SQLITE_DATA_SOURCE   = "9b.db"
	TLS_CERTIFICATE_FILE = ""
	TLS_PRIVATE_KEY_FILE = ""
)
