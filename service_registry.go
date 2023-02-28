/* 
 :: ServiceRegistry
 :: Contributed by @slimdestro
*/
package main

import (
	"fmt"
	"net/http"
)
type ServiceRegistry struct {
	services map[string]string
} 

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

func (s *ServiceRegistry) Register(name string, url string) {
	s.services[name] = url
}

func (s *ServiceRegistry) Get(name string) (string, bool) {
	url, ok := s.services[name]
	return url, ok
}

func main() {
	registry := NewServiceRegistry()
	registry.Register("foo", "http://server.com") 
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		url, ok := registry.Get(name)
		if !ok {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Service %s is at %s\n", name, url)
	})

	http.ListenAndServe(":8080", nil)
}