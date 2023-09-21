package iec61850

/*
#include "iec61850_server.h"
*/
import "C"

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
