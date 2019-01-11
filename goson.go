package goson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	key_generic_array_item = "*:"
)

type Node struct {
	Type  string
	Value interface{}
	Key   string
}

type MapNode map[string]*Node

type MapString map[string]interface{}

type ArrMapNode []MapNode

func (a ArrMapNode) generateJSONSchema(isrequired bool) MapString {
	res := map[string]interface{}{}
	res["type"] = "array"
	if len(a) > 0 {
		if key := getKeymapPrefix(a[0], key_generic_array_item); key != "" {
			// INDICATE GENERIC ARRAY-ITEMS
			node := a[0][key]
			itemsSchema := map[string]interface{}{
				"type": node.Type,
			}
			res["items"] = itemsSchema

		} else {
			res["items"] = a[0].generateJSONSchema(false)
		}

	}

	return res
}

func (m MapString) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func getKeymapPrefix(item MapNode, prefix string) string {
	for k, _ := range item {
		if strings.HasPrefix(k, prefix) {
			return k
		}
	}
	return ""
}

func (m MapNode) generateJSONSchema(isrequired bool) MapString {
	res := map[string]interface{}{}
	res["type"] = "object"
	mapNode := map[string]*Node(m)
	properties := map[string]interface{}{}
	required := []string{}
	for k, v := range mapNode {
		if v.Type == "object" {
			mapNode := v.Value.(MapNode)
			schema := mapNode.generateJSONSchema(false)
			properties[k] = schema
		} else if v.Type == "array" {
			arrMapNode := v.Value.([]MapNode)
			resMapNode := ArrMapNode(arrMapNode)
			properties[k] = resMapNode.generateJSONSchema(false)
		} else {
			properties[k] = map[string]interface{}{
				"type": v.Type,
			}
		}
		required = append(required, k)
	}
	res["properties"] = properties
	if isrequired {
		res["required"] = required
	}
	return res
}

func GenerateJSONSchema(jsonStr string) (MapString, error) {
	res := MapString{}

	if !strings.HasPrefix(jsonStr, "[") {
		raw := map[string]interface{}{}
		err := json.Unmarshal([]byte(jsonStr), &raw)
		if err != nil {
			return nil, err
		}
		mapRes := parseNonArray(raw)
		res = mapRes.generateJSONSchema(true)
	} else {
		// TODO: Array Documents
		raw := []interface{}{}
		err := json.Unmarshal([]byte(jsonStr), &raw)
		if err != nil {
			return nil, err
		}
		arrMapNode := parseArray(raw)
		arrMapRes := ArrMapNode(arrMapNode)
		res = arrMapRes.generateJSONSchema(true)
	}

	return res, nil
}

func GenerateJSONSchemaFromURL(url string) (MapString, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return GenerateJSONSchema(string(body))
}

func GenerateJSONSchemaFromFile(filePath string) (MapString, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return GenerateJSONSchema(string(body))
}

func parseArray(raw []interface{}) []MapNode {
	res := []MapNode{}
	for _, item := range raw {
		contentType := getType(item)
		temp := MapNode{}

		if contentType == "object" {
			mp := item.(map[string]interface{})
			temp = parseNonArray(mp)
		} else if contentType == "array" {
			// TODO: add logic here for nested array. For example: {"data":[[1,2,4],[3,6,1]]}
			// Currently it's only set the type is array, but the items is not defined
			// For example:
			/*
				{"properties": {
						"data": {
							"type": "array"
						}
					},
				"required": ["data"],
				"type": "object"}
			*/
		} else {
			node := &Node{
				Type:  contentType,
				Value: item,
				Key:   fmt.Sprintf("%s:%v", key_generic_array_item, item),
			}
			temp[node.Key] = node
		}
		if len(temp) > 0 {
			res = append(res, temp)
		}
	}
	return res
}

func parseNonArray(raw map[string]interface{}) MapNode {
	mapNode := map[string]*Node{}
	for key, val := range raw {
		contentType := getType(val)
		if contentType == "object" {
			item := val.(map[string]interface{})
			node := &Node{
				Key:   key,
				Value: parseNonArray(item),
				Type:  contentType,
			}
			mapNode[key] = node
		} else if contentType == "array" {
			node := &Node{
				Key:   key,
				Value: parseArray(val.([]interface{})),
				Type:  contentType,
			}
			mapNode[key] = node
		} else {
			node := &Node{
				Type:  getType(val),
				Value: val,
				Key:   key,
			}
			mapNode[key] = node
		}
	}
	return mapNode
}

func getType(item interface{}) string {
	switch item.(type) {
	case string:
		return "string"
	case map[string]interface{}:
		return "object"
	case int8, int, int16, int32, int64:
		return "integer"
	case float32, float64:
		return "number"
	case []interface{}:
		return "array"
	case nil:
		return "null"
	default:
		var r = reflect.TypeOf(item)
		return fmt.Sprintf("%s", r.String())
	}
}
