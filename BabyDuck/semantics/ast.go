package semantics

import (
	"errors"
	"fmt"
)

// Attrib es el tipo genérico usado por el parser
// para pasar valores semánticos entre reglas.
type Attrib interface{}

type VarInfo struct {
	Type    string
	Scope   string // global, local, param
	Address int
}

type FunctionInfo struct {
	ReturnType string
	Params     []VarInfo
	Vars       map[string]VarInfo
}

var (
	functionDirectory = make(map[string]*FunctionInfo)
	currentFunction   = ""
	tempParams        []VarInfo
)

// ------------------- FUNCIONES PARA FUNCIONES -------------------

func AddFunction(name string, returnType string) error {
	//fmt.Println("DEBUG AddFunction - agregando función:", name)

	if _, exists := functionDirectory[name]; exists {
		return fmt.Errorf("funcion '%s' ya declarada", name)
	}
	functionDirectory[name] = &FunctionInfo{
		ReturnType: returnType,
		Params:     tempParams,
		Vars:       make(map[string]VarInfo),
	}
	tempParams = []VarInfo{}
	return nil
}

func SetCurrentFunction(name string) error {
	//fmt.Println("DEBUG SetCurrentFunction - currentFunction ahora es:", name)
	if _, exists := functionDirectory[name]; !exists {
		return fmt.Errorf("SetCurrentFunction: función '%s' no encontrada", name)
	}
	currentFunction = name
	return nil
}

// ------------------- FUNCIONES PARA VARIABLES -------------------

func AddVariable(name string, typeAttrib Attrib) error {
	varType, ok := typeAttrib.(string)
	if !ok {
		return errors.New("AddVariable: typeAttrib no es un string")
	}

	f, ok := functionDirectory[currentFunction]
	if !ok || f == nil {
		return fmt.Errorf("AddVariable: función '%s' no encontrada", currentFunction)
	}
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("variable '%s' ya declarada en %s", name, currentFunction)
	}

	scope := "local"
	if currentFunction == "program" || currentFunction == "main" {
		scope = "global"
	}

	addr := Memory.Direccionar(scope, varType)

	f.Vars[name] = VarInfo{
		Type:    varType,
		Scope:   scope,
		Address: addr,
	}

	return nil
}

func AddParameter(name string, paramType string) error {
	f := functionDirectory[currentFunction]
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("parametro '%s' ya declarado", name)
	}
	param := VarInfo{Type: paramType, Scope: "param"}
	f.Params = append(f.Params, param)
	f.Vars[name] = param
	tempParams = append(tempParams, param)
	return nil
}

func LookupVariable(name string) (VarInfo, error) {
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

// ------------------- EXPRESIONES -------------------

func ArithmeticResultType(leftType, rightType string) (string, error) {
	if leftType == "float" || rightType == "float" {
		return "float", nil
	}
	if leftType == "int" && rightType == "int" {
		return "int", nil
	}
	return "", fmt.Errorf("tipos incompatibles '%s' y '%s'", leftType, rightType)
}

// ------------------- COMPATIBILIDAD DE TIPOS -------------------

func TypeCompatible(expectedAttrib, actualAttrib Attrib) bool {
	expected, ok1 := expectedAttrib.(string)
	actual, ok2 := actualAttrib.(string)
	if !ok1 || !ok2 {
		return false
	}
	return expected == actual
}

// ------------------- LLAMADAS A FUNCIONES -------------------

func ValidateFunctionCall(name string, argsAttrib Attrib) error {
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

func CollectIDs(id string, moreIDs Attrib) ([]string, error) {
	ids := []string{id}
	if rest, ok := moreIDs.([]string); ok {
		ids = append(ids, rest...)
		return ids, nil
	}
	return nil, errors.New("CollectIDs: moreIDs no es []string")
}

func AppendID(id string, rest Attrib) (Attrib, error) {
	ids := []string{id}
	if r, ok := rest.([]string); ok {
		ids = append(ids, r...)
		return ids, nil
	}
	return nil, errors.New("AppendID: rest no es []string")
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

// Construir lista de parámetros desde WP
func BuildParamList(id string, typeAttrib Attrib, rest Attrib) (Attrib, error) {
	typ, ok := typeAttrib.(string)
	if !ok {
		return nil, errors.New("BuildParamList: typeAttrib no es string")
	}

	params := []VarInfo{{Type: typ, Scope: "param"}}

	if r, ok := rest.([]VarInfo); ok {
		params = append(params, r...)
	}

	return params, nil
}

// Agregar más parámetros a la lista
func AppendParamList(newParam Attrib, rest Attrib) (Attrib, error) {
	param, ok := newParam.(VarInfo)
	if !ok {
		return nil, errors.New("AppendParamList: param no es VarInfo")
	}
	params := []VarInfo{param}
	if r, ok := rest.([]VarInfo); ok {
		params = append(params, r...)
	}
	return params, nil
}

// ------------ DEBUG: IMPRIMIR DIRECTORIOS y CUÁDRUPLOS -------------------

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

func PrintQuadruples() {
	fmt.Println("=== CUÁDRUPLOS GENERADOS ===")
	for i, q := range Quadruples {
		fmt.Printf("%03d: (%d, %d, %d, %d)\n", i, q.Op, q.Left, q.Right, q.Result)
	}
}
