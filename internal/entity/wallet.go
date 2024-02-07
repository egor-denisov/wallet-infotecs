package entity

// @Description Состояние кошелька
type Wallet struct {
	ID      string  `json:"id"       example:"5b53700ed469fa6a09ea72bb78f36fd9" description:"Уникальный ID кошелька" validate:"required"`
	Balance float64 `json:"balance"  example:"100.0"                            description:"Баланс кошелька"        validate:"required" minimum:"0.0" format:"float"`
}