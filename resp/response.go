package resp

import (
	"context"
	"errors"
	"net/http"
)

// ResponseOption applies only to the current response. Locale overrides
// Accept-Language; Location is useful with 201 Created.
type ResponseOption func(*responseConfig)
type responseConfig struct {
	locale   string
	location string
	internal string
}

func Locale(lang string) ResponseOption {
	return func(o *responseConfig) {
		o.locale = lang
	}
}

func Location(url string) ResponseOption {
	return func(o *responseConfig) {
		o.location = url
	}
}

func OK(ctx any, data any, opts ...ResponseOption) error {
	return current().OK(ctx, data, opts...)
}

func Created(ctx any, data any, opts ...ResponseOption) error {
	return current().Created(ctx, data, opts...)
}

func NoContent(ctx any) error {
	return current().NoContent(ctx)
}

func Fail(ctx any, err error, opts ...ResponseOption) error {
	return current().Fail(ctx, err, opts...)
}

func NotFound(ctx any, opts ...ResponseOption) error {
	return current().Fail(ctx, ErrNotFound, opts...)
}

func (e *Engine) OK(ctx any, data any, opts ...ResponseOption) error {
	return e.send(ctx, http.StatusOK, data, nil, opts...)
}

func (e *Engine) Created(ctx any, data any, opts ...ResponseOption) error {
	return e.send(ctx, http.StatusCreated, data, nil, opts...)
}

func (e *Engine) NoContent(ctx any) error {
	return e.send(ctx, http.StatusNoContent, nil, nil)
}

func (e *Engine) Fail(ctx any, err error, opts ...ResponseOption) error {
	if errors.Is(err, context.Canceled) {
		return ErrRequestCanceled
	}
	def, _ := e.Resolve(err)
	return e.send(ctx, def.Status, nil, &def, optsWithInternal(err, e.exposeInternal, opts)...)
}

// optsWithInternal is package-private so callers cannot accidentally enable
// diagnostic disclosure per request. Only Engine option controls it.
func optsWithInternal(err error, enabled bool, opts []ResponseOption) []ResponseOption {
	if !enabled || err == nil {
		return opts
	}
	result := make([]ResponseOption, 0, len(opts)+1)
	result = append(result, opts...)
	result = append(result, func(o *responseConfig) {
		o.internal = err.Error()
		if len(o.internal) > 8192 {
			o.internal = o.internal[:8192]
		}
	})
	return result
}

func configFor(opts []ResponseOption) responseConfig {
	var cfg responseConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}
