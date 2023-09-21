package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type IedClient struct {
	connection C.IedConnection
}

func NewIedClient() *IedClient {
	return &IedClient{
		connection: C.IedConnection_create(),
	}
}

func (client *IedClient) Connect(hostname string, tcpPort int) error {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))

	var clientError C.IedClientError
	C.IedConnection_connect(client.connection, &clientError, cHostname, C.int(tcpPort))
	if clientError != C.IED_ERROR_OK {
		return fmt.Errorf("failed to connect to %s:%d, clientError: %v", hostname, tcpPort, clientError)
	}
	return nil
}

func (client *IedClient) ReadObjectFloatValue(objectRef string, constraint FunctionalConstraint) (float64, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readFloatValue(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return 0, fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return float64(value), nil
}

func (client *IedClient) Close() {
	C.IedConnection_close(client.connection)
	C.IedConnection_destroy(client.connection)
}

func printSpaces(count int) {
	for i := 0; i < count; i++ {
		fmt.Print(" ")
	}
}

func (client *IedClient) BrowseDataAttributes(doRef string, spaces int) {
	var error C.IedClientError

	dataAttributes := C.IedConnection_getDataDirectory(client.connection, &error, C.CString(doRef))
	if dataAttributes != nil {
		for dataAttribute := C.LinkedList_getNext(dataAttributes); dataAttribute != nil; dataAttribute = C.LinkedList_getNext(dataAttribute) {
			dataAttributeName := C.GoString((*C.char)(dataAttribute.data))

			printSpaces(spaces) // Assuming you've a function that prints spaces
			fmt.Printf("DA: %s\n", dataAttributeName)

			daRef := fmt.Sprintf("%s.%s", doRef, dataAttributeName)
			client.BrowseDataAttributes(daRef, spaces+2)
		}
	}

	C.LinkedList_destroy(dataAttributes)
}

func (client *IedClient) BrowseModel() {
	var error C.IedClientError

	// Get Logical Device List
	deviceList := C.IedConnection_getLogicalDeviceList(client.connection, &error)
	defer C.LinkedList_destroy(deviceList)

	if error != 0 {
		fmt.Printf("Failed to retrieve logical device list. Error: %d\n", error)
		return
	}

	for device := C.LinkedList_getNext(deviceList); device != nil; device = C.LinkedList_getNext(device) {
		deviceName := C.GoString((*C.char)(device.data))
		fmt.Printf("LD: %s\n", deviceName)

		// Get Logical Node Directory
		logicalNodes := C.IedConnection_getLogicalDeviceDirectory(client.connection, &error, C.CString(deviceName))
		if error != 0 {
			fmt.Printf("Failed to retrieve logical nodes for device %v. Error: %v\n", deviceName, error)
			continue
		}

		for logicalNode := C.LinkedList_getNext(logicalNodes); logicalNode != nil; logicalNode = C.LinkedList_getNext(logicalNode) {
			logicalNodeName := C.GoString((*C.char)(logicalNode.data))
			fmt.Printf("  LN: %v\n", logicalNodeName)

			lnRef := fmt.Sprintf("%s/%s", deviceName, logicalNodeName)

			// Browse DataObjects
			dataObjects := C.IedConnection_getLogicalNodeDirectory(client.connection, &error, C.CString(lnRef), C.ACSI_CLASS_DATA_OBJECT)
			for dataObject := C.LinkedList_getNext(dataObjects); dataObject != nil; dataObject = C.LinkedList_getNext(dataObject) {
				dataObjectName := C.GoString((*C.char)(dataObject.data))
				fmt.Printf("    DO: %s\n", dataObjectName)

				doRef := fmt.Sprintf("%s/%s.%s", deviceName, logicalNodeName, dataObjectName)

				client.BrowseDataAttributes(doRef, 6)
			}

			// Cleanup for DataObjects
			C.LinkedList_destroy(dataObjects)
		}

		// Cleanup for each logical node
		C.LinkedList_destroy(logicalNodes)
	}
}
