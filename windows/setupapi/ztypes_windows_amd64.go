// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs windows/types.go

package setupapi

type SpDevinfoData struct {
	CbSize    uint32
	ClassGuid Guid
	DevInst   uint32
	Reserved  uint64
}
type Guid struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]uint8
}

const (
	SpDevinfoDataSz = 0x20
)
