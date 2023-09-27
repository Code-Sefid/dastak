package dto

type CreateProductsRequest struct {
	Title      string  `json:"title" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID *int    `json:"categoryId"`
	Inventory  int     `json:"inventory" binding:"required"`
}

type UpdateProductsRequest struct {
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	CategoryID int     `json:"categoryId"`
	Inventory  int     `json:"inventory"`
}

type ProductsResponse struct {
	ID        int                 `json:"id"`
	Title     string              `json:"title"`
	Price     float64             `json:"price"`
	Category  *CategoriesResponse `json:"category"`
	Inventory int                 `json:"inventory"`
}

type ProductFactorResponse struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}
