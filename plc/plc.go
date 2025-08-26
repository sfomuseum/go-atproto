package plc

const DID_PLC string = "did:plc"
const DID_KEY string = "did:key"

const MB_ED25519 string = "\xED\x01"

type CreatePlcService struct {
	Type     string `cbor:"type" json:"type"`
	Endpoint string `cbor:"endpoint" json:"endpoint"`
}

type CreatePlcOperation struct {
	Type                string                      `cbor:"type" json:"type"`
	VerificationMethods map[string]string           `cbor:"verificationMethods" json:"verificationMethods"`
	RotationKeys        []string                    `cbor:"rotationKeys" json:"rotationKeys"`
	AlsoKnownAs         []string                    `cbor:"alsoKnownAs" json:"alsoKnownAs"`
	Services            map[string]CreatePlcService `cbor:"services" json:"services"`
	Prev                interface{}                 `cbor:"prev" json:"prev"`
}

type CreatePlcOperationSigned struct {
	CreatePlcOperation
	Signature string `cbor:"sig" json:"sig"`
}
