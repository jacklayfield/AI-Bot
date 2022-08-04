package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"

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

	//the slack bot itself
    bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	//the client from Wit.ai
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	//the Wolfram client
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

    go printCmdEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition {
		Description: "send query to Wolfram",
		Example: "who is the CEO of Alphabet",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")

			fmt.Println(query)
			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			//get json
			data, _ := json.MarshalIndent(msg, "", "    ")
			//get string
			rough := string(data[:])
			//get the actual query to send to wolfram
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value") 
			//convert query to expected string value to send to wolfram
			answer := value.String()
			//get the result from wolfram
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("there is an error")
			}
			fmt.Println(value)
			response.Reply(res)

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