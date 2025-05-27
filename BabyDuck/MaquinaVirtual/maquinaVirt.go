package maquinavirtual

import (
	"BabyDuck/semantics"
	"fmt"
)

type MaquinaVirt struct {
	Quads  []semantics.Quadruple
	Memory map[int]interface{} // Dirección virtual → Valor
	IP     int                 // Instruction Pointer
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic(fmt.Sprintf("Tipo inválido: %T", val))
	}
}

func isFloat(val interface{}) bool {
	_, ok := val.(float64)
	return ok
}

func NuevaMV(quads []semantics.Quadruple, consts map[int]interface{}) *MaquinaVirt {
	mv := &MaquinaVirt{
		Quads:  quads,
		Memory: make(map[int]interface{}),
		IP:     0,
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
			mv.setValue(quad.Result, val)
		case 1: // +
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, isFloatL := left.(float64); isFloatL || isFloat(right) {
				result := toFloat(left) + toFloat(right)
				mv.setValue(quad.Result, result)
			} else {
				result := left.(int) + right.(int)
				mv.setValue(quad.Result, result)
			}
		case 2: // -
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, isFloatL := left.(float64); isFloatL || isFloat(right) {
				result := toFloat(left) - toFloat(right)
				mv.setValue(quad.Result, result)
			} else {
				result := left.(int) - right.(int)
				mv.setValue(quad.Result, result)
			}
		case 3: // *
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			if _, isFloatL := left.(float64); isFloatL || isFloat(right) {
				result := toFloat(left) * toFloat(right)
				mv.setValue(quad.Result, result)
			} else {
				result := left.(int) * right.(int)
				mv.setValue(quad.Result, result)
			}
		case 4: // /
			left := mv.getValue(quad.Left)
			right := mv.getValue(quad.Right)
			mv.setValue(quad.Result, left.(int)/right.(int))
		case 5: // <

		case 9: // GotoF
			condition := mv.getValue(quad.Left).(bool)
			if !condition {
				mv.IP = quad.Result
				continue
			}
		case 8: // Goto
			mv.IP = quad.Result
			continue
		case 11: // Print
			arg := quad.Left
			if arg != -1 {
				val := mv.getValue(arg)
				fmt.Println(val)
			} else {
				fmt.Println("WARNING: Cuádruplo PRINT con argumento -1")
			}
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
