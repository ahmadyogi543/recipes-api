// Recipes API
//
// This is a sample recipes API wrriten in Go with Gin! You can find out more about the API at https://github.com/ahmadyogi543/recipes-api
//
// Schemes: http
// Host: localhost:5000
// BasePath: /
// Version: 1.0.0
// Contact: Ahmad Yogi <ahmadyogi543@gmail.com> https://github.com/ahmadyogi543
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"encoding/json"
	"flag"
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
	addr := flag.String("addr", ":5000", "HTTP network adress")
	flag.Parse()

	router := gin.Default()
	router.GET("/recipes", getAllRecipesHandler)
	router.GET("/recipes/:id", getRecipeHandler)
	router.GET("/recipes/search", searchRecipesHandler)
	router.POST("/recipes", addRecipeHandler)
	router.PUT("/recipes/:id", updateHandler)
	router.DELETE("/recipes/:id", deleteRecipeHandler)

	router.Run(*addr)
}

// swagger:operation GET /recipes recipes getAllRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	  '200':
//		  description: Successful operation
func getAllRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation GET /recipes/{id} recipes getAllRecipes
// Returns recipe based on its id
// ---
// produces:
// - application/json
// responses:
//
//		 '200':
//			  description: Successful operation
//	   '404':
//	     description: Invalid recipe ID
func getRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe *Recipe

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			recipe = &recipes[i]
		}
	}

	if recipe == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation POST /recipes recipes addRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//	  '201':
//		  description: Successful operation
//		'400':
//		  description: Invalid input
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

// swagger:operation GET /recipes/search recipes searchRecipes
// Search recipe based on tags
// ---
// parameters:
//   - name: tag
//     in: query
//     description: Tags of the recipe
//     required: false
//     type: string
//
// produces:
// - application/json
// responses:
//
//	  '200':
//		  description: Successful operation
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

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	 '200':
//		  description: Successful operation
//	 '400':
//			description: Invalid input
//	 '404':
//	   description: Invalid recipe ID
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

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete a recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	 '200':
//		  description: Successful operation
//	 '404':
//	   description: Invalid recipe ID
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
