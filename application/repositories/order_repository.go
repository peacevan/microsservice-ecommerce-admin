package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

type OrderRepository interface {
	Insert(order *domain.Order) (*domain.Order, error)
	Find(id string) (*domain.Order, error)
	Update(order *domain.Order) (*domain.Order, error)
	FindAll() ([]domain.Order, error)
	Delete(order *domain.Order) (*domain.Order, error)
}

type OrderRepositoryDb struct {
	Db *gorm.DB
}

func (repo OrderRepositoryDb) Insert(order *domain.Order) (*domain.Order, error) {

	err := repo.Db.Create(order).Error

	if err != nil {
		return nil, err
	}

	return order, nil

}

func (repo OrderRepositoryDb) Find(id string) (*domain.Order, error) {
	var order domain.Order
	err := repo.Db.First(&order, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order with ID %s not found", id)
		}
		return nil, err
	}
	return &order, nil
}

func (repo OrderRepositoryDb) Update(order *domain.Order) (*domain.Order, error) {
	err := repo.Db.Save(&order).Error

	if err != nil {
		return nil, err
	}
	return order, nil
}
func (repo OrderRepositoryDb) Delete(id string) error {
	return repo.Db.Where("id = ?", id).Delete(&domain.Order{}).Error
}
func (repo OrderRepositoryDb) FindAll() ([]domain.Order, error) {
	var orders []domain.Order
	err := repo.Db.Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
