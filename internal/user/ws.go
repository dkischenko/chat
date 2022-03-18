package user

import (
	"context"
	"fmt"
	"github.com/dkischenko/chat/internal/user/models"
	"github.com/gorilla/websocket"
	"net/http"
)

func (s *Service) InitSocketConnection(w http.ResponseWriter, r *http.Request, u *models.User) error {
	s.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connection, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Entry.Errorf("WS connection error: %w", err)
	}
	s.clientAdd(connection)
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
		var (
			mes         = make(chan []byte, 1)
			errs        = make(chan error, 1)
			closeWS     = make(chan bool, 1)
			messageType int
			p           []byte
			err         error
		)
		go func() {
			messageType, p, err = conn.ReadMessage()
			if messageType == websocket.CloseMessage || err != nil {
				err = s.storage.UpdateOnline(ctx, u, false)
				if err != nil {
					s.logger.Entry.Errorf("Error with update online status: %w", err)
					errs <- fmt.Errorf("received error: %w", err)
				} else {
					s.logger.Entry.Printf("<%s> left chat", u.Username)
					closeWS <- true
				}
			}

			mes <- p

			close(closeWS)
			close(mes)
			close(errs)
		}()

		if cWS := <-closeWS; cWS {
			return nil
		}

		if err, ok := <-errs; ok {
			return err
		}

		s.rwMutex.RLock()
		defer s.rwMutex.RUnlock()

		if m, ok := <-mes; ok {
			s.handleMessages(ctx, m, u, conn)
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			return fmt.Errorf("received error: %w", err)
		}
	}
}

func (s Service) handleMessages(ctx context.Context, m []byte, u *models.User, conn *websocket.Conn) {
	msg := &models.Message{
		Text:  string(m),
		UFrom: u.ID,
	}
	mid, err := s.storage.StoreMessage(ctx, msg)
	if err != nil {
		s.logger.Entry.Errorf("Error with store message: %w", err)
	}

	mDB, err := s.storage.FindOneMessage(ctx, mid)
	if err != nil {
		s.logger.Entry.Errorf("Error with getting message: %w", err)
	}

	for connKey := range s.clients {
		if conn != connKey {
			err = connKey.WriteJSON(mDB)
			if err != nil {
				s.logger.Entry.Errorf("Error with sending message: %w", err)
			}
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
