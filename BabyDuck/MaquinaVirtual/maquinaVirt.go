package maquinavirtual

import (
	"BabyDuck/semantics"
	"fmt"
)

type MaquinaVirt struct {
	Quads       []semantics.Quadruple
	Memory      map[int]interface{}                // Dirección virtual → Valor
	IP          int                                // Instruction Pointer
	ReturnStack []int                              // Guarda IPs para regresar de funciones
	EraTemp     map[int]interface{}                // Memoria temporal de ERA
	ParamQueue  []interface{}                      //Fila de parámetros
	FuncDir     map[string]*semantics.FunctionInfo //Nombres de las funciones en los cuádruplos
}

func aFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic(fmt.Sprintf("Tipo inválido: %T", val))
	}
}

func esFloat(val interface{}) bool {
	_, ok := val.(float64)
	return ok
}

func NuevaMV(quads []semantics.Quadruple, consts map[int]interface{}, funcDir map[string]*semantics.FunctionInfo) *MaquinaVirt {
	mv := &MaquinaVirt{
		Quads:       quads,
		Memory:      make(map[int]interface{}),
		IP:          0,
		ReturnStack: []int{},
		EraTemp:     make(map[int]interface{}),
		ParamQueue:  []interface{}{},
		FuncDir:     funcDir,
	}
	for addr, val := range consts {
		mv.Memory[addr] = val
	}
	return mv
}

func (mv *MaquinaVirt) Run() error {
	for mv.IP < len(mv.Quads) {
		quad := mv.Quads[mv.IP]

		switch quad.Op {
		case 0:
			val := mv.getValue(quad.Left)
			mv.setValue(quad.Result.(int), val)
		case 1: // +
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, esFloatL := left.(float64); esFloatL || esFloat(right) {
				result := aFloat(left) + aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) + right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 2: // -
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, esFloatL := left.(float64); esFloatL || esFloat(right) {
				result := aFloat(left) - aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) - right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 3: // *
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, esFloatL := left.(float64); esFloatL || esFloat(right) {
				result := aFloat(left) * aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) * right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 4: // /
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, esFloatL := left.(float64); esFloatL || esFloat(right) {
				result := aFloat(left) / aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) / right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 5: // <
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if esFloat(left) || esFloat(right) {
				result := aFloat(left) < aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) < right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 6: // >
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if esFloat(left) || esFloat(right) {
				result := aFloat(left) > aFloat(right)
				mv.setValue(quad.Result.(int), result)
			} else {
				result := left.(int) > right.(int)
				mv.setValue(quad.Result.(int), result)
			}
		case 7: // !=
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			mv.setValue(quad.Result.(int), left != right)
		case 8: // Goto
			mv.IP = quad.Result.(int)
			continue
		case 9: // GotoF
			cond := mv.getValue(quad.Left)
			if b, ok := cond.(bool); ok && !b {
				mv.IP = quad.Result.(int)
				continue
			} else if !ok {
				panic(fmt.Sprintf("GotoF esperaba bool en %v pero recibió %T", quad.Left, cond))
			}
		case 11: // Print
			arg := quad.Left
			if arg != -1 {
				val := mv.getValue(arg)
				fmt.Println(val)
			} else {
				fmt.Println("WARNING: Cuádruplo PRINT con argumento -1")
			}
		case 12: //param
			val := mv.getValue(quad.Left)
			mv.ParamQueue = append(mv.ParamQueue, val)
		case 13: // gosub
			mv.ReturnStack = append(mv.ReturnStack, mv.IP+1)

			baseAddr := quad.Left

			for i, val := range mv.ParamQueue {
				mv.Memory[baseAddr+i] = val
			}
			mv.ParamQueue = []interface{}{}

			funcName, ok := quad.Result.(string)
			if !ok {
				panic("GOSUB esperaba nombre de función string en Result")
			}
			funcInfo, exists := mv.FuncDir[funcName]
			if !exists {
				panic(fmt.Sprintf("Función '%s' no encontrada en FuncDir", funcName))
			}
			mv.IP = funcInfo.StartQuad
			continue
		case 14: // ERA
			mv.EraTemp = make(map[int]interface{})
			mv.ParamQueue = []interface{}{}
		case 15: // EndFunc
			if len(mv.ReturnStack) == 0 {
				panic("EndFunc: pila de retorno vacía")
			}
			mv.IP = mv.ReturnStack[len(mv.ReturnStack)-1]
			mv.ReturnStack = mv.ReturnStack[:len(mv.ReturnStack)-1]
			continue
		case 16: //End
			return nil
		}

		mv.IP++
	}
	return nil
}

func (mv *MaquinaVirt) getValue(addr int) interface{} {
	val, ok := mv.Memory[addr]
	if !ok {
		panic(fmt.Sprintf("Memoria no inicializada en %d", addr))
	}
	return val
}

func (mv *MaquinaVirt) setValue(addr int, val interface{}) {
	mv.Memory[addr] = val
}
