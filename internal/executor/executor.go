package executor

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lucasepe/gptcli/internal/app"
	"github.com/lucasepe/gptcli/internal/shortcuts"
	"github.com/sashabaranov/go-openai"
)

func GPT(shc map[string]string) prompt.Executor {
	client := openai.NewClient(app.OpenAIApiKey())

	return func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		} else if s == "quit" || s == "exit" {
			fmt.Println("Bye!")
			os.Exit(0)
			return
		}

		query := shortcuts.Expand(shc, fmt.Sprintf("@@%s", s))
		if len(query) == 0 {
			return
		}

		resp, err := client.CreateChatCompletion(context.TODO(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: query,
					},
				},
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Got error: %s\n", err.Error())
			os.Exit(1)
			return
		}

		if len(resp.Choices) > 0 {
			content := strings.TrimSpace(resp.Choices[0].Message.Content)
			fmt.Println(content)
			fmt.Println()
		}
	}
}
