package filter

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	type args struct {
		input   io.Reader
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "nil input",
			args: args{
				input:   nil,
				pattern: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "two exact matches",
			args: args{
				input:   strings.NewReader("foo\nbar\nbaz\n"),
				pattern: "ba",
			},
			want:    []string{"bar", "baz"},
			wantErr: false,
		},
		{
			name: "two regexp matches",
			args: args{
				input:   strings.NewReader("foo\nbar\nbaz\n"),
				pattern: "a.+",
			},
			want:    []string{"bar", "baz"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Grep(tt.args.input, tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

// benchmark Grep
func BenchmarkGrepExact(b *testing.B) {
	input := strings.NewReader(text)
	for i := 0; i < b.N; i++ {
		Grep(input, "ci")
	}
}

func BenchmarkGrepRegexp(b *testing.B) {
	input := strings.NewReader(text)
	for i := 0; i < b.N; i++ {
		Grep(input, `\w*ci\w+`)
	}
}

func BenchmarkMatch(b *testing.B) {
	input := strings.NewReader(text)
	for i := 0; i < b.N; i++ {
		Match(input, "ci")
	}
}

const text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, 
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris 
nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in 
reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla 
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in 
culpa qui officia deserunt mollit anim id est laborum.`
