package goson

import (
	"encoding/json"
	"fmt"
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

func (m MapNode) generateJSONSchema() MapString {
	res := map[string]interface{}{}
	res["type"] = "object"
	mapNode := map[string]*Node(m)
	properties := map[string]interface{}{}
	for k, v := range mapNode {
		if v.Type == "object" {
			mapNode := v.Value.(MapNode)
			schema := mapNode.generateJSONSchema()
			properties[k] = schema
		} else if v.Type == "array" {
			itemsSchema := map[string]interface{}{}
			arrMapNode := v.Value.([]MapNode)
			if len(arrMapNode) > 0 {
				schema := map[string]interface{}{}
				if key := getKeymapPrefix(arrMapNode[0], key_generic_array_item); key != "" {
					// INDICATE GENERIC ARRAY-ITEMS
					node := arrMapNode[0][key]
					itemsSchema := map[string]interface{}{
						"type": node.Type,
					}
					schema = map[string]interface{}{
						"type":  "array",
						"items": itemsSchema,
					}
				} else {
					itemsSchema = arrMapNode[0].generateJSONSchema()
					fmt.Printf("MAPNODE: %+v\n", arrMapNode[0])

					schema = map[string]interface{}{
						"type":  v.Type,
						"items": itemsSchema,
					}
				}

				properties[k] = schema
			}
		} else {
			properties[k] = map[string]interface{}{
				"type": v.Type,
			}
		}
	}
	res["properties"] = properties
	return res
}

func GenerateJSONSchema(jsonStr string) (MapString, error) {
	raw := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &raw)
	if err != nil {
		return nil, err
	}
	mapRes := parseNonArray(raw)
	for k, v := range mapRes {
		fmt.Printf("key: %v, val: %+v \n", k, v)
	}
	return mapRes.generateJSONSchema(), nil
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
			// TODO
		} else {
			//  TODO
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
			// TODO If array is exists
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
	default:
		var r = reflect.TypeOf(item)
		fmt.Printf("Other:%v\n", r)
	}
	return ""
}
