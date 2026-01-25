package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Produk represents the data for a product
type Produk struct {
	ID    int    `json:"id" example:"1"`
	Nama  string `json:"nama" example:"Kopi Gadjah"`
	Harga int    `json:"harga" example:"2000"`
	Stok  int    `json:"stok" example:"10"`
}

// Category represents the data for a category
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Minuman"`
	Description string `json:"description" example:"Pelepas dahaga"`
}

var produk = []Produk{
	{ID: 1, Nama: "Kopi Gadjah", Harga: 2000, Stok: 10},
	{ID: 2, Nama: "Teh Tong Tji", Harga: 1500, Stok: 5},
	{ID: 3, Nama: "Indomie", Harga: 2500, Stok: 20},
	{ID: 4, Nama: "Super Bihun", Harga: 1400, Stok: 15},
	{ID: 5, Nama: "Panadol", Harga: 4700, Stok: 2},
}

var category = []Category{
	{ID: 1, Name: "Minuman", Description: "Pelepas dahaga"},
	{ID: 2, Name: "Makanan", Description: "Anti kelaparan"},
	{ID: 3, Name: "Obat", Description: "Obat mujarab"},
}

// health godoc
// @Summary Check API Health
// @Description Get the status of the API
// @Tags General
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API running well!",
	})
}

// getAllProduk godoc
// @Summary Get all products
// @Description Get list of all products
// @Tags Produk
// @Produce json
// @Success 200 {array} Produk
// @Failure 404 {string} string "Produk belum ada"
// @Router /api/produk [get]
func getAllProduk(w http.ResponseWriter, r *http.Request) {
	if len(produk) == 0 {
		http.Error(w, "Produk belum ada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produk)
}

// getProdukByID godoc
// @Summary Get product by ID
// @Description Get details of a single product by its ID
// @Tags Produk
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Produk
// @Failure 400 {string} string "Invalid produk ID"
// @Failure 404 {string} string "Produk belum ada"
// @Router /api/produk/{id} [get]
func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// addProduk godoc
// @Summary Add a new product
// @Description Create a new product with the provided data
// @Tags Produk
// @Accept json
// @Produce json
// @Param produk body Produk true "Product data"
// @Success 201 {object} Produk
// @Failure 400 {string} string "Invalid request"
// @Router /api/produk [post]
func addProduk(w http.ResponseWriter, r *http.Request) {
	var newProduk Produk
	err := json.NewDecoder(r.Body).Decode(&newProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// add new produk
	newProduk.ID = len(produk) + 1
	produk = append(produk, newProduk)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduk)
}

// updateProduk godoc
// @Summary Update an existing product
// @Description Update product details by ID
// @Tags Produk
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param produk body Produk true "Updated product data"
// @Success 200 {object} Produk
// @Failure 400 {string} string "Invalid request/ID"
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /api/produk/{id} [put]
func updateProduk(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// convert to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop over produk, find based on id and update data
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// deleteProduk godoc
// @Summary Delete a product
// @Description Remove a product from the list by its ID
// @Tags Produk
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "message: produk deleted"
// @Failure 400 {string} string "Invalid produk ID"
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /api/produk/{id} [delete]
func deleteProduk(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")

	// convert to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid produk ID", http.StatusBadRequest)
		return
	}

	// loop over produk
	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "produk deleted",
			})
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)

}

// getAllCategory godoc
// @Summary Get all categories
// @Description Get list of all categories
// @Tags Category
// @Produce json
// @Success 200 {array} Category
// @Failure 404 {string} string "Category belum ada"
// @Router /api/category [get]
func getAllCategory(w http.ResponseWriter, r *http.Request) {
	if len(category) == 0 {
		http.Error(w, "Category belum ada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// getCategoryByID godoc
// @Summary Get category by ID
// @Description Get details of a single category by its ID
// @Tags Category
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid Category ID"
// @Failure 404 {string} string "Produk belum ada"
// @Router /api/category/{id} [get]
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
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

// addCategory godoc
// @Summary Add a new category
// @Description Create a new category with the provided data
// @Tags Category
// @Accept json
// @Produce json
// @Param category body Category true "Category data"
// @Success 201 {object} Category
// @Failure 400 {string} string "Invalid request"
// @Router /api/category [post]
func addCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category
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

// updateCategory godoc
// @Summary Update an existing category
// @Description Update category details by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Updated category data"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid request/ID"
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /api/category/{id} [put]
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	// convert to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updateCategory Category
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

// deleteCategory godoc
// @Summary Delete a category
// @Description Remove a category from the list by its ID
// @Tags Category
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "message: category deleted"
// @Failure 400 {string} string "Invalid Category ID"
// @Failure 404 {string} string "Category tidak ditemukan"
// @Router /api/category/{id} [delete]
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id from req
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

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

// @title Kasir API
// @version 1.0
// @description Ini adalah API untuk sistem kasir sederhana.
// @host localhost:8080
// @BasePath /
func main() {
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/health", health)

	/*
		===== Produk Paths
	*/
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addProduk(w, r)
		case http.MethodGet:
			getAllProduk(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProdukByID(w, r)
		case http.MethodPut:
			updateProduk(w, r)
		case http.MethodDelete:
			deleteProduk(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	/*
		===== Category Paths
	*/
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addCategory(w, r)
		case http.MethodGet:
			getAllCategory(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Server running at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
