package test

import (
	"fmt"
	"github.com/jifanchn/go-libiec61850/iec61850"
	"testing"
	"time"
)

func TestIEC61850ClientReadValues(t *testing.T) {
	fmt.Println("CREATE SERVER....................")
	modelName := "test"
	tcpPort := 102

	// Create IedModel
	model := iec61850.NewIedModel(modelName)
	defer model.Destroy()
	// test/SENSORS/LLN0
	lDevice1 := model.CreateLogicalDevice("SENSORS")
	lln0 := lDevice1.CreateLogicalNode("LLN0")
	dataVSS := lln0.CreateDataObjectCDC_VSS("DATA")
	dataF := lln0.CreateDataObjectCDC_SAV("FLOAT", false)
	dataI := lln0.CreateDataObjectCDC_SAV("INT", true)
	fmt.Println(dataVSS)
	fmt.Println(dataF)
	fmt.Println(dataI)
	dataset := lln0.CreateDataSet("DATASET")
	fmt.Println(dataset)
	dataset.AddDataSetEntry("LLN0$ST$DATA$stVal")
	dataset.AddDataSetEntry("LLN0$MX$FLOAT$instMag$f")
	dataset.AddDataSetEntry("LLN0$MX$INT$instMag$i")

	server := iec61850.NewIedServer(model)

	// Start the server
	server.Start(tcpPort)

	server.LockDataModel()

	//		childString := dataVSS.GetChild("instMag.s")
	childFloat := dataF.GetChild("instMag.f")
	childInt := dataI.GetChild("instMag.i")
	childStr := dataVSS.GetChild("stVal")

	server.UpdateVisibleStringAttributeValue(childStr, "Hello World!")
	server.UpdateFloatAttributeValue(childFloat, 3.14159)
	server.UpdateInt32AttributeValue(childInt, 42)

	server.UnlockDataModel()

	client := iec61850.NewIedClient()
	err := client.Connect("localhost", 102)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	fmt.Println("BROWSER MODEL....................")
	client.BrowseModel()
	scl, err := client.BrowseModelToSCL()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("BROWSER SCL MODEL................")
	scl.Print()
	fmt.Println("READ VALUES......................")
	valueString, err := client.ReadString("testSENSORS/LLN0.DATA.stVal", iec61850.IEC61850_FC_ST)
	if err != nil {
		fmt.Println(err)
	}
	valInt, err := client.ReadInt32("testSENSORS/LLN0.INT.instMag.i", iec61850.IEC61850_FC_MX)
	if err != nil {
		fmt.Println(err)
	}
	valFloat, err := client.ReadFloat("testSENSORS/LLN0.FLOAT.instMag.f", iec61850.IEC61850_FC_MX)
	if err != nil {
		fmt.Println(err)
	}

	dataSet, err := client.ReadDataSetValues("testSENSORS/LLN0.DATASET", "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Result = %v, %v, %v\n", valueString, valInt, valFloat)
	fmt.Printf("Dataset = %v\n", dataSet)
	time.Sleep(time.Hour)
	server.Stop()
	server.Destroy()
}
