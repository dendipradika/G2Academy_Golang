package main

import (
	"bytes"
	"fmt"
	"net/url"
	"time"
)
import "net/http"
import "encoding/json"

var baseURL = "http://localhost:8080"

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

func main() {
	http.HandleFunc("/client/masuk", pengunjungMasuk)
	http.HandleFunc("/client/keluar", pengunjungKeluar)

	fmt.Println("starting web client at http://localhost:8089/")
	http.ListenAndServe(":8089", nil)
}

func ambilDataPengunjungMasuk() (ParkirMasuk, error) {
	var err error
	var client = &http.Client{}
	var karcis ParkirMasuk

	request, err := http.NewRequest("POST", baseURL+"/server/masuk", nil)
	if err != nil {
		return karcis, err
	}

	response, err := client.Do(request)
	if err != nil {
		return karcis, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&karcis)
	if err != nil {
		return karcis, err
	}

	return karcis,nil
}

func ambilDataPengunjungKeluar(idParkir string, tipeKendaraan string, platNo string) (ParkirKeluar, error) {
	var err error
	var client = &http.Client{}
	var struk ParkirKeluar

	var param = url.Values{}
	param.Set("id_parkir", idParkir)
	param.Set("tipe_kendaraan", tipeKendaraan)
	param.Set("plat_nomor", platNo)
	var payload = bytes.NewBufferString(param.Encode())

	request, err := http.NewRequest("POST", baseURL+"/server/keluar", payload)
	if err != nil {
		return struk, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return struk, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&struk)
	if err != nil {
		return struk, err
	}

	return struk, nil
}

func pengunjungMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	if r.Method == "POST" {
		var karcis, err = ambilDataPengunjungMasuk()

		if err != nil {
			fmt.Println("Error!", err.Error())
			return
		}

		fmt.Fprintln(w, "\tKarcis Masuk")
		fmt.Fprintln(w, "Id Parkir : ", karcis.IdParkir)
		fmt.Fprintln(w, "Waktu Masuk : ", karcis.WaktuMasuk.Format("2006-Jan-02 Monday 03:04:05"))
	}

	http.Error(w, "", http.StatusBadRequest)
}

func pengunjungKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var inpIdP = r.FormValue("id_parkir")
		var inpTipe = r.FormValue("tipe_kendaraan")
		var inpPlatNo = r.FormValue("plat_nomor")

		var struk, err = ambilDataPengunjungKeluar(inpIdP, inpTipe, inpPlatNo)

		if err != nil {
			fmt.Println("Error!", err.Error())
			return
		}

		fmt.Fprintln(w, "\tBiaya Parkir")
		fmt.Fprintln(w, "Id Parkir : ", struk.IdParkir)
		fmt.Fprintln(w, "Waktu Masuk : ", struk.WaktuMasuk.Format("2006-Jan-02 Monday 03:04:05"))
		fmt.Fprintln(w, "Waktu Keluar : ", struk.WaktuKeluar.Format("2006-Jan-02 Monday 03:04:05"))
		fmt.Fprintln(w, "Durasi : ", struk.DurasiParkir)
		fmt.Fprintln(w, "Tipe Kendaraan : ", struk.TipeKendaraan)
		fmt.Fprintln(w, "Plat Nomor : ", struk.PlatNo)
		fmt.Fprintln(w, "Biaya Parkir : Rp. ", struk.BiayaParkir)
	}

	http.Error(w, "", http.StatusBadRequest)

}