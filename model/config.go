package model

type ServerConfig struct {
	ServerName string `json:"server_name"`
	ServerVersion string `json:"server_version"`
	Config string `json:"config"`
}
