package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokensLimit = 3000

type NullWriter int

const filePath = "./conversation.txt"
const fileName = "conversation.txt"

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	apiKey := viper.GetString("API_KEY")

	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			fmt.Println("Es un placer volver a verla mi señorita... ¿De qué quiere conversar hoy?")

			for !quit {
				fmt.Print("('quit' to end):")
				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				appendToConversation("Tú: "+question, filePath)

				conversationHistory := GetHistory(filePath)

				prompt := strings.Join(conversationHistory, "\n")
				promptTokens := []byte(prompt)
				//token limit control
				if len(promptTokens) >= tokensLimit {
					prompt = reducePromptSize(filePath)
					ioutil.WriteFile(filePath, []byte(prompt), 0644)
				}

				switch question {
				case "quit":

					quit = true

				default:
					GetResponse(client, ctx, prompt)
				}
			}
		},
	}
	rootCmd.Execute()
}
