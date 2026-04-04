package models

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"created_at"`
}

type Review struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Language       string `json:"language"`
	OriginalCode   string `json:"original_code"`
	BugExplanation string `json:"bug_explanation"`
	FixedCode      string `json:"fixed_code"`
	Suggestions    string `json:"suggestions"`
	IsHelpful      *bool  `json:"is_helpful"`
	CreatedAt      string `json:"created_at"`
}

// what user sends when submitting code
type ReviewRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// what we parse back from openai
type AIResponse struct {
	BugExplanation string `json:"bug_explanation"`
	FixedCode      string `json:"fixed_code"`
	Suggestions    string `json:"suggestions"`
}