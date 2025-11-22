package process

import "fmt"

func New(port string) error {
	fmt.Printf("Process service listening on port %s\n", port)
	return nil
}