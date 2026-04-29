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
	defaultEndpoint = "http://localhost:11434/api/generate"
	defaultModel    = "llama3.2"
	timeout         = 30 * time.Second
)

type AIConfig struct {
	Enabled  bool
	Model    string
	Endpoint string
}

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

	prompt := fmt.Sprintf(`You are a security code analyzer. Determine if the following code snippet contains malicious behavior like credential theft, command injection, backdoor, or data exfiltration. Reply ONLY with a JSON object: {"malicious": true/false, "reason": "short explanation"}.

Code:
%s`, content)

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
		return false, "", fmt.Errorf("调用 AI 模型失败 (请确认 Ollama 已启动): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("读取 AI 响应失败: %w", err)
	}

	var oResp ollamaResponse
	if err := json.Unmarshal(body, &oResp); err != nil {
		return false, "", fmt.Errorf("解析 AI 响应失败: %w", err)
	}

	var result struct {
		Malicious bool   `json:"malicious"`
		Reason    string `json:"reason"`
	}

	response := oResp.Response
	if idx := strings.Index(response, "{"); idx >= 0 {
		if endIdx := strings.LastIndex(response, "}"); endIdx > idx {
			response = response[idx : endIdx+1]
		}
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		lower := strings.ToLower(oResp.Response)
		result.Malicious = strings.Contains(lower, "malicious") || strings.Contains(lower, "yes")
		result.Reason = "AI analysis result"
	}

	return result.Malicious, result.Reason, nil
}

func IsAvailable(endpoint string) bool {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(strings.TrimSuffix(endpoint, "/api/generate"))
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}
