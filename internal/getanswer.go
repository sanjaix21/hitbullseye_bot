package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (h *Handler) GetAnswers() map[int]int {
	// Format all questions into one prompt
	prompt := "You are a test-taking expert. Analyze these multiple choice questions and provide the best answers.\n\n"
	prompt += "Return ONLY a JSON object in this exact format: {\"1\":2,\"2\":1,\"3\":4}\n"
	prompt += "Where key is question number and value is option number (1,2,3,4).\n\n"

	for _, q := range h.QuestionBank {
		prompt += fmt.Sprintf("Question %d: %s\n", q.QuestionNo, q.Question)
		options := strings.Split(q.Options, " | ")
		for i, opt := range options {
			if strings.TrimSpace(opt) != "" {
				prompt += fmt.Sprintf("%d) %s\n", i+1, strings.TrimSpace(opt))
			}
		}
		prompt += "\n"
	}
	prompt += "Remember: respond with ONLY the JSON object, no other text."

	// Call Gemini API
	answers := callGemini(prompt)
	fmt.Printf("Got %d answers from AI\n", len(answers))
	return answers
}

func callGemini(prompt string) map[int]int {
	apiKey := os.Getenv("GEMINI_API_KEY") // Set this in your environment
	if apiKey == "" {
		fmt.Println("GEMINI_API_KEY not set, using random answers")
		return generateRandomAnswers()
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

	reqBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Gemini API call failed:", err)
		return generateRandomAnswers()
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var geminiResp GeminiResponse
	if json.Unmarshal(respBody, &geminiResp) != nil {
		fmt.Println("Failed to parse Gemini response")
		return generateRandomAnswers()
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		fmt.Println("Empty Gemini response")
		return generateRandomAnswers()
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	answers := parseResponse(responseText)

	if len(answers) == 0 {
		fmt.Println("Failed to parse AI response, using random answers")
		return generateRandomAnswers()
	}

	return answers
}

func parseResponse(response string) map[int]int {
	// Find JSON in response
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}") + 1

	if start == -1 || end <= start {
		return nil
	}

	jsonStr := response[start:end]
	var raw map[string]interface{}

	if json.Unmarshal([]byte(jsonStr), &raw) != nil {
		return nil
	}

	answers := make(map[int]int)
	for qStr, optVal := range raw {
		var q, opt int
		fmt.Sscanf(qStr, "%d", &q)

		switch v := optVal.(type) {
		case float64:
			opt = int(v)
		case string:
			fmt.Sscanf(v, "%d", &opt)
		}

		if opt >= 1 && opt <= 4 && q > 0 {
			answers[q] = opt
		}
	}

	return answers
}

func generateRandomAnswers() map[int]int {
	// Fallback: generate random answers (1-4)
	answers := make(map[int]int)
	for i := 1; i <= 50; i++ {
		answers[i] = (i % 4) + 1
	}
	return answers
}
