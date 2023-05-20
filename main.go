package main

import (
	"fmt"
	"os"
	"os/exec"
	"shereyascodes-tech/gh-back/gh"
	"shereyascodes-tech/gh-back/git"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var DEBUG = envOrDefault("DEBUG", "false") == "true"

func fail(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		if DEBUG {
			panic(err)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func env(key string) string {
	res := os.Getenv(key)
	if res == "" {
		fail(fmt.Errorf("error: ENV %s not found", key), fmt.Sprintf("error: ENV %s not found", key))
	}
	return res
}

func envOrDefault(key string, def string) string {
	res := os.Getenv(key)
	if res == "" {
		return def
	}
	return res
}

func main() {

	AUTH_TOKEN := env("AUTH_TOKEN")
	OWNER := envOrDefault("OWNER", "")
	REPO := env("REPO")
	PRIVATE := envOrDefault("PRIVATE", "false") == "true"
	CMD := strings.Split(env("CMD"), " ")
	DIR := env("DIR")

	user, err := gh.Get_user(AUTH_TOKEN)
	fail(err, "error: User not found")

	err = git.Config(user.Login, user.Email)
	fail(err, "error: Could not configure git")
	fmt.Println(user.Login, user.Email)

	if OWNER == "" {
		OWNER = user.Login
	}

	repo, err := gh.Get_Repo(AUTH_TOKEN, OWNER, REPO)
	if err != nil {
		fmt.Println("error: Repo not found")
		fmt.Println("Creating repo...")

		repo, err = gh.Create_repo(AUTH_TOKEN, user.Login, OWNER, REPO, PRIVATE)
	}
	fail(err, "error: Could not create repo")

	full_name := repo["full_name"].(string)

	// Clone the repo
	err = git.Exec(
		"",
		"clone",
		fmt.Sprintf("https://%s:%s@github.com/%s.git", user.Login, AUTH_TOKEN, full_name),
		DIR,
	)
	fail(err, "error: Could not clone repo")

	// install and track all files under git lfs
	err = git.Exec(
		DIR,
		"lfs",
		"install",
		"--skip-smudge",
	)
	fail(err, "error: Could not install git lfs")

	// Delete .gitattributes file
	os.Remove(fmt.Sprintf("%s/.gitattributes", DIR))

	err = git.Exec(
		DIR,
		"lfs",
		"track",
		"*",
	)
	fail(err, "error: Could not track files with git lfs")

	// Delete .gitattributes file
	err = os.Remove(fmt.Sprintf("%s/.gitattributes", DIR))
	fail(err, "error: Could not delete .gitattributes file")

	watcher, err := fsnotify.NewWatcher()
	fail(err, "error: Could not create watcher")
	defer watcher.Close()

	changed := false
	first := true
	timer := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				{

					if !ok {
						return
					}
					changed = true
				}
			case err, ok := <-watcher.Errors:
				{
					if !ok {
						return
					}
					fmt.Println("error:", err)
				}
			case <-timer.C:
				{
					if first {
						first = false
						continue
					}
					if changed {
						changed = false
						now := time.Now().In(time.FixedZone("IST", 19800)).Format("02-01-2006 03:04:05")
						fmt.Printf("Committing changes... %s\n", now)

						err = git.CommitAll(
							DIR,
							fmt.Sprintf("Committed at %s", now),
						)
						if err != nil {
							fmt.Println(err)
							fmt.Println("error: Could not commit files")
							continue
						}

						err = git.Push(DIR)
						fail(err, "error: Could not push files")
					}
				}
			}
		}
	}()

	watcher.Add(DIR)

	// Run the command
	cmd := exec.Command(CMD[0], CMD[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
