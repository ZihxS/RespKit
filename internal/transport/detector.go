package transport

import (
	"net/http"
	"reflect"
)

func Adapt(c any) (Adapter, error) {
	if c == nil {
		return nil, ErrUnsupportedContext
	}
	v := reflect.ValueOf(c)
	if (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) && v.IsNil() {
		return nil, ErrUnsupportedContext
	}
	switch x := c.(type) {
	case *HTTPContext:
		if x == nil || x.Writer == nil || x.Request == nil {
			return nil, ErrUnsupportedContext
		}
		return netAdapter{w: x.Writer, r: x.Request}, nil
	case ginLike:
		return ginAdapter{c: x, raw: c}, nil
	case echoLike:
		if x.Request() == nil {
			return nil, ErrUnsupportedContext
		}
		if _, ok := EchoResponseHeader(c); !ok {
			return nil, ErrUnsupportedContext
		}
		return echoAdapter{c: x, raw: c}, nil
	case fiberLike:
		if !validFiberStatus(c) {
			return nil, ErrUnsupportedContext
		}
		return FiberAdapter{C: x, Raw: c}, nil
	default:
		return nil, ErrUnsupportedContext
	}
}

func validFiberStatus(c any) bool {
	v := reflect.ValueOf(c).MethodByName("Status")
	return v.IsValid() && v.Type().NumIn() == 1 && v.Type().In(0).Kind() == reflect.Int
}

func GinWriterHeader(c any) (http.Header, bool) {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, false
	}
	w := v.FieldByName("Writer")
	if !w.IsValid() || !w.CanInterface() {
		return nil, false
	}
	if writer, ok := w.Interface().(http.ResponseWriter); ok {
		return writer.Header(), true
	}
	return nil, false
}

func GinRequest(c any) (*http.Request, bool) {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, false
	}
	r := v.FieldByName("Request")
	if r.IsValid() && r.CanInterface() {
		req, ok := r.Interface().(*http.Request)
		return req, ok
	}
	return nil, false
}

func EchoResponseHeader(c any) (http.Header, bool) {
	if c == nil {
		return nil, false
	}
	v := reflect.ValueOf(c)
	m := v.MethodByName("Response")
	if !m.IsValid() || m.Type().NumIn() != 0 || m.Type().NumOut() != 1 {
		return nil, false
	}
	r := m.Call(nil)[0]
	if (r.Kind() == reflect.Pointer || r.Kind() == reflect.Interface) && r.IsNil() {
		return nil, false
	}
	h := r.MethodByName("Header")
	if !h.IsValid() || h.Type().NumIn() != 0 || h.Type().NumOut() != 1 {
		return nil, false
	}
	result := h.Call(nil)[0]
	header, ok := result.Interface().(http.Header)
	return header, ok
}
