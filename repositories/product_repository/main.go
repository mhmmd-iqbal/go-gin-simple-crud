package product_repository

import (
	"github.com/mhmmd-iqbal/go-rest-gin/models"
)

func GetProducts(search string, page int, limit int) ([]models.Product, int64, error) {
	var total int64
	var products []models.Product

	query := models.DB.Model(&models.Product{})
	offset := (page - 1) * limit

	if search != "" {
		query = query.Where(
			"name LIKE ? OR sku LIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&products).Error; err != nil {

		return nil, 0, err
	}

	return products, total, nil
}

func GetProductBySKU(sku string) (*models.Product, error) {
	var product models.Product

	if err := models.DB.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func CreateProduct(product *models.Product) (*models.Product, error) {
	if err := models.DB.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func UpdateProductBySKU(sku string, product *models.Product) (*models.Product, error) {
	if err := models.DB.Where("sku = ?", sku).Updates(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func DeleteProductBySKU(sku string) error {
	var product models.Product

	if err := models.DB.Where("sku = ?", sku).Delete(&product).Error; err != nil {
		return err
	}

	return nil
}
