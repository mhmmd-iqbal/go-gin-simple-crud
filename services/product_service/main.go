package product_service

import (
	"errors"
	"fmt"

	"github.com/mhmmd-iqbal/go-rest-gin/models"
	"github.com/mhmmd-iqbal/go-rest-gin/repositories/product_repository"
	"github.com/mhmmd-iqbal/go-rest-gin/requests/product_request"
	"gorm.io/gorm"
)

func GetProducts(search string, page int, limit int) ([]models.Product, int64, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	products, total, err := product_repository.GetProducts(search, page, limit)

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func GetDetailProduct(sku string) (*models.Product, error) {
	product, err := product_repository.GetProductBySKU(sku)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}

		return nil, err
	}

	return product, nil
}

func CreateProduct(input product_request.CreateProductRequest) (*models.Product, error) {
	_, err := product_repository.GetProductBySKU(input.SKU)

	if err == nil {
		return nil, fmt.Errorf("%w: %s", errors.New("product already exists"), input.SKU)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		SKU:         input.SKU,
	}

	createdProduct, err := product_repository.CreateProduct(&product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func UpdateProduct(sku string, input product_request.UpdateProductRequest) (*models.Product, error) {
	product, err := product_repository.GetProductBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}

		return nil, err
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Description != nil {
		product.Description = *input.Description
	}

	if input.Price != nil && *input.Price > 0 {
		product.Price = *input.Price
	}

	updatedProduct, err := product_repository.UpdateProductBySKU(sku, product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func DeleteProduct(sku string) error {
	_, err := product_repository.GetProductBySKU(sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}

		return err
	}

	return product_repository.DeleteProductBySKU(sku)
}
