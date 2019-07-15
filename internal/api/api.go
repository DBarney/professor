package api

import (
	"fmt"
	"github.com/dbarney/professor/types"
	"github.com/dbarney/professor/www"
	"net/http"
)

type api struct {
	source    <-chan *types.Event
	listeners []chan *types.Event
	events    []*types.Event
}

type Store interface {
}

func Run(source <-chan *types.Event, store Store) {

	a := &api{
		source:    source,
		listeners: []chan *types.Event{},
		events:    []*types.Event{},
	}

	http.HandleFunc("/events.source", a.serveEvents)
	http.Handle("/", http.FileServer(www.AssetFile()))
	go a.record()
	go http.ListenAndServe(":8080", nil)
}

func (a *api) record() {
	for event := range a.source {
		if event.Status == types.Pending {
			a.events = []*types.Event{}
		}
		a.events = append(a.events, event)
		for _, l := range a.listeners {
			select {
			case l <- event:
			default:
			}
		}
	}
}

func (a *api) serveEvents(w http.ResponseWriter, r *http.Request) {
	c := make(chan *types.Event, 10)
	a.listeners = append(a.listeners, c)
	defer func() {
		//remove the channel from the listeners list
	}()

	w.Header().Set("Content-Type", "text/event-stream")

	for i := range a.events {
		event := a.events[i]
		_, err := fmt.Fprintf(w, "id:%v\nevent:%v\ndata:%v\n\n", event.Sha, event.Status.String(), string(event.Data))
		if err != nil {
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	for event := range c {
		_, err := fmt.Fprintf(w, "id:%v\nevent:%v\ndata:%v\n\n", event.Sha, event.Status.String(), string(event.Data))
		if err != nil {
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
