package handler

import (
	"errors"

	"echobackend/internal/service"
	"echobackend/pkg/response"

	"github.com/labstack/echo/v5"
)

type ExchangeRateHandler struct {
	exchangeRateService service.ExchangeRateService
}

func NewExchangeRateHandler(exchangeRateService service.ExchangeRateService) *ExchangeRateHandler {
	return &ExchangeRateHandler{exchangeRateService: exchangeRateService}
}

func (h *ExchangeRateHandler) GetRate(c *echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	result, err := h.exchangeRateService.GetRate(c.Request().Context(), from, to)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCurrencyPair) {
			return response.BadRequest(c, "Invalid currency pair", err)
		}
		return response.InternalServerError(c, "Failed to get exchange rate", err)
	}

	return response.Success(c, "Exchange rate fetched successfully", result)
}
