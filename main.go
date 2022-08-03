package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	
)

var (
	unsupportedCriteria = errors.New("unsupported criteria")
	emptyDeparture      = errors.New("empty departure station")
	emptyArrival        = errors.New("empty arrival station")
	wrongArrivalInput   = errors.New("bad arrival station input")
	wrongDepartureInput = errors.New("bad departure station input")
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationID"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func main() {
	var (
		departureStation, arrivalStation, criteria string
	)
	fmt.Println("Введіть номер станції відправлення:")
	_, err := fmt.Scanf("%s", &departureStation)
	if err != nil {
		log.Fatalf("input error %s", err.Error())
	}

	fmt.Println("Введіть номер станції прибуття:")
	_, err = fmt.Scanf("%s", &arrivalStation)
	if err != nil {
		log.Fatalf("input error %s", err.Error())
	}

	fmt.Println("Введіть критерій для сортування потягів(можливі критерії price, arrival-time, departure-time):")
	_, err = fmt.Scanf("%s", &criteria)
	if err != nil {
		log.Fatalf("input error %s", err.Error())
	}

	//	... запит даних від користувача
	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(result) > 3 {
		result = result[:3]
	}
	for _, v := range result {
		fmt.Printf("%+v\n", v)
	}
	//	... обробка помилки
	//	... друк result
}


func (t *Train) UnmarshalJSON(data []byte) error {
	timeFormat := "15:04:05"
	type TrainClone Train

	trn := struct {
		ArrivalTime   string `json:"arrivalTime"`
		DepartureTime string `json:"departureTime"`
		*TrainClone
	}{
		TrainClone: (*TrainClone)(t),
	}

	err := json.Unmarshal(data, &trn)
	if err != nil {
		return err
	}

	arrivalTime, err := time.Parse(timeFormat, trn.ArrivalTime)
	if err != nil {
		return err
	}

	departureTime, err := time.Parse(timeFormat, trn.DepartureTime)
	if err != nil {
		return err
	}

	t.ArrivalTime = arrivalTime
	t.DepartureTime = departureTime
	return nil
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	criteria = strings.ToLower(criteria)

	var everyTrain, result Trains

	content, err := os.ReadFile("data.json")
	if err != nil {
		return nil, fmt.Errorf("file reading error: %w", err)
	}

	err = json.Unmarshal(content, &everyTrain)
	if err != nil {
		return nil, err
	}

	if departureStation == "" {
		return nil, emptyDeparture
	}

	if arrivalStation == "" {
		return nil, emptyArrival
	}

	departStation, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, wrongDepartureInput
	}

	arrivStation, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, wrongArrivalInput
	}

	if arrivStation < 1 {
		return nil, wrongArrivalInput
	}

	if departStation < 1 {
		return nil, wrongDepartureInput
	}

	switch criteria {
	case "price":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Price < result[j].Price
		})
	case "arrival-time":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].ArrivalTime.Before(result[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].DepartureTime.Before(result[j].DepartureTime)
		})
	default:
		return nil, unsupportedCriteria
	}
	
	return result, nil 
}
