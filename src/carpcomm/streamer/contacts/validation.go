package contacts

func IsValidFrame(satellite_id string, frame []byte) bool {
	// Everyone is currently using AX.25 packets.
	// The AX.25 header is at least 15 bytes.
	if len(frame) < 15 {
		return false
	}

	return true
}