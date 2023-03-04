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

func GetResponse(client gpt3.Client, ctx context.Context, question string, conversationHistory []string) []string {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			"Nuestra conversación previa: \n" + "Mi actual mensaje a responder: \n" + question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Stream:      false,
		Temperature: gpt3.Float32Ptr(0.5),
	}, func(resp *gpt3.CompletionResponse) {

		response += resp.Choices[0].Text
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
	conversationHistory = append(conversationHistory, "GPT-3: "+response)
	return conversationHistory
}

type NullWriter int

var response string

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	//log.SetOutput(new(NullWriter))

	conversationHistory := []string{
		"Tú: Hola, ¿cómo estás?",
		"GPT-3: Hola, estoy bien. ¿Y tú?",
		"Tú: Estoy bien también. Conversemos como si fueras un mayordomo amable y calmado que esta dispuesto a ayudarme y escucharme atento? Quisiera que uses fraces como: 'Si, mi señorita', '¿Cómo se encuentra hoy mi damisela?', entre otras amables palabras.",
		"GPT-3: ¡Claro que si mi damisela! ¿En que puedo ayudarla hoy?",
		"Tú: No, no tan efusivo y alegre. Debes sonar sereno y calmado.",
		"GPT-3: De acuerdo dama mia, estoy atento a lo que necesite. ¿De qué quiere conversar hoy?",
	}

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
				conversationHistory = append(conversationHistory, "Tú: "+question)
				prompt := strings.Join(conversationHistory, "\n") + "\n"

				switch question {
				case "quit":

					quit = true

				default:
					conversationHistory = GetResponse(client, ctx, prompt, conversationHistory)
					//fmt.Println(conversationHistory)
				}
			}
		},
	}
	rootCmd.Execute()
}
