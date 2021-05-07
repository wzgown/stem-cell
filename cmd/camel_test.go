package cmd

import (
	"testing"
)

func Test_Camel(t *testing.T) {
	out := camel("hello_world_ni_hao")
	expect := "HelloWorldNiHao"
	if out != expect {
		t.Error(out, "; expect==>  ", expect)
	}
	t.Log(out)
}
