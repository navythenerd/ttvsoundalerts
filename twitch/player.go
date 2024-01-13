package twitch

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/navythenerd/lionrouter"
	"nhooyr.io/websocket"
)

func (s *Service) registerAlertPlayer() {
	s.router.Get("/player", http.HandlerFunc(s.alertPlayerHander))
	s.router.Get("/alerts", http.HandlerFunc(s.alertWebsocketHandler))
	s.router.Get("/static/*file", http.HandlerFunc(s.staticHandler))
}

func (s *Service) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 - Not Found!"))
}

func (s *Service) alertPlayerHander(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "player/default-player.html")
}

func (s *Service) alertWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	s.sockets = append(s.sockets, conn)
}

func (s *Service) staticHandler(w http.ResponseWriter, r *http.Request) {
	file := lionrouter.Param(r.Context(), "file")
	path := fmt.Sprintf("static/%s", file)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		s.notFoundHandler(w, r)
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("static/%s", file))
}
