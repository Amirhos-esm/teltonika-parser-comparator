package teltonica

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	Priority_low int = iota
	Priority_high
	Priority_panic
)

// Define the struct in Go
type GpsElement struct {
	Longitude  float64 // in decimal degrees
	Latitude   float64 // in decimal degrees
	Altitude   uint16  // uint16 in C
	Angle      uint16  // uint16 in C
	Satellites uint8   // uint8 in C
	Speed      uint16  // uint16 in C
}

// Define the struct in Go
type IoElement struct {
	Id    int
	Value int64
	Size  uint8
}

type IoElements struct {
	N1        map[int]IoElement
	N2        map[int]IoElement
	N4        map[int]IoElement
	N8        map[int]IoElement
	EventIoId uint8
}
type AvlData struct {
	Timestamp uint64 // Timestamp in milliseconds (1185345998335 →1185345998,335 in Unix Timestamp = 25 Jul 2007 06:46:38 UTC)
	Priority  uint8
	Gps       GpsElement
	Io        IoElements
}
type AvlPacket struct {
	CodecId      uint8
	numberOfData uint8
	AvlArray     []AvlData
}

// Function to add N1 element
func (u *IoElements) AddN1(ioID uint8, ioValue uint8) error {
	if u.N1 == nil {
		u.N1 = map[int]IoElement{}
	}
	if len(u.N1) >= 255 {
		return errors.New("too many N1 elements")
	}
	u.N1[int(ioID)] = IoElement{Id: int(ioID), Value: int64(ioValue), Size: 1}
	return nil
}

// Function to add N2 element
func (u *IoElements) AddN2(ioID uint8, ioValue uint16) error {
	if u.N2 == nil {
		u.N2 = map[int]IoElement{}
	}
	if len(u.N2) >= 255 {
		return errors.New("too many N2 elements")
	}
	u.N2[int(ioID)] = IoElement{Id: int(ioID), Value: int64(ioValue), Size: 2}
	return nil
}

// Function to add N4 element
func (u *IoElements) AddN4(ioID uint8, ioValue uint32) error {
	if u.N4 == nil {
		u.N4 = map[int]IoElement{}
	}
	if len(u.N4) >= 255 {
		return errors.New("too many N4 elements")
	}
	u.N4[int(ioID)] = IoElement{Id: int(ioID), Value: int64(ioValue), Size: 4}
	return nil
}

// Function to add N8 element
func (u *IoElements) AddN8(ioID uint8, ioValue uint64) error {
	if u.N8 == nil {
		u.N8 = map[int]IoElement{}
	}
	if len(u.N8) >= 255 {
		return errors.New("too many N8 elements")
	}
	u.N8[int(ioID)] = IoElement{Id: int(ioID), Value: int64(ioValue), Size: 8}
	return nil
}

// Serialize serializes AvlData into a byte slice
func (u *AvlData) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Set byte order to big-endian
	binary.Write(buf, binary.BigEndian, u.Timestamp)
	binary.Write(buf, binary.BigEndian, u.Priority)

	// Serialize GPS elements
	binary.Write(buf, binary.BigEndian, int32(u.Gps.Longitude*1e7))
	binary.Write(buf, binary.BigEndian, int32(u.Gps.Latitude*1e7))
	binary.Write(buf, binary.BigEndian, u.Gps.Altitude)
	binary.Write(buf, binary.BigEndian, u.Gps.Angle)
	binary.Write(buf, binary.BigEndian, u.Gps.Satellites)
	binary.Write(buf, binary.BigEndian, u.Gps.Speed)

	// Serialize IO elements
	binary.Write(buf, binary.BigEndian, u.Io.EventIoId)

	// Count of all IO elements
	totalElements := uint8(len(u.Io.N1) + len(u.Io.N2) + len(u.Io.N4) + len(u.Io.N8))
	binary.Write(buf, binary.BigEndian, totalElements)

	// Serialize N1 elements
	if len(u.Io.N1) > 255 {
		return nil, errors.New("too many N1 elements")
	}
	binary.Write(buf, binary.BigEndian, uint8(len(u.Io.N1)))
	for id, elem := range u.Io.N1 {
		binary.Write(buf, binary.BigEndian, uint8(id))
		binary.Write(buf, binary.BigEndian, uint8(elem.Value))
	}

	// Serialize N2 elements
	if len(u.Io.N2) > 255 {
		return nil, errors.New("too many N2 elements")
	}
	binary.Write(buf, binary.BigEndian, uint8(len(u.Io.N2)))
	for id, elem := range u.Io.N2 {
		binary.Write(buf, binary.BigEndian, uint8(id))
		binary.Write(buf, binary.BigEndian, uint16(elem.Value))
	}

	// Serialize N4 elements
	if len(u.Io.N4) > 255 {
		return nil, errors.New("too many N4 elements")
	}
	binary.Write(buf, binary.BigEndian, uint8(len(u.Io.N4)))
	for id, elem := range u.Io.N4 {
		binary.Write(buf, binary.BigEndian, uint8(id))
		binary.Write(buf, binary.BigEndian, uint32(elem.Value))
	}

	// Serialize N8 elements
	if len(u.Io.N8) > 255 {
		return nil, errors.New("too many N8 elements")
	}
	binary.Write(buf, binary.BigEndian, uint8(len(u.Io.N8)))
	for id, elem := range u.Io.N8 {
		binary.Write(buf, binary.BigEndian, uint8(id))
		binary.Write(buf, binary.BigEndian, uint64(elem.Value))
	}

	return buf.Bytes(), nil
}

