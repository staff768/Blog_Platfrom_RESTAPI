package models
import (
	"database/sql"
	"errors"
)

var ErrNoRecord = sql.ErrNoRows
var ErrDuplicateEmail = errors.New("duplicate email")
var ErrInvalidCredentials = errors.New("invalid credentials")
