package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"ALFARIUS/GPT"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
)

func main() {

	ctx := context.Background()
	client := gpt3.NewClient(GPT.GetApiKey())

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			fmt.Println(GPT.Intro) //start of talk

			for !quit {
				fmt.Print("('quit' to end):")
				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				GPT.AppendToConversation(GPT.NameUser+question, GPT.FilePath)

				conversationHistory := GPT.GetHistory(GPT.FilePath)

				prompt := strings.Join(conversationHistory, "\n")
				promptTokens := []byte(prompt)
				//fmt.Println(len(promptTokens))
				//token limit control
				if len(promptTokens) >= (GPT.TokensLimit - 700) {
					prompt = GPT.ReducePromptSize(GPT.FilePath)
					ioutil.WriteFile(GPT.FilePath, []byte(prompt), 0644)
				}

				switch question {
				case "quit":

					quit = true

				default:
					GPT.GetResponse(client, ctx, prompt)
				}
			}
		},
	}
	rootCmd.Execute()
}
