package fileutil

import (
	"fmt"
	"testing"
)

func TestListAllFilePathInDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "/home/lyr",
		},
	}
	for _, tt := range tests {
		fmt.Println("for path", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got := ListAllFilePathInDir(tt.args.dir)
			{
				fmt.Println(got)
			}
		})
	}
}
