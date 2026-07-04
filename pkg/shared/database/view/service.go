package view

import (
	"log"

	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/view/dto"
)

type Service[T any] interface {
	Add(view dto.View, f T) error
	Get(view dto.View) T
}

type service[T any] struct {
	views map[dto.View]T
}

func NewViewService[T any]() Service[T] {
	return &service[T]{
		views: make(map[dto.View]T),
	}
}

func (s *service[T]) Add(view dto.View, f T) error {
	// check if view already exists
	if _, exists := s.views[view]; exists {
		log.Fatalf("view %s already exists", view)
	}
	s.views[view] = f
	return nil
}

func (s *service[T]) Get(view dto.View) T {
	return s.views[view]
}
