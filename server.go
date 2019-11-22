package sm

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	DB map[string]*Trie
}

type SearchRequest struct {
	Name   string   `json:"name"`
	Key    []string `json:"key"`
	Option string   `json:"option"`
	Limit  int      `json:"limit"`
}

func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	var searchResponse map[string][]string

	searchResponse = make(map[string][]string)

	if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
		http.Error(w, err.Error(), 400)
	}

	if searchRequest.Name == "" {
		http.Error(w, "name is required", 400)
	}

	if searchRequest.Limit == 0 {
		searchRequest.Limit = 10
	}

	trie, ok := s.DB[searchRequest.Name]

	if !ok {
		http.Error(w, "no trie found", 400)
	}

	for _, key := range searchRequest.Key {
		searchResponse[key] = make([]string, 0)

		switch searchRequest.Option {
		case "forward":
			flags := trie.SeekBefore([]byte(key))
			for _, idx := range flags {
				searchResponse[key] = append(searchResponse[key], string(key[0:idx]))
			}
		case "backward":
			it := trie.SeekAfter([]byte(key))
			count := 0
			for it.HasNext() && count < searchRequest.Limit {
				k, _ := it.Next()
				searchResponse[key] = append(searchResponse[key], string(k))
			}
		}
	}

	if err := json.NewEncoder(w).Encode(searchResponse); err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func (s *Server) HandleKeyInsert(w http.ResponseWriter, r *http.Request) {
	var postData []string

	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, err.Error(), 400)
	}

	params := mux.Vars(r)
	name := params["name"]

	_, ok := s.DB[name]

	if !ok {
		s.DB[name] = NewTrie()
	}

	for _, key := range postData {
		s.DB[name].Insert([]byte(key), nil)
	}

	if err := json.NewEncoder(w).Encode(make(map[string]interface{})); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (s *Server) InitHTTPServer() {

	r := mux.NewRouter()
	r.HandleFunc("/api/search", s.HandleSearch).Methods(http.MethodPost)
	r.HandleFunc("/api/{name}", s.HandleKeyInsert).Methods(http.MethodPost)

	go func() {
		fmt.Println("Init HTTP Server...")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func NewServer() *Server {
	server := &Server{}
	server.DB = make(map[string]*Trie)
	return server
}
