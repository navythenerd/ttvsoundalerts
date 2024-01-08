package twitch

import (
	"fmt"
	"log"

	ttvirc "github.com/gempir/go-twitch-irc/v4"
)

type Service struct {
	irc    *ttvirc.Client
	config *Config
}

func New(cfg *Config) *Service {
	srv := &Service{
		config: cfg,
	}

	srv.irc = ttvirc.NewClient(cfg.User, fmt.Sprintf("oauth:%s", cfg.Token))

	srv.irc.OnConnect(func() {
		log.Printf("Bot joined twitch channel: %s\n", cfg.Channel)
		srv.irc.Say(cfg.Channel, cfg.JoinMessage)
	})

	srv.irc.OnPrivateMessage(srv.privateMessageHandler)
	srv.irc.Join(cfg.Channel)

	return srv
}

func (s *Service) Connect() {
	go func() {
		err := s.irc.Connect()

		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Service) Shutdown() {
	s.irc.Say(s.config.Channel, s.config.PartMessage)
}

func (s *Service) privateMessageHandler(message ttvirc.PrivateMessage) {
	log.Printf("%s: %s\n", message.User.DisplayName, message.Message)
}
