package user

import (
	"context"
	"fmt"
	"github.com/dkischenko/chat/internal/user/models"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func (s *Service) InitSocketConnection(w http.ResponseWriter, r *http.Request, u *models.User) error {
	s.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connection, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Entry.Errorf("WS connection error: %s", err)
	}
	s.clientAdd(connection)
	log.Printf("%s goes online.", u.Username)
	s.logger.Entry.Infof("%s goes online.", u.Username)
	err = s.storage.UpdateOnline(r.Context(), u, true)
	if err != nil {
		s.logger.Entry.Errorf("Error with update online status: %w", err)
	}

	cnt, err := s.storage.GetUnreadMessagesCount(r.Context(), u)
	if err != nil {
		s.logger.Entry.Errorf("Error with getting unread messages count: %w", err)
	}

	unreadMes, err := s.storage.GetUnreadMessages(r.Context(), u, cnt)
	if err != nil {
		s.logger.Entry.Errorf("Error with getting unread messages: %w", err)
	}

	for _, mes := range unreadMes {
		err := connection.WriteJSON(mes)
		if err != nil {
			s.logger.Entry.Errorf("Error with sending mesage: %w", err)
			err = connection.Close()
			if err != nil {
				s.logger.Entry.Errorf("error happens: %w", err)
			}
			s.clientDelete(connection)
		}
	}

	if err = s.reader(connection, r.Context(), u); err != nil {
		s.logger.Entry.Errorf("%w", err)
	}
	s.clientDelete(connection)
	err = connection.Close()
	if err != nil {
		s.logger.Entry.Errorf("error happens: %w", err)
	}
	return nil
}

func (s *Service) reader(conn *websocket.Conn, ctx context.Context, u *models.User) (err error) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			err = s.storage.UpdateOnline(ctx, u, false)
			log.Printf("<%s> left chat", u.Username)
			s.logger.Entry.Errorf("Error with update online status: %w", err)
			return fmt.Errorf("received error: %w", err)
		}

		s.rwMutex.RLock()
		defer s.rwMutex.RUnlock()
		var wg sync.WaitGroup

		msg := &models.Message{
			Text:  string(p),
			UFrom: u.ID,
		}
		mid, err := s.storage.StoreMessage(ctx, msg)
		if err != nil {
			s.logger.Entry.Errorf("Error with store message: %s", err)
		}

		m, err := s.storage.FindOneMessage(ctx, mid)
		if err != nil {
			s.logger.Entry.Errorf("Error with getting message: %s", err)
		}

		for connKey := range s.clients {
			wg.Add(1)
			go func(connKey *websocket.Conn, m *models.Message) {
				defer wg.Done()
				if conn != connKey {
					//log.Printf("<%s>: %s", u.Username, m.Text)
					err = connKey.WriteJSON(m)
					if err != nil {
						s.logger.Entry.Errorf("Error with sending message: %s", err)
					}
				}
			}(connKey, m)
		}
		wg.Wait()

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
