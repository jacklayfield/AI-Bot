package main

import("fmt", "context", "encoding/json", "log", "os", "github.com/joho/godotenv", "github.com/shomali11/slacker")

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

    go printCommandEvents(bot.CommandEvents())
}