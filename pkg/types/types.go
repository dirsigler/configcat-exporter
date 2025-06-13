package types

import "time"

// ConfigCat API response structures

// Organization represents a ConfigCat organization
type Organization struct {
	OrganizationID string `json:"organizationId"`
	Name           string `json:"name"`
}

// Product represents a ConfigCat product
type Product struct {
	ProductID string `json:"productId"`
	Name      string `json:"name"`
}

// Config represents a ConfigCat configuration
type Config struct {
	ConfigID string `json:"configId"`
	Name     string `json:"name"`
}

// Environment represents a ConfigCat environment
type Environment struct {
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
}

// FeatureFlag represents a ConfigCat feature flag/setting
type FeatureFlag struct {
	SettingID int    `json:"settingId"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	Hint      string `json:"hint"`
}

// ZombieFlag represents a ConfigCat zombieflag/setting
type ZombieFlag struct {
	ProductID string `json:"productId"`
	Name      string `json:"name"`
	Configs   []struct {
		ConfigID          string `json:"configId"`
		Name              string `json:"name"`
		EvaluationVersion string `json:"evaluationVersion"`
		HasCodeReferences bool   `json:"hasCodeReferences"`
		Settings          []struct {
			SettingID         int    `json:"settingId"`
			Name              string `json:"name"`
			Key               string `json:"key"`
			Hint              string `json:"hint"`
			HasCodeReferences bool   `json:"hasCodeReferences"`
			Tags              []struct {
				TagID        int `json:"tagId"`
				SettingTagID int `json:"settingTagId"`
			} `json:"tags"`
			SettingValues []struct {
				EnvironmentID string    `json:"environmentId"`
				UpdatedAt     time.Time `json:"updatedAt"`
				IsStale       bool      `json:"isStale"`
			} `json:"settingValues"`
		} `json:"settings"`
	} `json:"configs"`
	Environments []struct {
		EnvironmentID string `json:"environmentId"`
		Name          string `json:"name"`
	} `json:"environments"`
}
