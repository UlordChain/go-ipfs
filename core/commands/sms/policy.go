package sms

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const(
	SMSAddr	= "http://192.168.14.42:9090"
)
type policyObj struct{
	Ver int
	Expired int64
	Ext json.RawMessage
}


type respBody struct{
	Message string
}


func parsePolicy(token string) (*policyObj, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	tokenSlice := strings.Split(token, ":")
	if len(tokenSlice) != 3 {
		return nil, errors.New("token format error")
	}

	bs, err := base64.URLEncoding.DecodeString(tokenSlice[2])
	if err != nil {
		return nil, errors.Wrap(err, "decode token failed")
	}

	obj := &policyObj{}
	err = json.Unmarshal(bs, obj)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal policy string failed")
	}

	if time.Now().After(time.Unix(obj.Expired, 0)) {
		return nil, errors.Wrap(err, "token have out of date")
	}

	return obj, nil
}


func handleResp(resp *http.Response) error{
	defer resp.Body.Close()
	body := &respBody{}
	if resp.ContentLength > 0 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read response body failed")
		}

		err = json.Unmarshal(bs, body)
		if err != nil {
			return errors.Wrap(err, "unmarshal response body failed: " + string(bs))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("%d %s", resp.StatusCode, body.Message)
	}

	return nil
}