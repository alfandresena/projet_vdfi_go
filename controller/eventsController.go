package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"projet_vdfi/initializers"
	"projet_vdfi/models"
)

// 🟢 CRÉER UN ÉVÉNEMENT (ADMIN SEULEMENT)
func CreateEvent(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Vérifier si l'utilisateur est admin
	isAdmin := user.IsAdmin
	fmt.Println("isAdmin", isAdmin)

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé, administrateur requis"})
		return
	}

	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Ajouter l'événement en base de données
	if err := initializers.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la création de l'événement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Événement créé avec succès", "event": event})
}

// 🟡 LISTER TOUS LES ÉVÉNEMENTS
func GetEvents(c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de récupération des événements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// 🟠 OBTENIR UN ÉVÉNEMENT PAR ID
func GetEventByID(c *gin.Context) {
	id := c.Param("id")

	var event models.Event
	if err := initializers.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Événement introuvable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// 🔵 METTRE À JOUR UN ÉVÉNEMENT (ADMIN SEULEMENT)
func UpdateEvent(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Vérifier si l'utilisateur est admin
	isAdmin := user.IsAdmin
	fmt.Println("isAdmin", isAdmin)

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé, administrateur requis"})
		return
	}

	id := c.Param("id")
	var event models.Event

	// Vérifier si l'événement existe
	if err := initializers.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Événement introuvable"})
		return
	}

	// Mettre à jour l'événement
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Événement mis à jour avec succès", "event": event})
}

// 🔴 SUPPRIMER UN ÉVÉNEMENT (ADMIN SEULEMENT)
func DeleteEvent(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Vérifier si l'utilisateur est admin
	isAdmin := user.IsAdmin
	fmt.Println("isAdmin", isAdmin)

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé, administrateur requis"})
		return
	}

	id := c.Param("id")

	// Vérifier si l'événement existe
	if err := initializers.DB.Delete(&models.Event{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Échec de suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Événement supprimé avec succès"})
}
