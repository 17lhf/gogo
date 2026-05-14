package model

// UserStatus represents the status of a user account.
type UserStatus int16

const (
	UserStatusEnabled  UserStatus = 1
	UserStatusDisabled UserStatus = 2
	UserStatusLocked   UserStatus = 3
)

func (s UserStatus) String() string {
	switch s {
	case UserStatusEnabled:
		return "enabled"
	case UserStatusDisabled:
		return "disabled"
	case UserStatusLocked:
		return "locked"
	default:
		return "unknown"
	}
}
