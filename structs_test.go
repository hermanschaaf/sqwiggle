package sqwiggle

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	type dog struct {
		Name     string
		Password string
	}

	type house struct {
		Name string
		Age  int
	}

	// table-driven test of different structs
	tests := []struct {
		name    string
		a       interface{}
		b       interface{}
		want    difference
		wantErr error
	}{
		{
			name:    "nil pointers",
			a:       nil,
			b:       nil,
			want:    map[string]attr{},
			wantErr: errNilInterface,
		},
		{
			name:    "nil pointer vs struct",
			a:       nil,
			b:       dog{"woof", "wooooofle"},
			want:    map[string]attr{},
			wantErr: errNilInterface,
		},
		{
			name:    "equal empty string structs",
			a:       dog{},
			b:       dog{},
			want:    difference{},
			wantErr: nil,
		},
		{
			name:    "equal structs",
			a:       dog{"doge", "ilovewoof"},
			b:       dog{"doge", "ilovewoof"},
			want:    map[string]attr{},
			wantErr: nil,
		},
		{
			name:    "unequal structs",
			a:       dog{"doge", "ilovewoof"},
			b:       dog{"doge", "ihatewoof"},
			want:    map[string]attr{"Password": {"ilovewoof", "ihatewoof"}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		got, err := compare(tt.a, tt.b)
		if err != tt.wantErr {
			t.Fatalf("%q case: err = %v, want %v", tt.name, err, tt.wantErr)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q case: got %#v, want %#v", tt.name, got, tt.want)
		}
	}
}
