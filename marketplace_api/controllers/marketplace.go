package controllers

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/georgysavva/scany/pgxscan"
	"context"
)

type Produk struct {
	IDProduk			int    `json:"id_produk"`
	NamaProduk		string    `json:"nama_produk"`
	Harga		int    `json:"harga"`
	Stok		int		`json:"stok"`
	Lokasi 	string 	  `json:"lokasi"`
}

type DetailProduk struct {
	IDProduk			int    `json:"id_produk"`
	NamaProduk		string    `json:"nama_produk"`
	Harga		int    `json:"harga"`
	Stok		int		`json:"stok"`
	NamaToko	string 	  `json:"nama_toko"`
	Lokasi 	string 	  `json:"lokasi"`
}

type ProdukKeranjang struct {
	IDKeranjang			int    `json:"id_keranjang"`
	IDProduk			int    `json:"id_produk"`
	NamaProduk		string    `json:"nama_produk"`
	Harga		int    `json:"harga"`
	Stok		int		`json:"stok"`
	Lokasi 	string 	  `json:"lokasi"`
}

type IdProduk struct {
	IDProduk			int    `json:"id_produk"`
}

// INITIALIZE DB CONNECTION (TO AVOID TOO MANY CONNECTION)
var dbConnect *pgxpool.Pool
func InitiateDB(dbPool *pgxpool.Pool) {
	dbConnect = dbPool
}

func GetAllProduk(c *gin.Context) {
	var arrProduk []Produk
	err := pgxscan.Select(context.Background(), dbConnect, &arrProduk, `SELECT id_produk, nama_produk, harga, stok, lokasi FROM toko NATURAL JOIN produk;`)

	if err != nil {
		log.Printf("Error while getting all produk, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Semua Produk",
		"data": arrProduk,
	})
	return
}

func GetProdukDetail(c *gin.Context) {
	produkId := c.Param("produkId")
	var detailProduk DetailProduk
	err := pgxscan.Get(context.Background(), dbConnect, &detailProduk, `SELECT id_produk, nama_produk, harga, stok, nama_toko, lokasi FROM toko NATURAL JOIN produk WHERE id_produk=$1`, produkId)

	if err != nil {
		log.Printf("Error while getting produk detail, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Produk not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Detail produk",
		"data": detailProduk,
	})
	return
}

func GetAllKeranjang(c *gin.Context) {
	var arrKeranjang []ProdukKeranjang
	err := pgxscan.Select(context.Background(), dbConnect, &arrKeranjang, `SELECT id_keranjang, id_produk, nama_produk, harga, stok, lokasi 
	FROM keranjang NATURAL JOIN produk NATURAL JOIN toko;`)

	if err != nil {
		log.Printf("Error while getting all produk di keranjang, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Semua Produk di Keranjang",
		"data": arrKeranjang,
	})
	return
}

func AddKeranjang(c *gin.Context) {
	var newKeranjang IdProduk
    if err := c.Bind(&newKeranjang); err != nil {
        return
    }
	produkId := newKeranjang.IDProduk
	_, err := dbConnect.Exec(context.Background(), "INSERT INTO keranjang(id_produk) VALUES ($1)", produkId)
	if err != nil {
		log.Printf("Error while inserting new keranjang, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Item keranjang created Successfully",
	})
	return
}

func DeleteKeranjang(c *gin.Context) {
	keranjangId := c.Param("keranjangId")
	_, err := dbConnect.Exec(context.Background(), "DELETE FROM keranjang WHERE id_keranjang=$1", keranjangId)
	if err != nil {
		log.Printf("Error while deleting keranjang, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Item keranjang deleted Successfully",
	})
	return
}

func EditStok(c *gin.Context) {
	produkId := c.Param("produkId")
	_, err := dbConnect.Exec(context.Background(), "UPDATE produk SET stok = stok - 1 WHERE id_produk = $1;", produkId)
	if err != nil {
		log.Printf("Error while mengurangi stok produk, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Stok produk dikurangi Successfully",
	})
	return
}