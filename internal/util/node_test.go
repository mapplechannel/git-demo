package util_test

import (
	"encoding/json"
	"fmt"
	"hsm-scheduling-back-end/internal/util"
	"testing"
)

func TestGetNodeServer(t *testing.T) {
	res := util.GetNodeServer()
	fmt.Println(res)
}

func TestJson(t *testing.T) {

	req := make(map[string]string)
	requestBody, err1 := json.Marshal(req)

	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println(string(requestBody))
}
