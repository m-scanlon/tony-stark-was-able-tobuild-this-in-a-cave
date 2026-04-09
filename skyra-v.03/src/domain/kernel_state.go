package domain

import (
	"fmt"
	"strings"
)

type KernelState struct {
	beings map[string]*Being
}

func NewKernelState() *KernelState {
	return &KernelState{
		beings: make(map[string]*Being),
	}
}

func (s *KernelState) InsertBeing(being *Being) error {
	if being == nil {
		return ErrNilBeing
	}
	if err := being.Validate(); err != nil {
		return err
	}
	if _, exists := s.beings[being.Name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateBeingName, being.Name)
	}

	s.beings[being.Name] = being
	return nil
}

func (s *KernelState) BeingByName(name string) (*Being, bool) {
	being, ok := s.beings[strings.TrimSpace(name)]
	return being, ok
}

func (s *KernelState) SeedRelationship(leftName, rightName string) error {
	left, ok := s.beings[leftName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, leftName)
	}
	right, ok := s.beings[rightName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, rightName)
	}

	if err := left.SeedPeer(right.Name, right.Nature); err != nil {
		return err
	}
	if err := right.SeedPeer(left.Name, left.Nature); err != nil {
		return err
	}
	return nil
}
