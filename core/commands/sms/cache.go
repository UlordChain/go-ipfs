package sms

import (
	"fmt"
	"net/http"
)

func Cache(token, hash, peer string) (err error) {
	policy, err := parsePolicy(token)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v%d/resource/cache/%s/%s", SMSAddr, policy.Ver, peer, hash)

	r, err :=http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	r.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}

	err = handleResp(resp)
	return
}
