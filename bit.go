package main

func hasBit(data byte, pos uint) bool {
	val := data & (0b10000000 >> pos)
	return (val > 0)
}

func setBit(data byte, pos uint) byte {
	data |= byte(0b10000000 >> pos)
	return data
}

func setBitTo(data byte, pos uint, value bool) byte {
	if value {
		return setBit(data, pos)
	} else {
		return clearBit(data, pos)
	}
}

func clearBit(data byte, pos uint) byte {
	mask := byte(^(0b10000000 >> pos))
	data &= mask
	return data
}
