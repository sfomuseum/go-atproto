package plc

/*
func (d *DID) PublicKey(suffix string) ([]byte, error) {

	for _, m := range d.VerificationMethods {

		if !strings.HasSuffix(m.Id, suffix) {
			continue
		}

		key_mb := m.PublicKeyMultibase
		key_mb = strings.TrimLeft(key_mb, fmt.Sprintf("%s:", DID_KEY))

		_, body, err := multibase.Decode(key_mb)

		if err != nil {
			return nil, fmt.Errorf("Failed to decode multibase, %w", err)
		}

		return body[2:], nil
	}

	return nil, fmt.Errorf("Key not found")
}
*/
