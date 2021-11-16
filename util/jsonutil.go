package util

import (
	"encoding/json"
)

type jsonUtil struct {
}

func (c *jsonUtil) SliceToString(r []string) string {
	if r == nil {
		return ""
	}
	result, er := json.Marshal(r)
	if er != nil {
		panic(er)
	}
	return string(result)
}

func (c *jsonUtil) To2String(r interface{}) string {
	if r == nil {
		return ""
	}
	result, er := json.Marshal(r)
	if er != nil {
		panic(er)
	}
	return string(result)
}

func (c *jsonUtil) To2PrettyString(r interface{}) string {
	if r == nil {
		return ""
	}
	result, er := json.MarshalIndent(r, "", "    ")
	if er != nil {
		panic(er)
	}
	return string(result)
}

func (c *jsonUtil) BoolStr(status bool) interface{} {
	if status {
		return "success"
	} else {
		return "failure"
	}
}

var JsonUtil jsonUtil = jsonUtil{}
