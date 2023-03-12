package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	GPT "ALFARIUS/gpt3"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
	"github.com/sashabaranov/go-openai"
)

func newAudioIntBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  8000,
		},
	}
	for {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}

func VoiceRecord() {
	// Abre una nueva sesión de audio
	err := portaudio.Initialize()
	if err != nil {
		fmt.Println("Error al inicializar portaudio:", err)
		os.Exit(1)
	}
	defer portaudio.Terminate()

	// Configura la grabación
	sampleRate := 44100
	framesPerBuffer := 512
	duration := 5 * time.Second

	buffer := make([]int16, int(duration.Seconds()*float64(sampleRate)))

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), framesPerBuffer, func(in []int16) {
		copy(buffer, in)
	})
	if err != nil {
		fmt.Println("Error al abrir el stream de audio:", err)
		os.Exit(1)
	}
	defer stream.Close()

	// Comienza la grabación
	fmt.Println("Comenzando a grabar...")
	err = stream.Start()
	if err != nil {
		fmt.Println("Error al iniciar la grabación:", err)
		os.Exit(1)
	}

	// Espera la duración de la grabación
	time.Sleep(duration)

	// Detiene la grabación
	err = stream.Stop()
	if err != nil {
		fmt.Println("Error al detener la grabación:", err)
		os.Exit(1)
	}
	fmt.Println("Grabación finalizada.")

	// Guarda el archivo de audio WAV
	outputFile, err := os.Create("grabacion.wav")
	if err != nil {
		fmt.Println("Error al crear el archivo de salida:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Read raw PCM data from input file.
	in, err := os.Open("audio.pcm")
	if err != nil {
		log.Fatal(err)
	}

	// Output file.
	out, err := os.Create("output.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// 8 kHz, 16 bit, 1 channel, WAV.
	e := wav.NewEncoder(out, 8000, 16, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(in)
	if err != nil {
		log.Fatal(err)
	}
	// Write buffer to output file. This writes a RIFF header and the PCM chunks from the audio.IntBuffer.
	if err := e.Write(audioBuf); err != nil {
		log.Fatal(err)
	}
	if err := e.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Archivo de salida guardado en grabacion.wav.")
}

func main() {
	c := openai.NewClient(GPT.GetApiKey())
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: "test.m4a",
	}
	resp, err := c.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return
	}

	//responseForNow
	fmt.Println(resp.Text)
}
