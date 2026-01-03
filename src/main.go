package main

import (
	"fmt"
	"github.com/Niroloc/Temcshiki/v2/src/context"
)

func main() {
	user := context.NewUser(1, "John")
	fmt.Println(*user)
}