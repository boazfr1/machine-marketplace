package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/tmc/langchaingo/llms/openai"

	_ "embed"
)

//go:embed handle_error_prompt.tmpl
var handleErrorPrompt string

type (
	ErrorContext struct {
		Command string
		Output  string
		Error   string
	}

	ErrorFixResponse struct {
		CanFix      bool   `json:"canFix"`
		NewCommand  string `json:"newCommand"`
		Explanation string `json:"explanation"`
	}
)

func newLLM(openaiKey string, model string) (*openai.LLM, error) {
	return openai.New(
		openai.WithToken(openaiKey),
		openai.WithModel(model),
	)
}

func callLLM(ctx context.Context, openaiKey, prompt string, model string) (string, error) {
	llm, err := newLLM(openaiKey, model)
	if err != nil {
		return "", err
	}
	return llm.Call(ctx, prompt)
}

func HandleErrorWithLLM(ctx context.Context, openaiKey, command, output string, cmdError error, model string) (*ErrorFixResponse, error) {
	prompt, err := buildErrorPrompt(command, output, cmdError)
	if err != nil {
		return nil, err
	}

	llmResponse, err := callLLM(ctx, openaiKey, prompt, model)
	if err != nil {
		return nil, err
	}

	return parseErrorFixResponse(llmResponse)
}

func buildErrorPrompt(command, output string, cmdError error) (string, error) {
	tmpl, err := template.New("errorPrompt").Parse(handleErrorPrompt)
	if err != nil {
		return "", err
	}

	errorStr := ""
	if cmdError != nil {
		errorStr = cmdError.Error()
	}

	ctx := ErrorContext{
		Command: command,
		Output:  output,
		Error:   errorStr,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func parseErrorFixResponse(response string) (*ErrorFixResponse, error) {
	response = strings.TrimSpace(response)

	// Extract JSON from markdown code blocks if present
	if strings.Contains(response, "```json") {
		start := strings.Index(response, "```json") + 7
		end := strings.LastIndex(response, "```")
		if end > start {
			response = strings.TrimSpace(response[start:end])
		}
	} else if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		end := strings.LastIndex(response, "```")
		if end > start {
			response = strings.TrimSpace(response[start:end])
		}
	}

	var result ErrorFixResponse
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
