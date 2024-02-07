package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egor-denisov/wallet-infotecs/internal/entity"
	"github.com/egor-denisov/wallet-infotecs/internal/usecase"
	"github.com/egor-denisov/wallet-infotecs/pkg/logger"
)

type walletRoutes struct {
	w usecase.Wallet
	l logger.Interface
}

func newWalletRoutes(handler *gin.RouterGroup, w usecase.Wallet, l logger.Interface) {
	r := &walletRoutes{w, l}

	h := handler.Group("/wallet")
	{
		h.POST("", r.createNewWallet)
		h.POST("/:walletId/send", r.sendFunds)
		h.GET("/:walletId/history", r.getWalletHistoryById)
		h.GET("/:walletId", r.getWalletById)
	}
}

// @Summary     Создание кошелька
// @Description Создает новый кошелек с уникальным ID. Идентификатор генерируется сервером.
// @Description
// @Description Созданный кошелек должен иметь сумму 100.0 у.е. на балансе
// @Tags  	    Wallet
// @Success     200 {object} entity.Wallet "Кошелек создан"
// @Failure     400 "Ошибка в запросе"
// @Router      /wallet [post]
func (r *walletRoutes) createNewWallet(c *gin.Context) {
	wallet, err := r.w.CreateNewWalletWithDefaultBalance(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - createNewWallet")
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, wallet)
}

// @Summary     Перевод средств с одного кошелька на другой
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Param input body entity.TransactionRequest true "Запрос перевода средств"
// @Success     200 "Перевод успешно проведен"
// @Failure     404 "Исходящий кошелек не найден"
// @Failure     400 "Ошибка в пользовательском запросе или ошибка перевода"
// @Router      /wallet/{walletId}/send [post]
func (r *walletRoutes) sendFunds(c *gin.Context) {
	var TransactionRequest entity.TransactionRequest

	if err := c.BindJSON(&TransactionRequest); err != nil {
		r.l.Error(err, "http - v1 - sendFunds")
		c.Status(http.StatusBadRequest)

		return
	}

	err := r.w.SendFunds(c.Request.Context(), c.Param("walletId"), TransactionRequest.To, TransactionRequest.Amount)
	if errors.Is(err, entity.ErrWalletNotFound) {
		r.l.Error(err, "http - v1 - sendFunds")
		c.Status(http.StatusNotFound)

		return
	}
	if err != nil {
		r.l.Error(err, "http - v1 - sendFunds")
		c.Status(http.StatusBadRequest)

		return
	}

	c.Status(http.StatusOK)
}

// @Summary     Получение историй входящих и исходящих транзакций
// @Description Возвращает историю транзакций по указанному кошельку.
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Success     200 {object} []entity.Transaction "История транзакций получена"
// @Failure     404 "Указанный кошелек не найден"
// @Router      /wallet/{walletId}/history [get]
func (r *walletRoutes) getWalletHistoryById(c *gin.Context) {
	transactions, err := r.w.GetWalletHistoryById(c.Request.Context(), c.Param("walletId"))
	if err != nil {
		r.l.Error(err, "http - v1 - getWalletHistoryById")
		c.AbortWithStatus(http.StatusNotFound)

		return
	}
	
	c.JSON(http.StatusOK, transactions)
}

// @Summary     Получение текущего состояния кошелька
// @Tags  	    Wallet
// @Param walletId path string true "ID кошелька"
// @Success     200 {object} entity.Wallet "OK"
// @Failure     404 "Указанный кошелек не найден"
// @Router      /wallet/{walletId} [get]
func (r *walletRoutes) getWalletById(c *gin.Context) {
	wallet, err := r.w.GetWalletById(c.Request.Context(), c.Param("walletId"))
	if err != nil {
		r.l.Error(err, "http - v1 - getWalletById")
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	c.JSON(http.StatusOK, wallet)
}