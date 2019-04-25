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

func TestPlainText(t *testing.T) {
	c := Client{Users: map[string]string{
		"U03CP354N": "user",
		"U02J1PU37": "group",
	}}
	s := `&lt;123<@U03CP354N>456<subgroup!hogehoge|@U02J1PU37>789&gt;`
	if got, expected := c.PlainMessageText(s), "<123@user456@group789>"; got !=expected {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
