package log

import "fmt"

// GenerateKEYVALs generates a fixed amount of key value pairs.
func GenerateKEYVALs(n int) []interface{} {
	var kvs []interface{}
	for i := 0; i < n; i++ {
		kvs = append(kvs, fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
	return kvs
}
