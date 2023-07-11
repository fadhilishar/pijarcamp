package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fadhilishar/pijarcamp/config"
	"github.com/fadhilishar/pijarcamp/entity"
)

var countReq = 0

type ProdukFilter struct {
	ID         *int
	Nama       *string
	Keterangan *string
	HargaMin   *int
	HargaMax   *int
	Jumlah     *int
	Sort       *string
}

type ProdukModel struct {
	conn *sql.DB
	ID   int
}

func NewProdukModel() *ProdukModel {
	conn, err := config.DBCOnnection()
	if err != nil {
		panic(err)
	}

	return &ProdukModel{
		conn: conn,
	}
}

func (p *ProdukModel) Create(produk entity.Produk) bool {
	result, err := p.conn.Exec("insert into produk (nama, keterangan, harga, jumlah) values(?,?,?,?)", produk.Nama, produk.Keterangan, produk.Harga, produk.Jumlah)
	if err != nil {
		fmt.Println("error saat menambah data produk", err)
		return false
	}

	lastInsertId, _ := result.LastInsertId()
	p.ID = int(lastInsertId)
	return lastInsertId > 0
}

func (p *ProdukModel) FindAll(filter ProdukFilter) (res []entity.Produk, err error) {
	query := `select * from produk `

	// fmt.Println("SAMPAI ZINI")
	// args := []any{}
	mapFieldToValue := map[string]interface{}{}
	// if filter.ID != nil || filter.Nama != nil || filter.Keterangan != nil || filter.HargaMin != nil || filter.Jumlah != nil || filter.Sort != nil {
	// query += "where "
	if filter.ID != nil {
		mapFieldToValue["id"] = *filter.ID
		// query += "id=?"
		// args = append(args, filter.ID)
	}
	if filter.Nama != nil {
		// query += "nama=?"
		mapFieldToValue["nama"] = *filter.Nama
		// args = append(args, filter.Nama)
	}
	if filter.HargaMin != nil {
		// query += "harga>?"
		mapFieldToValue["hargaMin"] = *filter.HargaMin
		// args = append(args, filter.HargaMin)
	}
	if filter.HargaMax != nil {
		// query += "harga<?"
		mapFieldToValue["hargaMax"] = *filter.HargaMax
		// args = append(args, filter.HargaMax)
	}

	if filter.Jumlah != nil {
		mapFieldToValue["jumlah"] = *filter.Jumlah
	}

	// }
	// fmt.Println("DI ZNINI")

	// args := make([]any, len(mapFieldToValue))
	args := []any{}
	i := 0
	if len(mapFieldToValue) > 0 {
		// fmt.Println("MAPFIELDTOVALUE HERERRRR", mapFieldToValue)
		query += "where "
		for field, value := range mapFieldToValue {
			if i != 0 {
				query += " and "
			}

			switch field {
			case "hargaMin":
				query += "harga>=?"
				args = append(args, value)
			case "hargaMax":
				query += "harga<=?"
				args = append(args, value)
			case "nama":
				// query += "nama regexp '?'"
				// fmt.Println("VALUE", value)
				query += fmt.Sprintf("nama regexp '%v'", value)
			default:
				query += fmt.Sprintf("%v=?", field)
				args = append(args, value)
			}
			// if field == "hargaMin" {
			// 	query += "harga>=?"
			// } else if field == "hargaMax" {
			// 	query += "harga<=?"
			// } else {
			// 	query += fmt.Sprintf("%v=?", field)
			// }

			// args[i] = value
			// args = append(args, value)
			i++
		}
	}

	// fmt.Println("KALAU SINI?")
	// fmt.Println("FILTER", filter)
	if filter.Sort != nil {
		fields := strings.Split(*filter.Sort, "_")
		// query += "order by ? ?"
		query += "order by " + fields[0] + " " + fields[1]

		// mapFieldToValue["order by"] = fields[0] + " " + fields[1]
		// args = append(args, fields[0], fields[1])
	}
	countReq++
	// fmt.Println()
	// fmt.Println("REQUEST KE ", countReq)
	// fmt.Println()
	// fmt.Println("FILTER", filter)
	// fmt.Println("MAPFIELDTOVALUE", mapFieldToValue)
	// fmt.Println("QUERY ", query)
	// fmt.Println("ARGS", args)
	var rows *sql.Rows
	rows, err = p.conn.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	defer rows.Close()

	produk := entity.Produk{}
	for rows.Next() {
		rows.Scan(
			&produk.ID,
			&produk.Nama,
			&produk.Keterangan,
			&produk.Harga,
			&produk.Jumlah,
		)
		res = append(res, produk)
	}

	return
}

func (p *ProdukModel) FindOne(id int) (res *entity.Produk, err error) {
	res = &entity.Produk{}
	err = p.conn.QueryRow("select * from produk where id=?", id).Scan(
		&res.ID,
		&res.Nama,
		&res.Keterangan,
		&res.Harga,
		&res.Jumlah,
	)
	return
}

func (p *ProdukModel) Update(produk entity.Produk) bool {
	// fmt/.Println("produk di sini", produk.Nama, produk)
	result, err := p.conn.Exec(
		"update produk set nama=?, keterangan=?, harga=?, jumlah=? where id=?",
		produk.Nama, produk.Keterangan, produk.Harga, produk.Jumlah, produk.ID)
	if err != nil {
		// fmt.Println("error saat mengubah data produk", err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected >= 0
}

func (p *ProdukModel) Delete(id int) {
	p.conn.Exec("delete from produk where id=?", id)
}
