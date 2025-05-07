package semantics

import (
	"errors"
	"fmt"
)

// Attrib es el tipo genérico usado por el parser
// para pasar valores semánticos entre reglas.
type Attrib interface{}

type VarInfo struct {
	Type  string
	Scope string // global, local, param
}

type FunctionInfo struct {
	ReturnType string
	Params     []VarInfo
	Vars       map[string]VarInfo
}

var (
	functionDirectory = make(map[string]*FunctionInfo)
	currentFunction   = ""
)

// ------------------- FUNCIONES PARA FUNCIONES -------------------

func AddFunction(nameAttrib Attrib, returnType string) error {
	name, ok := nameAttrib.(string)
	if !ok {
		return errors.New("AddFunction: nameAttrib no es un string")
	}
	if _, exists := functionDirectory[name]; exists {
		return fmt.Errorf("funcion '%s' ya declarada", name)
	}
	functionDirectory[name] = &FunctionInfo{
		ReturnType: returnType,
		Params:     []VarInfo{},
		Vars:       make(map[string]VarInfo),
	}
	return nil
}

func SetCurrentFunction(nameAttrib Attrib) error {
	name, ok := nameAttrib.(string)
	if !ok {
		return errors.New("SetCurrentFunction: nameAttrib no es un string")
	}
	currentFunction = name
	return nil
}

// ------------------- FUNCIONES PARA VARIABLES -------------------

func AddVariable(nameAttrib Attrib, typeAttrib Attrib) error {
	name, ok := nameAttrib.(string)
	if !ok {
		return errors.New("AddVariable: nameAttrib no es un string")
	}
	varType, ok := typeAttrib.(string)
	if !ok {
		return errors.New("AddVariable: typeAttrib no es un string")
	}
	f := functionDirectory[currentFunction]
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("variable '%s' ya declarada en %s", name, currentFunction)
	}
	f.Vars[name] = VarInfo{Type: varType, Scope: "local"}
	return nil
}

func AddParameter(nameAttrib Attrib, paramType string) error {
	name, ok := nameAttrib.(string)
	if !ok {
		return errors.New("AddParameter: nameAttrib no es un string")
	}
	f := functionDirectory[currentFunction]
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("parametro '%s' ya declarado", name)
	}
	param := VarInfo{Type: paramType, Scope: "param"}
	f.Params = append(f.Params, param)
	f.Vars[name] = param
	return nil
}

func LookupVariable(nameAttrib Attrib) (VarInfo, error) {
	name, ok := nameAttrib.(string)
	if !ok {
		return VarInfo{}, errors.New("LookupVariable: nameAttrib no es un string")
	}
	if f, ok := functionDirectory[currentFunction]; ok {
		if v, ok := f.Vars[name]; ok {
			return v, nil
		}
	}
	if f, ok := functionDirectory["program"]; ok {
		if v, ok := f.Vars[name]; ok {
			return v, nil
		}
	}
	return VarInfo{}, fmt.Errorf("variable '%s' no declarada", name)
}

// ------------------- COMPATIBILIDAD DE TIPOS -------------------

func TypeCompatible(expected, actual string) bool {
	return expected == actual
}

// ------------------- LLAMADAS A FUNCIONES -------------------

func ValidateFunctionCall(nameAttrib Attrib, argsAttrib Attrib) error {
	name, ok := nameAttrib.(string)
	if !ok {
		return errors.New("ValidateFunctionCall: nameAttrib no es un string")
	}
	args, ok := argsAttrib.([]string)
	if !ok {
		return errors.New("ValidateFunctionCall: argsAttrib no es []string")
	}
	f, exists := functionDirectory[name]
	if !exists {
		return fmt.Errorf("funcion '%s' no declarada", name)
	}
	if len(args) != len(f.Params) {
		return fmt.Errorf("funcion '%s' espera %d argumentos, se dieron %d", name, len(f.Params), len(args))
	}
	for i, param := range f.Params {
		if !TypeCompatible(param.Type, args[i]) {
			return fmt.Errorf("tipo de argumento %d incorrecto: esperado '%s', recibido '%s'", i+1, param.Type, args[i])
		}
	}
	return nil
}

// ------------------- LISTAS AUXILIARES -------------------

func CollectIDs(id Attrib, moreIDs Attrib) []string {
	ids := []string{}
	if str, ok := id.(string); ok {
		ids = append(ids, str)
	}
	if rest, ok := moreIDs.([]string); ok {
		ids = append(ids, rest...)
	}
	return ids
}

func AppendID(id Attrib, rest Attrib) (Attrib, error) {
	ids := []string{}
	if s, ok := id.(string); ok {
		ids = append(ids, s)
	}
	if r, ok := rest.([]string); ok {
		ids = append(ids, r...)
	}
	return ids, nil
}

func EmptyIDList() (Attrib, error) {
	return []string{}, nil
}

// ------------------- LISTAS DE ARGUMENTOS -------------------

func EmptyArgList() (Attrib, error) {
	return []string{}, nil
}

func BuildArgList(argType Attrib, rest Attrib) (Attrib, error) {
	typeStr, ok := argType.(string)
	if !ok {
		return nil, errors.New("BuildArgList: argType no es string")
	}
	args := []string{typeStr}
	if r, ok := rest.([]string); ok {
		args = append(args, r...)
	}
	return args, nil
}

func AppendArg(argType Attrib, rest Attrib) (Attrib, error) {
	typeStr, ok := argType.(string)
	if !ok {
		return nil, errors.New("AppendArg: argType no es string")
	}
	args := []string{typeStr}
	if r, ok := rest.([]string); ok {
		args = append(args, r...)
	}
	return args, nil
}

// ------------------- DEBUG: IMPRIMIR DIRECTORIOS -------------------

func PrintSymbolTables() {
	fmt.Println("========= Directorio de Funciones =========")
	for name, info := range functionDirectory {
		fmt.Printf("Función: %s\n", name)
		fmt.Printf("  Tipo de retorno: %s\n", info.ReturnType)

		if len(info.Params) > 0 {
			fmt.Println("  Parámetros:")
			for i, param := range info.Params {
				fmt.Printf("    [%d] Nombre: (en Vars) / Tipo: %s / Scope: %s\n", i+1, param.Type, param.Scope)
			}
		} else {
			fmt.Println("  Parámetros: ninguno")
		}

		if len(info.Vars) > 0 {
			fmt.Println("  Variables:")
			for varName, varInfo := range info.Vars {
				fmt.Printf("    %s: Tipo: %s / Scope: %s\n", varName, varInfo.Type, varInfo.Scope)
			}
		} else {
			fmt.Println("  Variables: ninguna")
		}

		fmt.Println("----------------------------------")
	}
}
