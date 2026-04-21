package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	godotenv.Load(".env")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY not found in .env")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-flash-latest")
	resp, err := model.GenerateContent(ctx, genai.Text("Hello, are you working?"))
	if err != nil {
		log.Fatalf("failed to generate content: %v", err)
	}

	fmt.Printf("Response: %v\n", resp.Candidates[0].Content.Parts[0])
}
