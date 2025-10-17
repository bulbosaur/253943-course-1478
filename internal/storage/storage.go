package storage

import (
	"errors"
	"fmt"
	"sync"

	"lyceum/internal/models"
)

type OrderStorage struct {
	Orders map[string]models.Order
	nextID int32
	Mu           sync.Mutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		Orders: make(map[string]models.Order),
	}
}

func (s *OrderStorage) CreateOrder(item string, quantity int32) string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	ID := fmt.Sprintf("%d", s.nextID)

	order := models.Order{
		ID: ID,
		Item: item,
		Quantity: quantity,
	}
	s.nextID ++
	s.Orders[ID] = order

	return ID
}

func (s *OrderStorage) GetOrder(ID string) (models.Order, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	order, exist := s.Orders[ID]

	if !exist {
		return models.Order{}, errors.New("there is no order with this ID")
	}
	return order, nil
}

func (s *OrderStorage) DeleteOrder(ID string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	delete(s.Orders, ID)
}

func (s *OrderStorage) UpdateOrder(ID, item string, quantity int32) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	newOrder := models.Order{
		ID: ID,
		Item: item,
		Quantity: quantity,
	}

	s.Orders[ID] = newOrder
}

func (s *OrderStorage) ListOrders() map[string]models.Order {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	return s.Orders
}
