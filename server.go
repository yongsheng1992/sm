package sm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	DB  map[string]*Trie
	AOF *AOF
}

type SearchRequest struct {
	Name   string   `json:"name"`
	Key    []string `json:"key"`
	Option string   `json:"option"`
	Limit  int      `json:"limit"`
}

func (server *Server) Insert(name string, key []byte, value interface{}) error {
	trie, ok := server.DB[name]
	if !ok {
		return errors.New(fmt.Sprintf("trie name `%s` not found", name))
	}
	trie.Insert(key, value)
	server.AOF.Feed(ConvertInsert(name, string(key), ""))
	return nil
}

func (server *Server) Remove(name string, key string) error {
	trie, ok := server.DB[name]
	if !ok {
		return errors.New(fmt.Sprintf("trie name `%s` not found", name))
	}
	trie.Remove([]byte(key))
	server.AOF.Feed(ConvertRemove(name, key))
	return nil
}

func (server *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	var searchResponse map[string][]string

	searchResponse = make(map[string][]string)

	if err := json.NewDecoder(r.Body).Decode(&searchRequest); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if searchRequest.Name == "" {
		http.Error(w, "name is required", 400)
		return
	}

	if searchRequest.Limit == 0 {
		searchRequest.Limit = 10
	}

	trie, ok := server.DB[searchRequest.Name]

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
				k, node := it.Next()
				if node.IsKey {
					searchResponse[key] = append(searchResponse[key], string(k))
				}
			}
		}

	}

	if err := json.NewEncoder(w).Encode(searchResponse); err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func (server *Server) HandleKeyInsert(w http.ResponseWriter, r *http.Request) {
	var postData []string

	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	params := mux.Vars(r)
	name := params["name"]

	_, ok := server.DB[name]
	if !ok {
		server.DB[name] = NewTrie()
	}

	for _, key := range postData {
		if err := server.Insert(name, []byte(key), nil); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(make(map[string]interface{})); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (server *Server) HandleKeyRemove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	key := params["key"]

	if err := server.Remove(name, key); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(make(map[string]interface{})); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type KeyGetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (server *Server) HandleKeyGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	key := params["key"]

	_, ok := server.DB[name]

	if !ok {
		http.Error(w, "trie not found", 404)
		return
	}

	ret, value := server.DB[name].Find([]byte(key))

	if !ret {
		http.Error(w, "key not found", 404)
		return
	}

	var response KeyGetResponse
	response = KeyGetResponse{
		Key:   key,
		Value: value,
	}

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (server *Server) InitHTTPServer() {

	r := mux.NewRouter()
	r.HandleFunc("/api/search", server.HandleSearch).Methods(http.MethodPost)
	r.HandleFunc("/api/{name}", server.HandleKeyInsert).Methods(http.MethodPost)
	r.HandleFunc("/api/{name}/{key}", server.HandleKeyRemove).Methods(http.MethodDelete)
	r.HandleFunc("/api/{name}/{key}", server.HandleKeyGet).Methods(http.MethodGet)

	go func() {
		fmt.Println("Init HTTP Server...")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func NewServer() *Server {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := &Server{}
	server.DB = make(map[string]*Trie)
	// when receive a signal, server.AOF.Close() should be called.
	server.AOF = NewAOF("aof")

	return server
}
