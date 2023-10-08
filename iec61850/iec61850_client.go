package iec61850

// #include <iec61850_client.h>
import "C"
import (
	"errors"
	"fmt"
	"github.com/jifanchn/go-libiec61850/iec61850/scl_xml"
	"time"
	"unsafe"
)

type ClientState int

const (
	IED_STATE_CLOSED ClientState = iota
	IED_STATE_CONNECTING
	IED_STATE_CONNECTED
	IED_STATE_CLOSING
)

type GoMmsValue struct {
	Type  MMSType     // MMS_VALUE ENUM
	Value interface{} // The Go representation of the value
}

type Option func(client *IedClient)

type IedClient struct {
	connection C.IedConnection
}

func NewIedClient(options ...Option) *IedClient {
	client := &IedClient{
		connection: C.IedConnection_create(),
	}

	for _, op := range options {
		if op != nil {
			op(client)
		}
	}

	return client
}

func ConnectTimeout(timeout time.Duration) func(*IedClient) {
	// #define CONFIG_MMS_CONNECTION_DEFAULT_CONNECT_TIMEOUT 10000
	return func(c *IedClient) {
		// replace to c time
		C.IedConnection_setConnectTimeout(c.connection, C.uint(timeout/1e6))
	}
}

func RequestTimeout(timeout time.Duration) func(*IedClient) {
	// #define CONFIG_MMS_CONNECTION_DEFAULT_TIMEOUT 5000
	return func(c *IedClient) {
		C.IedConnection_setRequestTimeout(c.connection, C.uint(timeout/1e6))
	}
}

func (client *IedClient) Connect(hostname string, tcpPort int) error {
	cHostname := C.CString(hostname)
	defer C.free(unsafe.Pointer(cHostname))

	var clientError C.IedClientError
	C.IedConnection_connect(client.connection, &clientError, cHostname, C.int(tcpPort))
	if clientError == C.IED_ERROR_ALREADY_CONNECTED {
		return nil
	} else if clientError != C.IED_ERROR_OK {
		return fmt.Errorf("failed to connect to %s:%d, clientError: %v", hostname, tcpPort, Err(clientError))
	}
	return nil
}

func (client *IedClient) State() ClientState {
	state := C.IedConnection_getState(client.connection)
	return ClientState(state)
}

