package models

type RepoConfig struct {
	Url string `json:"url"`

	Auth Auth `json:"auth,omitempty"`
}
