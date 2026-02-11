package umqr

// CalculateCRC16CCITT calculates the CRC16-CCITT (False) checksum.
// Polynomial: 0x1021
// Initial Value: 0xFFFF
// Final XOR: 0x0000
// No reflection on input or output.
func CalculateCRC16CCITT(data []byte) uint16 {
	crc := uint16(0xFFFF)
	poly := uint16(0x1021)

	for _, b := range data {
		// XOR byte into the top byte of the CRC
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ poly
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}
