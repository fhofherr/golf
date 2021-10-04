package golfstdlib_test

import (
	"testing"

	"github.com/fhofherr/golf/golfstdlib"
	"github.com/stretchr/testify/assert"
)

var formatterTests = []formatterTest{
	{
		name: "empty kvs",
		kvs:  []interface{}{},
	},
	{
		name: "single value",
		kvs:  []interface{}{"just a message"},
	},
	{
		name: "single key value pair",
		kvs:  []interface{}{"key", "value"},
	},
	{
		name: "two key value pairs",
		kvs:  []interface{}{"key1", "value1", "key2", "value2"},
	},
	{
		name: "missing value",
		kvs:  []interface{}{"key", "value", "other-key"},
	},
	{
		name: "odd number of kvs",
		kvs:  []interface{}{"key1", "value1", "key2"},
	},
}

func TestPlainTextFormatter(t *testing.T) {
	expectations := map[string]string{
		"empty kvs":             "",
		"single value":          "just a message",
		"single key value pair": "key=value\n",
		"two key value pairs":   "key1=value1, key2=value2\n",
		"missing value":         "key=value, other-key=error: missing value\n",
		"odd number of kvs":     "key1=value1, key2=error: missing value\n",
	}
	exerciseFormatter(t, golfstdlib.PlainTextFormatter, expectations, func(t *testing.T, expected, actual string) {
		assert.Equal(t, expected, actual)
	})
}

func TestJSONFormatter(t *testing.T) {
	expectations := map[string]string{
		"empty kvs":             "{}",
		"single value":          `{"message": "just a message"}`,
		"single key value pair": `{"key": "value"}`,
		"two key value pairs":   `{"key1": "value1", "key2": "value2"}`,
		"missing value":         `{"key": "value", "other-key": "error: missing value"}`,
		"odd number of kvs":     `{"key1": "value1", "key2": "error: missing value"}`,
	}
	exerciseFormatter(t, golfstdlib.JSONFormatter, expectations, func(t *testing.T, expected, actual string) {
		assert.JSONEq(t, expected, actual)
	})
}

func TestJSONFormatter_MarshallingError(t *testing.T) {
	_, err := golfstdlib.JSONFormatter([]interface{}{"key", golfstdlib.JSONFormatter})
	assert.Error(t, err)
}

type formatterTest struct {
	name string
	kvs  []interface{}
}

type formatterAssertion func(*testing.T, string, string)

func exerciseFormatter(
	t *testing.T, formatter golfstdlib.Formatter, expectations map[string]string, assertion formatterAssertion,
) {
	for _, tt := range formatterTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			expected, ok := expectations[tt.name]
			if !ok {
				t.Fatalf("missing expectation: %s", tt.name)
			}
			bs, err := formatter(tt.kvs)
			assert.NoError(t, err)
			assertion(t, expected, string(bs))
		})
	}
}
