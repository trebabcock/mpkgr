package format

// MetaPayloadVersion is the payload version
const MetaPayloadVersion uint16 = 1

// MetaPayload is a meta payload
type MetaPayload struct {
	Pt     Payload
	Binary []byte
}

// AddRecord adds a record to the payload
func (m *MetaPayload) AddRecord() {

}

// DefaultMetaPayload is the default meta payload
func DefaultMetaPayload() *MetaPayload {
	r := &MetaPayload{
		Pt: Payload{
			Type:           payloadTypeMeta,
			Compression:    payloadCompressionNone,
			PayloadVersion: MetaPayloadVersion,
		},
		Binary: []byte{},
	}

	return r
}
