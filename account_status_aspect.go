package acct

import (
	"encoding/json"
	"errors"
)

// AccountStatus is the type of account status
type AccountStatus int

const (
	AccountStatusNormal AccountStatus = iota // normal
	AccountStatusLocked                      // locked
)

func (c AccountStatus) IsNormal() bool {
	return c == AccountStatusNormal
}

func (c AccountStatus) IsLocked() bool {
	return c == AccountStatusLocked
}

func (c AccountStatus) MarshalJSON() ([]byte, error) {
	var status string
	switch c {
	case AccountStatusNormal:
		status = "normal"
	case AccountStatusLocked:
		status = "locked"
	}

	return json.Marshal(status)
}

func (c *AccountStatus) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"normal"`:
		*c = AccountStatusNormal
	case `"locked"`:
		*c = AccountStatusLocked
	default:
		return errors.New("unknown account status")
	}
	return nil
}
