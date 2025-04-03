package main

import (
	"fmt"
)

// Main para probar la implementación
func main() {
	//Stack
	stack := Stack{} // Crear un stack vacío

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	val, err := stack.Peek()
	if err != nil {
		fmt.Println("Error al obtener el último elemento:", err)
	} else {
		fmt.Println("Último elemento (Peek):", val)
	} // Debería imprimir 30
	fmt.Println("Stack")
	stack.PrintS()
	val, _ = stack.Pop()
	fmt.Println("Elemento eliminado:", val) // Debería imprimir 30
	stack.PrintS()
	fmt.Println("Stack vacío?", stack.IsEmpty()) // Debería imprimir false

	//Queue
	myQueue := Queue{[]int{}}
	fmt.Println("\nQueue")
	fmt.Println(myQueue.IsEmpty())

	myQueue.Enqueue(1)
	myQueue.Enqueue(2)
	myQueue.Enqueue(3)

	myQueue.PrintQ()
	myQueue.Dequeue()
	myQueue.PrintQ()

	//Hash
	hashTable := NewHashTable()
	fmt.Println("\nHash Table")

	// Insertar valores
	hashTable.Put("uno", 1)
	hashTable.Put("dos", 2)
	hashTable.Put("tres", 3)

	// Buscar valores
	value, exists := hashTable.Get("dos")
	if exists {
		fmt.Println("Valor de 'dos':", value)
	} else {
		fmt.Println("'dos' no encontrado")
	}

	// Verificar si una clave existe
	fmt.Println("¿Contiene 'tres'?", hashTable.Contains("tres"))

	// Obtener todas las claves
	fmt.Println("Claves en la tabla:", hashTable.Keys())

	// Eliminar una clave
	hashTable.Remove("dos")
	fmt.Println("Después de eliminar 'dos':", hashTable.Keys())
}
