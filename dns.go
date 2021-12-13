package main

import (
	"encoding/binary"
	"strings"
)

// See also: https://courses.cs.duke.edu//fall16/compsci356/DNS/DNS-primer.pdf

/*
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DnsHeader struct {
	id     uint16 // A 16 bit identifier assigned by the program that generates any kind of query.
	qr     bool   // A one bit field that specifies whether this message is a query (0), or a response (1)
	opcode uint8  // A four bit field that specifies kind of query in this message.
	aa     bool   // Authoritative Answer
	tc     bool   // TrunCation - specifies that this message was truncated
	rd     bool   // Recursion Desired - this bit directs the name server to pursue the query recursively
	ra     bool   // RA Recursion Available - this be is set or cleared in a response, and denotes whether recursive query support is available in the name server. Recursive query support is optional
	rcode  uint8  // RCODE Response code - this 4 bit field is set as part of responses.
	/*
		The values have the following interpretation:
		0 No error condition
		1 Format error - The name server was unable to interpret the query.
		2
		2 Server failure - The name server was unable to process this query due to a problem with
		the name server.
		3 Name Error - Meaningful only for responses from an authoritative name server, this code
		signifies that the domain name referenced in the query does not exist.
		4 Not Implemented - The name server does not support the requested kind of query.
		5 Refused - The name server refuses to perform the specified operation for policy reasons.
	*/
	qdcount uint16 // an unsigned 16 bit integer specifying the number of entries in the question section. You should set this field to 1, indicating you have one question.
	ancount uint16 // an unsigned 16 bit integer specifying the number of resource records in the answer section. You should set this field to 0, indicating you are not providing any answers.
	nscount uint16 // an unsigned 16 bit integer specifying the number of name server resource records in the authority records section. You should set this field to 0, and should ignore any response entries in this section.
	arcount uint16 // an unsigned 16 bit integer specifying the number of resource records in the additional records section. You should set this field to 0, and should ignore any response entries in this section.
}

/*
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
/                    QNAME                      /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QTYPE                      |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QCLASS                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DnsQuestion struct {
	qname  string // A domain name represented as a sequence of labels, where each label consists of a length octet followed by that number of octets. The domain name terminates with the zero length octet for the null label of the root.
	qtype  uint16 // A two octet code which specifies the type of the query
	qclass uint16 // A two octet code that specifies the class of the query
}

/*
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
/                    NAME                       /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   ATYPE                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                   CLASS                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     TTL                       |
|                                               |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                  RDLENGTH                     | <- // The length of the RDATA field.
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
/                    RDATA                      /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type DnsAnswer struct {
	name   string // The domain name that was queried, in the same format as the QNAME in the questions.
	atype  uint16 // Two octets containing one of th type codes. This field specifies the meaning of the data in the RDATA field.
	aclass uint16 // Two octets which specify the class of the data in the RDATA field.
	ttl    uint32 // The number of seconds the results can be cached.
	rdata  []byte
	/*
		The data of the response. The format is dependent on the TYPE field: if the TYPE is 0x0001
		for A records, then this is the IP address (4 octets). If the type is 0x0005 for CNAMEs, then this
		is the name of the alias. If the type is 0x0002 for name servers, then this is the name of the
		server. Finally if the type is 0x000f for mail servers, the format is
		+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
		|                  PREFERENCE                   |
		+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
		/                   EXCHANGE                    /
		/                                               /
		+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
		where PREFERENCE is a 16 bit integer which specifies the preference of this mail server, and
		EXCHANGE is a domain name stored in the same format as QNAMEs. The latter two types are
		only relevant for the graduate version of this project.
	*/
}

/*
	DNS query
	Shown below is the hexdump (gathered via tcpdump and xxd) for an A-record query for www.northeastern.edu.

	0000000: db42 0100 0001 0000 0000 0000 0377 7777 .B...........www
	0000010: 0c6e 6f72 7468 6561 7374 6572 6e03 6564 .northeastern.ed
	0000020: 7500 0001 0001 u.....
*/
type DnsRequest struct {
	header   DnsHeader
	question DnsQuestion
}

