# go-pds

Experimental Go package implementing an ATProto Personal Data Server (pds).

## Motivation

This is experimental work to implement a "minimal viable" ATProto Personal Data Server (pds) such that SFO Museum can create accounts for things in its collection and participate on/with the ATProto network (read: Bluesky).

## Important

This work is incomplete and, in some notable cases, does not work. Specifically attempts to create, or register, new `did:plc` identifiers with the `https://plc.directory` fails with opaque "Bad Request" errors. For example:

```
$> go run cmd/pds-create-user/main.go -handle test939834748 -service https://atproto.sfomuseum.org

// START OF debugging output
{
  "type": "plc_operation",
  "verificationMethods": {
    "atproto": "did:key:zDnaevTSNUV3kMM5jn3i1rczD65J2qGio8vtPML4eWTedt6kH"
  },
  "rotationKeys": [
    "did:key:zDnaevTSNUV3kMM5jn3i1rczD65J2qGio8vtPML4eWTedt6kH"
  ],
  "alsoKnownAs": [
    "at://test939834748"
  ],
  "services": {
    "atproto_pds": {
      "type": "AtprotoPersonalDataServer",
      "endpoint": "https://atproto.sfomuseum.org"
    }
  },
  "sig": "4u0chguS8hKKywl__1_8pG3Jqy3ErmkWirUHrXXfLB3MV1QPXHr82ApUQK3pfSr8Yrs9S-w07HsJp4qohcX-eQ"
}
POST /did:plc:5NVS7CU7SB5FAUWH5GB7V5RZ HTTP/1.1
Host: plc.directory
User-Agent: Go-http-client/1.1
Content-Length: 434
Content-Type: application/json

{"type":"plc_operation","verificationMethods":{"atproto":"did:key:zDnaevTSNUV3kMM5jn3i1rczD65J2qGio8vtPML4eWTedt6kH"},"rotationKeys":["did:key:zDnaevTSNUV3kMM5jn3i1rczD65J2qGio8vtPML4eWTedt6kH"],"alsoKnownAs":["at://test939834748"],"services":{"atproto_pds":{"type":"AtprotoPersonalDataServer","endpoint":"https://atproto.sfomuseum.org"}},"sig":"4u0chguS8hKKywl__1_8pG3Jqy3ErmkWirUHrXXfLB3MV1QPXHr82ApUQK3pfSr8Yrs9S-w07HsJp4qohcX-eQ"}

// END OF debugging output

2025/08/26 18:04:04 Failed to run create user, Failed to create PLC for DID, Failed to execute request, Post "https://plc.directory/did:plc:5NVS7CU7SB5FAUWH5GB7V5RZ": dial tcp 3.129.34.168:443: connect: bad file descriptor
```

It is not clear to me what the problem is yet.
