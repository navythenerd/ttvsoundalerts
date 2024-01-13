package twitch

import (
	"log"

	ttvirc "github.com/gempir/go-twitch-irc/v4"
)

type Alert struct {
	Trigger     string   `json:"trigger"`
	Name        string   `json:"name"`
	Media       string   `json:"media"`
	Points      int      `json:"points"`
	FromUser    string   `json:"fromUser,omitempty"`
	Permissions []string `json:"permissions"`
}

func (s *Service) loadAlerts() {
	for _, alert := range s.config.Alerts {
		log.Printf("load alert: %s with trigger %s\n", alert.Name, alert.Trigger)
		s.alerts[alert.Trigger] = alert
	}
}

func (s *Service) alertMessageHandler(message ttvirc.PrivateMessage) {
	alert, ok := s.alerts[message.Message]

	if !ok {
		return
	}

	log.Printf("%s played: %s\n", message.User.DisplayName, alert.Name)

	userPermissionsMask := getPermissionsMask(message.User.Badges)
	requiredPermissions := getPermissionsMask(toPermissionsMap(alert.Permissions))

	if hasPermissions(userPermissionsMask, requiredPermissions) {
		userAlert := alert
		userAlert.FromUser = message.User.DisplayName

		s.alertQueue <- &userAlert
	}
}
