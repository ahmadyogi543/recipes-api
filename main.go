package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
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
	router.GET("/recipes/search", searchRecipesHandler)
	router.POST("/recipes", addRecipeHandler)
	router.PUT("/recipes/:id", updateHandler)
	router.DELETE("/recipes/:id", deleteRecipeHandler)

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

func searchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	filteredRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				filteredRecipes = append(filteredRecipes, recipes[i])
				break
			}
		}
	}

	c.JSON(http.StatusOK, filteredRecipes)
}

func updateHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}
	recipe.ID = id

	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

func deleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}
