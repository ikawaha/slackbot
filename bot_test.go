package slackbot

import (
	"testing"
)

const (
	slackAPIPath = "/api/rtm.start"
)

func TestBotUserNameEmpty(t *testing.T) {
	b := Bot{Users: map[string]string{}}
	if n := b.UserName(""); n != "" {
		t.Errorf("got %v, expected empty", n)
	}
}

func TestBotUserName(t *testing.T) {
	m := map[string]string{
		"U03CP354N": "foo",
		"U02J1PU37": "baa",
	}
	b := Bot{Users: m}
	for id, user := range m {
		if u := b.UserName(id); u != user {
			t.Errorf("got %v, expected empty", u, user)
		}
	}
}
func TestMessageTextBody(t *testing.T) {
	s := map[string]string{
		"<@U0ATYM9UZ1>:aaaa":       ":aaaa",
		"<@U0ATYM9UZ2>: bbbb":      "bbbb",
		"<@U0ATYM9UZ2>:      cccc": "cccc",
		"<@U0ATYM9UZ3>: ":          "",
		"<@U0ATYM9UZ4> ":           "",
		"<@U0ATYM9UZ5>:":           ":",
		"dddd":                     "dddd",
		"":                         "",
	}
	var m Message
	for txt, ans := range s {
		m.Text = txt
		b := m.TextBody()
		if b != ans {
			t.Errorf("input: %v, got %v, expected %v", txt, b, ans)
		}
	}
}

func TestMessageMentionID(t *testing.T) {
	s := map[string]string{
		"<@U0ATYM9UZ1>:aaaa":       "U0ATYM9UZ1",
		"<@U0ATYM9UZ2>: bbbb":      "U0ATYM9UZ2",
		"<@U0ATYM9UZ3>:      cccc": "U0ATYM9UZ3",
		"<@U0ATYM9UZ4>: ":          "U0ATYM9UZ4",
		"<@U0ATYM9UZ5> ":           "U0ATYM9UZ5",
		"<@U0ATYM9UZ6>:":           "U0ATYM9UZ6",
		"<@U0ATYM9UZ7|piyo>:":      "U0ATYM9UZ7|piyo",
		"dddd":                     "",
		"":                         "",
	}
	var m Message
	for txt, ans := range s {
		m.Text = txt
		id := m.MentionID()
		if id != ans {
			t.Errorf("input: %v, got %v, expected %v", txt, id, ans)
		}
	}
}
