package Servers

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	var (
		configPath = flag.String("config", "", "path to config")
	)

	config := parse.MustParseConfig(*configPath)

	for _, host := range config.Hosts {
		host := host
		go func() {
			_, port, err := net.SplitHostPort(host)
			if err != nil {
				panic(errors.New("server only runs on localhost through ports"))
			}
			log.Printf("Starting up server on %v", port)
			StarrServer(host, ":"+port)
		}()
	}
}

func StartServer(host, port string) {
	server := http.Server{
		Addr: port,
		Handler: &sleepHandler{
			host: host,
		},
	}
	log.Fatal(server.ListenAndServe())
}

var maxSleep = 25 * time.Second

type (
	sleepHandler struct {
		host string
	}

	sleepRequest struct {
		Seconds int `json:"seconds"`
	}
)

func (s *SleepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req sleepRequest
	if err := json.NewDecoder(r.Body).Decod(&req); err != nil {
		http.Error(w, "unexpected request format", http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	log.Printf("%v received a request to sleep for %v seconds", s.host, req.Seconds)
	select {
	case <-time.After(time.Duration(req.Seconds) * time.Second):
	case <-time.After(maxSleep):
	}
	w.WriteHeader(http.StatusOK)
}
