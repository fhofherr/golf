package log

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Formatter converts the passed slice of interface{} to a byte slice.
//
// A logger can then use the byte slice to write it to the log. A Formatter
// may return an error to indicate it could not write the slice to the log.
type Formatter func([]interface{}) ([]byte, error)

// PlainTextFormatter converts the passed kvs to strings and joins them
// separated by a coma and a space.
func PlainTextFormatter(kvs []interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	n := len(kvs)
	if n%2 != 0 {
		n = n - 1
	}
	for i := 0; i < n; i += 2 {
		s := fmt.Sprintf("%v=%v", kvs[i], kvs[i+1])
		b.Write([]byte(s))
		if i+2 < n {
			b.Write([]byte(", "))
		}
	}
	if n < len(kvs) {
		if b.Len() > 0 {
			b.Write([]byte(", "))
		}
		s := fmt.Sprintf("%v=error: missing value", kvs[n])
		b.Write([]byte(s))
	}
	return b.Bytes(), nil
}

func JSONFormatter(kvs []interface{}) ([]byte, error) {
	var (
		entry map[string]interface{}
		n     = len(kvs)
	)

	if n%2 == 0 {
		entry = make(map[string]interface{}, n/2)
	} else {
		entry = make(map[string]interface{}, (n+1)/2)
	}

	for i := 0; i+1 < n; i += 2 {
		k, v := fmt.Sprintf("%v", kvs[i]), kvs[i+1]
		entry[k] = v
	}
	if n%2 != 0 {
		k := fmt.Sprintf("%v", kvs[n-1])
		entry[k] = "error: missing value"
	}
	bs, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("marshal log entry: %v", err)
	}
	return bs, err
}
