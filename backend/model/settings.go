package model

type Settings struct {
	MockMode    bool   `json:"mockMode"`
	OpenAIURL   string `json:"openAIURL"`
	OpenAIModel string `json:"openAIModel"`
}

type SettingsUpdate struct {
	MockMode    *bool   `json:"mockMode"`
	OpenAIKey   string  `json:"openAIKey"`
	OpenAIURL   *string `json:"openAIURL"`
	OpenAIModel *string `json:"openAIModel"`
}
