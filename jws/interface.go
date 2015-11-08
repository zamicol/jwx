package jws

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

var (
	ErrInvalidCompactPartsCount = errors.New("compact JWS format must have three parts")
	ErrInvalidMac               = errors.New("invalid mac")
	ErrUnsupportedAlgorithm     = errors.New("unspported algorithm")
)

// Base64Encoder can encode itself into base64. But you can do more such as
// filling default values, validating them, etc. This is used in `Encode()`
// as both headers and payloads
type Base64Encoder interface {
	Base64Encode() ([]byte, error)
}

type Base64Decoder interface {
	Base64Decode([]byte) error
}

type EssentialHeader struct {
	Algorithm              jwa.SignatureAlgorithm `json:"alg,omitempty"`
	ContentType            string                 `json:"cty,omitempty"`
	Critical               []string               `json:"crit,omitempty"`
	Jwk                    jwk.JSONWebKey         `json:"jwk,omitempty"` // public key
	JwkSetURL              *url.URL               `json:"jku,omitempty"`
	KeyID                  string                 `json:"kid,omitempty"`
	Type                   string                 `json:"typ,omitempty"` // e.g. "JWT"
	X509Url                *url.URL               `json:"x5u,omitempty"`
	X509CertChain          []string               `json:"x5c,omitempty"`
	X509CertThumbprint     string                 `json:"x5t,omitempty"`
	X509CertThumbprintS256 string                 `json:"x5t#S256,omitempty"`
}

// Header represents a jws header.
type Header struct {
	*EssentialHeader `json:"-"`
	PrivateParams    map[string]interface{} `json:"-"`
}

// EncodedHeader represents a header value that is base64 encoded
// in JSON format
type EncodedHeader struct {
	Header
	encoded buffer.Buffer // sometimes our encoding and the source encoding don't match
}

// Signer generates signature for the given payload
type Signer interface {
	Jwk() jwk.JSONWebKey
	Kid() string
	Alg() jwa.SignatureAlgorithm
	Sign([]byte) ([]byte, error)
}

// Verifier is used to verify the signature against the payload
type Verifier interface {
	Verify([]byte, []byte) error
}

type RsaSign struct {
	Algorithm  jwa.SignatureAlgorithm
	JSONWebKey *jwk.RsaPublicKey
	KeyID      string
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type EcdsaSign struct {
	Algorithm  jwa.SignatureAlgorithm
	JSONWebKey *jwk.RsaPublicKey
	KeyID      string
	PrivateKey *ecdsa.PrivateKey
}

type MergedHeader struct {
	ProtectedHeader *EncodedHeader
	PublicHeader    *Header
}

type Signature struct {
	PublicHeader    Header        `json:"header"`              // Raw JWS Unprotected Heders
	ProtectedHeader EncodedHeader `json:"protected,omitempty"` // Base64 encoded JWS Protected Headers
	Signature       buffer.Buffer `json:"signature"`           // Base64 encoded signature
}

// Message represents a full JWS encoded message. Flattened serialization
// is not supported as a struct, but rather it's represented as a
// Message struct with only one `signature` element
type Message struct {
	Payload    buffer.Buffer `json:"payload"`
	Signatures []Signature   `json:"signatures"`
}

type MultiSigner interface {
	MultiSign(buffer.Buffer) (*Message, error)
	AddSigner(Signer)
}

type MultiSign struct {
	Signers []Signer
}

type HmacSign struct {
	Algorithm jwa.SignatureAlgorithm
	JSONWebKey *jwk.RsaPublicKey
	KeyID      string
	Key       []byte
}
