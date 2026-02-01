package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var produk = []models.Product{
	{ID: 1, Name: "Kopi Gadjah", Price: 2000, Stock: 10},
	{ID: 2, Name: "Teh Tong Tji", Price: 1500, Stock: 5},
	{ID: 3, Name: "Indomie", Price: 2500, Stock: 20},
	{ID: 4, Name: "Super Bihun", Price: 1400, Stock: 15},
	{ID: 5, Name: "Panadol", Price: 4700, Stock: 2},
}

var category = []models.Category{
	{ID: 1, Name: "Minuman", Description: "Pelepas dahaga"},
	{ID: 2, Name: "Makanan", Description: "Anti kelaparan"},
	{ID: 3, Name: "Obat", Description: "Obat mujarab"},
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API version 1.0 running well! ",
	})
}

/*
func getAllCategory(w http.ResponseWriter, r *http.Request) {
	if len(category) == 0 {
		http.Error(w, "Category belum ada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, c := range category {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory models.Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// add new category
	newCategory.ID = len(category) + 1
	category = append(category, newCategory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCategory)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// convert to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updateCategory models.Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop over category, find based on id and update data
	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// convert to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	// loop over category
	for i, c := range category {
		if c.ID == id {
			category = append(category[:i], category[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "category deleted",
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)

}
*/

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Println("DBConn:" + config.DBConn)

	// setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initiazed database: ", err)
	}

	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "API Kasir - Belajar Golang CodeWithUmam",
		})
	})

	// health check
	http.HandleFunc("/health", health)

	/*
		http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				addCategory(w, r)
			case http.MethodGet:
				getAllCategory(w, r)

			default:
				http.Error(w, "Method not allowed", http.StatusBadRequest)
			}
		})

		http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				getCategoryByID(w, r)
			case http.MethodPut:
				updateCategory(w, r)
			case http.MethodDelete:
				deleteCategory(w, r)

			default:
				http.Error(w, "Method not allowed", http.StatusBadRequest)
			}
		})
	*/

	fmt.Println("Server running at http://localhost:" + config.Port)
	errMsg := http.ListenAndServe(":"+config.Port, nil)
	if errMsg != nil {
		fmt.Println("Error starting server:", errMsg)
	}
}
