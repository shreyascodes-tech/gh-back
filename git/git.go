package git

import (
	"fmt"
	"os"
	"os/exec"
)

func Exec(path string, args ...string) error {
	new_args := []string{"git"}

	if path != "" {
		new_args = append(new_args, "-C", path)
	}

	new_args = append(new_args, args...)
	op, err := exec.Command(new_args[0], new_args[1:]...).CombinedOutput()
	if os.Getenv("DEBUG") == "true" {
		fmt.Print(string(op))
	}

	if err != nil {
		return err
	}

	return nil
}

func Clone(url string) error {
	return Exec("", "clone", url)
}

func Config(name string, email string) error {
	err := Exec("", "config", "--global", "user.name", name)
	if err != nil {
		return err
	}

	err = Exec("", "config", "--global", "user.email", email)
	return err
}

func CommitAll(path string, message string) error {
	err := Exec(path, "add", ".")

	if err != nil {
		fmt.Println("error: Could not add files to repo")
		return err
	}

	err = Exec(path, "pull")
	if err != nil {
		fmt.Println("error: Could not pull from repo")
		return err
	}

	err = Exec(path, "commit", "-m", fmt.Sprintf("\"%s\"", message))
	return err
}

func Push(path string) error {
	return Exec(path, "push")
}
