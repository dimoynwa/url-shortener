package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFount = errors.New("Redirect not found")
	ErrRedirectInvalid  = errors.New("Redirect invalid")
)

type redirectService struct {
	redirectRepository RedirectRepository
}

func NewRedirectService(repo RedirectRepository) RedirectService {
	return &redirectService{redirectRepository: repo}
}

func (service *redirectService) Find(code string) (*Redirect, error) {
	return service.redirectRepository.Find(code)
}

func (service *redirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return service.redirectRepository.Store(redirect)
}
