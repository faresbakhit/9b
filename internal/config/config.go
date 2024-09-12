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
	HTTP_SERVER_ADDR          = ":8000"
	HTTP_PUBLIC_DIRECTORY     = "./public/"
	HTTP_PUBLIC_ROUTE         = "/assets/"
	TEMPLATES_GLOB_PATTERN    = "./templates/*.html"
	TEMPLATES_LOAD_ON_RENDER  = true
	PASSWORD_BCRYPT_COST      = bcrypt.DefaultCost
	SESSION_TOKEN_COOKIE_NAME = "session_token"
	SESSION_TOKEN_LENGTH      = 16
	SQLITE_DATA_SOURCE        = "file:9b.db"
	TLS_CERTIFICATE_FILE      = ""
	TLS_PRIVATE_KEY_FILE      = ""
	ALLOW_HTTP_URLS           = false
)
