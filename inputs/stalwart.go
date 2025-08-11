package inputs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func StalwartHandler(ds func(id string, color string, token string, title *string, message *string, footer *string) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		type Payload struct {
			Events []struct {
				Id        string      `json:"id"`
				CreatedAt time.Time   `json:"createdAt"`
				Type      string      `json:"type"`
				Data      interface{} `json:"data"`
			} `json:"events"`
		}

		payload := Payload{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		errors := make([]error, 0)
		for i := 0; i < len(payload.Events); i++ {
			var bt []byte
			var out bytes.Buffer
			footer := fmt.Sprintf("%s | %s", payload.Events[i].Id, payload.Events[i].CreatedAt)

			bt, err = json.Marshal(payload.Events[i].Data)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			err = json.Indent(&out, bt, "", "  ")
			if err != nil {
				errors = append(errors, err)
				continue
			}

			bs := "```json\n" + string(out.Bytes()) + "\n```"
			if err = ds(vars["id"], "16731212", vars["token"], &payload.Events[i].Type, &bs, &footer); err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			e, err := json.Marshal(errors)
			if err != nil {
				e = []byte("Multiple errors detected")
			}
			http.Error(w, string(e), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
