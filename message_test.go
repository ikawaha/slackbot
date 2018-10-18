package slackbot

import (
	"testing"
)

func TestMessageTextBody(t *testing.T) {
	testdata := map[string]string{
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
	for input, expected := range testdata {
		m.Text = input
		if got := m.TextBody(); got != expected{
			t.Errorf("input: %v, got %v, expected %v", input, got, expected)
		}
	}
}

func TestMessageMentionID(t *testing.T) {
	testdata := map[string]string{
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
	for input, expected := range testdata {
		m.Text = input
		if got := m.MentionID();got != expected {
			t.Errorf("input: %v, got %v, expected %v", input, got, expected)
		}
	}
}
