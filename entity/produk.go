package entity

type Produk struct {
	ID         int    `json:"id"`
	Nama       string `json:"nama" validate:"required" label:"Nama Produk"`
	Keterangan string `json:"keterangan"`
	Harga      int    `json:"harga" label:"Harga Produk"`
	Jumlah     int    `json:"jumlah"`
}
