package test

import (
	"fmt"
	"github.com/jifanchn/go-libiec61850/iec61850"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestIEC61850ServerDynamicTest(t *testing.T) {
	modelName := "testmodel"
	tcpPort := 102

	// Create IedModel
	model := iec61850.NewIedModel(modelName)
	defer model.Destroy()

	// Create LogicalDevice
	lDevice1 := model.CreateLogicalDevice("SENSORS")

	// Create LogicalNodes
	lln0 := lDevice1.CreateLogicalNode("LLN0")
	ttmp1 := lDevice1.CreateLogicalNode("TTMP1")
	ggio1 := lDevice1.CreateLogicalNode("GGIO1")

	// Create DataObjects
	lln0.CreateDataObjectCDC_ENS("Mod")
	lln0.CreateDataObjectCDC_ENS("Health")
	ttmp1TmpSv := ttmp1.CreateDataObjectCDC_SAV("TmpSv", false)
	ggio1.CreateDataObjectCDC_APC("AnOut1", 2) // Assuming the constant CDC_CTL_MODEL_HAS_CANCEL | CDC_CTL_MODEL_SBO_ENHANCED equates to 2

	// Get DataAttributes
	temperatureValue := ttmp1TmpSv.GetChild("instMag.f")
	temperatureTimestamp := ttmp1TmpSv.GetChild("t")

	// Create IedServer
	server := iec61850.NewIedServer(model)
	defer server.Destroy()

	// Start the server
	server.Start(tcpPort)

	// Set up a signal handler to gracefully shut down the server on SIGINT (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		for range c {
			server.Stop()
			fmt.Println("Server stopped.")
			os.Exit(0)
		}
	}()

	// Sample update loop for attributes
	val := 0.0
	for i := 0; i < 10; i++ { // limited loop to 10 iterations for test
		server.LockDataModel()

		server.UpdateUTCTimeAttributeValue(temperatureTimestamp, time.Now().UnixNano()/int64(time.Millisecond))
		server.UpdateFloatAttributeValue(temperatureValue, float32(val))

		server.UnlockDataModel()

		val += 0.1

		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(time.Hour)
}
