package task

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"
)

type Task struct {
	Cmd string
	in  io.Reader
	Out io.ReadWriter
	Err io.ReadWriter

	next   *Task
	status int
	ctx    context.Context
}

func NewTask(cmd string) *Task {
	out := bytes.Buffer{}
	err := bytes.Buffer{}
	return &Task{
		Cmd: cmd,
		Err: &err,
		Out: &out,
	}
}

func (t *Task) Status() int {
	return t.status
}

func (t *Task) SetIn(in io.Reader) *Task {
	t.in = in
	return t
}
func (t *Task) SetErr(err io.ReadWriter) *Task {
	t.Err = err
	return t
}
func (t *Task) SetOut(out io.ReadWriter) *Task {
	t.Out = out
	return t
}

func (t *Task) Next(next *Task) *Task {
	t.next = next
	return t
}

func (t *Task) Run() error {
	split := strings.Split(t.Cmd, " ")
	exe := split[0]
	args := split[1:]
	cmd := exec.Command(exe, args...)
	cmd.Stderr = t.Err
	cmd.Stdout = t.Out

	if t.in != nil {
		cmd.Stdin = t.in
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	if t.next != nil {
		t.next.SetIn(t.Out)
		err = t.next.Run()
		if err != nil {
			return err
		}
		t.Out = t.next.Out
	}
	return nil
}
