package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"proyecto-treap/internal/database"
	"proyecto-treap/internal/treap"
)

// Instancias de la BD y el Treap entre endpoints.
type Servidor struct {
	db *sql.DB
	t  *treap.Treap
}

func main() {
	db, err := database.InicializarDB("./inventario.db")
	if err != nil {
		log.Fatalf("Error al inicializar la BD: %v", err)
	}
	defer db.Close()

	t := treap.NuevoTreap()

	productos, err := database.CargarProductos(db)
	if err != nil {
		log.Fatalf("Error al cargar productos iniciales: %v", err)
	}
	for _, p := range productos {
		t.Insertar(p.Codigo, p)
	}
	log.Printf("Treap cargado exitosamente con %d productos.", len(productos))

	srv := &Servidor{db: db, t: t}

	// Rutas del servidor
	http.Handle("/", http.FileServer(http.Dir("./frontend"))) // index.html
	http.HandleFunc("/api/productos", srv.manejarProductos)   // API REST

	log.Println("Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Rutear las peticiones según el método HTTP.
func (srv *Servidor) manejarProductos(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	// Obtener un producto por código, o todos si no hay código (GET).
	case http.MethodGet:
		codigo := req.URL.Query().Get("codigo")

		// Si no hay código, devolvemos todo el catálogo.
		if codigo == "" {
			productos := srv.t.ListarTodos()
			// Si el Treap está vacío, devolvemos un arreglo vacío en vez de null
			if productos == nil {
				productos = make([]any, 0)
			}
			json.NewEncoder(res).Encode(productos)
			return
		}

		// Búsqueda rápida de un producto específico O(log n)[cite: 2].
		val, encontrado := srv.t.Buscar(codigo)
		if !encontrado {
			http.Error(res, `{"error": "Producto no encontrado"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(res).Encode(val)

	// Insertar o actualizar un producto (POST).
	case http.MethodPost:
		var p database.Producto
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := database.GuardarProducto(srv.db, p); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		srv.t.Insertar(p.Codigo, p)
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(p)

	// Eliminar un producto por código (DELETE).
	case http.MethodDelete:
		codigo := req.URL.Query().Get("codigo")
		if codigo == "" {
			http.Error(res, `{"error": "Falta el parámetro 'codigo'"}`, http.StatusBadRequest)
			return
		}

		if err := database.EliminarProducto(srv.db, codigo); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		srv.t.Eliminar(codigo) // O(log n) esperado en memoria.
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"mensaje": "Producto eliminado"}`))

	default:
		http.Error(res, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
