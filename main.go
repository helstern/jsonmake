package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func composeJSON(jsonPaths []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, path := range jsonPaths {
		parts := strings.Split(path, ".")
		current := result
		for i, part := range parts[:len(parts)-1] {
			if strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
				indexStr := part[1 : len(part)-1]
				var index int
				if indexStr != "" {
					var err error
					index, err = strconv.Atoi(indexStr)
					if err != nil {
						panic(err)
					}
				} else {
					// Append to the array if index is not specified
					index = len(current)
				}
				if i == len(parts)-2 {
					current[index] = getDefaultValue(parts[len(parts)-1])
				} else {
					if _, ok := current[index]; !ok {
						current[index] = make(map[string]interface{})
					}
					current = current[index].(map[string]interface{})
				}
			} else {
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
		lastPart := parts[len(parts)-1]
		if strings.HasPrefix(lastPart, "[") && strings.HasSuffix(lastPart, "]") {
			indexStr := lastPart[1 : len(lastPart)-1]
			var index int
			if indexStr != "" {
				var err error
				index, err = strconv.Atoi(indexStr)
				if err != nil {
					panic(err)
				}
			} else {
				// Append to the array if index is not specified
				index = len(current)
			}
			current[index] = getDefaultValue(lastPart)
		} else {
			current[lastPart] = getDefaultValue(lastPart)
		}
	}
	return result
}

func getDefaultValue(part string) interface{} {
	if part == "null" {
		return nil
	}
	if i, err := strconv.Atoi(part); err == nil {
		return i
	}
	return part
}

func main() {
	var paths []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		paths = append(paths, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	jsonDocument := composeJSON(paths)
	jsonBytes, err := json.Marshal(jsonDocument)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
