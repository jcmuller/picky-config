// Package dmenu is a simple wrapper for the amazing dynamic menu for X (dmenu) in Go.
// It also can use rofi if it's symlinked to the dmenu binary.
package dmenu

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

// Dmenu holds the structure of this thing
type Dmenu struct {
	command   string
	arguments []string
}

// EmptySelectionError is returned if there is no selection
type EmptySelectionError struct{}

func (e *EmptySelectionError) Error() string {
	return "Nothing selected"
}

// Popup pops up the menu
func Popup(prompt string, options ...string) (selection string, err error) {
	//selection, err = defaultDmenu().Popup(prompt, options...)
	dmenu := defaultDmenu()
	selection, err = dmenu.Popup(prompt, options...)
	return
}

func defaultDmenu() *Dmenu {
	return New("dmenu", "-p", "%s")
}

// New instance of dmenu
func New(command string, arguments ...string) *Dmenu {
	program, err := exec.LookPath(command)

	if err != nil {
		log.Fatalf("%s not found", command)
	}

	return &Dmenu{program, arguments}
}

// NewZenityList constructs a dmenu for zenity
func NewZenityList() *Dmenu {
	return New("zenity", "--title=%s", `--list`, `--column`, `%s`)
}

// NewZenityYesNo constructs a dmenu for zenity yes/no questions
func NewZenityYesNo() *Dmenu {
	return New("zenity", `--title`, `%s`, `--question`, `--text`, `%s`)
}

// YesNo returns true of false
func (d *Dmenu) YesNo(prompt string) (bool, error) {
	processedArgs := []string{}
	for _, arg := range d.arguments {
		var parg string
		if strings.Contains(arg, "%s") {
			parg = fmt.Sprintf(arg, prompt)
		} else {
			parg = arg
		}

		processedArgs = append(processedArgs, parg)
	}
	cmd := exec.Command(d.command, processedArgs...)

	if err := cmd.Start(); err != nil {
		fmt.Printf("cmd.Start: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					return false, nil
				} else {
					return false, err
				}
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

// Popup pops up the menu
func (d *Dmenu) Popup(prompt string, options ...string) (selection string, err error) {
	processedArgs := []string{}
	for _, arg := range d.arguments {
		var parg string
		if strings.Contains(arg, "%s") {
			parg = fmt.Sprintf(arg, prompt)
		} else {
			parg = arg
		}

		processedArgs = append(processedArgs, parg)
	}
	cmd := exec.Command(d.command, processedArgs...)

	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatalf("Error getting pipe: %s", err)
	}

	go func(stdin io.WriteCloser) {
		defer stdin.Close()
		io.WriteString(stdin, strings.Join(options, "\n"))
	}(stdin)

	byteOut, err := cmd.Output()

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					err = &EmptySelectionError{}
				}
			}
		}

		return "", err
	}

	// Cast and trim
	selection = strings.TrimSpace(string(byteOut))

	return
}
