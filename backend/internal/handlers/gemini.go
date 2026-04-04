package handlers

import (
	"bytes"
	"codefix-ai/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getAIReview(code, language string) (models.AIResponse, error) {
	prompt := fmt.Sprintf(`You are an expert code reviewer. Analyze this %s code.

Find any bugs, errors, or bad practices. Respond ONLY in this exact JSON format with no extra text and no markdown:
{
  "bug_explanation": "explain what is wrong in simple terms",
  "fixed_code": "the complete corrected code here",
  "suggestions": "2-3 tips to improve the code beyond just fixing the bug"
}

Code to review:
%s`, language, code)

	reqBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return models.AIResponse{}, err
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return models.AIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.AIResponse{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AIResponse{}, err
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.Unmarshal(respBytes, &geminiResp)
	if err != nil {
		return models.AIResponse{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return models.AIResponse{}, fmt.Errorf("no response from gemini: %s", string(respBytes))
	}

	content := geminiResp.Candidates[0].Content.Parts[0].Text
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var aiResponse models.AIResponse
	err = json.Unmarshal([]byte(content), &aiResponse)
	if err != nil {
		return models.AIResponse{}, fmt.Errorf("could not parse AI response: %v", err)
	}

	return aiResponse, nil
}