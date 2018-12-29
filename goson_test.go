package goson_test

import (
	"fmt"
	"testing"

	"github.com/bxcodec/goson"
)

var (
	example = `{
		"id": "uk123",
		"name": "tom",
		"email": "tom@gmail.com",
		"user": {
			"address": "Sidikalang",
			"age": 23,
			"phones": [
				"08234239523",
				"08234239523"
			]
		},
		"arr": [{
			"site": "string",
			"url": "string"
		}],
		"tags": ["go", "js"],
		"numbers": [3, 4, 5] 
	}`
)

func TestJSONSchema(t *testing.T) {
	schema, err := goson.GenerateJSONSchema(example)
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jby, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA] ", string(jby))
}
