package commands

import (
	cmdkit "gx/ipfs/QmdE4gMduCKCGAcczM2F5ioYDfdeKuPix138wrES1YSr7f/go-ipfs-cmdkit"

	"github.com/pkg/errors"

	"bytes"
	"encoding/json"
	"net/http"

	"io/ioutil"

	cmds "github.com/ipfs/go-ipfs/commands"
)

const (
	uosUrlGetTableRows = "http://43.242.156.61:9008/v1/chain/get_table_rows"
	uosTableName       = "udfsuserrd"
	uosCode            = "ulorduosudfs"
)

var CheckCmd = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:          "check special user content hash valid on the uos.",
		ShortDescription: ``,
	},

	Arguments: []cmdkit.Argument{
		cmdkit.StringArg("user-account", true, false, "The account of user to valid for."),
		cmdkit.StringArg("content-hash", true, false, "The hash to be valid."),
	},
	Options: []cmdkit.Option{},
	Run: func(req cmds.Request, res cmds.Response) {
		account := req.StringArguments()[0]
		hashes := req.StringArguments()[1:]

		err := ValidOnUOS(account, hashes...)
		if err != nil {
			res.SetError(errors.Wrap(err, "valid failed"), cmdkit.ErrNormal)
			return
		}

		if len(hashes) > 0 {
			res.SetOutput("all valid")
		} else {
			res.SetOutput("valid")
		}

	},
}

// {"scope":"testertester" , "code":"ulorduosudfs" ,"table":"udfsuserrd" , "json":"true"}
type requestGetTableRows struct {
	Scope string `json:"scope"`
	Code  string `json:"code"`
	Table string `json:"table"`
	Json  bool   `json:"json"`
}

// {"rows":[{"inner_id":"2133026024183134724","size":30,"hash":"3adc120000000000000000000000000000000000000000000000000000000000",
// "folder":"a/b/d","string_name":"myudfs_2"},{"inner_id":"3538262513698496072","size":20,
// "hash":"1234000000000000000000000000000000000000000000000000000000000000","folder":"a/b/c","string_name":"myudfs"}],"more":false}

type tableRow struct {
	InnerID    string `json:"inner_id"`
	Size       int32  `json:"size"`
	Hash       string `json:"hash"`
	Folder     string `json:"folder"`
	StringName string `json:"string_name"`
}
type responseGetTableRows struct {
	Rows []*tableRow `json:"rows"`
}

func ValidOnUOS(account string, hashes ...string) error {
	rgtr := requestGetTableRows{
		Scope: account,
		Code:  uosCode,
		Table: uosTableName,
		Json:  true,
	}
	data, err := json.Marshal(rgtr)
	if err != nil {
		return errors.Wrap(err, "marshal getTableRowsBody error")
	}

	resp, err := http.Post(uosUrlGetTableRows, "application/json", bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "post failed")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read body from response error")
	}

	respTable := &responseGetTableRows{}
	err = json.Unmarshal(body, respTable)
	if err != nil {
		return errors.Wrap(err, "unmarshal body from response error")
	}

	count := len(respTable.Rows)
	if count == 0 {
		return errors.New("the user does`t have any content on the uos.")
	}

	for _, hash := range hashes {
		if inRows(respTable.Rows, hash) {
			continue
		}

		return errors.Errorf("can`t find hash=%s", hash)
	}

	return nil
}

func inRows(rows []*tableRow, hash string) bool {
	for _, row := range rows {
		if row.Hash == hash {
			return true
		}
	}

	return false
}
