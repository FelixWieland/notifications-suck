package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Notification struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

var notifications = []Notification{}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		/create?message=MYMESSAGE  	-- creates a notification
		/all				-- returns all created notifications
		/from?from=UNIXTIMESTAMP  	-- returns all created notifications after given timestamp
		`))
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["message"]
		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("message is missing"))
			return
		}
		message := keys[0]

		n := createNotification(message)
		b, err := json.Marshal(n)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	})

	http.HandleFunc("/from", func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["from"]
		if !ok || len(keys[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("from is missing"))
			return
		}

		fromRaw := keys[0]
		from, _ := strconv.Atoi(fromRaw)

		ns := []Notification{}
		for _, n := range notifications {
			if n.Time.Unix() > int64(from) {
				ns = append(ns, n)
			}
		}

		b, err := json.Marshal(ns)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(b)
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(notifications)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
	})

	if len(port) == 0 {
		port = "80"
	}
	http.ListenAndServe(":"+port, http.DefaultServeMux)
}

func createNotification(message string) Notification {
	n := Notification{
		Time:    time.Now(),
		Message: message,
	}
	notifications = append(notifications, n)
	return n
}
