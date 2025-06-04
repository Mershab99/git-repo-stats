package models

type CommitsRequest struct {
	Days int32 `json:"days"`

	Repositories []RepoConfig `json:"repositories"`
}
