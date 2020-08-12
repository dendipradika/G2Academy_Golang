package main

import (
	"context"
	"fmt"
	"log"

	model "Day04/model"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	URL = "localhost:8080"
)

func main() {
	con, err := grpc.Dial(URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Not Connect: %v", err)
	}
	defer con.Close()
	newserv := model.NewParkingServiceClient(con)
	var input int

	for input != 99 {
		fmt.Println("## MENU ##")
		fmt.Println("1. Parkir Masuk")
		fmt.Println("2. Parkir Keluar")
		fmt.Println("99. EXIT")
		fmt.Println("Pilih Menu:")
		fmt.Scanln(&input)
		switch input {
		case 1:
			newId, newTime := GetId(newserv)
			fmt.Println("Id :", newId)
			fmt.Println("Time In :", newTime)
			input = 0
		case 2:
			var (
				id, platno, tipe string
			)
			fmt.Println("ID Parkir")
			fmt.Scan(&id)
			fmt.Println("Jenis Kendaraan (Mobil/Motor)")
			fmt.Scan(&tipe)
			fmt.Println("Plat Nomor")
			fmt.Scan(&platno)
			datain := &model.InputData{
				Id:     id,
				Tipe:   tipe,
				Platno: platno,
			}
			result := ParkingOut(newserv, datain)
			fmt.Println(result)
			input = 0
		}
	}
}

func GetId(mod model.ParkingServiceClient) (string, string) {
	resp, err := mod.ParkIn(context.Background(), new(empty.Empty))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return resp.Id, resp.Time
}

func ParkingOut(c model.ParkingServiceClient, input *model.InputData) string {
	resp, err := c.ParkOut(context.Background(), input)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return "Id :" + resp.Id + "\n" +
		"PlatNo :" + resp.Platno + "\n" +
		"Duration :" + resp.Duration + "\n" +
		"Message :" + resp.Message
}
