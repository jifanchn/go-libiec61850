package scl_xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// IEC 61850 SCL (ICD) Data Structures

// DAValue interface for all possible data types
type DAValue interface{}

// Simple Data Types in Go for IEC 61850

type BOOLEAN bool
type INT32 int32
type INT16 int16
type INT8 int8
type FLOAT32 float32
type VisString255 string
type Unicode255 string

type SCL struct {
	IED               []IED             `scl_xml:"IED"`
	DataTypeTemplates DataTypeTemplates `scl_xml:"DataTypeTemplates"`
}

type IED struct {
	Name          string        `scl_xml:"name,attr"`
	Type          string        `scl_xml:"type,attr"`
	Desc          string        `scl_xml:"desc,attr"`
	ConfigVersion string        `scl_xml:"configVersion,attr"`
	AccessPoint   []AccessPoint `scl_xml:"AccessPoint"`
}

type AccessPoint struct {
	Name    string    `scl_xml:"name,attr"`
	LDevice []LDevice `scl_xml:"Server>LDevice"`
}

type LDevice struct {
	Inst string `scl_xml:"inst,attr"`
	LN   []LN   `scl_xml:"LN"`
	LN0  LN0    `scl_xml:"LN0"`
}

type LN0 struct {
	Inst     string    `scl_xml:"inst,attr"`
	LnType   string    `scl_xml:"lnType,attr"`
	LnClass  string    `scl_xml:"lnClass,attr"`
	DataSets []DataSet `scl_xml:"DataSet"`
}

type LN struct {
	Inst    string `scl_xml:"inst,attr"`
	Prefix  string `scl_xml:"prefix,attr"`
	LnType  string `scl_xml:"lnType,attr"`
	LnClass string `scl_xml:"lnClass,attr"`
	DOI     []DOI  `scl_xml:"DOI"`
}

type DOI struct {
	Desc string `scl_xml:"desc,attr"`
	Name string `scl_xml:"name,attr"`
	DAI  []DAI  `scl_xml:"DAI"`
	SDI  []SDI  `scl_xml:"SDI"`
}

type DAI struct {
	Name string `scl_xml:"name,attr"`
	Val  Val    `scl_xml:"Val"`
	SDI  []SDI  `scl_xml:"SDI"` // 新增
}

type SDI struct {
	Name string `scl_xml:"name,attr"`
	DAI  []DAI  `scl_xml:"DAI"`
	SDI  []SDI  `scl_xml:"SDI"` // 递归地包含SDI
}

type Val struct {
	Value string `scl_xml:",chardata"`
}

type DataSet struct {
	Name string      `scl_xml:"name,attr"`
	Desc string      `scl_xml:"desc,attr"`
	FCDA []FCDAEntry `scl_xml:"FCDA"`
}

type FCDAEntry struct {
	LDInst  string `scl_xml:"ldInst,attr,omitempty"`
	Prefix  string `scl_xml:"prefix,attr,omitempty"`
	LNClass string `scl_xml:"lnClass,attr"`
	LNInst  string `scl_xml:"lnInst,attr,omitempty"`
	DOName  string `scl_xml:"doName,attr"`
	DAName  string `scl_xml:"daName,attr,omitempty"`
	FC      string `scl_xml:"fc,attr"`
}

type DataTypeTemplates struct {
	LNodeType []LNodeType `scl_xml:"LNodeType"`
	DOType    []DOType    `scl_xml:"DOType"`
	DAType    []DAType    `scl_xml:"DAType"`
	EnumType  []EnumType  `scl_xml:"EnumType"`
}

type LNodeType struct {
	ID string `scl_xml:"id,attr"`
	DO []DO   `scl_xml:"DO"`
}

type DO struct {
	Name string `scl_xml:"name,attr"`
	Type string `scl_xml:"type,attr"`
	DA   []DA   `scl_xml:"DA"`
}

type DOType struct {
	ID  string `scl_xml:"id,attr"`
	DA  []DA   `scl_xml:"DA"`
	SDO []SDO  `scl_xml:"SDO"`
}

type DA struct {
	Name string  `scl_xml:"name,attr"`
	Type string  `scl_xml:"type,attr"`
	Val  DAValue `scl_xml:"Val"`
	DA   []DA    `scl_xml:"DA"`
}

type DAType struct {
	ID  string `scl_xml:"id,attr"`
	BDA []BDA  `scl_xml:"BDA"`
	DA  []DA   `scl_xml:"DA"`
}

type BDA struct {
	Name string  `scl_xml:"name,attr"`
	Type string  `scl_xml:"type,attr"`
	Val  DAValue `scl_xml:"Val"`
}

type EnumType struct {
	ID      string    `scl_xml:"id,attr"`
	EnumVal []EnumVal `scl_xml:"EnumVal"`
}

type EnumVal struct {
	Ord  int    `scl_xml:"ord,attr"`
	Name string `scl_xml:",chardata"`
}

