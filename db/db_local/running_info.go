package db_local

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe/constants"
	"github.com/turbot/steampipe/utils"
)

// RunningDBInstanceInfo contains data about the running process and it's credentials
type RunningDBInstanceInfo struct {
	Pid        int
	Port       int
	Listen     []string
	ListenType StartListenType
	Invoker    constants.Invoker
	Password   string
	User       string
	Database   string
}

func (r *RunningDBInstanceInfo) Save() error {
	if content, err := json.Marshal(r); err != nil {
		return err
	} else {
		return ioutil.WriteFile(runningInfoFilePath(), content, 0644)
	}
}

func (r *RunningDBInstanceInfo) String() string {
	writeBuffer := bytes.NewBufferString("")
	jsonEncoder := json.NewEncoder(writeBuffer)

	// redact the password from the string, so that it doesn't get printed
	// this should not affect the state file, since we use a json.Marshal there
	p := r.Password
	r.Password = "XXXX-XXXX-XXXX"

	jsonEncoder.SetIndent("", "")
	jsonEncoder.Encode(r)
	r.Password = p
	return writeBuffer.String()
}

func loadRunningInstanceInfo() (*RunningDBInstanceInfo, error) {
	utils.LogTime("db.loadRunningInstanceInfo start")
	defer utils.LogTime("db.loadRunningInstanceInfo end")

	if !helpers.FileExists(runningInfoFilePath()) {
		return nil, nil
	}

	fileContent, err := ioutil.ReadFile(runningInfoFilePath())
	if err != nil {
		return nil, err
	}
	var info = new(RunningDBInstanceInfo)
	err = json.Unmarshal(fileContent, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func removeRunningInstanceInfo() error {
	return os.Remove(runningInfoFilePath())
}

func runningInfoFilePath() string {
	return filepath.Join(constants.InternalDir(), "steampipe.json")
}
