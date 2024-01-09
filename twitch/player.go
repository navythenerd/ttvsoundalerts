package twitch

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

func (s *Service) registerAlertPlayer() {
	s.router.Get("/player", http.HandlerFunc(s.alertPlayerHander))
	s.router.Get("/alerts", http.HandlerFunc(s.alertWebsocketHandler))
}

func (s *Service) alertPlayerHander(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "player/index.html")
}

func (s *Service) alertWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	s.sockets = append(s.sockets, conn)
}
