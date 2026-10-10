// Package formatter constructs the stable public envelope independently of
// framework adapters, error detection, and language resolution.
package formatter

import "strings"

// Field converts RespKit-owned response keys; application data is untouched.
func Field(name, naming string) string {
	if naming == "snake_case" {
		return name
	}
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if part == "" || (i == 0 && naming == "camelCase") {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

// SafeRequestID allows only bounded, plain ASCII identifiers in API payloads.
func SafeRequestID(id string) string {
	if id == "" || len(id) > 128 {
		return ""
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' && r != '-' && r != '.' {
			return ""
		}
	}
	return id
}

type Detail struct{ Field, Code, Message string }

func Success(data any, naming string) map[string]any {
	return map[string]any{Field("success", naming): true, Field("data", naming): data}
}

func Failure(code, message, debug string, details []Detail, naming string) map[string]any {
	item := map[string]any{Field("code", naming): code, Field("message", naming): message}
	if debug != "" {
		item[Field("debug_message", naming)] = debug
	}
	if len(details) > 0 {
		rows := make([]map[string]any, 0, len(details))
		for _, d := range details {
			rows = append(rows, map[string]any{
				Field("field", naming):   d.Field,
				Field("code", naming):    d.Code,
				Field("message", naming): d.Message,
			})
		}
		item[Field("details", naming)] = rows
	}
	return map[string]any{Field("success", naming): false, Field("error", naming): item}
}

func Problem(status int, title, code, message, debug string) map[string]any {
	body := map[string]any{"type": "about:blank", "title": title, "status": status, "detail": message, "code": code}
	if debug != "" {
		body["debug_message"] = debug
	}
	return body
}
