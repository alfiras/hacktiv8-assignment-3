package data

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Data struct {
	Status StatusData `json:"status"`
}

type StatusData struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

var Dep = make([]func(Data), 0)

func GetRandomData() Data {
	return Data{
		Status: StatusData{
			Water: rand.Intn(100),
			Wind:  rand.Intn(100),
		},
	}
}

func ReadFromJson() Data {
	file, err := os.ReadFile("data/status.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	data := Data{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}

func WriteToJson() {
	data := GetRandomData()
	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	_ = os.WriteFile("data/status.json", json, 0644)
	RunDep(data)
}

func AddDep(f func(Data)) {
	Dep = append(Dep, f)
}

func RunDep(data Data) {
	if len(Dep) > 0 {
		for _, d := range Dep {
			d(data)
		}
	}
}

func RunEvery() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Err", fmt.Sprintf("%v", r))
		}
	}()

	for {
		time.Sleep(15 * time.Second)
		WriteToJson()
	}
}
