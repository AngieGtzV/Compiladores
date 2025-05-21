package semantics

//"errors"
import "fmt"

//"strings"

var (
	FunctionDirectory map[string]*FunctionInfo
	CurrentFunction   string
	OperandStack      []Operand
	OperatorStack     []int
	TypeStack         []string
	Quadruples        []Quadruple
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
	GlobalIntStart   = 1000
	GlobalFloatStart = 2000

	// Locales
	LocalIntStart   = 3000
	LocalFloatStart = 4000

	// Temporales
	TempIntStart   = 5000
	TempFloatStart = 6000
	TempBoolStart  = 7000

	// Constantes
	ConstIntStart   = 8000
	ConstFloatStart = 9000
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

// ------------------- DIRECCIONES VIRTUALES -------------------

var Memory = NewMemoryManager()

type MemoryManager struct {
	globalInt, globalFloat       int
	localInt, localFloat         int
	tempInt, tempFloat, tempBool int
	constInt, constFloat         int
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		globalInt:   GlobalIntStart,
		globalFloat: GlobalFloatStart,
		localInt:    LocalIntStart,
		localFloat:  LocalFloatStart,
		tempInt:     TempIntStart,
		tempFloat:   TempFloatStart,
		tempBool:    TempBoolStart,
		constInt:    ConstIntStart,
		constFloat:  ConstFloatStart,
	}
}

func (mm *MemoryManager) Allocate(segment, typ string) int {
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
		}
	}
	panic(fmt.Sprintf("Allocate: segmento o tipo inválido (%s, %s)", segment, typ))
}

// ------------------- Tabla de Constantes -------------------

var ConstTab *ConstTable

type ConstTable struct {
	table map[string]int
	mm    *MemoryManager
}

func NewConstTable(mm *MemoryManager) *ConstTable {
	return &ConstTable{
		table: make(map[string]int),
		mm:    mm,
	}
}

func (ct *ConstTable) GetOrAddConstant(lit string, typ string) int {
	if addr, ok := ct.table[lit]; ok {
		return addr
	}
	addr := ct.mm.Allocate("const", typ)
	ct.table[lit] = addr
	return addr
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
