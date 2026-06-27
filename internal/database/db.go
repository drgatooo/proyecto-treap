package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Producto struct {
	Codigo string
	Nombre string
	Precio float64
	Stock  int
}

// InicializarDB crea la tabla de inventario si no existe.
func InicializarDB(ruta string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ruta)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS inventario (
		codigo TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		precio REAL NOT NULL,
		stock INTEGER NOT NULL
	);`

	_, err = db.Exec(query)
	return db, err
}

// GuardarProducto inserta o actualiza un producto en la base de datos (UPSERT).
func GuardarProducto(db *sql.DB, p Producto) error {
	query := `
	INSERT INTO inventario (codigo, nombre, precio, stock) 
	VALUES (?, ?, ?, ?)
	ON CONFLICT(codigo) DO UPDATE SET
		nombre = excluded.nombre,
		precio = excluded.precio,
		stock = excluded.stock;`

	_, err := db.Exec(query, p.Codigo, p.Nombre, p.Precio, p.Stock)
	return err
}

// EliminarProducto borra un registro por su código.
func EliminarProducto(db *sql.DB, codigo string) error {
	_, err := db.Exec("DELETE FROM inventario WHERE codigo = ?", codigo)
	return err
}

// CargarProductos obtiene todos los productos para inicializar el Treap.
func CargarProductos(db *sql.DB) ([]Producto, error) {
	rows, err := db.Query("SELECT codigo, nombre, precio, stock FROM inventario")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.Codigo, &p.Nombre, &p.Precio, &p.Stock); err != nil {
			return nil, err
		}
		productos = append(productos, p)
	}
	return productos, nil
}