/*
	DNS response
	Shown below is the hexdump (gathered via tcpdump and xxd) for the query above.

	0000000: db42 8180 0001 0001 0000 0000 0377 7777 .B...........www
	0000010: 0c6e 6f72 7468 6561 7374 6572 6e03 6564 .northeastern.ed
	0000020: 7500 0001 0001 c00c 0001 0001 0000 0258 u..............X
	0000030: 0004 9b21 1144 ...!.D
*/
type DnsResponse struct {
	header   DnsHeader
	question DnsQuestion
	answer   DnsAnswer
}

type DnsPacket struct {
	header    DnsHeader
	questions []DnsQuestion
	answers   []DnsAnswer
}

func DecodeHeader(headerData []byte) DnsHeader {
	return DnsHeader{
		id:      binary.BigEndian.Uint16(headerData[0:2]),
		qr:      hasBit(headerData[2], 0),
		opcode:  uint8(clearBit(headerData[2], 0) >> 3),
		aa:      hasBit(headerData[2], 5),
		tc:      hasBit(headerData[2], 6),
		rd:      hasBit(headerData[2], 7),
		ra:      hasBit(headerData[3], 0),
		rcode:   uint8(clearBit(headerData[3], 0)),
		qdcount: binary.BigEndian.Uint16(headerData[4:6]),
		ancount: binary.BigEndian.Uint16(headerData[6:8]),
		nscount: binary.BigEndian.Uint16(headerData[8:10]),
		arcount: binary.BigEndian.Uint16(headerData[10:12]),
	}
}

func EncodeHeader(header DnsHeader) []byte {
	result := make([]byte, 12)
	binary.BigEndian.PutUint16(result[0:2], header.id)
	firstFlagByte := byte(header.opcode << 3)
	firstFlagByte = setBitTo(firstFlagByte, 0, header.qr)
	firstFlagByte = setBitTo(firstFlagByte, 5, header.aa)
	firstFlagByte = setBitTo(firstFlagByte, 6, header.tc)
	firstFlagByte = setBitTo(firstFlagByte, 7, header.rd)
	result[2] = firstFlagByte
	secondFlagByte := byte(header.rcode)
	secondFlagByte = setBitTo(secondFlagByte, 0, header.ra)
	result[3] = secondFlagByte
	binary.BigEndian.PutUint16(result[4:6], header.qdcount)
	binary.BigEndian.PutUint16(result[6:8], header.ancount)
	binary.BigEndian.PutUint16(result[8:10], header.nscount)
	binary.BigEndian.PutUint16(result[10:12], header.arcount)
	return result
}

