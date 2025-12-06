package process

import (
	"fmt"
	db "machine-marketplace/internal/DB/generated"
)

func New(port string, queries *db.Queries) error {
	fmt.Printf("Process service listening on port %s\n", port)
	return nil
}