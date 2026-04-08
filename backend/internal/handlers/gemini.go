package handlers// we always need to write in go in which package/folder it is in

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
	}//building request body to send to gemini api

	bodyBytes, err := json.Marshal(reqBody)//converts go object/map to json bytes format
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

	resp, err := http.DefaultClient.Do(req)//do() sends request to network
	if err != nil {
		return models.AIResponse{}, err//returns empty structure and error separately
	}
	defer resp.Body.Close()//close the stream of data and free resources like closing the tap 
	//if we didn't write defer keyword their might be possibility that it skipped after error occured but defer make sure us that it will
	//run despite of getting error or response
	//prevents memory leak, crashing server
	respBytes, err := io.ReadAll(resp.Body)//read data and store data in bytes
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
	}//go needs structure to fill up the response but in case of nodejs it directly converts json to object
// Node.js → dynamic (no predefined shape needed)
// Go → static (must define structure)
	err = json.Unmarshal(respBytes, &geminiResp)//json bytes to go object
	if err != nil {
		return models.AIResponse{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return models.AIResponse{}, fmt.Errorf("no response from gemini: %s", string(respBytes))
	}//it checks whether the api returned valid data or not, if not then it returns error rather than crashing

	content := geminiResp.Candidates[0].Content.Parts[0].Text
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
//It extracts the AI response text and removes unwanted markdown (like ```json) so it becomes clean JSON.
	var aiResponse models.AIResponse
	err = json.Unmarshal([]byte(content), &aiResponse)
	if err != nil {
		return models.AIResponse{}, fmt.Errorf("could not parse AI response: %v", err)
	}
//It converts the cleaned JSON string into your final structured Go object and returns it.
	return aiResponse, nil
}