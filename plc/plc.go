package plc

const DID_PLC string = "did:plc"
const DID_KEY string = "did:key"

type PlcOperationService struct {
	Type     string `cbor:"type" json:"type"`
	Endpoint string `cbor:"endpoint" json:"endpoint"`
}

type PlcOperation struct {
	Type                string                         `cbor:"type" json:"type"`
	VerificationMethods map[string]string              `cbor:"verificationMethods" json:"verificationMethods"`
	RotationKeys        []string                       `cbor:"rotationKeys" json:"rotationKeys"`
	AlsoKnownAs         []string                       `cbor:"alsoKnownAs" json:"alsoKnownAs"`
	Services            map[string]PlcOperationService `cbor:"services" json:"services"`
	Prev                interface{}                    `cbor:"prev,omitempty" json:"prev,omitempty"`
}

type PlcOperationSigned struct {
	PlcOperation
	Signature string `cbor:"sig" json:"sig"`
}
