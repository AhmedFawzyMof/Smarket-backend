package tables

import (
	"alwadi_markets/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type Favourite struct {
	Id      int
	User    string
	Product int
}

func (f Favourite) Add(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	FindProduct := db.QueryRow("SELECT * FROM Favourite WHERE (user, product) = (?, ?)", f.User, f.Product)

	err := FindProduct.Scan(&f.Id, &f.User, &f.Product)
	if err != nil {
		_, err := db.Exec("INSERT INTO `Favourite`(`user`, `product`) VALUES (?, ?)", f.User, f.Product)
		if err != nil {
			Response["Error"] = true

			middleware.SendResponse(response, Response)
			return
		}

		Response["Error"] = false

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false
	Response["Message"] = "المنتج موجود بالفعل في المفضلة"

	middleware.SendResponse(response, Response)
}

func (f Favourite) Get(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	var Products []int
	var FavProducts []Product
	getProducts, err := db.Query("SELECT * FROM Favourite WHERE user = ?", f.User)

	if err != nil {
		panic(err.Error())
	}

	defer getProducts.Close()

	for getProducts.Next() {
		var Fav Favourite

		if err = getProducts.Scan(&Fav.Id, &Fav.User, &Fav.Product); err != nil {
			panic(err.Error())
		}

		Products = append(Products, Fav.Product)
	}

	ids := ""

	for i, id := range Products {

		if len(Products) > 0 {
			if i >= 0 {
				ids += fmt.Sprintf("%d", id)
			}
			if i == len(Products)-1 {
				ids += fmt.Sprintf(",%d", id)
			}
		}

		var stmt string = fmt.Sprintf("SELECT * FROM Products WHERE id IN (%s)", ids)
		if ids != "" {
			products, err := db.Query(stmt)

			if err != nil {
				fmt.Println(stmt, err)
			}

			defer products.Close()

			for products.Next() {
				var Product Product

				if err := products.Scan(&Product.Id, &Product.Name, &Product.Description, &Product.Company, &Product.Subcategories, &Product.Category, &Product.Image, &Product.Available, &Product.Price, &Product.Offer); err != nil {
					panic(err.Error())
				}

				FavProducts = append(FavProducts, Product)
			}

		}

	}
	products, err := json.Marshal(FavProducts)
	if err != nil {
		fmt.Println(err.Error())
	}

	response <- products
	wg.Done()
}

func (f Favourite) Delete(db *sql.DB, response chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})
	_, err := db.Exec("DELETE FROM Favourite WHERE (user, product) = (?, ?)", f.User, f.Product)

	if err != nil {
		Response["Error"] = true

		middleware.SendResponse(response, Response)
		return
	}

	Response["Error"] = false

	middleware.SendResponse(response, Response)
}
