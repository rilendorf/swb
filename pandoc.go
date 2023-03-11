package main

import (
	"fmt"
	"os/exec"
	"strings"

	"os"
)

func pandoc(out string, args ...string) error {
	f, err := os.OpenFile(out, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer f.Close()

	fmt.Printf("Exec: pandoc %s\n", strings.Join(args, " "))
	cmd := exec.Command("pandoc", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = f

	return cmd.Run()
}
