package info

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/mangalorg/libmangal"
)

func TestNew(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantInfo Info
		wantErr  bool
	}{
		{
			name: "bundle",
			args: args{
				r: strings.NewReader(`
type = "lua"
id = "some-id"
`),
			},
			wantInfo: Info{
				ProviderInfo: libmangal.ProviderInfo{
					ID: "some-id",
				},
				Type: TypeLua,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := New(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("New() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}
