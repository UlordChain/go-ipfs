package commands

import (
	"bytes"
	"encoding/json"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"
	"github.com/pkg/errors"
	"gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	"gx/ipfs/Qmde5VP1qUkyQXKCfmEUA7bP64V2HAptbJ7phuPp7jXWwg/go-ipfs-cmdkit"
	"net/http"

	"io/ioutil"

	"gx/ipfs/QmSXUokcP4TJpFfqozT69AVAYRtzXVMUjzQVkYX41R9Svs/go-ipfs-cmds"
)
//
//const (
//	uosUrlGetTableRows = "http://43.242.156.61:9008/v1/chain/get_table_rows"
//	uosTableName       = "udfsuserrd"
//	uosCode            = "ulorduosudfs"
//)

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
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		account := req.Arguments[0]
		hash := req.Arguments[1]

		node, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		cfg, _ := node.Repo.Config()
		if !cfg.UOSCheck.Enable {
			return errors.New("disabled")
		}

		_, err = ValidOnUOS(&cfg.UOSCheck, account, hash)
		if err != nil {
			return errors.Wrap(err, "valid failed")
		}

		if len(hash) > 0 {
			res.Emit("all valid")
		} else {
			res.Emit("valid")
		}

		return nil
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

func ValidOnUOS(uoscheck *config.UOSCheck, account string, hash string) (uint64, error) {
	rgtr := requestGetTableRows{
		Scope: account,
		Code:  uoscheck.Code,
		Table: uoscheck.TableName,
		Json:  true,
	}
	data, err := json.Marshal(rgtr)
	if err != nil {
		return 0, errors.Wrap(err, "marshal getTableRowsBody error")
	}

	resp, err := http.Post(uoscheck.UrlForGetTableRows, "application/json", bytes.NewReader(data))
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
