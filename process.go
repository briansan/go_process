// Helper module for forking and monitoring process
package process

import (
	"os"
	"os/exec"
	"sync"
)

// Process monitor
type ProcessMonitor struct {
	CmdName *string
	CmdArgs *[]string
	Process *os.Process
	Cmd     *exec.Cmd
	Output  *[]byte
	Err     error
}

// Process state listener interface
type ProcessStateListener interface {
	OnComplete(processMonitor *ProcessMonitor)
	OnError(processMonitor *ProcessMonitor, err error)
}

// Method to fork a process for given command
// and return ProcessMonitor
func Fork(processStateListener ProcessStateListener, cmdName string, cmdArgs ...string) (*ProcessMonitor, *sync.WaitGroup) {
	processMonitor := &ProcessMonitor{
		CmdName: &cmdName,
		CmdArgs: &cmdArgs,
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		command := exec.Command(cmdName, cmdArgs...)
		processMonitor.Process = command.Process
		processMonitor.Cmd = command

		output, err := command.CombinedOutput()
		if err != nil {
			processMonitor.Err = err
			processStateListener.OnError(processMonitor, err)
		}
		processMonitor.Output = &output
		processStateListener.OnComplete(processMonitor)
		wg.Done()
	}()
	return processMonitor, wg
}
