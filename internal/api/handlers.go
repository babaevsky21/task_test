package api

import (
	"crypto-watcher/internal/models"
	"crypto-watcher/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	cryptoService *service.CryptoService
}

func SetupRouter(cryptoService *service.CryptoService) *gin.Engine {
	r := gin.Default()

	handler := &Handler{cryptoService: cryptoService}

	// Swagger документация
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := r.Group("/currency")
	{
		api.POST("/add", handler.AddCoin)
		api.DELETE("/remove", handler.RemoveCoin)
		api.GET("/price", handler.GetPrice)
	}

	// Запуск фоновой коллекции цен
	go cryptoService.SimulatePriceCollection()

	return r
}

// AddCoin добавляет криптовалюту в список наблюдения
// @Summary Добавить криптовалюту в watchlist
// @Description Добавляет криптовалюту в список наблюдения
// @Tags currency
// @Accept json
// @Produce json
// @Param request body models.AddCoinRequest true "Данные криптовалюты"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /currency/add [post]
func (h *Handler) AddCoin(c *gin.Context) {
	var req models.AddCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.cryptoService.AddCoin(req.Coin); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coin added successfully"})
}

// RemoveCoin удаляет криптовалюту из списка наблюдения
// @Summary Удалить криптовалюту из watchlist
// @Description Удаляет криптовалюту из списка наблюдения
// @Tags currency
// @Accept json
// @Produce json
// @Param request body models.AddCoinRequest true "Данные криптовалюты"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /currency/remove [delete]
func (h *Handler) RemoveCoin(c *gin.Context) {
	var req models.AddCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.cryptoService.RemoveCoin(req.Coin); err != nil {
		if err.Error() == "coin not found in watchlist" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coin removed successfully"})
}

// GetPrice получает цену криптовалюты на определенный момент времени
// @Summary Получить цену криптовалюты
// @Description Получает цену криптовалюты на указанный момент времени
// @Tags currency
// @Accept json
// @Produce json
// @Param coin query string true "Название криптовалюты"
// @Param timestamp query int true "Unix timestamp"
// @Success 200 {object} models.PriceHistory
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /currency/price [get]
func (h *Handler) GetPrice(c *gin.Context) {
	coin := c.Query("coin")
	timestampStr := c.Query("timestamp")

	if coin == "" || timestampStr == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "coin and timestamp parameters are required"})
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid timestamp format"})
		return
	}

	price, err := h.cryptoService.GetPrice(coin, timestamp)
	if err != nil {
		if err.Error() == "coin not found in watchlist" || err.Error() == "price not found for the given timestamp" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, price)
}
