package main

import (
	"log"

	"github.com/dath-241/coin-price-be-go/services/price-service/routers"
	"github.com/dath-241/coin-price-be-go/utils"
)

func main() {
    if err := utils.ConnectMongoDB("mongodb://localhost:27017"); err != nil {
        log.Fatal(err.Error())  
    }

    r := routers.SetupRouter()
    r.Run(":3000")
}