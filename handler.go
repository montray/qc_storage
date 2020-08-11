package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"sync/atomic"
)

type StoreHandler struct {
	rCount         uint64
	storageService StorageService
	v              *validator.Validate
}

func NewStoreHandler(service StorageService) *StoreHandler {
	return &StoreHandler{
		storageService: service,
		v:              validator.New(),
	}
}

func (sh *StoreHandler) Store(c *gin.Context) {
	type StoreRequest struct {
		ProductId int `json:"product_id" binding:"required" validate:"gte=1"`
		Quantity  int `json:"quantity" validate:"required,gte=1"`
	}

	atomic.AddUint64(&sh.rCount, 1)
	defer atomic.StoreUint64(&sh.rCount, sh.rCount-1)

	if sh.rCount >= 4 {
		c.AbortWithStatus(429)
		return
	}

	sr := StoreRequest{}
	c.ShouldBindJSON(&sr)

	err := sh.v.Struct(sr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "product_id": nil, "quantity": nil})
		return
	}

	qty, err := sh.storageService.Take(sr.ProductId, sr.Quantity)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "product_id": sr.ProductId, "quantity": qty})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "product_id": sr.ProductId, "quantity": qty})
}

func (sh *StoreHandler) Add(c *gin.Context) {
	type AddRequest struct {
		ProductId int `json:"product_id" binding:"required" validate:"gte=1"`
		Quantity  int `json:"quantity" validate:"required,gte=1"`
	}

	ar := AddRequest{}
	c.ShouldBindJSON(&ar)

	err := sh.v.Struct(ar)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "product_id": nil, "quantity": nil})
		return
	}

	qty := sh.storageService.Add(ar.ProductId, ar.Quantity)

	c.JSON(http.StatusOK, gin.H{"message": "success", "product_id": ar.ProductId, "quantity": qty})
}

func (sh *StoreHandler) Get(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil || productId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "product_id must be integer and greater than 0"})
	}

	qty := sh.storageService.Get(productId)

	c.JSON(http.StatusOK, gin.H{"message": "success", "product_id": productId, "quantity": qty})
}
