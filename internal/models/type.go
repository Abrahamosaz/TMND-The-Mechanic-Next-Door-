package models


type PaginationQuery struct {
	Page *int
	Limit *int
}

type FilterQuery struct {
	Search *string
	Status *string
	*PaginationQuery
}


type PaginationResponse[T any] struct {
	Data       	[]T   `json:"data"`        // Holds the list of items
	Total      	*int64 `json:"total"`       // Total number of records
	Page       	*int   `json:"page"`        // Current page number
	Limit      	*int   `json:"limit"`       // Number of items per page
	TotalPages 	*int   `json:"totalPages"`  // Total pages available
}