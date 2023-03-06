package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

func appendToConversation(message string, filePath string) {
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

func reducePromptSize(filePath string) string {

	conversationHistory, err := loadConversation("./personality.txt")

	//conversation
	var conversation []string
	conversation, err = loadConversation(filePath)

	//pesonality and conversation
	conversationHistory = append(conversationHistory, conversation[1:]...)

	if err != nil {
		log.Fatal("Failed to load file when reducePromptSize:", err)
	}

	newPrompt := strings.Join(conversationHistory, "\n")

	return newPrompt
}

func GetHistory(filePath string) []string {
	//personality
	conversationHistory, err := loadConversation("./personality.txt")

	//conversation
	var conversation []string
	conversation, err = loadConversation(filePath)

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
		MaxTokens:   gpt3.IntPtr(tokensLimit),
		Temperature: gpt3.Float32Ptr(0.5),
	}, func(resp *gpt3.CompletionResponse) {

		sb.WriteString(resp.Choices[0].Text)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	defer func() {
		appendToConversation(sb.String(), filePath) // save the response.
		beautifyText(fileName)
	}()

	fmt.Printf("\n")
}

func beautifyText(fileName string) {
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
