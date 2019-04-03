package sms

import (
	"fmt"
	"net/http"
)

func CheckAdd(token string, size int64) (err error) {
	policy, err := parsePolicy(token)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v%d/resource/check/%d", SMSAddr, policy.Ver, size)

	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	r.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}

	return handleResp(resp)
}

func FinishAdd(token string, size uint64, hash string) (err error) {
	policy, err := parsePolicy(token)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v%d/resource/finish/%d/%s", SMSAddr, policy.Ver, size, hash)

	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return
	}
	r.Header.Set("token", token)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}

	return handleResp(resp)
}
