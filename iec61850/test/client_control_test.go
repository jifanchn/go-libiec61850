package test

import (
	"fmt"
	"github.com/jifanchn/go-libiec61850/iec61850"
	"testing"
	"time"
)

func TestIEC61850DirectWithNormalSecurity(t *testing.T) {
	client := iec61850.NewIedClient(iec61850.ConnectTimeout(time.Second * 5))
	err := client.Connect("localhost", 10102)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	r1, err := client.ReadBoolean("BMS01DL_TRANS0/ykGAPC0.SPC00.stVal", iec61850.IEC61850_FC_ST)
	if err != nil {
		return
	}
	fmt.Printf("read success, r1 value: %+v\n", r1)

	err = client.DirectWithNormalSecurity("BMS01DL_TRANS0/ykGAPC0.SPC00", false)
	if err != nil {
		fmt.Println(err)
		return
	}

	r2, err := client.ReadBoolean("BMS01DL_TRANS0/ykGAPC0.SPC00.stVal", iec61850.IEC61850_FC_ST)
	if err != nil {
		return
	}
	fmt.Printf("read success, r2 value: %+v\n", r2)
}
