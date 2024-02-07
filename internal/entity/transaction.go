package entity

import "time"

// @Description Денежный перевод
type Transaction struct {
	Time   time.Time `json:"time"   example:"2024-02-04T17:25:35.448Z"         description:"Дата и время перевода"  validate:"required" format:"date-time"`
	From   string    `json:"from"   example:"5b53700ed469fa6a09ea72bb78f36fd9" description:"ID исходящего кошелька" validate:"required" pg:"from_wallet_id"`
	To     string    `json:"to"     example:"eb376add88bf8e70f80787266a0801d5" description:"ID входящего кошелька"  validate:"required" pg:"to_wallet_id"`
	Amount float64   `json:"amount" example:"30.0"                             description:"Сумма перевода"         validate:"required" format:"float" minimum:"0.0"`
}

// @Description Запрос перевода средств
type TransactionRequest struct {
	To     string  `json:"to"     example:"eb376add88bf8e70f80787266a0801d5" description:"ID кошелька, куда нужно перевести деньги" validate:"required"`
	Amount float64 `json:"amount" example:"100.0"                            description:"Сумма перевода"                           validate:"required" format:"float" minimum:"0.0"`
}