package systeminteraction

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"saidvandeklundert/agrippa/agrippalogger"
	"time"
)

/*
Contains the command that was executed using exec as well as the stdout and stderr
*/
type CommandResponse struct {
	Command []string
	Stdout  string
	Stderr  string
}

/*
Runs target command on the local system and returns a CommandResponse.
*/
func RunCommand(command ...string) (CommandResponse, error) {
	log := agrippalogger.GetLogger()
	log.Debug("running command ", command)

	if len(command) == 0 {
		log.Error("no command given")
		return CommandResponse{}, errors.New("RunCommand requires a command but no command given")
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	// accomodate commands with and without arguments
	var cmd *exec.Cmd
	if len(command) == 1 {
		cmd = exec.Command(command[0])

	} else {
		cmd = exec.Command(command[0], command[1:]...)
	}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Run(); err != nil {
		return CommandResponse{}, err
	}
	stdoutString := stdoutBuf.String()
	stderrString := stderrBuf.String()
	log.Info("stdout ", stdoutString)
	log.Info("stderr ", stderrString)
	response := CommandResponse{
		Command: command,
		Stdout:  stdoutString,
		Stderr:  stderrString,
	}
	return response, nil
}

func runCommand(ctx context.Context, command ...string) (CommandResponse, error) {
	log := agrippalogger.GetLogger()
	log.Debug("running command ", command)

	if len(command) == 0 {
		log.Error("no command given")
		return CommandResponse{}, errors.New("RunCommand requires a command but no command given")
	}

	var stdoutBuff, stderrBuff bytes.Buffer

	// accomodate commands with or without arguments
	var cmd *exec.Cmd
	if len(command) == 1 {
		cmd = exec.CommandContext(ctx, command[0])
	} else {
		cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	}
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Start()
	if err != nil {
		return CommandResponse{}, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			log.Error("failed to kill the process, ", err)
		}
		log.Error("failed to complete the command in time")
		return CommandResponse{}, context.DeadlineExceeded
	case err := <-done:
		if err != nil {
			return CommandResponse{}, err
		}
	}

	stdoutString := stdoutBuff.String()
	stderrString := stderrBuff.String()
	log.Info("stdout ", stdoutString)
	log.Info("stderr ", stderrString)
	response := CommandResponse{
		Command: command,
		Stdout:  stdoutString,
		Stderr:  stderrString,
	}
	return response, nil

}

/*
Runs a command on the system guarded by a timeout. If the execution of the
command takes longer then the given time.Duration, the command is killed
*/
func RunCommandWithTimeout(timeout time.Duration, command ...string) (CommandResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return runCommand(ctx, command...)

}
