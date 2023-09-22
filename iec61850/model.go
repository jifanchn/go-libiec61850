package iec61850

/*
#include "iec61850_server.h"
*/
import "C"
import "unsafe"

type IedModel struct {
	model *C.IedModel
}

func NewIedModel(name string) *IedModel {
	return &IedModel{
		model: C.IedModel_create(C.CString(name)),
	}
}

func (m *IedModel) Destroy() {
	C.IedModel_destroy(m.model)
}

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

// ENS: EnumerationString
// VSS: Visible String Setting
// SAV: Sampled Value
// APC: Analogue Process Control

func (n *LogicalNode) CreateDataObjectCDC_ENS(name string) *DataObject {
	return &DataObject{
		object: C.CDC_ENS_create(C.CString(name), (*C.ModelNode)(n.node), 0),
	}
}

func (n *LogicalNode) CreateDataObjectCDC_VSS(name string) *DataObject {
	return &DataObject{
		object: C.CDC_VSS_create(C.CString(name), (*C.ModelNode)(n.node), 0),
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

type DataSet struct {
	dataSet *C.DataSet
}

// CreateDataSet creates a new DataSet under this LogicalNode.
func (ln *LogicalNode) CreateDataSet(name string) *DataSet {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cDataSet := C.DataSet_create(cName, ln.node)
	return &DataSet{dataSet: cDataSet}
}

// AddDataSetEntry adds a new DataSetEntry to this DataSet.
func (ds *DataSet) AddDataSetEntry(ref string) {
	cRef := C.CString(ref)
	defer C.free(unsafe.Pointer(cRef))

	C.DataSetEntry_create(ds.dataSet, cRef, -1, nil)
}
