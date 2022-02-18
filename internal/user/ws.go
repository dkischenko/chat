package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *service) InitSocketConnection(w http.ResponseWriter, r *http.Request, u *User) error {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Entry.Errorf("WS connection error: %s", err)
	}
	clients[connection] = true
	log.Println("Client successfully connected.")
	if err = s.reader(connection, r.Context(), u); err != nil {
		s.logger.Entry.Errorf("%s", err)
	}
	delete(clients, connection)
	err = connection.Close()
	if err != nil {
		s.logger.Entry.Errorf("error happens: %s", err)
	}
	return nil
}

func (s *service) reader(conn *websocket.Conn, ctx context.Context, u *User) (err error) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			// remove from online
			err = s.storage.UpdateOnline(ctx, u, map[string]bool{"is_online": false})
			if err != nil {
				s.logger.Entry.Errorf("Error with update online status: %s", err)
			}
			return fmt.Errorf("Received error: %s", err)
		}
		err = s.storage.UpdateOnline(ctx, u, map[string]bool{"is_online": true})
		if err != nil {
			s.logger.Entry.Errorf("Error with update online status: %s", err)
		}
		log.Printf("Received message: %s", string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			return fmt.Errorf("Received error: %s", err)
		}
	}
}
