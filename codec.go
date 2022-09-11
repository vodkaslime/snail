package main

import "encoding/binary"

func encodePayload(p *[]byte) []byte {
	payloadLen := len(*p)
	totalLen := INT_LEN + int(payloadLen)
	buf := make([]byte, totalLen)

	// Assign length of the payload
	binary.LittleEndian.PutUint64(buf, uint64(payloadLen))

	// Assign the actual payload into buf
	copy(buf[INT_LEN:], *p)

	return buf
}
