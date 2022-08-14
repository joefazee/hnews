package models

import (
	"errors"
	"github.com/upper/db/v4"
	"testing"
)

func Test_convertUpperIDtoInt(t *testing.T) {
	type args struct {
		id db.ID
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"upper db.ID to int", args{id: db.ID(1)}, 1},
		{"upper db.ID to int", args{id: int64(1)}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertUpperIDtoInt(tt.args.id); got != tt.want {
				t.Errorf("convertUpperIDtoInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errHasDuplicate(t *testing.T) {
	type args struct {
		err error
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ERROR: duplicate key value", args{err: errors.New(`ERROR: duplicate key value violates unique constraint "users_email_key"`), key: "users_email_key"}, true},
		{"key not in string", args{err: errors.New(`regular error without the key`), key: "some random keu=y"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errHasDuplicate(tt.args.err, tt.args.key); got != tt.want {
				t.Errorf("errHasDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
