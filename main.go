package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func composeJSON(keys []string, values []string) interface{} {
	result := make(map[string]interface{})
	for i, key := range keys {
		value := getDefaultValue(values[i])
		parts := strings.Split(key, ".")
		current := result
		for j, part := range parts[:len(parts)-1] {

			// append to list
			if strings.HasSuffix(part, "[]") {
				arrKey := strings.TrimSuffix(part, "[]")
				if _, ok := current[arrKey]; !ok {
					current[arrKey] = make([]map[string]interface{}, 0)
					
				}
				if arr, ok := current[arrKey].([]map[string]interface{}); ok {
					next := make(map[string]interface{})
					current[arrKey] = append(arr, next)
					current = next
				} else {
					panic(fmt.Sprintf("expected an array, but found %T", current[arrKey]))
				}
			} else if strings.HasSuffix(part, "[-1]") {
				arrKey := strings.TrimSuffix(part, "[-1]")
				if arr, ok := current[arrKey].([]map[string]interface{}); ok {				
					current = arr[len(arr)- 1 ]
				} else {
					panic(fmt.Sprintf("expected an array, but found %T", current[arrKey]))
				}
			} else {
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]interface{})
				}
				if m, ok := current[part].(map[string]interface{}); ok {
					current = m
				} else {
					panic(fmt.Sprintf("expected a map, but found %T", current[part]))
				}
			}
			if j == len(parts)-2 {
				lastPart := parts[len(parts)-1]
				if strings.HasSuffix(lastPart, "[]") {
					arrKey := strings.TrimSuffix(lastPart, "[]")
					if _, ok := current[arrKey]; !ok {
						current[arrKey] = make([]interface{}, 0)
					}
					if arr, ok := current[arrKey].([]interface{}); ok {
						current[arrKey] = append(arr, value)
					} else {
						panic(fmt.Sprintf("expected an array, but found %T", current[arrKey]))
					}
				} else {
					current[lastPart] = value
				}
			}
		}
	}
	return result
}

func getDefaultValue(value string) interface{} {
	if value == "null" {
		return nil
	}
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return value
}

func main() {
	var keys []string
	var values []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pair := strings.Split(scanner.Text(), "=")
		keys = append(keys, pair[0])
		values = append(values, pair[1])
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	jsonDocument := composeJSON(keys, values)
	jsonBytes, err := json.Marshal(jsonDocument)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
