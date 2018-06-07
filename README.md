```
package main

import (
	"fmt"

	"github.com/xaionaro-go/pinentry"
)

func main() {
	client, _ := pinentry.NewPinentryClient()
	client.SetTitle("Some title here")
	client.SetDesc("Some description here")
	client.SetPrompt("Enter the passphrase, please:")
	client.SetOK("Ok")
	p, _ := client.GetPin()
	fmt.Println(string(p))
	client.Close()
}
```
