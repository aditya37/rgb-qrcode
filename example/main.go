package main

import (
	"encoding/json"
	"log"
	"os"

	qr "github.com/aditya37/rgb-qrcode"
)

type JsonData struct {
	LocationId    int `json:"location_id"`
	SubLocationId int `json:"sub_location_id"`
}

func main() {
	// open file
	f, _ := os.Open("fleet.png")
	defer f.Close()

	j := JsonData{
		LocationId:    1,
		SubLocationId: 1,
	}
	bytejs, _ := json.Marshal(j)

	q, err := qr.New(qr.GenerateParam{
		LogoPath: f,
		QrValue:  string(bytejs),
		QrSize:   256,
	})
	if err != nil {
		log.Println(err)
		return
	}
	res, err := q.Encode()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Base64 => ", res.Base64)
}
