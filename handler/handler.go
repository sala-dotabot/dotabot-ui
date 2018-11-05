package handler 

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	local_telegram "dotabot-ui/telegram"
)

type Handler struct {
	commands []Command
}

func CreateHandler(commands []Command) *Handler {
	return &Handler{
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

	handled := false
	for _, command := range this.commands {
		if command.CanHandle(u) {
			handled = true

			err = command.Handle(u)
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
