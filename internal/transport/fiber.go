package transport

import (
	"reflect"
)

// Fiber v2 and v3 satisfy this small structural interface. Status() has
// different return types, so reflection is limited to that method alone.
// No Fiber dependency or runtime plugins are required in the core module.
type fiberLike interface {
	Get(string, ...string) string
	Set(string, string)
	Send([]byte) error
}

type FiberAdapter struct {
	C   fiberLike
	Raw any
}

func (a FiberAdapter) Header(key string) string {
	return a.C.Get(key)
}

func (a FiberAdapter) SetHeader(key, value string) {
	a.C.Set(key, value)
}

func (a FiberAdapter) AddVary(v string) {
	old := ""
	if c, ok := a.Raw.(interface {
		GetRespHeader(string, ...string) string
	}); ok {
		old = c.GetRespHeader("Vary")
	}
	if old == "" {
		a.C.Set("Vary", v)
	} else if !HasVary(old, v) {
		a.C.Set("Vary", old+", "+v)
	}
}

func (a FiberAdapter) Method() string {
	m := reflect.ValueOf(a.Raw).MethodByName("Method")
	if !m.IsValid() {
		return ""
	}
	t := m.Type()
	if t.NumOut() != 1 || t.Out(0).Kind() != reflect.String || t.NumIn() > 1 {
		return ""
	}
	// Fiber v2 Method(...string) and v3 Method() both work with no args.
	out := m.Call(nil)
	return out[0].String()
}

func (a FiberAdapter) Write(code int, body []byte, ctype string) error {
	m := reflect.ValueOf(a.Raw).MethodByName("Status")
	if !m.IsValid() || m.Type().NumIn() != 1 || m.Type().In(0).Kind() != reflect.Int {
		return ErrUnsupportedContext
	}
	m.Call([]reflect.Value{reflect.ValueOf(code)})
	if ctype != "" {
		a.C.Set("Content-Type", ctype)
	}
	return a.C.Send(body)
}
