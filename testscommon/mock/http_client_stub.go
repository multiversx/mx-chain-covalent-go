package mock

import "net/http"

type HTTPClientStub struct {
	GetCalled func(url string) (resp *http.Response, err error)
}

func (hcs *HTTPClientStub) Get(url string) (resp *http.Response, err error) {
	if hcs.GetCalled != nil {
		return hcs.GetCalled(url)
	}

	return nil, nil
}
