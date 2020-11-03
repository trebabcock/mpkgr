package format

// RecordType is the type of record for a payload
type RecordType uint8

const (
	recordTypeUnknown RecordType = 0
	recordTypeInt8    RecordType = 1
	recordTypeUint8   RecordType = 2
	recordTypeInt16   RecordType = 3
	recordTypeUint16  RecordType = 4
	recordTypeInt32   RecordType = 5
	recordTypeUint32  RecordType = 6
	recordTypeInt64   RecordType = 7
	recordTypeUint64  RecordType = 8
	recordTypeString  RecordType = 9
)

// RecordTag is a tag for the type of record
type RecordTag RecordType

const (
	recordTagUnknown RecordTag = 0
)
