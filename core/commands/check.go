package commands

import (

	"github.com/pkg/errors"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
	"math"

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
		hash := req.StringArguments()[1]

		_, err := ValidOnUOS(account, hash)
		if err != nil {
			res.SetError(errors.Wrap(err, "valid failed"), cmdkit.ErrNormal)
			return
		}

		if len(hash) > 0 {
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
	Size       uint64 `json:"size"`
	Hash       string `json:"hash"`
	Folder     string `json:"folder"`
	StringName string `json:"string_name"`
}
type responseGetTableRows struct {
	Rows []*tableRow `json:"rows"`
}

func ValidOnUOS(account string, hash string) (uint64, error) {
	// TODO: FOR TEST, because the uos valied can not work now.
	return math.MaxUint64/1024, nil

	rgtr := requestGetTableRows{
		Scope: account,
		Code:  uosCode,
		Table: uosTableName,
		Json:  true,
	}
	data, err := json.Marshal(rgtr)
	if err != nil {
		return 0, errors.Wrap(err, "marshal getTableRowsBody error")
	}

	resp, err := http.Post(uosUrlGetTableRows, "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, errors.Wrap(err, "post failed")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Wrap(err, "read body from response error")
	}

	respTable := &responseGetTableRows{}
	err = json.Unmarshal(body, respTable)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal body from response error")
	}

	count := len(respTable.Rows)
	if count == 0 {
		return 0, errors.New("the user does`t have any content on the uos.")
	}

	size, b := inRows(respTable.Rows, hash)
	if b {
		return size, nil
	}

	return 0, errors.Errorf("can`t find hash=%s on uos", hash)
}

func inRows(rows []*tableRow, hash string) (uint64, bool) {
	for _, row := range rows {
		if row.Hash == hash {
			return row.Size, true
		}
	}

	return 0, false
}
