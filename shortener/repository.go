package shortener

type RedirectRepository interface {
	Find(string) (*Redirect, error)
	Store(*Redirect) error
}
