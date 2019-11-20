package api

type Product struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	Price    float32 `json:"price"`
	Category string  `json:"category"`
	Image    string  `json:"image"`
	Rating   Rating  `json:"rating"`
}

type Rating struct {
	// out of 5 (e.g. 3.9)
	Stars float32 `json:"stars"`
	// number of customers who reviewed the product
	Customers int `json:"customers"`
}
