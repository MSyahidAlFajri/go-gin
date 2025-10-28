package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MSyahidAlFajri/go-gin/database"
	"github.com/MSyahidAlFajri/go-gin/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateProductRequest untuk membuat produk baru
type CreateProductRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=255"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Category string  `json:"category" binding:"required,min=1,max=100"`
}

// UpdateProductRequest untuk update produk
type UpdateProductRequest struct {
	Name     *string  `json:"name"`
	Price    *float64 `json:"price"`
	Category *string  `json:"category"`
}

// GetAllProducts fetch semua produk, dengan optional filter category
func GetAllProducts(c *gin.Context) {
	var products []models.Product
	db := database.DB

	// optional: filter berdasarkan category
	if cat := c.Query("category"); cat != "" {
		db = db.Where("category = ?", strings.TrimSpace(cat))
	}

	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID fetch produk berdasarkan ID
func GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not found",
				"message": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"message": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct handler membuat produk baru
func CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	// trim inputs
	req.Name = strings.TrimSpace(req.Name)
	req.Category = strings.TrimSpace(req.Category)

	prod := models.Product{
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
	}

	if err := database.DB.Create(&prod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Create failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, prod)
}

// UpdateProduct handler update produk
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not found",
				"message": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"message": err.Error(),
			})
		}
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	if req.Name != nil {
		product.Name = strings.TrimSpace(*req.Name)
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Category != nil {
		product.Category = strings.TrimSpace(*req.Category)
	}

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Update failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct handler untuk menghapus produk
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	result := database.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Delete failed",
			"message": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not found",
			"message": "Product not found or already deleted",
		})
		return
	}

	// Atau jika Anda ingin mengembalikan pesan sukses:
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
