package main

import (
	"fmt"
	"gitlab.weiheng-tech.com/SystemIntegrated/battery-pack-tools/61850-test/iec61850"
)

func main() {
	client := iec61850.NewIEDClient()
	err := client.Connect("localhost", 102)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	value, err := client.ReadObject("BMSBC3/batModMMXN1.Vol279.mag.f")
	if err != nil {
		fmt.Println("Error reading object:", err)
	} else {
		fmt.Println("Read value:", value)
	}

	client.Close()
}
