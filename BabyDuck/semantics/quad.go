package semantics

//"errors"
import (
	"fmt"
	"strconv"
	"strings"
)

var (
	FunctionDirectory map[string]*FunctionInfo
	CurrentFunction   string
	OperandStack      []Operand
	OperatorStack     []int
	TypeStack         []string
	Quadruples        []Quadruple
	JumpStack         []int
)

type Quadruple struct {
	Op     int
	Left   int
	Right  int
	Result int
}

func InitGlobals() {
	OperandStack = []Operand{}
	OperatorStack = []int{}
	TypeStack = []string{}
	Quadruples = []Quadruple{}

	functionDirectory = make(map[string]*FunctionInfo)
	currentFunction = ""
	tempParams = []VarInfo{}
}

type Operand struct {
	Type  string
	Addr  int
	Value string // solo útil en constantes
}

const (
	// Globales
	GlobalIntInicio   = 1000
	GlobalFloatInicio = 2000

	// Locales
	LocalIntInicio   = 3000
	LocalFloatInicio = 4000

	// Temporales
	TempIntInicio   = 5000
	TempFloatInicio = 6000
	TempBoolInicio  = 7000

	// Constantes
	ConstIntInicio    = 8000
	ConstFloatInicio  = 9000
	ConstStringInicio = 10000
)

// ------------------- FUNCIONES STACK  -------------------
func PushOperand(op Operand) {
	OperandStack = append(OperandStack, op)
	TypeStack = append(TypeStack, op.Type)
}

func PopOperand() Operand {
	last := OperandStack[len(OperandStack)-1]
	OperandStack = OperandStack[:len(OperandStack)-1]
	return last
}

func PushType(typ string) {
	TypeStack = append(TypeStack, typ)
}

func PopType() string {
	last := TypeStack[len(TypeStack)-1]
	TypeStack = TypeStack[:len(TypeStack)-1]
	return last
}

func PopOperator() int {
	if len(OperatorStack) == 0 {
		panic("PopOperator: pila de operadores vacía")
	}
	op := OperatorStack[len(OperatorStack)-1]
	OperatorStack = OperatorStack[:len(OperatorStack)-1]
	return op
}

func PushJump(pos int) {
	JumpStack = append(JumpStack, pos)
}

func PopJump() (int, error) {
	if len(JumpStack) == 0 {
		return -1, fmt.Errorf("JumpStack vacío")
	}
	top := JumpStack[len(JumpStack)-1]
	JumpStack = JumpStack[:len(JumpStack)-1]
	return top, nil
}

// ------------------- DIRECCIONES VIRTUALES -------------------

var Memory = NewDirecVirtuales()

type DirecVirtuales struct {
	globalInt, globalFloat       int
	localInt, localFloat         int
	tempInt, tempFloat, tempBool int
	constInt, constFloat         int
	constString                  int
}

func NewDirecVirtuales() *DirecVirtuales {
	return &DirecVirtuales{
		globalInt:   GlobalIntInicio,
		globalFloat: GlobalFloatInicio,
		localInt:    LocalIntInicio,
		localFloat:  LocalFloatInicio,
		tempInt:     TempIntInicio,
		tempFloat:   TempFloatInicio,
		tempBool:    TempBoolInicio,
		constInt:    ConstIntInicio,
		constFloat:  ConstFloatInicio,
		constString: ConstStringInicio,
	}
}

func (mm *DirecVirtuales) Direccionar(segment, typ string) int {
	switch segment {
	case "global":
		if typ == "int" {
			addr := mm.globalInt
			mm.globalInt++
			return addr
		} else if typ == "float" {
			addr := mm.globalFloat
			mm.globalFloat++
			return addr
		}
	case "local":
		if typ == "int" {
			addr := mm.localInt
			mm.localInt++
			return addr
		} else if typ == "float" {
			addr := mm.localFloat
			mm.localFloat++
			return addr
		}
	case "temp":
		if typ == "int" {
			addr := mm.tempInt
			mm.tempInt++
			return addr
		} else if typ == "float" {
			addr := mm.tempFloat
			mm.tempFloat++
			return addr
		} else if typ == "bool" {
			addr := mm.tempBool
			mm.tempBool++
			return addr
		}
	case "const":
		if typ == "int" {
			addr := mm.constInt
			mm.constInt++
			return addr
		} else if typ == "float" {
			addr := mm.constFloat
			mm.constFloat++
			return addr
		} else if typ == "string" {
			addr := mm.constString
			mm.constString++
			return addr
		}
	}
	panic(fmt.Sprintf("Direccionar: segmento o tipo inválido (%s, %s)", segment, typ))
}

// ------------------- Tabla de Constantes -------------------

var ConstTab *ConstTable

type ConstTable struct {
	table map[string]int
	mm    *DirecVirtuales
}

func NewConstTable(mm *DirecVirtuales) *ConstTable {
	return &ConstTable{
		table: make(map[string]int),
		mm:    mm,
	}
}

func (ct *ConstTable) GetOrAddConstant(lit string, typ string) int {
	if addr, ok := ct.table[lit]; ok {
		return addr
	}
	addr := ct.mm.Direccionar("const", typ)
	ct.table[lit] = addr
	fmt.Println("=== Tabla de Constantes ===")
	for lit, addr := range ct.table {
		fmt.Printf("Constante: %q -> Dirección: %d\n", lit, addr)
	}
	return addr
}

// Devuelve un mapa que asocia direcciones virtuales con sus valores reales parseados
func (ct *ConstTable) GetAddrValueMap() map[int]interface{} {
	result := make(map[int]interface{})
	for lit, addr := range ct.table {
		// Determinar tipo del literal
		if strings.HasPrefix(lit, "\"") { // string
			result[addr] = strings.Trim(lit, "\"")
		} else if strings.Contains(lit, ".") { // float
			if f, err := strconv.ParseFloat(lit, 64); err == nil {
				result[addr] = f
			}
		} else { // int
			if i, err := strconv.Atoi(lit); err == nil {
				result[addr] = i
			}
		}
	}
	return result
}

// ------------------------ CUÁDRUPLOS -------------------------

func AddQuadruple(op int, left int, right int, result int) {
	Quadruples = append(Quadruples, Quadruple{
		Op:     op,
		Left:   left,
		Right:  right,
		Result: result,
	})
}
