package api

type Service struct {
	Name     string `json:"name"`
	Base_url string `json:"base_url"`
}

type AppState struct {
	services     []Service
	internal_key string
}
