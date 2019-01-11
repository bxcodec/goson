package goson_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
		"numbers": [3, 4, 5],
		"nested":[[1,2], [2,4]],
		"null": null
	}`
	malformedJSON = `{
		"id": "uk123",
		"name": "tom",
		"email": "tom@gmail.com",
		"user": {
			"address": "Sidikalang"
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
	malformedArray = `[{
		"id": "uk123",
		"name": "tom",
		"email": "tom@gmail.com",
		"user": {
			"address": "Sidikalang"
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
	exampleArr2   = `[1 ,3 ,4 ,5]`
	filePath      = "./testfile.json"
	wrongFilePath = "./testdata/testfile.json"
	wrongEndpoint = "test.qwerty.com:1234/post/1/a/1"
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

func TestJSONSchema_MalformedJSON(t *testing.T) {
	_, err := goson.GenerateJSONSchema(malformedJSON)
	if err == nil {
		t.Error("Expected nil; but got: ", err)
	}
}

func TestJSONSchema_MalformedArray(t *testing.T) {
	_, err := goson.GenerateJSONSchema(malformedArray)
	if err == nil {
		t.Error("Expected nil; but got: ", err)
	}
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

func TestJSONSchemaFromURL_WrongEndpoint(t *testing.T) {
	_, err := goson.GenerateJSONSchemaFromURL(wrongEndpoint)
	if err == nil {
		t.Error("Expected nil; but got: ", err)
	}
}

func TestJSONSchemaFromFile(t *testing.T) {
	schema, err := goson.GenerateJSONSchemaFromFile(filePath)
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	jbyt, err := schema.ToJSON()
	if err != nil {
		t.Error("Expected nil; but got: ", err)
	}
	fmt.Println("[SCHEMA] ", string(jbyt))
}

func TestJSONSchemaFromFile_WrongPath(t *testing.T) {
	_, err := goson.GenerateJSONSchemaFromFile(wrongFilePath)
	if err == nil {
		t.Error("Expected nil; but got: ", err)
	}
}
