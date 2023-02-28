/* 
 :: ServiceMonitoring
 :: Very basic. will add more to this file later
 :: Contributed by @slimdestro
*/

package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Destro!")
	})

	// Start the server
	go http.ListenAndServe(":8080", nil)
 
	for {
		resp, err := http.Get("http://localhost:8080")
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Status:", resp.StatusCode)
		}
		time.Sleep(time.Second * 5)
	}
}