package dto

type CreateCategoriesRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoriesRequest struct {
	Name string `json:"name"`
}

type CategoriesResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
