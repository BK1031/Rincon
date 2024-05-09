package utils

import (
	"math/rand"
	"rincon/config"
	"strconv"
)

func GenerateID(length int) int {
	if length == 0 {
		length, _ = strconv.Atoi(config.ServiceIDLength)
	}
	var id string
	for i := 0; i < length; i++ {
		id += strconv.Itoa(rand.Intn(10))
	}
	idInt, _ := strconv.Atoi(id)
	return idInt
}
