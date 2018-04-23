// Simple web server that will listen for requests and then will execute git pull
// on the host machine. Used as an API endpoint for GitHub webhooks that will effectively
// deploy code when a merged pull request is completed on the origin repo into master.

package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	addr := require("ADDR")
	tlsKey := require("TLSKEY")
	tlsCert := require("TLSCERT")

	http.HandleFunc("/update", UpdateHandler)
	log.Println("update server is listenting on https://" + addr + "...")
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, nil))
}

// require gets and returns the environment variable value with the given name.
// If the given name is unset then the process will exit
func require(name string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		log.Fatalf("no value set for '%s'", name)
	}
	return val
}

// UpdateHandler handles requests to update the currently served client code.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = "/site/Website"
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "error executing command: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
