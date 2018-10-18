package slackbot

import (
	"testing"
)

func TestClientUserNameEmpty(t *testing.T) {
	c := Client{Users: map[string]string{}}
	if got := c.UserName(""); got != "" {
		t.Errorf("got %v, expected empty", got)
	}
}

func TestClientUserName(t *testing.T) {
	m := map[string]string{
		"U03CP354N": "foo",
		"U02J1PU37": "baa",
	}
	c := Client{Users: m}
	for id, expected := range m {
		if got := c.UserName(id); got != expected {
			t.Errorf("got %v, expected %v", got, expected)
		}
	}
}
