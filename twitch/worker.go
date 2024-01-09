package twitch

import (
	"nhooyr.io/websocket"
)

func (s *Service) StartAlertServiceWorker() {
	go func() {
		for {
			message := <-s.alerts

			for _, conn := range s.sockets {
				conn.Write(s.ctx, websocket.MessageText, []byte(message))
			}
		}

	}()
}
