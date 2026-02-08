package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}
	tools := []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "Read",
			Description: openai.String("Read and return the content of the file"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"file_path": map[string]any{
						"type":        "string",
						"description": "the path of the file to read",
					},
				},
				"required": []string{"file_path"},
			},
		}),
	}

	result, err := gameLoop(client, messages, tools)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, result)
}

func gameLoop(client openai.Client, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) (string, error) {
	for {
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools:    tools,
			},
		)
		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

		// 1
		fmt.Print(resp.Choices[0].Message.Content)

		// 2 workaround for test.
		messages = append(messages,
			openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(resp.Choices[0].Message.Content),
					},
				},
			},
		)

		var data []byte

		// rangeover toolcalls
		if toolCalls := resp.Choices[0].Message.ToolCalls; len(toolCalls) > 0 {

			messages = append(messages,
				openai.ChatCompletionMessageParamUnion{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						ToolCalls: toolCalls,
					},
				},
			)

			for i := range toolCalls {
				toolCall := toolCalls[i]
				toolCallFunction := toolCall.Function
				toolCallFunctionArgs := toolCallFunction.Arguments

				// only doing Read here as for now without even looking the name of function
				var jsonArgs struct {
					FilePath string `json:"file_path"`
				}
				err = json.Unmarshal([]byte(toolCallFunctionArgs), &jsonArgs)
				if err != nil {
					return "", err
				}

				data, err = os.ReadFile(jsonArgs.FilePath)
				if err != nil {
					return "", err
				}

				messages = append(messages,
					openai.ChatCompletionMessageParamUnion{
						OfTool: &openai.ChatCompletionToolMessageParam{
							ToolCallID: toolCall.ID,
							Content: openai.ChatCompletionToolMessageParamContentUnion{
								OfString: openai.String(string(data)),
							},
						},
					},
				)

			}
		} else {
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(resp.Choices[0].Message.Content),
					},
				},
			},
			)
			return string(data), nil
		}
	}
}
