package cmd

import "machine-marketplace/pkg/order"



type OrderService struct {
	Port string `short:"p" long:"port" description:"Port to listen on" default:"3001"`
}

func (o *OrderService) Execute(args []string) error {
	err := order.New(o.Port)
	if err != nil {
		return err
	}
	return err
}