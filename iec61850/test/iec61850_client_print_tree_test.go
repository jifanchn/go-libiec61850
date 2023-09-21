package test

import (
	"fmt"
	"gitlab.weiheng-tech.com/SystemIntegrated/battery-pack-tools/61850-test/iec61850"
	"testing"
)

func TestIEC61850ClientPrintTree(t *testing.T) {
	client := iec61850.NewIedClient()
	err := client.Connect("localhost", 102)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	client.BrowseModel()

	client.Close()
}
