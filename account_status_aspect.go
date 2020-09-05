package acct

// AccountStatus is the type of account status
type AccountStatus int

const (
	AccountStatusNormal = iota // normal
	AccountStatusLocked        // locked
)

func (c AccountStatus) IsNormal() bool {
	return c == AccountStatusNormal
}

func (c AccountStatus) IsLocked() bool {
	return c == AccountStatusLocked
}
