package main

import (
	"fmt"
	"github.com/go-faker/faker/v4"
)

func main() {
	faker.SetServerLang("id")
	fmt.Println(faker.Name())
	fmt.Println(faker.Name())
}
