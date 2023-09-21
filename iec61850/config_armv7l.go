//go:build linux && arm

package iec61850

// #cgo CFLAGS: -I./inc/hal/inc -I./inc/common/inc -I./inc/goose -I./inc/iec61850/inc -I./inc/iec61850/inc_private -I./inc/logging -I./inc/mms/inc -I./inc/mms/inc_private -I./inc/mms/iso_mms/asn1c
// #cgo LDFLAGS: -static-libgcc -static-libstdc++ -L./lib/linux_armv7l -liec61850 -lhal
import "C"
