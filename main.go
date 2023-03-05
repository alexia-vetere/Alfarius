package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func appendToConversation(message string, filePath string) {
	// Open file with append mode
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	}
	defer f.Close()

	// Add message to file with a newline separator
	if _, err := f.WriteString(message + "\n"); err != nil {
	}
	if err != nil {
		log.Fatal("Failed to save file:", err)
	}
}

func loadConversation(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conversation []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conversation = append(conversation, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return conversation, nil
}

func GetResponse(client gpt3.Client, ctx context.Context, prompt string) {
	var sb strings.Builder
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			prompt,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0.5),
	}, func(resp *gpt3.CompletionResponse) {

		//response += resp.Choices[0].Text
		sb.WriteString(resp.Choices[0].Text)
		//fmt.Print(response)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	defer func() {
		appendToConversation("\n"+sb.String(), filePath) // se guarda la respuesta completa en el archivo
	}()

	fmt.Printf("\n")
}

var response string

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
				appendToConversation("\nTú: "+question, filePath)

				conversationHistory, err := loadConversation(filePath)

				if err != nil {
					log.Fatal("Failed to load file:", err)
				}

				//conversationHistory = append(conversationHistory, "\nTú: "+question)
				prompt := strings.Join(conversationHistory, "\n")

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
