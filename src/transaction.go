package main

type Transaction struct {
	Id 					int 		`json:"id"`
	Description string	`json:"description"`
	BuyMoney 	bool			`json:"buyMoney"`
	Money 	float32			`json:"money"`
	Owner_Id 	float64		`json:"owner_id"`
	Match_Id 	float64		`json:"match_id"`
}