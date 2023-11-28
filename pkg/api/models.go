package api

type Service struct {
	Name                       string `json:"name"`
	Base_url                   string `json:"base_url"`
	AuthorizationKeyIdentifier string `json:"authorization_key_identifier"`
}

type AppState struct {
	services    []Service
	internalKey string
}
