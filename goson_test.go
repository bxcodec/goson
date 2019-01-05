package goson_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/goson"
	"github.com/spf13/afero"
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
		"numbers": [3, 4, 5],
		"nested":[[1,2], [2,4]] 
	}`
	exampleArr = `[{
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
	}]`
	exampleArr2 = `[1 ,3 ,4 ,5]`
)

func TestJSONSchema(t *testing.T) {
	schema, err := goson.GenerateJSONSchema(example)
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jbyt, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA] ", string(jbyt))
}

func TestJSONSchemaArray(t *testing.T) {
	schema, err := goson.GenerateJSONSchema(exampleArr2)
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jbyt, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA ARR] ", string(jbyt))
}

func TestJSONSchemaFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(example))

			if r.Method != "GET" {
				t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
			}

			if r.URL.EscapedPath() != "/posts/1" {
				t.Errorf("Expected request to ‘/posts/1’, got ‘%s’", r.URL.EscapedPath())
			}
		},
	))

	defer ts.Close()

	schema, err := goson.GenerateJSONSchemaFromURL(ts.URL + "/posts/1")
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jbyt, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA] ", string(jbyt))
}

func TestJSONSchemaFromFile(t *testing.T) {
	appFS := afero.NewMemMapFs()
	appFS.MkdirAll("testdata", 0755)
	afero.WriteFile(appFS, "testdata/data.json", []byte(example), 0644)
	file, _ := appFS.Open("testdata/data.json")

	schema, err := goson.GenerateJSONSchemaFromFile(file.(io.Reader))
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jbyt, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA] ", string(jbyt))
}
