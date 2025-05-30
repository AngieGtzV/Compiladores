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

	ParamCount int
	VarCount   int
	TempCount  int
	StartQuad  int
}

var (
	FunctionDirectory = make(map[string]*FunctionInfo)
	CurrentFunction   = ""
	CurrentCall       = ""
	tempParams        []VarInfo
	ProgramName       = ""
	PrintArgs         []int
)

// ------------------- FUNCIONES PARA FUNCIONES -------------------

func AddFunction(name string, returnType string) error {
	//fmt.Println("DEBUG AddFunction - agregando función:", name)

	if _, exists := FunctionDirectory[name]; exists {
		return fmt.Errorf("funcion '%s' ya declarada", name)
	}

	FunctionDirectory[name] = &FunctionInfo{
		ReturnType: returnType,
		Params:     tempParams,
		Vars:       make(map[string]VarInfo),
		ParamCount: len(tempParams),
		VarCount:   0,
		TempCount:  0,
		StartQuad:  0,
	}
	tempParams = []VarInfo{}
	return nil
}

func SetCurrentFunction(name string) error {
	//fmt.Println("DEBUG SetCurrentFunction - currentFunction ahora es:", name)
	if _, exists := FunctionDirectory[name]; !exists {
		return fmt.Errorf("SetCurrentFunction: función '%s' no encontrada", name)
	}
	CurrentFunction = name
	return nil
}

// ------------------- FUNCIONES PARA VARIABLES -------------------

func AddVariable(name string, typeAttrib Attrib) error {
	varType, ok := typeAttrib.(string)
	if !ok {
		return errors.New("AddVariable: typeAttrib no es un string")
	}

	f, ok := FunctionDirectory[CurrentFunction]
	if !ok || f == nil {
		return fmt.Errorf("AddVariable: función '%s' no encontrada", CurrentFunction)
	}
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("variable '%s' ya declarada en %s", name, CurrentFunction)
	}

	scope := "local"
	if FunctionDirectory[CurrentFunction].ReturnType == "program" {
		scope = "global"
	}

	addr := Memory.Direccionar(scope, varType)

	f.Vars[name] = VarInfo{
		Type:    varType,
		Scope:   scope,
		Address: addr,
	}

	if scope == "local" {
		f.VarCount++
	}

	return nil
}

func AddParameter(name string, paramType string) error {
	f := FunctionDirectory[CurrentFunction]
	if _, exists := f.Vars[name]; exists {
		return fmt.Errorf("parametro '%s' ya declarado", name)
	}
	address := Memory.Direccionar("local", paramType)

	param := VarInfo{
		Type:    paramType,
		Scope:   "param",
		Address: address,
	}

	// Guardar en lista de parámetros y también en tabla de variables
	f.Params = append(f.Params, param)
	f.Vars[name] = param

	// También agregarlo a los parámetros temporales (por si aún no se ha llamado AddFunction)
	tempParams = append(tempParams, param)

	return nil
}

func LookupVariable(name string) (VarInfo, error) {
	if f, ok := FunctionDirectory[CurrentFunction]; ok {
		if v, ok := f.Vars[name]; ok {
			return v, nil
		}
	}
	if f, ok := FunctionDirectory[ProgramName]; ok {
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
	f, exists := FunctionDirectory[name]
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

func AppendArg(argType Attrib, rest Attrib) (Attrib, error) {
	typeStr, ok := argType.(string)
	if !ok {
		return nil, errors.New("AppendArg: argType no es string")
	}
	args := []string{typeStr}
	if r, ok := rest.([]string); ok {
		args = append(args, r...)
	}

	// Sacar argumento de PilaO
	arg := PopOperand()
	argTypeFromStack := PopType()

	if argTypeFromStack != typeStr {
		return nil, fmt.Errorf("AppendArg: tipo en stack '%s' no coincide con '%s'", argTypeFromStack, typeStr)
	}

	// Obtener función actual
	callFunc := CurrentCall
	paramIndex := FunctionDirectory[callFunc].ParamCount
	paramName := fmt.Sprintf("param%d", paramIndex+1) // empieza desde 1

	AddQuadruple(12, arg.Addr, -1, paramName)

	FunctionDirectory[callFunc].ParamCount++ // para el siguiente argumento

	return args, nil
}

// ------------ DEBUG: IMPRIMIR CUÁDRUPLOS -------------------

func PrintQuadruples() {
	fmt.Println("=== CUÁDRUPLOS GENERADOS ===")
	for i, q := range Quadruples {
		fmt.Printf("%03d: (%d, %d, %d, %v)\n", i, q.Op, q.Left, q.Right, q.Result)
	}
}
