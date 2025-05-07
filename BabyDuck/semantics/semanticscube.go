package semantics

type SemanticCube map[string]map[string]map[string]string

func NewSemanticCube() SemanticCube {
	cube := make(SemanticCube)

	ops := []string{"+", "-", "*", "/", "<", ">", "!="}
	types := []string{"int", "float", "string"}

	// Inicializar
	for _, op := range ops {
		cube[op] = make(map[string]map[string]string)
		for _, left := range types {
			cube[op][left] = make(map[string]string)
			for _, right := range types {
				cube[op][left][right] = "error" // default
			}
		}
	}

	// +
	cube["+"]["int"]["int"] = "int"
	cube["+"]["int"]["float"] = "float"
	cube["+"]["float"]["int"] = "float"
	cube["+"]["float"]["float"] = "float"
	cube["+"]["string"]["string"] = "error"

	// -
	cube["-"]["int"]["int"] = "int"
	cube["-"]["int"]["float"] = "float"
	cube["-"]["float"]["int"] = "float"
	cube["-"]["float"]["float"] = "float"

	// *
	cube["*"]["int"]["int"] = "int"
	cube["*"]["int"]["float"] = "float"
	cube["*"]["float"]["int"] = "float"
	cube["*"]["float"]["float"] = "float"

	// /
	cube["/"]["int"]["int"] = "float"
	cube["/"]["int"]["float"] = "float"
	cube["/"]["float"]["int"] = "float"
	cube["/"]["float"]["float"] = "float"

	// <
	cube["<"]["int"]["int"] = "int"
	cube["<"]["int"]["float"] = "int"
	cube["<"]["float"]["int"] = "int"
	cube["<"]["float"]["float"] = "int"

	// >
	cube[">"]["int"]["int"] = "int"
	cube[">"]["int"]["float"] = "int"
	cube[">"]["float"]["int"] = "int"
	cube[">"]["float"]["float"] = "int"

	// !=
	cube["!="]["int"]["int"] = "int"
	cube["!="]["int"]["float"] = "int"
	cube["!="]["float"]["int"] = "int"
	cube["!="]["float"]["float"] = "int"

	return cube
}

// LookupTypeResult devuelve el tipo resultante o "error"
func (c SemanticCube) LookupTypeResult(op, leftType, rightType string) string {
	if opMap, ok := c[op]; ok {
		if leftMap, ok := opMap[leftType]; ok {
			if result, ok := leftMap[rightType]; ok {
				return result
			}
		}
	}
	return "error"
}
