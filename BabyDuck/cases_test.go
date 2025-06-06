package main

import (
	"BabyDuck/lexer"
	"BabyDuck/maquinavirtual"
	"BabyDuck/parser"
	"BabyDuck/semantics"
	"fmt"
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
				program test;
				var x: int;
				main {
					x = 10 / 2;
					print(x);
				}
				end
			`,
		valid: true,
	},
	{
		src: `
				program test;
				var x, y: float;
				main {
					y = 3.7;
					x = 1.5 + y;
				}
				end
			`,
		valid: true,
	},
	{
		src: `
				program test;
				var x: int;
				main {
					x = 6;
					while (x < 10) do
	    				{x = x + 1;}
						;
				}
				end
			`,
		valid: true,
	},
	{
		src: `
				program test5;
				var x: int;
				main {
					x = 6;
					print("Hola","mundo");
				}
				end
			`,
		valid: true,
	},
	{
		src: `
				program test;
				var x: int;
				main {
					x = 7;
					if (x < 10) {
						x = x + 1;
						print(x);
					} else {
						x = x - 2;
						print(x);
					};
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
				};
			}
			end
		`,
		valid: true,
	},
	{
		src: `program withCycle;
	         var i: int;
	         main {
	            i = 0;
	            while (i < 10/5) do {
	               print(i);
	               i = i + 1;
	            };
	         }
	         end`,

		valid: true,
	},
	{
		src: `
			program test8;
			void funcion() [ var x: int;
				{print("func");}
			];
			main {
				funcion();
			}
			end
		`,
		valid: true,
	},
	{
		src: `
						program test8;
						var y: int;
						void funcion() [ var x: int;
							{print("func", y);}
						];
						void second() [ var x: int;
							{funcion(); 
							x = 9;}
						];
						main {
							y = 10;
							
							second();
						}
						end
					`,
		valid: true,
	},
	{
		src: `program testFibonacci;
		var n, resultado: int;

		void fib(num: int)
		[
			var a, b, i, temp: int;
			{
				a = 0;
				b = 1;
				i = 0;
				while (i < num) do {
					temp = b;
					b = a + b;
					a = temp;
					i = i + 1;
				};
				resultado = a;
			}
		];

		main {
			n = 10;
			resultado = 0;
			fib(n);
			print("Resultado de Fibonacci", n, "es:", resultado);
		}
	end`,
		valid: true,
	},

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
	mm := semantics.NewDirecVirtuales()
	semantics.Memory = mm
	semantics.ConstTab = semantics.NewConstTable(mm)

	p := parser.NewParser()
	pass := true
	for _, ts := range testData {
		semantics.InitGlobals()
		s := lexer.NewLexer([]byte(ts.src))

		_, err := p.Parse(s)
		semantics.PrintQuadruples()

		fmt.Println("=== FunctionDirectory ===")
		for name := range semantics.FunctionDirectory {
			fmt.Println("FuncDir contiene:", name)
		}

		fmt.Println("=== MÁQUINA VIRTUAL ===")
		mv := maquinavirtual.NuevaMV(semantics.Quadruples, semantics.ConstTab.GetAddrValueMap(), semantics.FunctionDirectory)
		mv.Run()

		fmt.Println("|=============================|")

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
