package utils

import (
	"errors"
	"reflect"
	"testing"
)

func TestUtils_TemplateToString(t *testing.T) { // TODO
	validFile := []byte("valid")
	type fields struct{}
	type args struct {
		name string
	}
	type mockArgs struct {
		obj     string
		method  string
		args    []interface{}
		returns []interface{}
	}
	testCases := []struct {
		name     string
		fields   fields
		args     args
		want     []byte
		wantErr  bool
		mockArgs []mockArgs
	}{
		{
			name: "Should return a valid file when file are read successfully",
			args: args{
				name: "anyfile",
			},
			mockArgs: []mockArgs{
				{
					obj:     "os",
					method:  "ReadFile",
					args:    []interface{}{"anyfile"},
					returns: []interface{}{validFile, nil},
				},
			},
			want:    validFile,
			wantErr: false,
		},
		{
			name: "Should return an error if reading the file fails",
			args: args{
				name: "anyfile",
			},
			mockArgs: []mockArgs{
				{
					obj:     "os",
					method:  "ReadFile",
					args:    []interface{}{"anyfile"},
					returns: []interface{}{[]byte{}, errors.New("read file error")},
				},
			},
			want:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := Utils{}
			for _, mockArgs := range tt.mockArgs {
				switch mockArgs.obj {
				case "os":
					osReadFile = func(_ string) ([]byte, error) {
						err, _ := mockArgs.returns[1].(error)
						return mockArgs.returns[0].([]byte), err
					}
				}
			}
			got, err := u.ReadFile(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Utils.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
