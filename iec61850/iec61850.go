package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type IEDClient struct {
	connection C.IedConnection
}

func NewIEDClient() *IEDClient {
	return &IEDClient{
		connection: C.IedConnection_create(),
	}
}

func (client *IEDClient) Connect(hostname string, tcpPort int) error {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))

	var clientError C.IedClientError
	C.IedConnection_connect(client.connection, &clientError, cHostname, C.int(tcpPort))
	if clientError != C.IED_ERROR_OK {
		return fmt.Errorf("failed to connect to %s:%d, clientError: %v", hostname, tcpPort, clientError)
	}
	return nil
}

func (client *IEDClient) ReadObject(objectRef string) (float64, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readObject(client.connection, &clientError, cObjectRef, C.IEC61850_FC_MX)

	if clientError != C.IED_ERROR_OK {
		return 0, fmt.Errorf("failed to read object %s, clientError: %v", objectRef, clientError)
	}

	if C.MmsValue_getType(value) == C.MMS_FLOAT {
		floatVal := float64(C.MmsValue_toFloat(value))
		return floatVal, nil
	} else if C.MmsValue_getType(value) == C.MMS_DATA_ACCESS_ERROR {
		return 0, errors.New("data access clientError")
	}
	return 0, errors.New("unknown clientError")
}

func (client *IEDClient) Close() {
	C.IedConnection_close(client.connection)
	C.IedConnection_destroy(client.connection)
}
