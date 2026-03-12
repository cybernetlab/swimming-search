package test_tg

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Message struct {
	ChatID string
	Text   string
}

type Server struct {
	Messages chan Message
	Server   *httptest.Server
}

func New(t *testing.T) Server {
	var s Server

	s.Messages = make(chan Message, 10)
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		switch {
		case strings.HasSuffix(r.URL.Path, "getMe"):
			_, err := w.Write([]byte(`{
            "ok": true, 
            "result": {
                "id": 111, 
                "is_bot": true, 
                "first_name": "TestBot", 
                "username": "test_bot"
            }
        }`))
			require.NoError(t, err)

		case strings.HasSuffix(r.URL.Path, "sendMessage"):
			err := r.ParseForm()
			require.NoError(t, err)
			s.Messages <- Message{ChatID: r.PostFormValue("chat_id"), Text: r.PostFormValue("text")}
			_, err = w.Write([]byte(`{"ok": true, "result": {"message_id": 123}}`))
			require.NoError(t, err)

		default:
			_, err := w.Write([]byte(`{"ok": false, "description": "unknown method"}`))
			require.NoError(t, err)
		}
	}))
	return s
}
