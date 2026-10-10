package resp

import (
	"errors"
	"testing"
)

func TestNewReturnsTranslationLoaderError(t *testing.T) {
	wantErr := errors.New("catalog unavailable")
	engine, err := newEngine(nil, func() (map[string]map[string]string, error) {
		return nil, wantErr
	})
	if engine != nil {
		t.Fatalf("newEngine() engine = %#v, want nil", engine)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("newEngine() error = %v, want wrapped %v", err, wantErr)
	}
}

func TestMustNewEnginePanicsOnInitializationError(t *testing.T) {
	wantErr := errors.New("catalog unavailable")
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("mustNewEngine() did not panic")
		}
		if !errors.Is(got.(error), wantErr) {
			t.Fatalf("mustNewEngine() panic = %v, want wrapped %v", got, wantErr)
		}
	}()

	mustNewEngine(func(...Option) (*Engine, error) {
		return nil, wantErr
	})
}
