package main

import (
	"fmt"

	"github.com/xaionaro-go/pinentry"
)

func main() {
	client, _ := pinentry.NewPinentryClient()
	client.SetTitle("Some title here")
	client.SetDesc("Some description here")
	client.SetPrompt("PIN")
	client.SetOK("OK")
	client.SetCancel("Cancel")
	p, _ := client.GetPin()
	fmt.Println(string(p))
	client.Confirm()
	client.Close()
}
