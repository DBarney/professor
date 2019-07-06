package api

import (
	"fmt"
	"github.com/dbarney/professor/www"
	"net/http"
)

type Set struct {
	Name   string
	Events <-chan *Event
}

type Event struct {
	Name string
	Data string
}

type api struct {
	source    <-chan *Set
	set       *Set
	listeners []chan *Event
	events    []*Event
}

func Run(source <-chan *Set) {

	a := &api{
		source:    source,
		listeners: []chan *Event{},
		events:    []*Event{},
	}

	http.HandleFunc("/events.source", a.serveEvents)
	http.Handle("/", http.FileServer(www.AssetFile()))
	go a.record()
	go http.ListenAndServe(":8080", nil)
}

func (a *api) record() {
	for set := range a.source {
		a.set = set
		a.events = []*Event{}
		a.forward(&Event{Name: "new", Data: set.Name})
		for event := range set.Events {
			a.forward(event)
		}
		a.forward(&Event{Name: "done", Data: set.Name})
	}
}

func (a *api) forward(e *Event) {
	a.events = append(a.events, e)
	for _, l := range a.listeners {
		select {
		case l <- e:
		default:
		}
	}
}

func (a *api) serveEvents(w http.ResponseWriter, r *http.Request) {
	c := make(chan *Event, 10)
	a.listeners = append(a.listeners, c)
	defer func() {
		//remove the channel from the listeners list
	}()

	w.Header().Set("Content-Type", "text/event-stream")

	for _, event := range a.events {
		fmt.Fprintf(w, "event:%v\ndata:%v\n\n", event.Name, event.Data)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	for event := range c {
		_, err := fmt.Fprintf(w, "event:%v\ndata:%v\n\n", event.Name, event.Data)
		if err != nil {
			return
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
