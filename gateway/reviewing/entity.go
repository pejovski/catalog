package reviewing

type Rating struct {
	// out of 5 (e.g. 3.9)
	Stars float32 `json:"rating"`
	// number of customers who reviewed the product
	Customers int `json:"customers"`
}
