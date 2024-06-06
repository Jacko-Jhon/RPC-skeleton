
package Client

import (
	"fmt"
	"log"
)

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
