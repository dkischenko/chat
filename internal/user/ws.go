package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *Service) InitSocketConnection(w http.ResponseWriter, r *http.Request, u *User) error {
	s.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connection, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Entry.Errorf("WS connection error: %s", err)
	}
	s.clientAdd(connection)
	log.Printf("%s goes online.", u.Username)
	if err = s.reader(connection, r.Context(), u); err != nil {
		s.logger.Entry.Errorf("%s", err)
	}
	s.clientDelete(connection)
	err = connection.Close()
	if err != nil {
		s.logger.Entry.Errorf("error happens: %s", err)
	}
	return nil
}

func (s *Service) reader(conn *websocket.Conn, ctx context.Context, u *User) (err error) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			err = s.storage.UpdateOnline(ctx, u, false)
			log.Printf("<%s> left chat", u.Username)
			s.logger.Entry.Errorf("Error with update online status: %s", err)
			return fmt.Errorf("received error: %s", err)
		}
		err = s.storage.UpdateOnline(ctx, u, true)
		if err != nil {
			s.logger.Entry.Errorf("Error with update online status: %s", err)
		}

		s.rwMutex.RLock()
		defer s.rwMutex.RUnlock()
		for connKey := range s.clients {
			if conn != connKey {
				log.Printf("<%s>: %s", u.Username, string(p))
				err := connKey.WriteMessage(messageType, p)
				if err != nil {
					s.logger.Entry.Errorf("Error with sending message: %s", err)
				}
			}
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			return fmt.Errorf("received error: %s", err)
		}
	}
}

func (s *Service) clientAdd(conn *websocket.Conn) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	s.clients[conn] = true
}

func (s *Service) clientDelete(conn *websocket.Conn) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()
	delete(s.clients, conn)
}
