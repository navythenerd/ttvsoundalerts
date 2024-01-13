package twitch

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"nhooyr.io/websocket"
)

func (s *Service) startAlertServiceWorker() {
	go func() {
		for {
			alert := <-s.alertQueue
			jsonPayload, err := json.Marshal(alert)

			if err != nil {
				log.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(s.ctx, time.Second*30)

			for _, conn := range s.sockets {
				conn.Write(ctx, websocket.MessageText, jsonPayload)
			}

			cancel()
		}
	}()
}
