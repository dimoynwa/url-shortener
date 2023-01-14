package shortener

type RedirectSerializer interface {
	Decode([]byte) (*Redirect, error)
	Encode(*Redirect) ([]byte, error)
}
