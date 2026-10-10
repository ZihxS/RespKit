package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ZihxS/RespKit/internal/formatter"
	"github.com/ZihxS/RespKit/internal/transport"
)

func (e *Engine) field(name string) string {
	return formatter.Field(name, string(e.naming))
}

func safeRequestID(id string) string {
	return formatter.SafeRequestID(id)
}

func (e *Engine) send(ctx any, status int, data any, failure *Definition, opts ...ResponseOption) error {
	out, err := transport.Adapt(ctx)
	if err != nil {
		return err
	}
	cfg := configFor(opts)
	if status == http.StatusNoContent {
		return out.Write(status, nil, "")
	}
	language := e.effectiveLocale(e.chooseLocale(out.Header("Accept-Language"), cfg.locale), failure)
	requestID := safeRequestID(out.Header("X-Request-ID"))
	if cfg.location != "" && status == 201 {
		out.SetHeader("Location", cfg.location)
	}
	out.SetHeader("Content-Language", language)
	out.AddVary("Accept-Language")

	mediaType := "application/json; charset=utf-8"
	var body map[string]any
	if failure == nil {
		body = formatter.Success(data, string(e.naming))
	} else {
		for key, values := range failure.Headers {
			if len(values) > 0 {
				out.SetHeader(key, values[0])
			}
		}
		// HTTP 401 requires an authentication challenge. The default Bearer
		// challenge can be overridden through WithHeader on an AppError.
		if status == http.StatusUnauthorized && (failure.Headers == nil || failure.Headers.Get("WWW-Authenticate") == "") {
			out.SetHeader("WWW-Authenticate", "Bearer")
		}
		message := e.publicMessage(language, failure.Code, status)
		if e.format == ProblemDetails {
			mediaType = "application/problem+json"
			debug := ""
			if e.exposeInternal {
				debug = cfg.internal
			}
			body = formatter.Problem(status, http.StatusText(status), failure.Code, message, debug)
			if requestID != "" {
				body["request_id"] = requestID
			}

		} else {
			detailRows := make([]formatter.Detail, 0, len(failure.Details))
			for _, d := range failure.Details {
				params := map[string]string{"field": d.Field}
				for k, v := range d.Params {
					if len(v) > 256 || k == "field" {
						continue
					}
					params[k] = v
				}
				detailRows = append(detailRows, formatter.Detail{
					Field:   d.Field,
					Code:    d.Code,
					Message: e.fieldMessage(language, d.Code, params),
				})
			}
			debug := ""
			if e.exposeInternal {
				debug = cfg.internal
			}
			body = formatter.Failure(failure.Code, message, debug, detailRows, string(e.naming))
		}
	}
	if requestID != "" && (failure == nil || e.format != ProblemDetails) {
		body[e.field("meta")] = map[string]any{e.field("request_id"): requestID}
	}
	buf, err := json.Marshal(body)
	if err != nil {
		// Handle failed success serialization BEFORE writing a status/body.
		// Do not attempt to marshal unsafe data again.
		return e.writeMarshalFallback(out, language, requestID, err)
	}
	if strings.EqualFold(out.Method(), http.MethodHead) {
		buf = nil
	}
	if err := out.Write(status, buf, mediaType); err != nil {
		return fmt.Errorf("resp: response write: %w", err)
	}
	return nil
}

func (e *Engine) writeMarshalFallback(out transport.Adapter, locale, requestID string, cause error) error {
	debug := ""
	if e.exposeInternal && cause != nil {
		debug = cause.Error()
	}
	body := formatter.Failure("INTERNAL_ERROR", e.translate(locale, "INTERNAL_ERROR", nil), debug, nil, string(e.naming))
	mediaType := "application/json; charset=utf-8"
	if e.format == ProblemDetails {
		// RFC 9457 defines errors, including serialization failures.
		mediaType = "application/problem+json"
		body = formatter.Problem(
			500,
			"Internal Server Error",
			"INTERNAL_ERROR",
			e.translate(locale, "INTERNAL_ERROR", nil),
			debug,
		)
	}
	if requestID != "" {
		if e.format == ProblemDetails {
			body["request_id"] = requestID
		} else {
			body[e.field("meta")] = map[string]any{e.field("request_id"): requestID}
		}
	}
	b, _ := json.Marshal(body)
	if strings.EqualFold(out.Method(), http.MethodHead) {
		b = nil
	}
	return out.Write(http.StatusInternalServerError, b, mediaType)
}

// publicMessage uses a generic status-appropriate message for custom error
// codes without a registered translation. Raw errors are never serialized.
func (e *Engine) publicMessage(locale, code string, status int) string {
	if e.hasTranslation(locale, code) || e.hasTranslation(e.defaultLocale, code) {
		return e.translate(locale, code, nil)
	}
	fallbackCodes := map[int]string{
		400: "BAD_REQUEST",
		401: "UNAUTHORIZED",
		403: "FORBIDDEN",
		404: "RESOURCE_NOT_FOUND",
		409: "RESOURCE_CONFLICT",
		422: "VALIDATION_FAILED",
		429: "RATE_LIMITED",
		503: "SERVICE_UNAVAILABLE",
	}
	fallback := fallbackCodes[status]
	if fallback == "" {
		fallback = "INTERNAL_ERROR"
	}
	return e.translate(locale, fallback, nil)
}

func (e *Engine) fieldMessage(locale, code string, params map[string]string) string {
	if !e.hasTranslation(locale, code) && !e.hasTranslation(e.defaultLocale, code) {
		return e.translate(locale, "VALIDATION_FAILED", nil)
	}
	return e.translate(locale, code, params)
}
