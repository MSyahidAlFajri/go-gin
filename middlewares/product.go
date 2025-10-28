package middlewares

import (
	"log"

	"github.com/MSyahidAlFajri/go-gin/database"
	"github.com/MSyahidAlFajri/go-gin/models"
	"github.com/gin-gonic/gin"
)

// ProductMigrationMiddleware menjalankan migrasi tabel produk jika belum ada.
func ProductMigrationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.DB
		// AutoMigrate dari GORM — membuat tabel, kolom, index jika belum ada
		if err := db.AutoMigrate(&models.Product{}); err != nil {
			log.Printf("Failed to auto-migrate Product model: %v", err)
			// bisa tetap lanjut tapi log sebagai warning
		}
		c.Next()
	}
}
