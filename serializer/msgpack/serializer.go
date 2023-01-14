package msgpack

import (
	"github.com/dimoynwa/url-shortener/shortener"
	errs "github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(shortener.ErrRedirectInvalid, "msgpack.Serialized.Decode")
	}

	return redirect, nil
}

func (r *Redirect) Encode(redirect *shortener.Redirect) ([]byte, error) {
	data, err := msgpack.Marshal(redirect)

	if err != nil {
		return nil, errs.Wrap(err, "msgpack.Serializer.Decode")
	}

	return data, nil
}
