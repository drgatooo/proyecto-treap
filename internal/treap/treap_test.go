package treap

import (
	"testing"
)

// TestInsertarYBuscar verifica que los elementos se inserten y recuperen correctamente.
func TestInsertarYBuscar(t *testing.T) {
	treap := NuevoTreap()

	// Insertar datos de prueba
	treap.Insertar("7750001", "Gatorade")
	treap.Insertar("7750002", "Doritos")

	// Caso 1: Buscar elemento existente
	val, encontrado := treap.Buscar("7750001")
	if !encontrado {
		t.Error("Se esperaba encontrar el producto '7750001'")
	}
	if val.(string) != "Gatorade" {
		t.Errorf("Se esperaba 'Gatorade', pero se obtuvo '%v'", val)
	}

	// Caso 2: Buscar elemento inexistente
	_, encontrado = treap.Buscar("9999999")
	if encontrado {
		t.Error("No se esperaba encontrar el producto inexistente '9999999'")
	}
}

// TestEliminar verifica que el nodo se remueva manteniendo la integridad del árbol.
func TestEliminar(t *testing.T) {
	treap := NuevoTreap()

	treap.Insertar("7750003", "Agua San Luis")

	// Eliminar el elemento
	treap.Eliminar("7750003")

	// Verificar que ya no exista
	_, encontrado := treap.Buscar("7750003")
	if encontrado {
		t.Error("El producto '7750003' no debió ser encontrado tras su eliminación")
	}
}

// TestListarTodosInOrder comprueba que el catálogo devuelva los datos ordenados alfabéticamente.
func TestListarTodosInOrder(t *testing.T) {
	treap := NuevoTreap()

	// Insertar de forma desordenada
	treap.Insertar("C", "Producto C")
	treap.Insertar("A", "Producto A")
	treap.Insertar("B", "Producto B")

	productos := treap.ListarTodos()

	if len(productos) != 3 {
		t.Fatalf("Se esperaban 3 productos, se obtuvieron %d", len(productos))
	}

	// Al ser in-orden, deben salir estrictamente en secuencia A -> B -> C[cite: 2]
	if productos[0].(string) != "Producto A" ||
		productos[1].(string) != "Producto B" ||
		productos[2].(string) != "Producto C" {
		t.Error("El recorrido in-orden no devolvió los elementos en el orden alfabético correcto de sus claves")
	}
}
