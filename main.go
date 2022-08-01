package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	// "github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	//"github.com/tidwall/gjson"

	witai "github.com/wit-ai/wit-go/v2"

)

func printCmdEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Cmd Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
    godotenv.Load(".env")
    fmt.Println("Test, nothing yet")

    bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

    go printCmdEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition {
		Description: "send query to Wolfram",
		Example: "who is the the Pittsburgh Penguins captain",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")

			fmt.Println(query)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			fmt.Println(msg)

			response.Reply("Success")
		}, 
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}