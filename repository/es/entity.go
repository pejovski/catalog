package es

type Document struct {
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	Price    float32 `json:"price"`
	Category string  `json:"category"`
	Image    string  `json:"image"`
}

type Update struct {
	Doc *Document `json:"doc"`
}

type Hit struct {
	Id     string   `json:"_id"`
	Source Document `json:"_source"`
}

type Result struct {
	Hits struct {
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}
