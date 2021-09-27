package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func JsonToStruct(filename string, output interface{}) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}

	json.Unmarshal(b, &output)
}
