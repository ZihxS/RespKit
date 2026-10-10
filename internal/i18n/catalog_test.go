package i18n

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
)

type catalogTestFS struct {
	files     map[string][]byte
	filenames []string
	globErr   error
	readErr   error
}

func (f catalogTestFS) Open(string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (f catalogTestFS) Glob(string) ([]string, error) {
	return f.filenames, f.globErr
}

func (f catalogTestFS) ReadFile(name string) ([]byte, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}
	data, ok := f.files[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return data, nil
}

func TestLoadTranslations(t *testing.T) {
	tests := []struct {
		name    string
		fs      catalogTestFS
		want    map[string]map[string]string
		wantErr bool
	}{
		{
			name:    "glob error",
			fs:      catalogTestFS{globErr: errors.New("glob failed")},
			wantErr: true,
		},
		{
			name:    "no catalogs",
			fs:      catalogTestFS{},
			wantErr: true,
		},
		{
			name:    "read error",
			fs:      catalogTestFS{filenames: []string{"locales/en.json"}, readErr: errors.New("read failed")},
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			fs:      catalogTestFS{filenames: []string{"locales/en.json"}, files: map[string][]byte{"locales/en.json": []byte("{")}},
			wantErr: true,
		},
		{
			name:    "empty catalog",
			fs:      catalogTestFS{filenames: []string{"locales/en.json"}, files: map[string][]byte{"locales/en.json": []byte("{}")}},
			wantErr: true,
		},
		{
			name:    "blank code",
			fs:      catalogTestFS{filenames: []string{"locales/en.json"}, files: map[string][]byte{"locales/en.json": []byte(`{" ":"Message"}`)}},
			wantErr: true,
		},
		{
			name:    "blank message",
			fs:      catalogTestFS{filenames: []string{"locales/en.json"}, files: map[string][]byte{"locales/en.json": []byte(`{"CODE":" "}`)}},
			wantErr: true,
		},
		{
			name: "valid catalog",
			fs: catalogTestFS{
				filenames: []string{"locales/en.json"},
				files:     map[string][]byte{"locales/en.json": []byte(`{"CODE":"Message"}`)},
			},
			want: map[string]map[string]string{"en": {"CODE": "Message"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadTranslations(tt.fs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadTranslations() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("loadTranslations() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
