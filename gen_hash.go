package main
import (
	"fmt"
	"alhikmah-attendance-api/pkg/utils"
)
func main() {
	hash, _ := utils.HashPassword("password123")
	fmt.Println(hash)
}
