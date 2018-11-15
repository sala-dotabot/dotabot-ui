package handler 

import (
	"dotabot-ui/state"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	local_telegram "dotabot-ui/telegram"
)

const INIT_STATE = ""

type Handler struct {
	stateRepository state.StateRepository
	commands []Command
}

func CreateHandler(commands []Command, stateRepository state.StateRepository) *Handler {
	return &Handler{
		stateRepository: stateRepository,
		commands: commands,
	}
}

func (this Handler) Handle(w http.ResponseWriter, r *http.Request) {
    if (r.URL == nil || r.URL.Path != "/webhook") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var u local_telegram.Update
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%#v", u.Message)

	log.Printf("Loading state")
	state, err := this.stateRepository.LoadState(u.Message.Chat.Id)
	if (err != nil) {
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
		log.Printf("Handler was not found")
	}

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}
