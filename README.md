# goson
⚡️[WIP] Golang SIMPLE JSON Schema extractor from raw-json string

## Roadmap
Currently with MVP release ( target: stable v1), there are a few things we need to finish.
You can view the list in the: https://waffle.io/bxcodec/goson

## Contributions
I really appreciate for everyone who willing to help to finish this, you can put any task from https://waffle.io/bxcodec/goson.

The task with label `epic` means it's need to finish ASAP.

## Example

```go
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
schema,err:=goson.GenerateJSONSchema(example)
if err != nil {
    log.Fatal(err)
}

jbyt, err := schema.ToJSON()

if err != nil {
	t.Error("Expected nil; but got: ", err)
}

fmt.Println(string(jbyt))
/*
 Will Print: 
{"properties":{"arr":{"items":{"properties":{"site":{"type":"string"},"url":{"type":"string"}},"type":"object"},"type":"array"},"email":{"type":"string"},"id":{"type":"string"},"name":{"type":"string"},"numbers":{"items":{"type":"number"},"type":"array"},"tags":{"items":{"type":"string"},"type":"array"},"user":{"properties":{"address":{"type":"string"},"age":{"type":"number"},"phones":{"items":{"type":"string"},"type":"array"}},"type":"object"}},"type":"object"}
*/

// Or in Formatted data:
/*
{
 "properties": {
	"arr": {
		"items": {
			"properties": {
				"site": {
					"type": "string"
				},
				"url": {
					"type": "string"
				}
			},
			"type": "object"
		},
		"type": "array"
	},
	"email": {
		"type": "string"
	},
	"id": {
		"type": "string"
	},
	"name": {
		"type": "string"
	},
	"numbers": {
		"items": {
			"type": "number"
		},
		"type": "array"
	},
	"tags": {
		"items": {
			"type": "string"
		},
		"type": "array"
	},
	"user": {
		"properties": {
			"address": {
				"type": "string"
			},
			"age": {
				"type": "number"
			},
			"phones": {
				"items": {
					"type": "string"
				},
				"type": "array"
			}
		},
		"type": "object"
	}
	},
	"type": "object"
}
*/
```
