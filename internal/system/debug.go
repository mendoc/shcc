package system

import (
	"fmt"
	"os"
)

// Debug affiche un message formaté si la variable d'environnement DEBUG est à "true"
func Debug(format string, a ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf("DEBUG: "+format+"\n", a...)
	}
}
