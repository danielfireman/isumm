package isumm

const (
	Currency = "R$"
)

// Those are essentially sets.
var AllowedUsers = map[string]struct{}{
	"danielfireman@gmail.com": struct{}{},
}
var AllowedTestUsers = map[string]struct{}{
	"test@example.com": struct{}{},
}
