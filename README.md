# slackbot

Tiny slack bot lib.

# Prepare

Please get the token for your slack bot.
https://[YOURTEAM].slack.com/services/new/bot

see.https://api.slack.com/

# Interface

Echo bot sample.

```
bot, err := slackbot.New(token) // set token of your bot!
if err != nil {
    log.Fatal(err)
}
bot.Close()
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

Samples:

|Bot |Src            |Usage            |
|:---|:---           |:---             |
|echo|cmd/echobot.go | Mention your bot|
|haiku|cmd/haikubot.go| Passive        |
|dajare|cmd/dajarebot.go| Passive      |
|ja-morph|cmd/kagomebot.go| Mention your bot|

# Lisence

MIT

# Awesome

* github.com/mattn/go-haiku
* github.com/kurehajime/dajarep

# Author

ikawaha
