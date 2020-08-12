package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	model "Day04/model"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	port = ": 8080"
)

type ParkingService struct{}

var Parking = make(map[string]time.Time)

func main() {
	log.Println("Server Running --> PORT", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	var parkingService ParkingService
	model.RegisterParkingServiceServer(s, parkingService)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Print(s.Serve(lis))
}

func generateID() string {
	var acakID string

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		acakID += fmt.Sprint(rand.Intn(7))
	}
	return acakID
}

func (ParkingService) ParkIn(ctx context.Context, empty *empty.Empty) (*model.DataIn, error) {
	newId := generateID()
	newTime := time.Now()
	TimeToString := newTime.String()
	Parking[newId] = newTime
	log.Println("Id Parking:", newId, "\t", "Time In :", newTime)
	message := &model.DataIn{
		Id:   newId,
		Time: TimeToString,
	}
	return message, nil
}

func (ParkingService) ParkOut(ctx context.Context, DataInput *model.InputData) (*model.DataOut, error) {
	id := DataInput.Id
	setNewPark := Parking[id]
	tempcounttime := time.Now().Unix() - setNewPark.Unix()
	counttime := int(tempcounttime)
	strcounttime := strconv.Itoa(counttime)
	var tarif int
	tipe := DataInput.Tipe
	platno := DataInput.Platno
	msg := ""

	if _, found := Parking[id]; found {
		switch tipe {
		case "Motor":
			if counttime >= 1 {
				tarif = 3000
			}
			tarif = tarif + ((counttime - 1) * 2000)
			tariffix := strconv.Itoa(tarif)
			msg = "Biaya Parkir anda: Rp." + tariffix
			delete(Parking, id)

		case "Mobil":
			if counttime >= 1 {
				tarif = 5000
			}
			tarif = tarif + ((counttime - 1) * 3000)
			tariffix := strconv.Itoa(tarif)
			msg = "Biaya Parkir anda: Rp." + tariffix
			delete(Parking, id)
		}
	} else {
		msg = "Parking Id not found please input another Id"
	}
	DataOut := &model.DataOut{
		Id:       id,
		Platno:   platno,
		Duration: strcounttime,
		Message:  msg,
	}
	return DataOut, nil
}
