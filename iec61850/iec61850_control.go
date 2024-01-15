package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// DirectWithNormalSecurity Get the value of the specified data set
func (client *IedClient) DirectWithNormalSecurity(controlReference string, val bool) error {
	//var clientError C.IedClientError

	cDataSetReference := C.CString(controlReference)
	defer C.free(unsafe.Pointer(cDataSetReference))

	control := C.ControlObjectClient_create(cDataSetReference, client.connection)
	if control == nil {
		return fmt.Errorf("error creating control object client")
	}

	defer C.ControlObjectClient_destroy(control)

	ctlVal := C.MmsValue_newBoolean(C._Bool(val))
	defer C.MmsValue_delete(ctlVal)

	C.ControlObjectClient_setOrigin(control, nil, 3)

	if bool(C.ControlObjectClient_operate(control, ctlVal, 0)) {
		return nil
	} else {
		return fmt.Errorf("failed to operate %s\n", controlReference)
	}

}
