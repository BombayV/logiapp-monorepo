package product

// GetProducts is the service function to retrieve all products.
func GetProducts() ([]Product, error) {
	// In a real app, this would call a repository to fetch data from the DB.
	return []Product{
		{ID: "prod_1", Name: "Super Widget", Description: "A very fine widget.", Price: 19.99},
		{ID: "prod_2", Name: "Mega Gadget", Description: "An impressive gadget.", Price: 29.99},
	}, nil
}
