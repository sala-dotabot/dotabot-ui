package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/saladinkzn/dotabot-ui/state"

	local_telegram "github.com/saladinkzn/dotabot-ui/telegram"
)

const INIT_STATE = ""

type Handler struct {
	stateRepository       state.StateRepository
	commands              []Command
	totalCounter          prometheus.Counter
	notFoundErrorCounter  prometheus.Counter
	decodeErrorCounter    prometheus.Counter
	loadStateErrorCounter prometheus.Counter
	unhandledCounter      prometheus.Counter
	errorCounter          prometheus.Counter
}

func CreateHandler(commands []Command,
	stateRepository state.StateRepository,
	totalCounter prometheus.Counter,
	notFoundErrorCounter prometheus.Counter,
	decodeErrorCounter prometheus.Counter,
	loadStateErrorCounter prometheus.Counter,
	unhandledCounter prometheus.Counter,
	errorCounter prometheus.Counter) *Handler {
	return &Handler{
		stateRepository:       stateRepository,
		commands:              commands,
		notFoundErrorCounter:  notFoundErrorCounter,
		totalCounter:          totalCounter,
		decodeErrorCounter:    decodeErrorCounter,
		loadStateErrorCounter: loadStateErrorCounter,
		unhandledCounter:      unhandledCounter,
		errorCounter:          errorCounter,
	}
}

func (this Handler) Handle(w http.ResponseWriter, r *http.Request) {
	this.totalCounter.Inc()

	if r.URL == nil || r.URL.Path != "/webhook" {
		this.notFoundErrorCounter.Inc()
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var u local_telegram.Update
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		this.decodeErrorCounter.Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%#v", u)
	log.Printf("%#v", u.Message)

	if u.Message == nil {
		fmt.Fprint(w, "OK")
		return
	}

	log.Printf("Loading state")
	state, err := this.stateRepository.LoadState(u.Message.Chat.Id)
	if err != nil {
		this.loadStateErrorCounter.Inc()
		log.Printf("Error has occurred while loading state: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Loaded state: %s", state)

	handled := false
	for _, command := range this.commands {
		if command.CanHandle(u, state) {
			handled = true

			err = command.Handle(u, state)
		}
	}

	if !handled {
		this.unhandledCounter.Inc()
		log.Printf("Handler was not found")
	}

	if err != nil {
		this.errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}