/*
  0  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
| 1  1|                                  OFFSET |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
func DecodeLengthOrPointer(data []byte) (length int8, pointer int16) {
	firstByte := data[0]
	if hasBit(firstByte, 0) && hasBit(firstByte, 1) {
		// it's a pointer...
		pointerData := make([]byte, 2)
		copy(pointerData, data[0:2])
		pointerData[0] = clearBit(pointerData[0], 0)
		pointerData[0] = clearBit(pointerData[0], 1)
		return 0, int16(binary.BigEndian.Uint16(pointerData))
	} else {
		// it's a length
		return int8(firstByte), 0
	}
}

func DecodeName(data []byte) (name string, nameLength int) {
	i := 0

	var sb strings.Builder
	for {
		partLength, pointer := DecodeLengthOrPointer(data[i : i+2])
		if pointer != 0 {
			// TODO: add support for pointers
			return "", 2
		}

		if partLength == 0 {
			return sb.String(), i + 1
		}

		lastCharIndex := i + int(partLength)
		if i != 0 {
			sb.WriteString(".")
		}
		sb.Write(data[i+1 : lastCharIndex+1])
		i = lastCharIndex + 1
	}
}

func EncodeName(name string) []byte {
	parts := strings.Split(name, ".")
	size := len(name) + 2
	var result = make([]byte, size)
	i := 0
	for _, part := range parts {
		partBytes := []byte(part)
		result[i] = byte(len(partBytes))
		copy(result[i+1:], partBytes)
		i = i + 1 + len(partBytes)
	}
	return result
}

func DecodeQuestion(questionData []byte) (question DnsQuestion, length int) {
	name, namePartLength := DecodeName(questionData)
	return DnsQuestion{
		qname:  name,
		qtype:  binary.BigEndian.Uint16(questionData[namePartLength : namePartLength+2]),
		qclass: binary.BigEndian.Uint16(questionData[namePartLength+2 : namePartLength+4]),
	}, namePartLength + 4
}

func EncodeQuestion(question DnsQuestion) []byte {
	encodedName := EncodeName(question.qname)
	everythingElse := make([]byte, 4)
	binary.BigEndian.PutUint16(everythingElse[0:2], question.qtype)
	binary.BigEndian.PutUint16(everythingElse[2:4], question.qclass)

	return append(encodedName, everythingElse...)
}

func DecodeAnswer(data []byte) (answer DnsAnswer, offset int) {
	name, nameOffset := DecodeName(data)
	rdlength := binary.BigEndian.Uint16(data[nameOffset+8 : nameOffset+10])
	return DnsAnswer{
		name:   name,
		atype:  binary.BigEndian.Uint16(data[nameOffset : nameOffset+2]),
		aclass: binary.BigEndian.Uint16(data[nameOffset+2 : nameOffset+4]),
		ttl:    binary.BigEndian.Uint32(data[nameOffset+4 : nameOffset+8]),
		rdata:  data[nameOffset+10 : nameOffset+10+int(rdlength)],
	}, nameOffset + 10 + int(rdlength)
}

func EncodeAnswer(answer DnsAnswer) []byte {
	encodedName := EncodeName(answer.name)

	data := make([]byte, 8+len(answer.rdata))
	binary.BigEndian.PutUint16(data[0:2], answer.atype)
	binary.BigEndian.PutUint16(data[2:4], answer.aclass)
	binary.BigEndian.PutUint32(data[4:8], answer.ttl)
	binary.BigEndian.PutUint16(data[8:10], uint16(len(answer.rdata)))
	copy(data[10:], answer.rdata)

	return append(encodedName, data...)
}

func DecodeRequest(packet []byte) DnsRequest {
	headerData := packet[0:12]
	questionData := packet[12:]
	question, _ := DecodeQuestion(questionData)
	return DnsRequest{
		header:   DecodeHeader(headerData),
		question: question,
	}
}

func EncodeRequest(request DnsRequest) []byte {
	return append(EncodeHeader(request.header), EncodeQuestion(request.question)...)
}

func DecodeResponse(packet []byte) DnsResponse {
	header := DecodeHeader(packet[0:12])
	question, questionLength := DecodeQuestion(packet[12:])
	answer, _ := DecodeAnswer(packet[12+questionLength:])

	return DnsResponse{
		header:   header,
		question: question,
		answer:   answer,
	}
}

func EncodeResponse(response DnsResponse) []byte {
	headerData := EncodeHeader(response.header)
	questionData := EncodeQuestion(response.question)
	answerData := EncodeAnswer(response.answer)
	return append(append(headerData, questionData...), answerData...)
}

func DecodePacket(packet []byte) DnsPacket {
	header := DecodeHeader(packet[0:12])
	questionsCount := header.qdcount
	answersCount := header.ancount

	questions := make([]DnsQuestion, questionsCount)
	answers := make([]DnsAnswer, answersCount)

	offset := 12
	for q := 0; q < int(questionsCount); q++ {
		question, questionLength := DecodeQuestion(packet[offset:])
		questions[q] = question
		offset += questionLength
	}
	for a := 0; a < int(answersCount); a++ {
		answer, answerLength := DecodeAnswer(packet[offset:])
		answers[a] = answer
		offset += answerLength
	}

	return DnsPacket{
		header:    header,
		questions: questions,
		answers:   answers,
	}
}

func EncodePacket(packet DnsPacket) []byte {
	headerData := EncodeHeader(packet.header)
	data := headerData

	for _, question := range packet.questions {
		data = append(data, EncodeQuestion(question)...)
	}
	for _, answer := range packet.answers {
		data = append(data, EncodeAnswer(answer)...)
	}

	return data
}
