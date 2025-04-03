package main

import (
	"fmt"
)

// Stack representa la estructura de la pila
type Stack struct {
	items []int
}

// Push agrega un elemento al tope del stack
func (s *Stack) Push(value int) {
	s.items = append(s.items, value) // Agrega el elemento al final
}

// Pop elimina y retorna el último elemento del stack
func (s *Stack) Pop() (int, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("stack vacío") // Devuelve un error si la pila está vacía
	}

	lastIndex := len(s.items) - 1
	value := s.items[lastIndex]   // Obtiene el último elemento
	s.items = s.items[:lastIndex] // Elimina el último elemento

	return value, nil
}

// Peek retorna el último elemento sin eliminarlo
func (s *Stack) Peek() (int, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("stack vacío")
	}
	return s.items[len(s.items)-1], nil
}

// IsEmpty verifica si el stack está vacío
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Imprime el stack
func (s *Stack) PrintS() {
	for _, item := range s.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
