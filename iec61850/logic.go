package iec61850

/*
#include "iec61850_server.h"
*/
import "C"
import "unsafe"

type LogicalDevice struct {
	device *C.LogicalDevice
}

func (m *IedModel) CreateLogicalDevice(name string) *LogicalDevice {
	return &LogicalDevice{
		device: C.LogicalDevice_create(C.CString(name), m.model),
	}
}

type LogicalNode struct {
	node *C.LogicalNode
}

func (d *LogicalDevice) CreateLogicalNode(name string) *LogicalNode {
	return &LogicalNode{
		node: C.LogicalNode_create(C.CString(name), d.device),
	}
}

type DataObject struct {
	object *C.DataObject
}

func (n *LogicalNode) CreateDataObjectCDC_ENS(name string) *DataObject {
	return &DataObject{
		object: C.CDC_ENS_create(C.CString(name), (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_SAV(name string, isInteger bool) *DataObject {
	return &DataObject{
		object: C.CDC_SAV_create(C.CString(name), (*C.ModelNode)(n.node), 0, C.bool(isInteger)),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_APC(name string, ctlModel int) *DataObject {
	return &DataObject{
		object: C.CDC_APC_create(C.CString(name), (*C.ModelNode)(n.node), 0, C.uint(ctlModel), C.bool(false)),
	}
}

type DataAttribute struct {
	attribute *C.DataAttribute
}

func (do *DataObject) GetChild(name string) *DataAttribute {
	return &DataAttribute{
		attribute: (*C.DataAttribute)(unsafe.Pointer(C.ModelNode_getChild((*C.ModelNode)(unsafe.Pointer(do.object)), C.CString(name)))),
	}
}
