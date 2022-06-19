package controllers

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/abdillahzakkie/silkroad/helpers"
// 	"github.com/abdillahzakkie/silkroad/models"
// )

// // POST "/products/{seller_id}/new"
// // CreateNewProduct creates new product
// func CreateNewProduct(w http.ResponseWriter, r *http.Request) {
// 	var productService models.ProductService
// 	var product models.Product

// 	// parse url-encoded form data into struct
// 	err := helpers.ParseForm(r, &product); if err != nil {
// 		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	user := models.User{
// 		ID: product.UserID,
// 	}
// 	// lookup seller by Id
// 	var userService models.UserService
// 	err = userService.IsExistingUser(user); if err != nil {
// 		switch err {
// 			case models.ErrNotFound:
// 				helpers.RespondWithError(w, http.StatusNotFound, "seller not found")
// 				return
// 			default:
// 				helpers.RespondWithError(w, http.StatusInternalServerError, models.ErrInternalServerError.Error())
// 				return
// 		}
// 	}

// 	// checks if category ID exists
// 	var categoryService models.CategoryService
// 	_, err = categoryService.GetCategoryById(product.CategoryID); if err != nil {
// 		helpers.RespondWithError(w, http.StatusNotFound, "category does not exist")
// 		return
// 	}

// 	// checks if product has already existed
// 	// discard error if product does not exist
// 	err = productService.IsExistingProduct(product)
// 	if err == nil {
// 		helpers.RespondWithError(w, http.StatusBadRequest, models.ErrorProductAlreadyExists.Error())
// 		return
// 	}

// 	// insert new product record
// 	product, err = productService.CreateNewProduct(product); if err != nil {
// 		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(product)
// }

// // GET "/products"
// // GetAllProducts get all products
// func GetAllProducts(w http.ResponseWriter, r *http.Request) {
// 	var productService models.ProductService

// 	products, err := productService.GetAllProducts(); if err != nil {
// 		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(products)
// }

// // GET "/products/{product_id}"
// // GetProductById get product by Product ID
// func GetProductById(w http.ResponseWriter, r *http.Request) {
// 	var product models.Product
// 	var err error
// 	vars := mux.Vars(r)

// 	product.ID, err = strconv.Atoi(vars["product_id"]); if err != nil {
// 		helpers.RespondWithError(w, http.StatusBadRequest, "invalid product ID")
// 		return
// 	}

// 	err = product.GetProductById(); if err != nil {
// 		helpers.RespondWithError(w, http.StatusNotFound, "product does not exist")
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(product)
// }

// // GET "/products/{seller_id}"
// // GetProductsBySellerId get product by Seller ID
// func GetProductsBySellerId(w http.ResponseWriter, r *http.Request) {
// 	var product models.Product
// 	var user models.User
// 	var err error
// 	vars := mux.Vars(r)

// 	product.UserID, err = strconv.Atoi(vars["seller_id"]); if err != nil {
// 		helpers.RespondWithError(w, http.StatusBadRequest, "invalid seller ID")
// 		return
// 	}

// 	// checks if seller exists
// 	err = user.GetUser(); if err != nil {
// 		helpers.RespondWithError(w, http.StatusNotFound, "seller does not exist")
// 		return
// 	}

// 	// get all products by seller ID
// 	products, err := product.GetProductsBySellerId(); if err != nil {
// 		helpers.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("%v", err))
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(products)
// }