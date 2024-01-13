package twitch

import (
	"context"
	"fmt"
	"log"
	"net/http"

	ttvirc "github.com/gempir/go-twitch-irc/v4"
	"github.com/navythenerd/lionrouter"
	"nhooyr.io/websocket"
)

type Service struct {
	ctx        context.Context
	cancelCtx  context.CancelFunc
	irc        *ttvirc.Client
	config     *Config
	router     *lionrouter.Router
	sockets    []*websocket.Conn
	alertQueue chan *Alert
	alerts     map[string]Alert
}

func New(cfg *Config) *Service {
	srv := &Service{
		config:     cfg,
		router:     lionrouter.New(),
		sockets:    make([]*websocket.Conn, 0),
		alertQueue: make(chan *Alert),
		alerts:     make(map[string]Alert),
	}

	ctx, cancel := context.WithCancel(context.Background())
	srv.ctx = ctx
	srv.cancelCtx = cancel

	srv.irc = ttvirc.NewClient(cfg.User, fmt.Sprintf("oauth:%s", cfg.Token))

	srv.irc.OnConnect(func() {
		log.Printf("alert bot joined twitch channel: %s\n", cfg.Channel)
		srv.irc.Say(cfg.Channel, cfg.JoinMessage)
	})

	srv.irc.OnPrivateMessage(srv.alertMessageHandler)
	srv.irc.Join(cfg.Channel)

	srv.loadAlerts()
	srv.registerAlertPlayer()
	srv.startAlertServiceWorker()

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

	for _, conn := range s.sockets {
		conn.Close(websocket.StatusNormalClosure, "alert service is shutting down")
	}

	s.cancelCtx()
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
