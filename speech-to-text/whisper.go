package main

import (
	"context"
	"fmt"

	"ALFARIUS/GPT"

	"github.com/sashabaranov/go-openai"
)

func voiceRecord() {

}

func voiceToText() string {
	c := openai.NewClient(GPT.GetApiKey())
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "./records/audio.wav",
	}
	resp, err := c.CreateTranscription(ctx, req)

	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return resp.Text
	}

	return resp.Text
}

func main() {

	//voiceRecord()

	//defer func() {
	//	recordInText := voiceToText()
	//responseForNow
	//	fmt.Println(recordInText)
	//}()

}
