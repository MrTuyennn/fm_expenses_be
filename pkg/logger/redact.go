package logger

import (
	"encoding/json"
	"strings"
)

// sensitiveKeys lists JSON field names (case-insensitive) whose value must
// never be printed to logs in full.
var sensitiveKeys = map[string]bool{
	"password":      true,
	"access_token":  true,
	"refresh_token": true,
	"token_hash":    true,
	"password_hash": true,
}

// redactJSON parses body as JSON and replaces the value of any sensitive
// key (anywhere in the structure, including nested objects/arrays) with
// "***". If body isn't valid JSON, it's returned unchanged — this keeps
// plain-text bodies (e.g. the "/health" response) printable too.
func redactJSON(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return string(body)
	}

	redactValue(v)

	// 2-space indent, matching JSON.stringify(v, null, 2).
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(body)
	}
	return string(out)
}

func redactValue(v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, vv := range val {
			if sensitiveKeys[strings.ToLower(k)] {
				val[k] = "***"
				continue
			}
			redactValue(vv)
		}
	case []interface{}:
		for _, item := range val {
			redactValue(item)
		}
	}
}
