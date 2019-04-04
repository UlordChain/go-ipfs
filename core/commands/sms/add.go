package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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


type finishAddRequestBody struct {
	Peer     string
	FileName string
}


func FinishAdd(token string, size uint64, hash, peer, filename string) (err error) {
	policy, err := parsePolicy(token)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v%d/resource/finish/%s/%d", SMSAddr, policy.Ver, hash, size)

	farb := &finishAddRequestBody{
		Peer: peer,
		FileName: filename,
	}
	bs , err := json.Marshal(farb)
	if err != nil {
		err = errors.Wrap(err, "unmarshal finish add request body failed")
		return
	}
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
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

type BackupRequestBody struct {
	Status   int8
	Description string
	Peers    []string
}

func BackupResult(token string, hash string, brb *BackupRequestBody) (err error) {
	policy, err := parsePolicy(token)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v%d/resource/backup/%s", SMSAddr, policy.Ver, hash)

	bs , err := json.Marshal(brb)
	if err != nil {
		err = errors.Wrap(err, "unmarshal finish add request body failed")
		return
	}
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
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