// Simple web server that will listen for requests and then will execute git pull
// on the host machine. Used as an API endpoint for GitHub webhooks that will effectively
// deploy code when a merged pull request is completed on the origin repo into master.

package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ContentDir is a custom string type so we can use
// a string as a receiver parameter for HTTP handlers.
type ContentDir string

// String returns the ContentDir as a string.
func (cd ContentDir) String() string {
	return string(cd)
}

func main() {

	addr := require("ADDR")
	tlsKey := require("TLSKEY")
	tlsCert := require("TLSCERT")

	gitRepo := require("AUTO_UPDATE_GIT_REPO")
	workDir := ContentDir(require("AUTO_UPDATE_CONTENT_DIR"))

	ensureGitRepo(workDir.String(), gitRepo)

	http.HandleFunc("/update", workDir.UpdateHandler)
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

// ensureGitRepo ensures that the given directory is
// ready to pull from the given git repo.
func ensureGitRepo(dir, gitRepo string) {
	dir = strings.TrimSuffix(dir, "/")

	// Ensure the directory is a git repo.
	if _, err := os.Stat(dir + "/.git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		_, err = cmd.Output()
		if err != nil {
			log.Fatalf("error initializing empty git repo: %v", err)
		}
	}

	// Check the git remote properties.
	remoteBytes, err := execCommand(dir, "git", "remote", "-v")
	if err != nil {
		log.Fatalf("error running git remote -v: %v", err)
	}

	remoteString := string(remoteBytes)

	// Make sure an origin remtoe exists.
	if !strings.Contains(remoteString, "origin") {
		_, err = execCommand(dir, "git", "remote", "add", "origin", gitRepo)
		if err != nil {
			log.Fatalf("error adding remote origin: %v", err)
		}

	}

	// Get the value of the origin remote.
	getOriginOutput, err := execCommand(dir, "git", "remote", "get-url", "origin")
	if err != nil {
		log.Fatalf("error getting url of origin remote")
	}

	// If the origin remote is not the git repo, set it to be the git repo.
	if string(getOriginOutput) != gitRepo {
		_, err = execCommand(dir, "git", "remote", "set-url", "origin", gitRepo)
		if err != nil {
			log.Fatalf("error setting remote origin url: %v", err)
		}
	}

	// Fetch the latest master objects and refs from origin.
	_, err = execCommand(dir, "git", "fetch", "origin", "master")
	if err != nil {
		log.Fatalf("error fetching origin master: %v", err)
	}

	// Reset the current git repo to the latest commit on origin/master.
	// this is necessary since the local git repo could be created over existing
	// content and it would cause the entire directory to be included in a diff
	// and cause all further git pulls to fail. This will also force
	_, err = execCommand(dir, "git", "reset", "--hard", "origin/master")
	if err != nil {
		log.Fatalf("error resetting local git repo to latest origin master sha: %v", err)
	}
}

// UpdateHandler handles requests to update the currently served client code.
func (cd ContentDir) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = cd.String()
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "error executing command: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(output)
	if err != nil {
		log.Printf("error writing output to client: %v", err)
	}
}

// execCommand runs the given command args in the given dir
// and returns the output. Returns an error if one occurred.
func execCommand(dir string, args ...string) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("args must have length of at least 1")
	}
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}
	cmd.Dir = dir
	return cmd.Output()
}
