package json

import (
	"encoding/json"

	"github.com/dimoynwa/url-shortener/shortener"
	errs "github.com/pkg/errors"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(shortener.ErrRedirectInvalid, "json.Serialized.Decode")
	}

	return redirect, nil
}

func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	data, err := json.Marshal(redirect)

	if err != nil {
		return nil, errs.Wrap(err, "json.Serializer.Decode")
	}

	return data, nil
}
