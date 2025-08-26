package plc

const DID_PLC string = "did:plc"
const DID_KEY string = "did:key"

type GenesisOperationService struct {
	Type     string `cbor:"type" json:"type"`
	Endpoint string `cbor:"endpoint" json:"endpoint"`
}

type GenesisOperation struct {
	Type                string                             `cbor:"type" json:"type"`
	VerificationMethods map[string]string                  `cbor:"verificationMethods" json:"verificationMethods"`
	RotationKeys        []string                           `cbor:"rotationKeys" json:"rotationKeys"`
	AlsoKnownAs         []string                           `cbor:"alsoKnownAs" json:"alsoKnownAs"`
	Services            map[string]GenesisOperationService `cbor:"services" json:"services"`
	Prev                interface{}                        `cbor:"prev" json:"prev"`
}

type GenesisOperationSigned struct {
	GenesisOperation
	Signature string `cbor:"sig" json:"sig"`
}

type CreateOperation struct {
	Type        string `json:"type"`
	SigningKey  string `json:"signingKey"`
	RecoveryKey string `json:"recoveryKey"`
	Handle      string `json:"handle"`
	Service     string `json:"service"`
	Prev        any    `json:"prev"`
}

type CreateOperationSigned struct {
	CreateOperation
	Signature string `json:"sig"`
}
