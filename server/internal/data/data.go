package data

type CreditCard struct {
	Id             int    `json:"id"`
	OwnerId        int    `json:"owner-id"`
	Number         int    `json:"credit-number"`
	ExpirationDate string `json:"expiration-date"`
	Secret         int    `json:"secret"`
}

type Machine struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	BuyerId int    `json:"buyer-id"`
	OwnerId int    `json:"owner-id"`
	RAM     int    `json:"RAM"`
	CPU     int    `json:"CPU"`
	Memory  int    `json:"memory"`
	Key     string `json:"key"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
