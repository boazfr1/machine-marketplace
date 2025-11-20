package order

import "fmt"

func New(port string) error {
	fmt.Printf("Order service listening on port %s\n", port)
	return nil
}