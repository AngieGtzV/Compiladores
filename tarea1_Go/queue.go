package main

import (
	"fmt"
)

// Queue
type Queue struct {
	items []int
}

// Enqueue -> insertar o push
func (q *Queue) Enqueue(data int) {
	q.items = append(q.items, data)
}

// Dequeue -> eliminar o pop
func (q *Queue) Dequeue() {
	q.items = q.items[1:]
}

// Front regresa un error si la Queue está vacía. Si no, regresa el primer elemento [0].
func (q *Queue) Front() (int, error) {
	if len(q.items) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	return q.items[0], nil
}

// IsEmpty verifica si la Queue está vacío
func (q *Queue) IsEmpty() bool {
	if len(q.items) == 0 {
		return true
	} else {
		return false
	}
}

// Imprime la Queue
func (q *Queue) PrintQ() {
	for _, item := range q.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
