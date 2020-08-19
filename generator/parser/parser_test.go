package parser

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	type fields struct {
		Name   string
		Fields []Field
	}
	type args struct {
		name   string
		fields string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "name should not be empty",
			args:    args{},
			wantErr: true,
		},
		{
			name:    "name is not empty and fields can't be empty",
			args:    args{name: "some name", fields: ""},
			wantErr: true,
		},
		{
			name: "name is not empty; fields are not correct",
			args: args{
				name:   "some name",
				fields: "key1",
			},
			wantErr: true,
		},
		{
			name: "name is not empty and fields are correct",
			args: args{
				name:   "some name",
				fields: "key:value,key1:value1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Parser{
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			if err := c.Parse(tt.args.name, tt.args.fields, ""); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
