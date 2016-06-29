package isumm

const (
	Currency = "R$"
)

// Those are essentially sets.
var AllowedUsers = map[string]struct{}{
	"isumm.demo.staging@gmail.com": struct{}{},
	"danielfireman@gmail.com":      struct{}{},
	"contato@diasbruno.com":        struct{}{},
	"idnotfound@gmail.com":         struct{}{},
	"marco.rosner@gmail.com":        struct{}{},
}
var AllowedTestUsers = map[string]struct{}{
	"test@example.com": struct{}{},
}
