package pinentry

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

var (
	PinentryUtilityName = "pinentry"
)

type PinentryClient interface {
	SetDesc(desc string)
	SetPrompt(prompt string)
	SetTitle(title string)
	SetOK(ok string)
	SetCancel(cancel string)
	SetError(errorMsg string)
	SetQualityBar()
	SetQualityBarTT(tt string)
	GetPin() (pin []byte, err error)
	Confirm() bool
	Close()
}

type pinentryClient struct {
	in   io.WriteCloser
	pipe *bufio.Reader
}

// set descriptive text to display
func (c *pinentryClient) SetDesc(desc string) {
	c.in.Write([]byte("SETDESC " + desc + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

// set desciption for user
func (c *pinentryClient) SetPrompt(prompt string) {
	c.in.Write([]byte("SETPROMPT " + prompt + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetTitle(title string) {
	c.in.Write([]byte("SETTITLE " + title + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetOK(okLabel string) {
	c.in.Write([]byte("SETOK " + okLabel + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetCancel(cancelLabel string) {
	c.in.Write([]byte("SETCANCEL " + cancelLabel + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetError(errorMsg string) {
	c.in.Write([]byte("SETERROR " + errorMsg + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetQualityBar() {
	c.in.Write([]byte("SETQUALITYBAR\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) SetQualityBarTT(tt string) {
	c.in.Write([]byte("SETQUALITYBAR_TT" + tt + "\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) != 0 {
		panic(string(ok))
	}
}

func (c *pinentryClient) Confirm() bool {
	confirmed := false
	c.in.Write([]byte("CONFIRM\n"))
	// ok
	ok, _, _ := c.pipe.ReadLine()
	if bytes.Compare(ok, []byte("OK")) == 0 {
		confirmed = true
	}
	return confirmed
}

func (c *pinentryClient) GetPin() (pin []byte, err error) {
	c.in.Write([]byte("GETPIN\n"))
	// D pin
	d_pin, _, err := c.pipe.ReadLine()
	if bytes.Compare(d_pin[:2], []byte("D ")) == 0 {
		ok, _, _ := c.pipe.ReadLine()
		if bytes.Compare(ok, []byte("OK")) != 0 {
			panic(string(ok))
		}
		return d_pin[2:], nil
	} else if bytes.Compare(d_pin[:2], []byte("OK")) == 0 {
		return nil, nil
	}
	return nil, fmt.Errorf("unexpected response for GetPin: %s", d_pin)
}

func (c *pinentryClient) Close() {
	c.in.Close()
	return
}

func startProcess(cmdName string) (io.WriteCloser, *bufio.Reader, error) {
	if runtime.GOOS == "windows" {
		cmdName += ".exe"
	}
	cmd := exec.Command(cmdName, "-T", "/dev/tty")
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	bufout := bufio.NewReader(out)
	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}
	// welcome
	welcome, _, err := bufout.ReadLine()
	if err != nil {
		return nil, nil, err
	}

	if bytes.Compare(welcome[:2], []byte("OK")) != 0 {
		return nil, nil, fmt.Errorf("Invalid welcome message: %v", string(welcome))
	}

	return in, bufout, nil
}

func NewPinentryClient() (PinentryClient, error) {
	in, bufout, err := startProcess(PinentryUtilityName)
	if err != nil {
		return nil, err
	}

	pinentry := &pinentryClient{in, bufout}

	//Setup default layout
	pinentry.SetTitle("pinentry")
	pinentry.SetDesc("")
	pinentry.SetPrompt("Enter the passphrase, please:")
	pinentry.SetOK("Ok")

	return pinentry, nil
}