func (client *IedClient) ReadBoolean(objectRef string, constraint FunctionalConstraint) (bool, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readBooleanValue(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return false, fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return bool(value), nil
}

func (client *IedClient) ReadFloat(objectRef string, constraint FunctionalConstraint) (float64, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readFloatValue(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return float64(0), fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return float64(value), nil
}

func (client *IedClient) ReadString(objectRef string, constraint FunctionalConstraint) (string, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readStringValue(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return "", fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return C.GoString(value), nil
}

func (client *IedClient) ReadInt32(objectRef string, constraint FunctionalConstraint) (int32, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readInt32Value(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return int32(0), fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return int32(value), nil
}

func (client *IedClient) ReadInt64(objectRef string, constraint FunctionalConstraint) (int64, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readInt64Value(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return int64(0), fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return int64(value), nil
}

func (client *IedClient) ReadUnsigned32(objectRef string, constraint FunctionalConstraint) (uint32, error) {
	cObjectRef := C.CString(objectRef)
	defer C.free(unsafe.Pointer(cObjectRef))

	var clientError C.IedClientError
	value := C.IedConnection_readUnsigned32Value(client.connection, &clientError, cObjectRef, C.FunctionalConstraint(constraint))

	if clientError != C.IED_ERROR_OK {
		return uint32(0), fmt.Errorf("failed to read object %s, clientError: %v", objectRef, Err(clientError))
	}

	return uint32(value), nil
}

func (client *IedClient) resolveValue(value *C.MmsValue, valueType MMSType) interface{} {
	goValue := interface{}(nil)

	// Refer to https://support.mz-automation.de/doc/libiec61850/c/latest/group__MMS__VALUE.html

	switch valueType {
	case MMS_BOOLEAN:
		realValue := bool(C.MmsValue_getBoolean(value))
		goValue = realValue
	case MMS_FLOAT:
		realValue := float64(C.MmsValue_toDouble(value))
		goValue = realValue
	case MMS_INTEGER:
		realValue := int64(C.MmsValue_toInt64(value))
		goValue = realValue
	case MMS_UNSIGNED:
		realValue := int64(C.MmsValue_toInt64(value))
		goValue = realValue
	case MMS_STRING:
		realValue := C.GoString(C.MmsValue_toString(value))
		goValue = realValue
	case MMS_VISIBLE_STRING:
		realValue := C.GoString(C.MmsValue_toString(value))
		goValue = realValue
	case MMS_STRUCTURE:
		goValue = client.digIntoStructure(value)
	case MMS_ARRAY:
		goValue = client.digIntoStructure(value)
	case MMS_BIT_STRING:
		goValue = uint32(C.MmsValue_getBitStringAsInteger(value))
	case MMS_UTC_TIME:
		goValue = uint32(C.MmsValue_toUnixTimestamp(value))
	}

	return goValue
}

func (client *IedClient) digIntoStructure(mms *C.MmsValue) []GoMmsValue {
	mmsType := C.MmsValue_getType(mms)
	if MMSType(mmsType) != MMS_STRUCTURE {
		return nil
	}
	goValues := make([]GoMmsValue, 0)
	index := 0
	for {
		value := C.MmsValue_getElement(mms, C.int(index))
		if value == nil {
			return goValues
		}
		valueType := C.MmsValue_getType(value)
		var goValue GoMmsValue
		goValue.Value = client.resolveValue(value, MMSType(valueType))
		goValue.Type = (MMSType)(valueType)
		goValues = append(goValues, goValue)
		index++
	}
}

func (client *IedClient) ReadDataSetValues(dataSetReference string, identifier string) ([]GoMmsValue, error) {
	var clientError C.IedClientError

	cDataSetReference := C.CString(dataSetReference)
	defer C.free(unsafe.Pointer(cDataSetReference))

	clientDataSet := C.IedConnection_readDataSetValues(client.connection, &clientError, cDataSetReference, nil)

	if clientError != C.IED_ERROR_OK {
		return nil, fmt.Errorf("failed to read dataset values, error code: %s", Err(clientError))
	}

	defer C.ClientDataSet_destroy(clientDataSet)

	// 获取数据集中的值
	values := C.ClientDataSet_getValues(clientDataSet)
	size := int(C.ClientDataSet_getDataSetSize(clientDataSet))

	goValues := make([]GoMmsValue, size)

	// 遍历数据集中的值
	for i := 0; i < size; i++ {
		value := C.MmsValue_getElement(values, C.int(i))
		valueType := MMSType(C.MmsValue_getType(value))
		// Refer to https://support.mz-automation.de/doc/libiec61850/c/latest/group__MMS__VALUE.html
		goValues[i].Type = valueType

		goValues[i].Value = client.resolveValue(value, valueType)
	}

	return goValues, nil
}

func (client *IedClient) ExplainDataSetValues(values []GoMmsValue, dSetScl *scl_xml.DataSetDetail) (map[string]interface{}, error) {
	if len(dSetScl.FCDA) != len(values) {
		return nil, errors.New("error dataset scl")
	}

	ret := make(map[string]interface{})
	for idx, entity := range dSetScl.FCDA {
		ref := fmt.Sprintf("%s/%s%s%s.%s", entity.LDInst, entity.Prefix, entity.LNClass, entity.LNInst, entity.DOName)
		val := values[idx]
		if entity.DAName != "" {
			ref += fmt.Sprintf(".%s", entity.DAName)
			ret[ref] = val.Value
		} else {
			if valueList, ok := val.Value.([]GoMmsValue); ok {
				doTyp := dSetScl.GetDOType(entity.Prefix, entity.LNClass, entity.DOName)
				for i, v := range valueList {
					var refNew string
					if len(doTyp.DA) > i+1 {
						refNew = fmt.Sprintf("%s.%s", ref, doTyp.DA[i].Name)
					} else {
						refNew = fmt.Sprintf("%s.%d", ref, i)
					}
					ret[refNew] = v.Value
				}
			}
		}
	}

	return ret, nil
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
	var clientError C.IedClientError

	dataAttributes := C.IedConnection_getDataDirectory(client.connection, &clientError, C.CString(doRef))
	if dataAttributes != nil {
		for dataAttribute := C.LinkedList_getNext(dataAttributes); dataAttribute != nil; dataAttribute = C.LinkedList_getNext(dataAttribute) {
			dataAttributeName := C.GoString((*C.char)(dataAttribute.data))

			printSpaces(spaces) // Assuming you've a function that prints spaces
			fmt.Printf("DA: %s\n", string(dataAttributeName))

			daRef := fmt.Sprintf("%s.%s", doRef, string(dataAttributeName))
			client.BrowseDataAttributes(daRef, spaces+2)
		}
	}

	C.LinkedList_destroy(dataAttributes)
}

func (client *IedClient) BrowseModel() {
	var clientError C.IedClientError

	// Get Logical Device List
	deviceList := C.IedConnection_getLogicalDeviceList(client.connection, &clientError)
	defer C.LinkedList_destroy(deviceList)

	if clientError != 0 {
		fmt.Printf("Failed to retrieve logical device list. Error: %s\n", Err(clientError))
		return
	}

	for device := C.LinkedList_getNext(deviceList); device != nil; device = C.LinkedList_getNext(device) {
		deviceName := C.GoString((*C.char)(device.data))
		fmt.Printf("LD: %s\n", string(deviceName))

		// Get Logical Node Directory
		logicalNodes := C.IedConnection_getLogicalDeviceDirectory(client.connection, &clientError, C.CString(deviceName))
		if clientError != 0 {
			fmt.Printf("Failed to retrieve logical nodes for device %v. Error: %s\n", deviceName, Err(clientError))
			continue
		}

		for logicalNode := C.LinkedList_getNext(logicalNodes); logicalNode != nil; logicalNode = C.LinkedList_getNext(logicalNode) {
			logicalNodeName := C.GoString((*C.char)(logicalNode.data))
			fmt.Printf("  LN: %v\n", logicalNodeName)

			lnRef := fmt.Sprintf("%s/%s", string(deviceName), string(logicalNodeName))

			// Browse DataObjects
			dataObjects := C.IedConnection_getLogicalNodeDirectory(client.connection, &clientError, C.CString(lnRef), C.ACSI_CLASS_DATA_OBJECT)
			for dataObject := C.LinkedList_getNext(dataObjects); dataObject != nil; dataObject = C.LinkedList_getNext(dataObject) {
				dataObjectName := C.GoString((*C.char)(dataObject.data))
				fmt.Printf("    DO: %s\n", string(dataObjectName))

				doRef := fmt.Sprintf("%s/%s.%s", string(deviceName), string(logicalNodeName), string(dataObjectName))

				client.BrowseDataAttributes(doRef, 6)
			}

			// Cleanup for DataObjects
			C.LinkedList_destroy(dataObjects)
		}

		// Cleanup for each logical node
		C.LinkedList_destroy(logicalNodes)
	}
}

func (client *IedClient) BrowseDataAttributesSCL(ref string) []scl_xml.DAI {
	var dais []scl_xml.DAI
	var clientError C.IedClientError

	attributes := C.IedConnection_getDataDirectory(client.connection, &clientError, C.CString(ref))
	defer C.LinkedList_destroy(attributes)

	if clientError != 0 {
		fmt.Printf("Failed to retrieve DAs for reference %s. Error: %s\n", ref, Err(clientError))
		return dais
	}

	for attribute := C.LinkedList_getNext(attributes); attribute != nil; attribute = C.LinkedList_getNext(attribute) {
		attributeName := C.GoString((*C.char)(attribute.data))
		childRef := fmt.Sprintf("%s.%s", ref, string(attributeName))

		dai := scl_xml.DAI{
			Name: attributeName,
			Val:  scl_xml.Val{Value: "SomeValue"}, // 这里简化了，实际上可能需要从远程设备读取属性值
		}

		// 递归获取SDI
		dai.SDI = client.BrowseSDISCL(childRef)

		dais = append(dais, dai)
	}

	return dais
}

func (client *IedClient) BrowseSDISCL(ref string) []scl_xml.SDI {
	var sdis []scl_xml.SDI
	var clientError C.IedClientError

	subdataObjects := C.IedConnection_getDataDirectory(client.connection, &clientError, C.CString(ref))
	defer C.LinkedList_destroy(subdataObjects)

	if clientError != 0 {
		fmt.Printf("Failed to retrieve SDIs for reference %s. Error: %s\n", ref, Err(clientError))
		return sdis
	}

	for sdo := C.LinkedList_getNext(subdataObjects); sdo != nil; sdo = C.LinkedList_getNext(sdo) {
		sdoName := C.GoString((*C.char)(sdo.data))
		childRef := fmt.Sprintf("%s.%s", ref, string(sdoName))

		sdi := scl_xml.SDI{
			Name: sdoName,
		}

		// 递归获取DAI和SDI
		sdi.DAI = client.BrowseDataAttributesSCL(childRef)
		sdi.SDI = client.BrowseSDISCL(childRef)

		sdis = append(sdis, sdi)
	}

	return sdis
}

func (client *IedClient) BrowseModelToSCL() (*scl_xml.SCL, error) {
	scl := &scl_xml.SCL{}
	var clientError C.IedClientError

	deviceList := C.IedConnection_getLogicalDeviceList(client.connection, &clientError)
	defer C.LinkedList_destroy(deviceList)

	if clientError != 0 {
		return nil, fmt.Errorf("failed to retrieve logical device list. Error: %s", Err(clientError))
	}

	for device := C.LinkedList_getNext(deviceList); device != nil; device = C.LinkedList_getNext(device) {
		deviceName := C.GoString((*C.char)(device.data))

		lDevice := scl_xml.LDevice{
			Inst: deviceName,
		}

		logicalNodes := C.IedConnection_getLogicalDeviceDirectory(client.connection, &clientError, C.CString(deviceName))
		if clientError != 0 {
			fmt.Printf("Failed to retrieve logical nodes for device %v. Error: %s\n", deviceName, Err(clientError))
			continue
		}

		for logicalNode := C.LinkedList_getNext(logicalNodes); logicalNode != nil; logicalNode = C.LinkedList_getNext(logicalNode) {
			logicalNodeName := C.GoString((*C.char)(logicalNode.data))

			ln := scl_xml.LN{
				Inst: logicalNodeName,
			}

			lnRef := fmt.Sprintf("%s/%s", string(deviceName), string(logicalNodeName))

			dataObjects := C.IedConnection_getLogicalNodeDirectory(client.connection, &clientError, C.CString(lnRef), C.ACSI_CLASS_DATA_OBJECT)
			for dataObject := C.LinkedList_getNext(dataObjects); dataObject != nil; dataObject = C.LinkedList_getNext(dataObject) {
				dataObjectName := C.GoString((*C.char)(dataObject.data))

				doi := scl_xml.DOI{
					Name: dataObjectName,
					DAI:  client.BrowseDataAttributesSCL(fmt.Sprintf("%s/%s.%s", string(deviceName), string(logicalNodeName), string(dataObjectName))),
				}

				ln.DOI = append(ln.DOI, doi)
			}

			C.LinkedList_destroy(dataObjects)

			lDevice.LN = append(lDevice.LN, ln)
		}

		C.LinkedList_destroy(logicalNodes)

		ied := scl_xml.IED{
			Name: deviceName,
			AccessPoint: []scl_xml.AccessPoint{
				{
					Name:    deviceName + "_AP",
					LDevice: []scl_xml.LDevice{lDevice},
				},
			},
		}

		scl.IED = append(scl.IED, ied)
	}

	return scl, nil
}
