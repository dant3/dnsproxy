package main

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func StripWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func BinaryString(str string) []byte {
	data, err := hex.DecodeString(StripWhitespace(str))
	if err != nil {
		log.Fatal("Failure: ", err)
	}
	return data
}

func TestDecodeRequest(t *testing.T) {
	data := BinaryString(`
		db42 0100 0001 0000 0000 0000 0377 7777
		0c6e 6f72 7468 6561 7374 6572 6e03 6564
		7500 0001 0001
	`)
	dnsRequest := DecodeRequest(data)

	assert.Equal(t, uint16(0xdb42), dnsRequest.header.id)
	assert.Equal(t, false, dnsRequest.header.qr)
	assert.Equal(t, uint8(0), dnsRequest.header.opcode)
	assert.Equal(t, false, dnsRequest.header.aa)
	assert.Equal(t, false, dnsRequest.header.tc)
	assert.Equal(t, true, dnsRequest.header.rd)
	assert.Equal(t, false, dnsRequest.header.ra)
	assert.Equal(t, uint8(0), dnsRequest.header.rcode)
	assert.Equal(t, uint16(1), dnsRequest.header.qdcount)
	assert.Equal(t, uint16(0), dnsRequest.header.ancount)
	assert.Equal(t, uint16(0), dnsRequest.header.nscount)
	assert.Equal(t, uint16(0), dnsRequest.header.arcount)

	assert.Equal(t, "www.northeastern.edu", dnsRequest.question.qname)
	assert.Equal(t, uint16(0x0001), dnsRequest.question.qtype)
	assert.Equal(t, uint16(0x0001), dnsRequest.question.qclass)
}

func TestEncodeRequest(t *testing.T) {
	payload := "db42 0100 0001 0000 0000 0000 0377 7777" +
		"0c6e 6f72 7468 6561 7374 6572 6e03 6564" +
		"7500 0001 0001"
	data, _ := hex.DecodeString(strings.ReplaceAll(payload, " ", ""))
	dnsRequest := DecodeRequest(data)

	encodeResult := EncodeRequest(dnsRequest)

	assert.Equal(t, data, encodeResult)
}

func TestDecodeResponse(t *testing.T) {
	data := BinaryString(`db42 8180 0001 0001 0000 0000 0377 7777
		0c6e 6f72 7468 6561 7374 6572 6e03 6564
		7500 0001 0001 c00c 0001 0001 0000 0258
		0004 9b21 1144
	`)

	dnsResponse := DecodeResponse(data)

	assert.Equal(t, uint16(0xdb42), dnsResponse.header.id)
	assert.Equal(t, true, dnsResponse.header.qr)
	assert.Equal(t, uint8(0), dnsResponse.header.opcode)
	assert.Equal(t, false, dnsResponse.header.aa)
	assert.Equal(t, false, dnsResponse.header.tc)
	assert.Equal(t, true, dnsResponse.header.rd)
	assert.Equal(t, true, dnsResponse.header.ra)
	assert.Equal(t, uint8(0), dnsResponse.header.rcode)
	assert.Equal(t, uint16(1), dnsResponse.header.qdcount)
	assert.Equal(t, uint16(1), dnsResponse.header.ancount)
	assert.Equal(t, uint16(0), dnsResponse.header.nscount)
	assert.Equal(t, uint16(0), dnsResponse.header.arcount)

	assert.Equal(t, "www.northeastern.edu", dnsResponse.question.qname)
	assert.Equal(t, uint16(0x0001), dnsResponse.question.qtype)
	assert.Equal(t, uint16(0x0001), dnsResponse.question.qclass)

	// TODO: to read name from answer we should support DNS compression schema
	assert.Equal(t, "", dnsResponse.answer.name)
	assert.Equal(t, uint16(0x0001), dnsResponse.answer.atype)
	assert.Equal(t, uint16(0x0001), dnsResponse.answer.aclass)
	assert.Equal(t, uint32(600), dnsResponse.answer.ttl)
	assert.Equal(t, BinaryString("9b21 1144"), dnsResponse.answer.rdata)
}

func TestEncodeResponse(t *testing.T) {
	data := BinaryString(`db42 8180 0001 0001 0000 0000 0377 7777
		0c6e 6f72 7468 6561 7374 6572 6e03 6564
		7500 0001 0001 c00c 0001 0001 0000 0258
		0004 9b21 1144
	`)

	dnsResponse := DecodeResponse(data)
	encodedResponse := EncodeResponse(dnsResponse)
	assert.NotEmpty(t, encodedResponse)
}
