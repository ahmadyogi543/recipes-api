package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"published_at"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)

	file, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(file, &recipes); err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := gin.Default()
	router.GET("/recipes", getAllRecipesHandler)
	router.POST("/recipes", addRecipeHandler)

	router.Run()
}

func getAllRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func addRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	c.JSON(http.StatusCreated, recipe)
}
