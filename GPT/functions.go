package GPT

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
)

func AppendToConversation(message string, filePath string) {
	// Open file with append mode
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	}
	defer f.Close()

	// Add message to file with a newline separator
	if _, err := f.WriteString("\n" + message); err != nil {
	}
	if err != nil {
		log.Fatal("Failed to save file:", err)
	}
}

func LoadConversation(filePath string) ([]string, error) {
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

func ReducePromptSize(filePath string) string {

	conversationHistory, err := LoadConversation("./gpt/personality.txt")

	//conversation
	var conversation []string
	conversation, err = LoadConversation(filePath)

	//pesonality and conversation
	if len(conversation) > 10 {
		conversationHistory = append(conversationHistory, conversation[11:]...)
	} else {
		conversationHistory = append(conversationHistory, conversation[1:]...)
	}

	if err != nil {
		log.Fatal("Failed to load file when reducePromptSize:", err)
	}

	newPrompt := strings.Join(conversationHistory, "\n")

	return newPrompt
}

func GetHistory(filePath string) []string {
	//personality
	conversationHistory, err := LoadConversation("./gpt/personality.txt")

	//conversation
	var conversation []string
	conversation, err = LoadConversation(filePath)

	//pesonality and conversation
	conversationHistory = append(conversationHistory, conversation...)

	if err != nil {
		log.Fatal("Failed to load file when GetHistory:", err)
	}

	return conversationHistory
}

func GetResponse(client gpt3.Client, ctx context.Context, prompt string) {
	var sb strings.Builder
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			prompt,
		},
		MaxTokens:   gpt3.IntPtr(TokensLimit),
		Temperature: gpt3.Float32Ptr(0.2),
	}, func(resp *gpt3.CompletionResponse) {

		sb.WriteString(resp.Choices[0].Text)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	defer func() {
		AppendToConversation(sb.String(), FilePath) // save the response.
		BeautifyText(FileName)
	}()

	fmt.Printf("\n")
}

func BeautifyText(fileName string) {
	// delete the lines breaks of the saved conversation
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// save the lines of the conversation in this slice
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// overwrite the conversation without lines breaks
	for _, line := range lines {
		fmt.Fprintln(outFile, line)
	}
}

func GetApiKey() string {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}
	return apiKey
}

type NullWriter int

var TokensLimit = 3000

const FilePath = "./gpt/conversation.txt"
const FileName = "conversation.txt"
const Intro = "Es un placer volver a verla mi señorita... ¿De qué quiere conversar hoy?"
const NameUser = "Alex: "

func (NullWriter) Write([]byte) (int, error) { return 0, nil }
