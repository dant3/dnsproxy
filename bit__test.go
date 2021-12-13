package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasBit(t *testing.T) {
	meaningfulByte := byte(0b00100001)

	assert.Equal(t, false, hasBit(meaningfulByte, 0))
	assert.Equal(t, false, hasBit(meaningfulByte, 1))
	assert.Equal(t, true, hasBit(meaningfulByte, 2))
	assert.Equal(t, false, hasBit(meaningfulByte, 3))
	assert.Equal(t, false, hasBit(meaningfulByte, 4))
	assert.Equal(t, false, hasBit(meaningfulByte, 5))
	assert.Equal(t, false, hasBit(meaningfulByte, 6))
	assert.Equal(t, true, hasBit(meaningfulByte, 7))
}

func TestHasBit0xc(t *testing.T) {
	byte := byte(0xc0)
	assert.True(t, hasBit(byte, 0))
	assert.True(t, hasBit(byte, 1))
}

func TestSetBit(t *testing.T) {
	meaningfulByte := setBit(byte(0b00000001), 3)

	assert.Equal(t, false, hasBit(meaningfulByte, 0))
	assert.Equal(t, false, hasBit(meaningfulByte, 1))
	assert.Equal(t, false, hasBit(meaningfulByte, 2))
	assert.Equal(t, true, hasBit(meaningfulByte, 3))
	assert.Equal(t, false, hasBit(meaningfulByte, 4))
	assert.Equal(t, false, hasBit(meaningfulByte, 5))
	assert.Equal(t, false, hasBit(meaningfulByte, 6))
	assert.Equal(t, true, hasBit(meaningfulByte, 7))
}

func TestClearBit(t *testing.T) {
	meaningfulByte := clearBit(byte(0b00010101), 7)

	assert.Equal(t, false, hasBit(meaningfulByte, 0))
	assert.Equal(t, false, hasBit(meaningfulByte, 1))
	assert.Equal(t, false, hasBit(meaningfulByte, 2))
	assert.Equal(t, true, hasBit(meaningfulByte, 3))
	assert.Equal(t, false, hasBit(meaningfulByte, 4))
	assert.Equal(t, true, hasBit(meaningfulByte, 5))
	assert.Equal(t, false, hasBit(meaningfulByte, 6))
	assert.Equal(t, false, hasBit(meaningfulByte, 7))
}
