package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultEndpoint = "https://chat.zxhkuav.cn/v1"
	defaultModel    = "gemma-4-26b-a4b-it"
	timeout         = 60 * time.Second
)

type AIConfig struct {
	Enabled  bool
	Model    string
	Endpoint string
}

// OpenAI-compatible chat completion request
type chatRequest struct {
	Model     string        `json:"model"`
	Messages  []chatMessage `json:"messages"`
	Stream    bool          `json:"stream"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Ollama API request/response (fallback)
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func Analyze(content string, cfg *AIConfig) (bool, string, error) {
	if !cfg.Enabled {
		return false, "", nil
	}

	model := cfg.Model
	if model == "" {
		model = defaultModel
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	systemPrompt := "You are a security code analyzer. Determine if the following code snippet contains malicious behavior like credential theft, command injection, backdoor, or data exfiltration. Reply ONLY with a JSON object: {\"malicious\": true/false, \"reason\": \"short explanation\"}."
	userContent := fmt.Sprintf("Code:\n%s", content)

	// Determine if endpoint is OpenAI-compatible or Ollama
	if strings.Contains(endpoint, "/v1") || strings.Contains(endpoint, "chat/completions") {
		return analyzeOpenAI(endpoint, model, systemPrompt, userContent)
	}
	return analyzeOllama(endpoint, model, systemPrompt, userContent)
}

func analyzeOpenAI(endpoint, model, systemPrompt, userContent string) (bool, string, error) {
	chatEndpoint := strings.TrimRight(endpoint, "/")
	if !strings.HasSuffix(chatEndpoint, "/chat/completions") {
		chatEndpoint += "/chat/completions"
	}

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userContent},
		},
		Stream: false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return false, "", fmt.Errorf("序列化请求失败: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(chatEndpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return false, "", fmt.Errorf("调用 AI 模型失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("读取 AI 响应失败: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return false, "", fmt.Errorf("解析 AI 响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return false, "", fmt.Errorf("AI 模型返回错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return false, "", fmt.Errorf("AI 响应为空")
	}

	response := chatResp.Choices[0].Message.Content
	return parseAIResponse(response)
}

func analyzeOllama(endpoint, model, systemPrompt, userContent string) (bool, string, error) {
	prompt := systemPrompt + "\n\n" + userContent

	reqBody := ollamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return false, "", fmt.Errorf("序列化请求失败: %w", err)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return false, "", fmt.Errorf("调用 Ollama 模型失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("读取 Ollama 响应失败: %w", err)
	}

	var oResp ollamaResponse
	if err := json.Unmarshal(body, &oResp); err != nil {
		return false, "", fmt.Errorf("解析 Ollama 响应失败: %w", err)
	}

	return parseAIResponse(oResp.Response)
}

func parseAIResponse(response string) (bool, string, error) {
	var result struct {
		Malicious bool   `json:"malicious"`
		Reason    string `json:"reason"`
	}

	// Try to extract JSON from response
	if idx := strings.Index(response, "{"); idx >= 0 {
		if endIdx := strings.LastIndex(response, "}"); endIdx > idx {
			response = response[idx : endIdx+1]
		}
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Fallback: heuristic check
		lower := strings.ToLower(response)
		result.Malicious = strings.Contains(lower, "malicious") || strings.Contains(lower, "true")
		result.Reason = "AI analysis result (heuristic)"
	}

	return result.Malicious, result.Reason, nil
}

func IsAvailable(endpoint string) bool {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	client := &http.Client{Timeout: 5 * time.Second}

	if strings.Contains(endpoint, "/v1") {
		// OpenAI-compatible: check models endpoint
		baseURL := strings.TrimRight(endpoint, "/")
		baseURL = strings.TrimSuffix(baseURL, "/chat/completions")
		modelsURL := baseURL + "/models"
		resp, err := client.Get(modelsURL)
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return resp.StatusCode == 200
	}

	// Ollama: check root endpoint
	resp, err := client.Get(strings.TrimSuffix(endpoint, "/api/generate"))
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == 200
}
