package main

import (
	"BabyDuck/lexer"
	"BabyDuck/parser"
	"testing"
)

type TI struct {
	src   string
	valid bool // true si debe ser aceptado, false si no
}

var testData = []*TI{
	//Casos válidos simples
	{
		src: `
			program test;
			main {
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test2;
			var x: int;
			main {
				x = 10;
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test3;
			var x, y: float;
			main {
				x = 1.5 + 2.0;
			}
			end
		`,
		valid: true,
	},
	/*{
		src: `
			program test4;
			main {
				print("Hola!");
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test5;
			main {
				print("Hola");
				print("Mundo");
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test6;
			main {
				if (1 < 2) {
					print("Si");
				} else {
					print("No");
				}
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test7;
			main {
				while (1 < 2) do {
					print("loop");
				};
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test8;
			void foo() [ var x: int;
				{print("func");}
			];
			main {
				foo();
			}
			end
		`,
		valid: true,
	},
	{
		src: `
			program test8;
			void foo() [ var x: int;
				{print("func");}
			];
			void second() [ var x: int;
				{print("func");}
			];
			main {
				foo();
				second();
			}
			end
		`,
		valid: true,
	},*/

	/*//Casos inválidos
	{
		src: `
			program test
			main {
			}
			end
		`,
		valid: false, // falta ';' después de "program test"
	},
	{
		src: `
			program test;
			var x int;
			main {
			}
			end
		`,
		valid: false, // falta ':' entre x e int
	},
	{
		src: `
			program test;
			var x: int
			main {
			}
			end
		`,
		valid: false, // falta ';' después de declaración de variable
	},
	{
		src: `
			program test;
			main {
				x = ;
			}
			end
		`,
		valid: false, // falta expresión en asignación
	},
	{
		src: `
			program test;
			main {
				if 1 < 2 {
					print("si");
				}
			}
			end
		`,
		valid: false, // falta paréntesis en `if (1 < 2)`
	},
	{
		src: `
			program test;
			main {
				print(Hola);
			}
			end
		`,
		valid: false, // falta comillas en el string
	},
	{
		src: `
			program test;
			main {
				while 1 < 2 do {
					print("loop");
				};
			}
			end
		`,
		valid: false, // falta paréntesis en while
	},
	{
		src: `
			program test;
			void func ( ) [ var x int; ] {
				print("error");
			};
			main {
				func();
			}
			end
		`,
		valid: false, // falta ':' entre x e int en parámetros
	},
	{
		src: `
			program test;
			var x, y: int
			main {
				print("hola");
			}
			end
		`,
		valid: false, // falta ';' al final de declaración de variables
	},
	{
		src: `
			program test;
			main {
				1 = x;
			}
			end
		`,
		valid: false, // asignación inválida: no se puede asignar a un literal
	},
	{
		src: `
			1 = x;
		`,
		valid: false, // no empieza por program ni el resto de la secuencia
	},*/
}

func Test1(t *testing.T) {
	p := parser.NewParser()
	pass := true
	for _, ts := range testData {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)

		if (err == nil) != ts.valid {
			pass = false
			if err != nil {
				t.Logf("Expected valid: %v, but got error: %s", ts.valid, err.Error())
			} else {
				t.Logf("Expected invalid, but parsed successfully:\n%s", ts.src)
			}
		}
	}

	if !pass {
		t.Fail()
	}
}
