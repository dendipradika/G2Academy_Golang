package main

import (
	"fmt"
  "math/rand"
	"time"
)

type KarcisMasuk struct {
	id string
	time time.Time
}

type KarcisKeluar struct {
	id_ string
	platNomor_ string
	time_ time.Time
	jenisKendaraan_ string
}

func generateID() string{
	var acakID string
	
  rand.Seed(time.Now().UnixNano())
  for i:=0;i<10;i++{
		acakID += fmt.Sprint(rand.Intn(7))
  }
	return acakID
}

func generateKarcis() (string,time.Time){
	id := generateID()
	time := time.Now()

	return id,time
}

func main() {
	var kendaraan []KarcisMasuk
	var kendaraan_ KarcisKeluar
	var menu int = 0

	for menu != 99 {
		fmt.Println("\n## MENU PARKIR ##")
		fmt.Println("1. Parkir Masuk")
		fmt.Println("2. Parkir Keluar")
		fmt.Println("99. Exit")
		fmt.Print("--> Pilih Menu : ")
		fmt.Scan(&menu)
		fmt.Println()
		switch menu {
			case 1:
				id,time := generateKarcis()
				simpanKarcis := KarcisMasuk{id,time}

				kendaraan = append(kendaraan,simpanKarcis)
				for i:=0; i<len(kendaraan); i++ {
					fmt.Println("ID Parkir ke-", i+1, " : ", kendaraan[i].id)
				}
			case 2:
				now := time.Now()
					var subMenu int = 0
					for subMenu != 99 {
						fmt.Println("## SUB MENU ##")
						fmt.Println("1. Input ID Parkir")
						fmt.Println("99. Kembali ke Menu Utama")
						fmt.Print("--> Pilih Sub Menu : ")
						fmt.Scan(&subMenu)
						if subMenu == 1 {
							for i:=0; i<len(kendaraan); i++ {
							fmt.Print("\nInput ID Parkir : ")
							fmt.Scan(&kendaraan_.id_)
							if kendaraan_.id_ == kendaraan[i].id {
								waktuParkir := now.Sub(kendaraan[i].time).Seconds()

								fmt.Print("Jenis Kendaraan : ")
								fmt.Scan(&kendaraan_.jenisKendaraan_)
								fmt.Print("Plat Nomor Kendaraan : ")
								fmt.Scan(&kendaraan_.platNomor_)
								fmt.Println("---------------------------- >>")
								if kendaraan_.jenisKendaraan_ == "Mobil" {
									if waktuParkir > 1 {
										fmt.Println("ID Parkir : ", kendaraan_.id_)
										fmt.Println("Biaya Parkir Rp. : ", int(((waktuParkir-1)*3000)+5000))
										fmt.Println("Lama Parkir : ", waktuParkir)
										fmt.Println("---------------------------- >>\n")
									} else {
										fmt.Println("Biaya Parkir Rp. : ", 5000)
										fmt.Println("---------------------------- >>\n")
									}
								} else if kendaraan_.jenisKendaraan_ == "Motor" {
									if waktuParkir > 1 {
										fmt.Println("ID Parkir : ", kendaraan_.id_)
										fmt.Println("Biaya Parkir Rp. : ", int(((waktuParkir-1)*2000)+3000))
										fmt.Println("Lama Parkir : ", waktuParkir)
										fmt.Println("---------------------------- >>\n")
									} else {
										fmt.Println("Biaya Parkir Rp. : ", 3000)
										fmt.Println("---------------------------- >>\n")
									}
								}
								kendaraan = append(kendaraan[:i], kendaraan[i+1:]...)
							} else {
								fmt.Println("ID Parkir yang dimasukkan salah.\n")
							}
						}
						} else {
							fmt.Println("Anda salah input.\n")
						}
					}
			case 99:
				fmt.Println("Terima kasih\n")
				break
			default:
				fmt.Println("Anda salah input.\n")
				break
		}
	}
}