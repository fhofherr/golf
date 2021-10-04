package golf

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
	if len(kvs) == 0 {
		return nil, nil
	}

	b := bytes.Buffer{}
	n := len(kvs)
	if n == 1 {
		fmt.Fprintf(&b, "%s", kvs[0])
		return b.Bytes(), nil
	}
	if n%2 != 0 {
		n--
	}
	for i := 0; i < n; i += 2 {
		fmt.Fprintf(&b, "%v=%v", kvs[i], kvs[i+1])
		if i+2 < n {
			fmt.Fprint(&b, ", ")
		}
	}
	if n < len(kvs) {
		if b.Len() > 0 {
			fmt.Fprint(&b, ", ")
		}
		fmt.Fprintf(&b, "%v=error: missing value", kvs[n])
	}
	fmt.Fprint(&b, "\n")

	return b.Bytes(), nil
}

// JSONFormatter marshals the passed key-value pairs into a JSON object.
func JSONFormatter(kvs []interface{}) ([]byte, error) {
	var (
		n = len(kvs)
	)

	entry := make(map[string]interface{}, (n+1)/2)
	if n == 1 {
		entry["message"] = kvs[0]
	} else {
		for i := 0; i+1 < n; i += 2 {
			k, v := fmt.Sprintf("%v", kvs[i]), kvs[i+1]
			entry[k] = v
		}
		if n%2 != 0 {
			k := fmt.Sprintf("%v", kvs[n-1])
			entry[k] = "error: missing value"
		}
	}

	bs, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("marshal log entry: %v", err)
	}
	return bs, err
}
