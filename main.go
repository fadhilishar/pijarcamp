package main

import (
	"net/http"

	"github.com/fadhilishar/pijarcamp/controller/produkcontroller"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// http.HandleFunc("/", produkcontroller.Index)
	http.HandleFunc("/produk", produkcontroller.Index)
	// http.HandleFunc("/produk/index", produkcontroller.Index)
	http.HandleFunc("/produk/tambah", produkcontroller.Add)
	http.HandleFunc("/produk/ubah", produkcontroller.Edit)
	http.HandleFunc("/produk/hapus", produkcontroller.Delete)

	http.ListenAndServe(":3000", nil)

}
