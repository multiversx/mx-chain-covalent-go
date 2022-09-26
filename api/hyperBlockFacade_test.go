package api

import (
	"fmt"
	"testing"
)

func TestSmth(t *testing.T) {
	path := buildUrlWithBlockQueryOptions(fmt.Sprintf("https://gateway.elrond.com/%s/%d", hyperBlockPathByNonce, 3333), options)
	fmt.Println(path)
}
