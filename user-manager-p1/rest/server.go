package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"user-manager/db"
	"user-manager/handlers"
	"user-manager/helper"
	"user-manager/services"

	"github.com/gorilla/mux"
)

type Server struct {
	Address   string
	Port      int
	channel   chan string
	guard     chan struct{}
	debugMode bool
	handler   *handlers.Handler
	wg        *sync.WaitGroup
}

func NewServer(addr string, port int) *Server {
	ch := make(chan string)

	max_goroutines, _ := strconv.Atoi(helper.Getenv("MAX_GOROUTINES", "5000"))

	guard := make(chan struct{}, max_goroutines)

	log.Printf("Using %d goroutines...\n", max_goroutines)

	srv := &Server{
		Address: addr,
		Port:    port,
		channel: ch,
		wg:      &sync.WaitGroup{},
		guard:   guard,
	}

	if os.Getenv("DEBUG") == "true" {
		srv.debugMode = true
	}
	// Initialize the store
	db := db.NewDBManager()
	db.Connect()
	svc := services.UserService{UserStore: db}
	srv.handler = handlers.NewUserHandler(svc, ch, srv.wg)

	lp, _ := strconv.Atoi(helper.Getenv("LOGGER_PROCESSORS", "10"))

	log.Println("Starting logging processes with", lp, "instances...")

	for i := 0; i < lp; i++ {
		go func() {
			for logMsg := range ch {
				log.Println(logMsg)
			}
		}()
	}

	return srv
}

// Start HTTP Server
func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/users", s.handler.AddUser).Methods("POST")
	router.HandleFunc("/users/{id}", s.handler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", s.handler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", s.handler.DeleteUser).Methods("DELETE")

	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)

	log.Println("Starting REST server on", addr)

	log.Fatal(http.ListenAndServe(addr, router))

	// Wait for all requests to complete
	s.wg.Wait()
	close(s.channel) // Close the channel after all requests are done
}
