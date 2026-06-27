package treap

import (
	"math/rand"
	"time"
)

// Un Nodo representa un elemento individual dentro del Treap.
type Nodo struct {
	Clave     string
	Valor     any
	Prioridad float64
	Izq       *Nodo
	Der       *Nodo
}

// Un Treap maneja la referencia al nodo raíz y el generador de números aleatorios.
type Treap struct {
	Raiz *Nodo
	rnd  *rand.Rand
}

// NuevoTreap inicializa la estructura con una semilla de tiempo para las prioridades.
func NuevoTreap() *Treap {
	return &Treap{
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CrearNodo genera un nuevo nodo con una prioridad aleatoria continua entre [0.0, 1.0).
func (t *Treap) CrearNodo(clave string, valor any) *Nodo {
	return &Nodo{
		Clave:     clave,
		Valor:     valor,
		Prioridad: t.rnd.Float64(), // Prioridades continuas i.i.d. como indica el paper.
	}
}

// rotarIzquierda mueve el nodo 'x' hacia abajo y sube su hijo derecho 'y'.
func rotarIzquierda(x *Nodo) *Nodo {
	y := x.Der
	x.Der = y.Izq
	y.Izq = x
	return y // 'y' pasa a ser la nueva raíz de este subárbol.
}

// rotarDerecha mueve el nodo 'y' hacia abajo y sube su hijo izquierdo 'x'.
func rotarDerecha(y *Nodo) *Nodo {
	x := y.Izq
	y.Izq = x.Der
	x.Der = y
	return x // 'x' pasa a ser la nueva raíz de este subárbol.
}

// Insertar expone la función pública para agregar un par clave-valor.
func (t *Treap) Insertar(clave string, valor any) {
	t.Raiz = t.insertar(t.Raiz, clave, valor)
}

// insertar realiza la búsqueda binaria y rebalancea usando rotaciones de forma recursiva.
func (t *Treap) insertar(raiz *Nodo, clave string, valor any) *Nodo {
	if raiz == nil {
		return t.CrearNodo(clave, valor)
	}

	if clave < raiz.Clave {
		raiz.Izq = t.insertar(raiz.Izq, clave, valor)
		// Si el hijo izquierdo tiene mayor prioridad, rotamos a la derecha.
		if raiz.Izq.Prioridad > raiz.Prioridad {
			raiz = rotarDerecha(raiz)
		}
	} else if clave > raiz.Clave {
		raiz.Der = t.insertar(raiz.Der, clave, valor)
		// Si el hijo derecho tiene mayor prioridad, rotamos a la izquierda.
		if raiz.Der.Prioridad > raiz.Prioridad {
			raiz = rotarIzquierda(raiz)
		}
	} else {
		// Clave duplicada: actualizamos el valor en nuestro caso de uso.
		raiz.Valor = valor
	}

	return raiz
}

// Buscar expone la función pública para obtener un valor por su clave.
func (t *Treap) Buscar(clave string) (any, bool) {
	return t.buscar(t.Raiz, clave)
}

// buscar desciende por el árbol comparando claves en tiempo O(log n) esperado.
func (t *Treap) buscar(raiz *Nodo, clave string) (any, bool) {
	if raiz == nil {
		return nil, false
	}

	if clave < raiz.Clave {
		return t.buscar(raiz.Izq, clave)
	} else if clave > raiz.Clave {
		return t.buscar(raiz.Der, clave)
	}

	return raiz.Valor, true
}

// Eliminar expone la función pública para borrar un elemento por su clave.
func (t *Treap) Eliminar(clave string) {
	t.Raiz = t.eliminar(t.Raiz, clave)
}

// eliminar busca el nodo y lo rota hacia abajo hasta que es una hoja para borrarlo.
func (t *Treap) eliminar(raiz *Nodo, clave string) *Nodo {
	if raiz == nil {
		return nil
	}

	if clave < raiz.Clave {
		raiz.Izq = t.eliminar(raiz.Izq, clave)
	} else if clave > raiz.Clave {
		raiz.Der = t.eliminar(raiz.Der, clave)
	} else {
		// Encontramos el nodo a eliminar.
		if raiz.Izq == nil {
			return raiz.Der // Si no tiene hijo izquierdo, el derecho toma su lugar.
		} else if raiz.Der == nil {
			return raiz.Izq // Si no tiene hijo derecho, el izquierdo toma su lugar.
		}

		// Si tiene ambos hijos, rotamos basándonos en la prioridad más alta.
		if raiz.Izq.Prioridad > raiz.Der.Prioridad {
			raiz = rotarDerecha(raiz)
			raiz.Der = t.eliminar(raiz.Der, clave)
		} else {
			raiz = rotarIzquierda(raiz)
			raiz.Izq = t.eliminar(raiz.Izq, clave)
		}
	}

	return raiz
}
