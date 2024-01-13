package twitch

import (
	"encoding/json"
	"os"
)

type Config struct {
	User        string  `json:"user"`
	Channel     string  `json:"channel"`
	Token       string  `json:"token"`
	JoinMessage string  `json:"joinMessage"`
	PartMessage string  `json:"partMessage"`
	Alerts      []Alert `json:"alerts"`
}

func ReadConfig(cfg *Config, file string) error {
	rawFile, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(rawFile, cfg)
	return err
}
