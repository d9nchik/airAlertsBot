package puller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"airAlertsBot/pkg/core"
)

type Puller struct {
	sender        core.Sender
	previousState State
}

func NewPuller(sender core.Sender) *Puller {
	return &Puller{sender: sender}
}

func (p *Puller) Run(ctx context.Context) {
	ticker := time.NewTicker(1)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ticker.Reset(time.Minute)

			notification, err := getStatuses()
			if err != nil {
				log.Printf("Problem with retrieving statuses: %v", err.Error())
				continue
			}

			state := notification.States["Волинська область"]

			if p.previousState.Equal(state) {
				continue
			}

			if ok := p.sender.SendMessage(state.toMessage()); ok {
				p.previousState = state
			}

		case <-ctx.Done():
			return
		}
	}
}

func getStatuses() (*Notifications, error) {
	resp, err := http.Get("https://emapa.fra1.cdn.digitaloceanspaces.com/statuses.json")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var notofication Notifications
	err = json.Unmarshal(data, &notofication)
	if err != nil {
		return nil, err
	}

	return &notofication, nil
}

type Notifications struct {
	States map[string]State `json:"states"`
}

type State struct {
	IsEnabled  bool      `json:"enabled"`
	EnabledAt  time.Time `json:"enabled_at"`
	DisabledAt time.Time `json:"disabled_at"`
}

func (s *State) toMessage() string {
	loc, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		loc = time.Local
	}

	if s.IsEnabled {
		return fmt.Sprint("🔴 ", s.DisabledAt.In(loc).Hour(), ":", s.DisabledAt.Minute(), " Повітряна тривоги")
	}
	return fmt.Sprint("🟢 ", s.DisabledAt.In(loc).Hour(), ":", s.DisabledAt.Minute(), " Відбій тривоги")
}

func (s *State) Equal(s2 State) bool {
	if s.IsEnabled != s2.IsEnabled {
		return false
	}
	if !s.EnabledAt.Equal(s2.EnabledAt) {
		return false
	}
	return s.DisabledAt.Equal(s2.DisabledAt)
}
