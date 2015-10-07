package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type LimitedWriter struct {
	N   int
	buf []byte
}

func (w *LimitedWriter) Write(data []byte) (int, error) {
	writable := len(data)
	if writable > w.N {
		writable = w.N
	}
	after := len(w.buf) + writable
	discard := 0
	if after > w.N {
		discard = after - w.N
	}
	w.buf = append(w.buf[discard:], data[len(data)-writable:]...)
	return len(data), nil
}

func (w *LimitedWriter) Bytes() []byte {
	return w.buf
}

func quiet(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("command argument was not supplied")
	}
	command := args[0]
	args = args[1:]
	if command == "-h" || command == "--help" {
		fmt.Printf(`quiet COMMAND [ARGS...]

quiet executes supplied command with its arguments, keeping the last QUIET_MAX
bytes of combined standard output and prints them only if the command fails.
QUIET_MAX defaults to 1024.

Originally used in crontab.
`)
		return nil
	}
	n := 1024
	nEnv := os.Getenv("QUIET_MAX")
	if nEnv != "" {
		nn, err := strconv.ParseInt(nEnv, 10, 32)
		if err != nil {
			return err
		}
		n = int(nn)
	}
	w := &LimitedWriter{
		N: n,
	}
	cmd := exec.Command(command, args...)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, string(w.Bytes()))
	}
	return err
}

func main() {
	err := quiet(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