// Deserialize deserializes a byte slice into AvlData
// Deserialize deserializes a byte slice into AvlData and returns the number of bytes read
func (u *AvlData) Deserialize(buf *bytes.Reader) (int, error) {

	startLen := buf.Len() // Capture the initial length of the buffer
	// Initialize the bytes read counter

	// Set byte order to big-endian and read fields
	if err := binary.Read(buf, binary.BigEndian, &u.Timestamp); err != nil {
		return 0, err
	}

	if err := binary.Read(buf, binary.BigEndian, &u.Priority); err != nil {
		return 0, err
	}

	// Deserialize GPS elements
	var longitude, latitude int32
	if err := binary.Read(buf, binary.BigEndian, &longitude); err != nil {
		return 0, err
	}

	if err := binary.Read(buf, binary.BigEndian, &latitude); err != nil {
		return 0, err
	}

	u.Gps.Longitude = float64(longitude) / 1e7
	u.Gps.Latitude = float64(latitude) / 1e7

	if err := binary.Read(buf, binary.BigEndian, &u.Gps.Altitude); err != nil {
		return 0, err
	}

	if err := binary.Read(buf, binary.BigEndian, &u.Gps.Angle); err != nil {
		return 0, err
	}

	if err := binary.Read(buf, binary.BigEndian, &u.Gps.Satellites); err != nil {
		return 0, err
	}

	if err := binary.Read(buf, binary.BigEndian, &u.Gps.Speed); err != nil {
		return 0, err
	}

	// Deserialize IO elements
	if err := binary.Read(buf, binary.BigEndian, &u.Io.EventIoId); err != nil {
		return 0, err
	}

	var totalElements uint8
	if err := binary.Read(buf, binary.BigEndian, &totalElements); err != nil {
		return 0, err
	}

	// Deserialize N1 elements
	var n1Count uint8
	if err := binary.Read(buf, binary.BigEndian, &n1Count); err != nil {
		return 0, err
	}

	u.Io.N1 = make(map[int]IoElement, n1Count)
	for i := 0; i < int(n1Count); i++ {
		var id, value uint8
		if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
			return 0, err
		}

		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return 0, err
		}

		u.Io.N1[int(id)] = IoElement{Id: int(id), Value: int64(value), Size: 1}
	}

	// Deserialize N2 elements
	var n2Count uint8
	if err := binary.Read(buf, binary.BigEndian, &n2Count); err != nil {
		return 0, err
	}

	u.Io.N2 = make(map[int]IoElement, n2Count)
	for i := 0; i < int(n2Count); i++ {
		var id uint8
		var value uint16
		if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
			return 0, err
		}

		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return 0, err
		}

		u.Io.N2[int(id)] = IoElement{Id: int(id), Value: int64(value), Size: 2}
	}

	// Deserialize N4 elements
	var n4Count uint8
	if err := binary.Read(buf, binary.BigEndian, &n4Count); err != nil {
		return 0, err
	}

	u.Io.N4 = make(map[int]IoElement, n4Count)
	for i := 0; i < int(n4Count); i++ {
		var id uint8
		var value uint32
		if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
			return 0, err
		}

		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return 0, err
		}

		u.Io.N4[int(id)] = IoElement{Id: int(id), Value: int64(value), Size: 4}
	}

	// Deserialize N8 elements
	var n8Count uint8
	if err := binary.Read(buf, binary.BigEndian, &n8Count); err != nil {
		return 0, err
	}

	u.Io.N8 = make(map[int]IoElement, n8Count)
	for i := 0; i < int(n8Count); i++ {
		var id uint8
		var value uint64
		if err := binary.Read(buf, binary.BigEndian, &id); err != nil {
			return 0, err
		}

		if err := binary.Read(buf, binary.BigEndian, &value); err != nil {
			return 0, err
		}

		u.Io.N8[int(id)] = IoElement{Id: int(id), Value: int64(value), Size: 8}
	}

	return startLen - buf.Len(), nil
}

func (gps *GpsElement) IsValid() bool {
	return gps.Altitude == 0 && gps.Satellites == 0 && gps.Speed == 0
}

func (u *AvlPacket) SerializeTcp() ([]byte, error){
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(0))
	data,err :=u.Serialize()
	if err != nil{
		return nil,err
	}
	binary.Write(buf, binary.BigEndian, uint32(len(data)))
	binary.Write(buf, binary.BigEndian, data)
	binary.Write(buf, binary.BigEndian, uint32(calculateCRC16(data)))
	return buf.Bytes(),nil
}

