package iec61850

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
	IED               []IED             `xml:"IED"`
	DataTypeTemplates DataTypeTemplates `xml:"DataTypeTemplates"`
}

type IED struct {
	Name          string        `xml:"name,attr"`
	Type          string        `xml:"type,attr"`
	Desc          string        `xml:"desc,attr"`
	ConfigVersion string        `xml:"configVersion,attr"`
	AccessPoint   []AccessPoint `xml:"AccessPoint"`
}

type AccessPoint struct {
	Name    string    `xml:"name,attr"`
	LDevice []LDevice `xml:"Server>LDevice"`
}

type LDevice struct {
	Inst string `xml:"inst,attr"`
	LN   []LN   `xml:"LN"`
	LN0  LN     `xml:"LN0"`
}

type LN struct {
	Inst    string `xml:"inst,attr"`
	Prefix  string `xml:"prefix,attr"`
	LnType  string `xml:"lnType,attr"`
	LnClass string `xml:"lnClass,attr"`
	DOI     []DOI  `xml:"DOI"`
}

type DOI struct {
	Desc string `xml:"desc,attr"`
	Name string `xml:"name,attr"`
	DAI  []DAI  `xml:"DAI"`
}

type DAI struct {
	Name string `xml:"name,attr"`
	Val  Val    `xml:"Val"`
}

type Val struct {
	Value string `xml:",chardata"`
}

type DataSet struct {
	Name string      `xml:"name,attr"`
	Desc string      `xml:"desc,attr"`
	FCDA []FCDAEntry `xml:"FCDA"`
}

type FCDAEntry struct {
	LDInst  string `xml:"ldInst,attr,omitempty"`
	Prefix  string `xml:"prefix,attr,omitempty"`
	LNClass string `xml:"lnClass,attr"`
	LNInst  string `xml:"lnInst,attr,omitempty"`
	DOName  string `xml:"doName,attr"`
	DAName  string `xml:"daName,attr,omitempty"`
	FC      string `xml:"fc,attr"`
}

type DataTypeTemplates struct {
	LNodeType []LNodeType `xml:"LNodeType"`
	DOType    []DOType    `xml:"DOType"`
	DAType    []DAType    `xml:"DAType"`
	EnumType  []EnumType  `xml:"EnumType"`
}

type LNodeType struct {
	ID string `xml:"id,attr"`
	DO []DO   `xml:"DO"`
}

type DO struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}

type DOType struct {
	ID  string `xml:"id,attr"`
	DA  []DA   `xml:"DA"`
	SDO []SDO  `xml:"SDO"`
}

type DA struct {
	Name string  `xml:"name,attr"`
	Type string  `xml:"type,attr"`
	Val  DAValue `xml:"Val"`
	DA   []DA    `xml:"DA"`
}

type DAType struct {
	ID  string `xml:"id,attr"`
	BDA []BDA  `xml:"BDA"`
	DA  []DA   `xml:"DA"`
}

type BDA struct {
	Name string  `xml:"name,attr"`
	Type string  `xml:"type,attr"`
	Val  DAValue `xml:"Val"`
}

type EnumType struct {
	ID      string    `xml:"id,attr"`
	EnumVal []EnumVal `xml:"EnumVal"`
}

type EnumVal struct {
	Ord  int    `xml:"ord,attr"`
	Name string `xml:",chardata"`
}

type SDO struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	DA   []DA   `xml:"DA"`
}

func (scl *SCL) PrintHierarchy() {
	for _, ied := range scl.IED {
		ied.printHierarchy(0)
	}
	scl.DataTypeTemplates.printDataTypeTemplates(0)
}

func (ied *IED) printHierarchy(depth int) {
	fmt.Printf("%sName: %s, Type: %s, Desc: %s\n", getIndentation(depth), ied.Name, ied.Type, ied.Desc)

	for _, ap := range ied.AccessPoint {
		fmt.Printf("%sAccessPoint: %s\n", getIndentation(depth+1), ap.Name)
		for _, ld := range ap.LDevice {
			fmt.Printf("%sLDevice: %s\n", getIndentation(depth+2), ld.Inst)
			ld.LN0.printHierarchy(depth + 3)
			for _, ln := range ld.LN {
				ln.printHierarchy(depth + 3)
			}
		}
	}
}

func (ln *LN) printHierarchy(depth int) {
	if ln.Inst != "" {
		fmt.Printf("%sLN: %s, Prefix: %s, LnType: %s, LnClass: %s\n", getIndentation(depth), strings.TrimSpace(ln.Inst), strings.TrimSpace(ln.Prefix), strings.TrimSpace(ln.LnType), strings.TrimSpace(ln.LnClass))
	}

	for _, doi := range ln.DOI {
		if doi.Name != "" {
			fmt.Printf("%sDOI: %s, Desc: %s\n", getIndentation(depth+1), strings.TrimSpace(doi.Name), strings.TrimSpace(doi.Desc))
			for _, dai := range doi.DAI {
				if dai.Name != "" {
					fmt.Printf("%sDAI: %s\n", getIndentation(depth+2), strings.TrimSpace(dai.Name))
					dai.Val.printHierarchy(depth + 3)
				}
			}
		}
	}
}

func (val *Val) printHierarchy(depth int) {
	fmt.Printf("%sValue: %s\n", getIndentation(depth), val.Value)
}

func (dt DataTypeTemplates) printDataTypeTemplates(depth int) {
	fmt.Printf("%sDataTypeTemplates:\n", getIndentation(depth))
	for _, lnt := range dt.LNodeType {
		fmt.Printf("%sLNodeType: %s\n", getIndentation(depth+1), lnt.ID)
		for _, dt := range lnt.DO {
			printDO(dt, depth+2)
		}
	}
	for _, dot := range dt.DOType {
		fmt.Printf("%sDOType: %s\n", getIndentation(depth+1), dot.ID)
		for _, da := range dot.DA {
			printDA(da, depth+2)
		}
	}
	for _, dat := range dt.DAType {
		fmt.Printf("%sDAType: %s\n", getIndentation(depth+1), dat.ID)
		for _, bda := range dat.BDA {
			printBDA(bda, depth+2)
		}
		for _, da := range dat.DA {
			printDA(da, depth+2)
		}
	}
	for _, et := range dt.EnumType {
		fmt.Printf("%sEnumType: %s\n", getIndentation(depth+1), et.ID)
		for _, ev := range et.EnumVal {
			fmt.Printf("%sEnumVal: Ord: %d, Name: %s\n", getIndentation(depth+2), ev.Ord, ev.Name)
		}
	}
}

func printDO(do DO, depth int) {
	if do.Name != "" && do.Type != "" {
		fmt.Printf("%sDO: %s, Type: %s\n", getIndentation(depth), do.Name, do.Type)
	}
}

func printDA(da DA, depth int) {
	if da.Name != "" && da.Type != "" {
		fmt.Printf("%sDA: %s, Type: %s\n", getIndentation(depth), da.Name, da.Type)
		for _, subDA := range da.DA {
			printDA(subDA, depth+1)
		}
	}
}

func printBDA(bda BDA, depth int) {
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
