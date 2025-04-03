package main

// Hash
type HashTable struct {
	table map[string]int // Mapa para almacenar los datos
}

// NewHashTable crea una nueva tabla hash vacía
func NewHashTable() *HashTable {
	return &HashTable{
		table: make(map[string]int),
	}
}

// Put inserta un valor con una clave en la tabla hash
func (h *HashTable) Put(key string, value int) {
	h.table[key] = value
}

// Get busca un valor por su clave en la tabla hash
func (h *HashTable) Get(key string) (int, bool) {
	value, exists := h.table[key]
	return value, exists
}

// Remove elimina una clave de la tabla hash
func (h *HashTable) Remove(key string) {
	delete(h.table, key)
}

// Contains verifica si una clave existe en la tabla hash
func (h *HashTable) Contains(key string) bool {
	_, exists := h.table[key]
	return exists
}

// Keys devuelve todas las claves de la tabla hash
func (h *HashTable) Keys() []string {
	keys := make([]string, 0, len(h.table))
	for key := range h.table {
		keys = append(keys, key)
	}
	return keys
}
