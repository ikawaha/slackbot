# slackbot

Tiny slack bot client.

# Prepare

Please get the token for your slack bot.

`https://[YOURTEAM].slack.com/services/new/bot`

see. https://api.slack.com/

# Interface

Echo bot sample.

```Go
// your bot
type Bot struct {
	*slackbot.Client
}

func NewBot(token string) (*Bot, error) {
	c, err := slackbot.New(token)
	if err != nil {
		return nil, err
	}
	return &Bot{Client: c}, err
}

bot, err := NewBot(token) // set your bot token!
if err != nil {
    log.Fatal(err)
}
defer bot.Close()
fmt.Println("^C exits")

for {
    msg, err := bot.GetMessage()
    if err != nil {
       log.Printf("receive error, %v", err)
    }
    if bot.ID == msg.MentionID() && msg.Type == "message" && msg.SubType == "" {
        go func(m slackbot.Message) {
           m.Text = m.TextBody()
           bot.PostMessage(m)
        }(msg)
    }
}

```

# Lisence

MIT
