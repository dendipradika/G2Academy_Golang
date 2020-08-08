package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ParkirMasuk struct {
	IdParkir  string 	  `json:"id_parkir"`
	WaktuMasuk time.Time  `json:"waktu_masuk"`
}

type ParkirKeluar struct {
	IdParkir      string 	 `json:"id_parkir"`
	WaktuMasuk    time.Time  `json:"waktu_masuk"`
	WaktuKeluar   time.Time  `json:"waktu_keluar"`
	DurasiParkir  int        `json:"durasi_parkir"`
	TipeKendaraan string     `json:"tipe_kendaraan"`
	PlatNo        string     `json:"plat_nomor"`
	BiayaParkir   int        `json:"biaya_parkir"`
}

var dataParkirMasuk []ParkirMasuk
var dataParkirKeluar []ParkirKeluar

func main() {
	http.HandleFunc("/server/masuk", customerMasuk)
	http.HandleFunc("/server/keluar", customerKeluar)
	http.HandleFunc("/server/semuaParkirMasuk", parkirMasuk)
	http.HandleFunc("/server/semuaParkirKeluar", parkirKeluar)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func generateID() string{
	var acakID string

	rand.Seed(time.Now().UnixNano())
	for i:=0;i<10;i++{
		acakID += fmt.Sprint(rand.Intn(7))
	}
	return acakID
}

func simpanDataParkirMasuk(idParkir string, waktuMasuk time.Time) {
	parkirMasuk := ParkirMasuk{IdParkir: idParkir, WaktuMasuk: waktuMasuk}
	dataParkirMasuk = append(dataParkirMasuk, parkirMasuk)
}

func simpanDataParkirKeluar(idParkir string, waktuMasuk time.Time, waktuKeluar time.Time, tipeKendaraan string,
	platNo string, durasiParkir int, biayaParkir int) {
	parkirKeluar := ParkirKeluar{IdParkir: idParkir, WaktuMasuk: waktuMasuk, WaktuKeluar: waktuKeluar,
		DurasiParkir: durasiParkir, TipeKendaraan: tipeKendaraan, PlatNo: platNo,
		BiayaParkir: biayaParkir}
	dataParkirKeluar = append(dataParkirKeluar,parkirKeluar)
}

func hitungBiayaParkir(idParkir string, waktuMasuk time.Time, waktuKeluar time.Time, inpTipe string,
	inpPlatNo string) {
	if inpTipe == "Motor" {
		parkirMotor(idParkir, waktuMasuk, waktuKeluar, inpTipe, inpPlatNo)
	} else if inpTipe == "Mobil" {
		parkirMobil(idParkir, waktuMasuk, waktuKeluar, inpTipe, inpPlatNo)
	} else {
		fmt.Println("Tipe kendaraan yang anda masukkan salah.")
	}
}

func parkirMotor(idParkir string, waktuMasuk time.Time, waktuKeluar time.Time, tipeKendaraan string, platNo string) {
	diff := waktuKeluar.Sub(waktuMasuk)
	durasiParkir := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasiParkir == 1 {
		biayaParkir := 3000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	} else if durasiParkir >= satuHari{
		biayaParkir := ((satuHari - 1) * 2000) + 3000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	} else {
		biayaParkir := ((durasiParkir - 1) * 2000) + 3000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	}
}

func parkirMobil(idParkir string, waktuMasuk time.Time, waktuKeluar time.Time, tipeKendaraan string, platNo string) {
	diff := waktuKeluar.Sub(waktuMasuk)
	durasiParkir := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasiParkir == 1 {
		biayaParkir := 5000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	} else if durasiParkir >= satuHari{
		biayaParkir := ((satuHari - 1) * 3000) + 5000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	} else {
		biayaParkir := ((durasiParkir - 1) * 3000) + 5000
		simpanDataParkirKeluar(idParkir, waktuMasuk, waktuKeluar, tipeKendaraan, platNo,
			durasiParkir, biayaParkir)
	}
}

func customerMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	if r.Method == "POST" {
		idParkir := generateID()
		waktuMasuk := time.Now()
		karcis := ParkirMasuk{idParkir,  waktuMasuk}

		var result, err = json.Marshal(karcis)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(result)
		if err != nil {
			log.Fatal(err)
		}

		simpanDataParkirMasuk(idParkir, waktuMasuk)

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func customerKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var inpIdP = r.FormValue("id_parkir")
		var inpTipe = r.FormValue("tipe_kendaraan")
		var inpPlatNo = r.FormValue("plat_nomor")
		waktuKeluar := time.Now()

		var result []byte
		var err error

		for i := 0; i < len(dataParkirMasuk); i++ {
			if dataParkirMasuk[i].IdParkir == inpIdP {
				hitungBiayaParkir(dataParkirMasuk[i].IdParkir, dataParkirMasuk[i].WaktuMasuk, waktuKeluar, inpTipe, inpPlatNo)
				for _, dpk := range dataParkirKeluar {
					if dpk.IdParkir == dataParkirMasuk[i].IdParkir {
						result, err = json.Marshal(dpk)

						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						dataParkirMasuk = append(dataParkirMasuk[:i], dataParkirMasuk[i+1:]...)

						w.Write(result)
						return
					}
				}
			}
		}

		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func parkirMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(dataParkirMasuk)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func parkirKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(dataParkirKeluar)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}