type SDO struct {
	Name string `scl_xml:"name,attr"`
	Type string `scl_xml:"type,attr"`
	DA   []DA   `scl_xml:"DA"`
}

func (scl *SCL) Print() {
	for _, ied := range scl.IED {
		ied.Print(0)
	}
	scl.DataTypeTemplates.Print(0)
}

func (ied *IED) Print(depth int) {
	fmt.Printf("%sIED Name: %s, Type: %s, Desc: %s\n", getIndentation(depth), ied.Name, ied.Type, ied.Desc)
	for _, ap := range ied.AccessPoint {
		fmt.Printf("%sAccessPoint: %s\n", getIndentation(depth+1), ap.Name)
		for _, ld := range ap.LDevice {
			fmt.Printf("%sLDevice: %s\n", getIndentation(depth+2), ld.Inst)
			for _, ln := range ld.LN {
				ln.Print(depth + 3)
			}
		}
	}
}

func (ln *LN) Print(depth int) {
	fmt.Printf("%sLN Inst: %s, Prefix: %s, LnType: %s, LnClass: %s\n", getIndentation(depth), ln.Inst, ln.Prefix, ln.LnType, ln.LnClass)
	for _, doi := range ln.DOI {
		doi.Print(depth + 1)
	}
}

func (doi *DOI) Print(depth int) {
	fmt.Printf("%sDOI Name: %s, Desc: %s\n", getIndentation(depth), doi.Name, doi.Desc)
	for _, dai := range doi.DAI {
		dai.Print(depth + 1)
	}
	// Here, assuming you also want to print SDIs if they are included in your model
	for _, sdi := range doi.SDI {
		sdi.Print(depth + 1)
	}
}

func (dai *DAI) Print(depth int) {
	fmt.Printf("%sDAI Name: %s, Value: %s\n", getIndentation(depth), dai.Name, dai.Val.Value)
	for _, sdi := range dai.SDI {
		sdi.Print(depth + 1)
	}
}

func (sdi *SDI) Print(depth int) {
	// Print SDI and its related DAIs
	// Note: If SDIs can have nested SDIs, this function will need recursion
	fmt.Printf("%sSDI: %s\n", getIndentation(depth), sdi.Name)
	for _, dai := range sdi.DAI {
		fmt.Printf("%sDAI: %s\n", getIndentation(depth+1), dai.Name)
		dai.Val.Print(depth + 2)
	}
}

func (val *Val) Print(depth int) {
	fmt.Printf("%sValue: %s\n", getIndentation(depth), val.Value)
}

func (dt DataTypeTemplates) Print(depth int) {
	fmt.Printf("%sDataTypeTemplates:\n", getIndentation(depth))
	for _, lnt := range dt.LNodeType {
		fmt.Printf("%sLNodeType: %s\n", getIndentation(depth+1), lnt.ID)
		for _, dt := range lnt.DO {
			dt.Print(depth + 2)
		}
	}
	for _, dot := range dt.DOType {
		fmt.Printf("%sDOType: %s\n", getIndentation(depth+1), dot.ID)
		for _, da := range dot.DA {
			da.Print(depth + 2)
		}
	}
	for _, dat := range dt.DAType {
		fmt.Printf("%sDAType: %s\n", getIndentation(depth+1), dat.ID)
		for _, bda := range dat.BDA {
			bda.Print(depth + 2)
		}
		for _, da := range dat.DA {
			da.Print(depth + 2)
		}
	}
	for _, et := range dt.EnumType {
		fmt.Printf("%sEnumType: %s\n", getIndentation(depth+1), et.ID)
		for _, ev := range et.EnumVal {
			fmt.Printf("%sEnumVal: Ord: %d, Name: %s\n", getIndentation(depth+2), ev.Ord, ev.Name)
		}
	}
}

func (do DO) Print(depth int) {
	if do.Name != "" && do.Type != "" {
		fmt.Printf("%sDO: %s, Type: %s\n", getIndentation(depth), do.Name, do.Type)
	}
}

func (da DA) Print(depth int) {
	if da.Name != "" && da.Type != "" {
		fmt.Printf("%sDA: %s, Type: %s\n", getIndentation(depth), da.Name, da.Type)
		for _, subDA := range da.DA {
			subDA.Print(depth + 1)
		}
	}
}

func (bda BDA) Print(depth int) {
	if bda.Name != "" && bda.Type != "" {
		fmt.Printf("%sBDA: %s, Type: %s\n", getIndentation(depth), bda.Name, bda.Type)
	}
}

func getIndentation(depth int) string {
	return strings.Repeat("  ", depth)
}

func GetSCL(path string) (SCL, error) {
	// 打开并读取ICD文件
	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return SCL{}, errors.New("open file failed")
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var scl SCL
	err = xml.Unmarshal(byteValue, &scl)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return SCL{}, errors.New("unmarshall failed")
	}

	return scl, nil
}
