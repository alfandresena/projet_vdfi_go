package controllers

import (
	"net/http"

	"projet_vdfi/initializers"
	"projet_vdfi/models"

	"github.com/gin-gonic/gin"
)

// 🟡 LISTER TOUTES LES PAROLES (PUBLIC)
func GetLyrics(c *gin.Context) {
	var lyrics []models.Lyric
	if err := initializers.DB.Find(&lyrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de récupération des paroles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

// 🟠 OBTENIR LES PAROLES PAR ID (PUBLIC)
func GetLyricsByID(c *gin.Context) {
	id := c.Param("id")

	var lyrics models.Lyric
	if err := initializers.DB.First(&lyrics, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paroles introuvables"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

// 🟢 CRÉER UNE CHANSON (ADMIN SEULEMENT - Middleware appliqué)
func CreateLyrics(c *gin.Context) {
	var lyrics models.Lyric
	if err := c.ShouldBindJSON(&lyrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	if err := initializers.DB.Create(&lyrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la création des paroles"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Paroles ajoutées avec succès", "lyrics": lyrics})
}

// 🔵 METTRE À JOUR UNE CHANSON (ADMIN SEULEMENT - Middleware appliqué)
func UpdateLyrics(c *gin.Context) {
	id := c.Param("id")
	var lyrics models.Lyric

	// Vérifier si la chanson existe
	if err := initializers.DB.First(&lyrics, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paroles introuvables"})
		return
	}

	// Mettre à jour uniquement les champs modifiés
	var updatedLyrics models.Lyric
	if err := c.ShouldBindJSON(&updatedLyrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	initializers.DB.Model(&lyrics).Updates(updatedLyrics)

	c.JSON(http.StatusOK, gin.H{"message": "Paroles mises à jour avec succès", "lyrics": lyrics})
}

// 🔴 SUPPRIMER UNE CHANSON (ADMIN SEULEMENT - Middleware appliqué)
func DeleteLyrics(c *gin.Context) {
	id := c.Param("id")

	// Vérifier si la chanson existe avant suppression
	var lyrics models.Lyric
	if err := initializers.DB.First(&lyrics, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paroles introuvables"})
		return
	}

	// Supprimer la chanson
	if err := initializers.DB.Delete(&lyrics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paroles supprimées avec succès"})
}
