package api

import (
	"github.com/JadenRazo/Project-Website/backend/internal/messaging/domain"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, wordFilterUseCase domain.WordFilterUseCase) {
	wordFilterHandler := NewWordFilterHandler(wordFilterUseCase)

	// Word filter routes
	wordFilterGroup := router.Group("/api/word-filters")
	{
		wordFilterGroup.POST("", wordFilterHandler.CreateWordFilter)
		wordFilterGroup.PUT("/:id", wordFilterHandler.UpdateWordFilter)
		wordFilterGroup.DELETE("/:id", wordFilterHandler.DeleteWordFilter)
		wordFilterGroup.GET("/:id", wordFilterHandler.GetWordFilter)
		wordFilterGroup.GET("/server/:server_id", wordFilterHandler.ListWordFilters)
		wordFilterGroup.GET("/server/:server_id/scope/:scope", wordFilterHandler.GetWordFiltersByScope)
	}

	// ... existing routes ...
}
