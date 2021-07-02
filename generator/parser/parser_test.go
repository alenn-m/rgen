package parser

import (
	"testing"

	"github.com/alenn-m/rgen/util/misc"
	"github.com/stretchr/testify/assert"
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

func TestParser_ParseActions__Success(t *testing.T) {
	a := assert.New(t)

	p := new(Parser)
	err := p.Parse("User", "first_name:string", "index, show, delete")
	a.Nil(err)
	a.Len(p.Actions, 3)
	a.Equal([]string{misc.ACTION_INDEX, misc.ACTION_SHOW, misc.ACTION_DELETE}, p.Actions)
}

func TestParser_ParseActions__WrongAction(t *testing.T) {
	a := assert.New(t)

	p := new(Parser)
	err := p.Parse("User", "first_name:string", "index, show, wrong_action")
	a.NotNil(err)
	a.Equal(ErrInvalidAction, err)
}
