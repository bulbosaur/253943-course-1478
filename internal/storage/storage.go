package storage

import (
	"errors"
	"fmt"
	"sync"

	pb "lyceum/pkg/api/test"
)

type OrderStorage struct {
	Orders map[string]*pb.Order
	nextID int32
	Mu     sync.Mutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		Orders: make(map[string]*pb.Order),
	}
}

func (s *OrderStorage) CreateOrder(item string, quantity int32) string {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	ID := fmt.Sprintf("%d", s.nextID)

	order := &pb.Order{
		Id:       ID,
		Item:     item,
		Quantity: quantity,
	}
	s.nextID++
	s.Orders[ID] = order

	return ID
}

func (s *OrderStorage) GetOrder(ID string) (*pb.Order, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	order, exist := s.Orders[ID]

	if !exist {
		return nil, errors.New("there is no order with this ID")
	}
	return order, nil
}

func (s *OrderStorage) DeleteOrder(ID string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	delete(s.Orders, ID)

	if _, exist := s.Orders[ID]; exist {
		return false
	}
	return true
}

func (s *OrderStorage) UpdateOrder(ID, item string, quantity int32) *pb.Order {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	newOrder := &pb.Order{
		Id:       ID,
		Item:     item,
		Quantity: quantity,
	}

	s.Orders[ID] = newOrder

	return s.Orders[ID]
}

func (s *OrderStorage) ListOrders() []*pb.Order {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if len(s.Orders) == 0 {
		return nil
	}

	slice := make([]*pb.Order, 0, len(s.Orders))
	for _, order := range s.Orders {
		slice = append(slice, order)
	}

	return slice
}