func (u *AvlPacket) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Set byte order to big-endian

	u.numberOfData = uint8(len(u.AvlArray))
	binary.Write(buf, binary.BigEndian, u.CodecId)
	binary.Write(buf, binary.BigEndian, u.numberOfData)
	for _, value := range u.AvlArray {
		output, err := value.Serialize()
		if err != nil {
			return nil, err
		}
		binary.Write(buf, binary.BigEndian, output)
	}
	binary.Write(buf, binary.BigEndian, u.numberOfData)
	return buf.Bytes(), nil
}
func NewAvlPacket(data []byte) (*AvlPacket, error) {
	u := AvlPacket{}
	buf := bytes.NewReader(data)

	// Read CodecId (uint8)
	if err := binary.Read(buf, binary.BigEndian, &u.CodecId); err != nil {
		return nil, err
	}

	// Read number of data (uint8)
	if err := binary.Read(buf, binary.BigEndian, &u.numberOfData); err != nil {
		return nil, err
	}
	startIndex := 2
	// Deserialize each AVL data
	u.AvlArray = make([]AvlData, u.numberOfData) // Assuming AvlData is the type stored in AvlArray
	for i := 0; i < int(u.numberOfData); i++ {
		var avlData AvlData = AvlData{}
		// chunk := data[startIndex:]
		offset, err := avlData.Deserialize(buf)
		if err != nil {
			return nil, err
		}
		startIndex += offset
		u.AvlArray[i] = avlData
	}

	// Read number of data again (uint8) for consistency check
	var numberOfDataAgain uint8
	if err := binary.Read(buf, binary.BigEndian, &numberOfDataAgain); err != nil {
		return nil, err
	}
	// if numberOfDataAgain != u.numberOfData {
	// 	return errors.New(
	// 		fmt.Sprint("number of data mismatch: %d vs %d", u.numberOfData, numberOfDataAgain),
	// 	)
	// }
	if buf.Len() != 0 {
		return nil, errors.New(fmt.Sprint("buff len is not zeor : ", buf.Len()))
	}

	return &u, nil
}

// Print method for AvlData
func (avlData *AvlData) Print() {
	fmt.Printf("  Timestamp: %d\n", avlData.Timestamp)
	fmt.Printf("  Priority: %d\n", avlData.Priority)
	fmt.Printf("  GPS - Latitude: %f, Longitude: %f, Altitude: %d, Angle: %d, Satellites: %d, Speed: %d\n",
		avlData.Gps.Latitude, avlData.Gps.Longitude, avlData.Gps.Altitude, avlData.Gps.Angle, avlData.Gps.Satellites, avlData.Gps.Speed)
	fmt.Printf("  IO Elements:\n")
	for n, ioMap := range map[string]map[int]IoElement{
		"N1": avlData.Io.N1,
		"N2": avlData.Io.N2,
		"N4": avlData.Io.N4,
		"N8": avlData.Io.N8,
	} {
		fmt.Printf("    %s:\n", n)
		for id, ioElement := range ioMap {
			fmt.Printf("      ID: %d, Value: %d, Size: %d\n", id, ioElement.Value, ioElement.Size)
		}
	}
	fmt.Printf("  EventIoId: %d\n", avlData.Io.EventIoId)
}

// Print method for AvlPacket
func (packet *AvlPacket) Print() {
	fmt.Printf("CodecId: %d\n", packet.CodecId)
	fmt.Printf("NumberOfData: %d\n", packet.numberOfData)
	for i, avlData := range packet.AvlArray {
		fmt.Printf("AvlData %d:\n", i)
		avlData.Print() // Call the Print method of AvlData
	}
}
func calculateCRC16(data []byte) uint16 {
	var crc uint16 = 0x0000 // Initial CRC value
	if data == nil {
		return 0
	}
	for _, b := range data {
		crc ^= uint16(b)
		for k := 0; k < 8; k++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc = crc >> 1
			}
		}
	}
	return crc
}
func NewAvlPacketTcp(data []byte, calculateCrc bool) (*AvlPacket, error) {
	buf := bytes.NewReader(data)
	var preamble, length, crc uint32
	// Read CodecId (uint8)
	if err := binary.Read(buf, binary.BigEndian, &preamble); err != nil {
		return nil, err
	}
	if preamble != 0 {
		return nil, errors.New(fmt.Sprint("bad preamble.require ", 0, " got ", preamble))

	}

	// Read number of data (uint8)
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if uint32(len(data)) != length+12 {
		return nil, errors.New(fmt.Sprint("bad lentgh.require ", length+12, " got ", len(data)))
	}

	avlpacket := data[8 : len(data)-4]
	if calculateCrc {
		buf = bytes.NewReader(data[len(data)-4:])
		// Read number of data (uint8)
		if err := binary.Read(buf, binary.BigEndian, &crc); err != nil {
			return nil, err
		}
		temp := uint32(calculateCRC16(avlpacket))
		if crc != temp {
			return nil, errors.New(fmt.Sprint("bad crc.require ",temp, " got ", crc))
		}
	}
	return NewAvlPacket(avlpacket)
}
