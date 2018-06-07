package pinentry

import "testing"

func TestPinentry(t *testing.T) {
	c, err := NewPinentryClient()
	if err != nil {
		panic(err)
	}
	c.SetDesc("Type your passphrase:")
	c.SetPrompt("PIN:")
	c.GetPin()
	c.Close()
}
