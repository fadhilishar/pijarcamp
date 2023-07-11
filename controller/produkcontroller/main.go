package produkcontroller

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/fadhilishar/pijarcamp/entity"
	"github.com/fadhilishar/pijarcamp/libraries"
	"github.com/fadhilishar/pijarcamp/model"
)

var validation = libraries.NewValidation()
var produkModel = model.NewProdukModel()

func Index(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("masuk index")
	queryString := req.URL.Query()
	// sort := req.Form.Get("sort")
	// req.GetBody()
	// fmt.Printf("req %#v\n", req)
	// fmt.Println("filterSort", sort)
	// fmt.Println("sort", sort)
	// fmt.Println("queryString: ", queryString)

	filterProduk := model.ProdukFilter{}
	// fmt.Println()
	if len(queryString) > 0 {
		// fmt.Println("queryString.Get(sort): ", queryString.Get("sort"))
		sort := queryString.Get("sort")
		if sort != "" {
			filterProduk.Sort = &sort
		}

		nama := queryString.Get("nama")
		if nama != "" {
			filterProduk.Nama = &nama
		}

		jumlahStr := queryString.Get("jumlah")
		if jumlahStr != "" {
			jumlah, _ := strconv.Atoi(jumlahStr)
			filterProduk.Jumlah = &jumlah
		}

		hargaMinStr := queryString.Get("hargaMin")
		if hargaMinStr != "" {
			hargaMin, _ := strconv.Atoi(hargaMinStr)
			filterProduk.HargaMin = &hargaMin
		}

		hargaMaxStr := queryString.Get("hargaMax")
		if hargaMaxStr != "" {
			hargaMax, _ := strconv.Atoi(hargaMaxStr)
			filterProduk.HargaMax = &hargaMax
		}
	}

	produks, err := produkModel.FindAll(filterProduk)
	if err != nil {
		return
	}

	data := map[string]interface{}{
		"produks": produks,
	}

	temp, err := template.ParseFiles("views/produk/index.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(res, data)
}

func Add(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/produk/tambah.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(res, nil)
	} else if req.Method == http.MethodPost {
		req.ParseForm()
		produk := entity.Produk{}
		produk.Nama = req.Form.Get("nama")
		// hargaStr:= req.Form.Get("harga")
		produk.Harga, _ = strconv.Atoi(req.Form.Get("harga"))
		produk.Keterangan = req.Form.Get("keterangan")
		produk.Jumlah, _ = strconv.Atoi(req.Form.Get("jumlah"))

		vErrors := validation.Struct(produk)

		var data = map[string]interface{}{}
		if vErrors != nil {
			data = map[string]interface{}{
				"validasi": vErrors,
				"produk":   produk,
			}
		} else {
			produkModel.Create(produk)
			data = map[string]interface{}{
				"pesan": "Data produk berhasil disimpan",
				"id":    produkModel.ID,
			}
		}

		temp, _ := template.ParseFiles("views/produk/tambah.html")
		temp.Execute(res, data)

	}

}

func Edit(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		queryString := req.URL.Query()
		id, _ := strconv.Atoi(queryString.Get("id"))

		produk, err := produkModel.FindOne(id)
		if err != nil {
			return
		}
		if produk == nil {
			return
		}

		data := map[string]interface{}{
			"produk": produk,
		}

		temp, err := template.ParseFiles("views/produk/ubah.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(res, data)
	} else if req.Method == http.MethodPost {
		req.ParseForm()
		produk := entity.Produk{}
		produk.ID, _ = strconv.Atoi(req.Form.Get("id"))
		produk.Nama = req.Form.Get("nama")
		produk.Harga, _ = strconv.Atoi(req.Form.Get("harga"))
		produk.Keterangan = req.Form.Get("keterangan")
		produk.Jumlah, _ = strconv.Atoi(req.Form.Get("jumlah"))

		vErrors := validation.Struct(produk)

		var data = map[string]interface{}{}
		if vErrors != nil {
			data = map[string]interface{}{
				"validasi": vErrors,
				"produk":   produk,
			}
		} else {
			produkModel.Update(produk)
			data = map[string]interface{}{
				"pesan": "Data produk berhasil diubah",
			}
		}

		temp, _ := template.ParseFiles("views/produk/ubah.html")
		temp.Execute(res, data)
	}
}

func Delete(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("ID HERE", id)
	queryString := req.URL.Query()
	id, _ := strconv.Atoi(queryString.Get("id"))
	produkModel.Delete(id)
	http.Redirect(res, req, "/produk", http.StatusSeeOther)
}
