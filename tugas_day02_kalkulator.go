package main

import (
	"fmt"
	"math"
)

func main() {
	var angkaKe1,angkaKe2,angkaKe3 int
	var input int = 0
	for input != 99 {
		fmt.Println("## Menu ##")
		fmt.Println("1. Pertambahan")
		fmt.Println("2. Pengurangan")
		fmt.Println("3. Pembagian")
		fmt.Println("4. Perkalian")
		fmt.Println("5. Akar")
		fmt.Println("6. Pangkat")
		fmt.Println("7. Luas Persegi")
		fmt.Println("8. Luas Lingkaran")
		fmt.Println("9. Volume tabung")
		fmt.Println("10. Volume Balok")
		fmt.Println("11. Volume Prisma")
		fmt.Println("99. Exit")
		fmt.Print("--> Pilih Menu : ")
		fmt.Scan(&input)
		switch input {
		case 1:
			fmt.Println("\nPertambahan")
			fmt.Print("Angka ke - 1 : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Angka ke - 2 : ")
			fmt.Scan(&angkaKe2)
			tambahkan := tambah(angkaKe1, angkaKe2)
			fmt.Println("Hasil Pertambahan :", tambahkan, "\n")
		case 2:
			fmt.Println("\nPengurangan")
			fmt.Print("Angka ke - 1 : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Angka ke - 2 : ")
			fmt.Scan(&angkaKe2)
			kurangi := kurang(angkaKe1,angkaKe2)
			fmt.Println("Hasil Pengurangan:", kurangi, "\n")
		case 3:
			fmt.Println("\nPembagian")
			fmt.Print("Angka ke - 1 : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Angka ke - 2 : ")
			fmt.Scan(&angkaKe2)
			dibagi := bagi(angkaKe1,angkaKe2)
			fmt.Println("Hasil Pembagian:", dibagi, "\n")
		case 4:
			fmt.Println("\nPerkalian")
			fmt.Print("Angka ke - 1 : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Angka ke - 2 : ")
			fmt.Scan(&angkaKe2)
			dikali := kali(angkaKe1,angkaKe2)
			fmt.Println("Hasil Perkalian:", dikali, "\n")
		case 5:
			fmt.Println("\nAkar")
			fmt.Print("Masukkan Angka : ")
			fmt.Scan(&angkaKe1)
			akar_ := akar(float64(angkaKe1))
			fmt.Println("Hasil Perkalian:", akar_, "\n")
		case 6:
			fmt.Println("\nPerpangkatan")
			fmt.Print("Angka ke - 1 : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Angka ke - 2 : ")
			fmt.Scan(&angkaKe2)
			perpangkatan := kali(angkaKe1,angkaKe2)
			fmt.Println("Hasil Perpangkatan:", perpangkatan, "\n")
		case 7:
			fmt.Println("\nLuas Persegi")
			fmt.Print("Masukkan Angka : ")
			fmt.Scan(&angkaKe1)
			persegi := Persegi{float64(angkaKe1)}
			fmt.Println("Hasil Luas Persegi :",persegi.luaspersegi(), "\n")
		case 8:
			fmt.Println("\nLuas lingkaran")
			fmt.Print("Masukkan Angka : ")
			fmt.Scan(&angkaKe1)
			fmt.Println("Hasil Luas lingkaran :",luasLingkaran(float64(angkaKe1)), "\n")
		case 9:
			fmt.Println("\nVolume Tabung")
			fmt.Print("Luas Alas : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Tinggi : ")
			fmt.Scan(&angkaKe2)
			fmt.Println("Hasil Volume Tabung :",volumeTabung(float64(angkaKe1), float64(angkaKe2)), "\n")
		case 10:
			fmt.Println("\nVolume Balok")
			fmt.Print("Panjang : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Lebar : ")
			fmt.Scan(&angkaKe2)
			fmt.Print("Tinggi : ")
			fmt.Scan(&angkaKe3)
			fmt.Println("Hasil Volume Balok :",volumeBalok(float64(angkaKe1), float64(angkaKe2), float64(angkaKe3)), "\n")
		case 11:
			fmt.Println("\nVolume Prisma")
			fmt.Print("Luas Alas : ")
			fmt.Scan(&angkaKe1)
			fmt.Print("Luas Tutup : ")
			fmt.Scan(&angkaKe2)
			fmt.Print("Luas Selimut : ")
			fmt.Scan(&angkaKe3)
			fmt.Println("Hasil Volume Prisma :", volumePrisma(float64(angkaKe1), float64(angkaKe2), float64(angkaKe3)), "\n")
		case 99:
			fmt.Println("Terima kasih\n")
			break
		default:
			fmt.Println("Anda salah input.\n")
			break
		}
	}
}

type Persegi struct {
	sisi float64
}

type Segitiga struct {
	alas float32
	tinggi float32
}

func input() (int,int) {
	var angkaKe1,angkaKe2 int
	fmt.Scan(&angkaKe1)
	fmt.Scan(&angkaKe2)
	return angkaKe1,angkaKe2
}

func change(original *int, value int) {
	*original = value
}

func tambah(angkaKe1, angkaKe2 int) int {
	return angkaKe1 + angkaKe2;
}

func kurang(angkaKe1, angkaKe2 int) int {
	return angkaKe1 - angkaKe2;
}

func bagi(angkaKe1, angkaKe2 int) int {
	return angkaKe1 / angkaKe2;
}

func kali(angkaKe1, angkaKe2 int) int {
	return angkaKe1 * angkaKe2;
}

func akar(angkaKe1 float64) float64 {
	return math.Sqrt(angkaKe1);
}

func pangkat(angkaKe1, angkaKe2 float64) float64 {
	pangkat_ := math.Pow(angkaKe1, angkaKe2)
	return pangkat_
}

func (angkaKe1 Persegi) luaspersegi() float64 {
	return math.Pow(angkaKe1.sisi,2)
}

func luasLingkaran(angkaKe1 float64) float64 {
	return angkaKe1*angkaKe1*math.Pi;
}

func volumeTabung(r, t float64) float64 {
	return math.Pi*r*r*t;
}

func volumeBalok(p, l, t float64) float64{
	return p*l*t;
}

func volumePrisma(a, t, s float64) float64 {
	return ((a*t)/2)*s;
}