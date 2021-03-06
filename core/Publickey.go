package core

import (
	"github.com/btcboost/secp256k1-go/secp256k1"
	"github.com/pkg/errors"
	"reflect"
)

var (
	errPublicKeySerialize = errors.New("secp256k1 public key serialize error")
)

type PublicKey secp256k1.PublicKey

func ParsePubKey(pubKeyStr []byte) (*PublicKey, error) {
	_, pubKey, err := secp256k1.EcPubkeyParse(secp256k1Context, pubKeyStr)
	return (*PublicKey)(pubKey), err
}

func (publicKey *PublicKey) ToSecp256k() *secp256k1.PublicKey {
	return (*secp256k1.PublicKey)(publicKey)
}

func (publicKey *PublicKey) SerializeUncompressed() []byte {
	_, serializedComp, err := secp256k1.EcPubkeySerialize(secp256k1Context, publicKey.ToSecp256k(), secp256k1.EcUncompressed)
	if err != nil {
		panic(errPublicKeySerialize)
	}
	return serializedComp
}

func (publicKey *PublicKey) SerializeCompressed() []byte {
	_, serializedComp, err := secp256k1.EcPubkeySerialize(secp256k1Context, publicKey.ToSecp256k(), secp256k1.EcCompressed)
	if err != nil {
		panic(errPublicKeySerialize)
	}
	return serializedComp
}

func (publicKey *PublicKey) IsEqual(otherPublicKey *PublicKey) bool {
	publicKeyBytes := publicKey.SerializeUncompressed()
	otherBytes := otherPublicKey.SerializeUncompressed()
	return reflect.DeepEqual(publicKeyBytes, otherBytes)
}

func IsCompressedOrUncompressedPubKey(bytes []byte) bool {
	if len(bytes) < 33 {
		return false
	}
	if bytes[0] == 0x04 {
		if len(bytes) != 65 {
			return false
		}
	} else if bytes[0] == 0x02 || bytes[0] == 0x03 {
		if len(bytes) != 33 {
			return false
		}
	} else {
		return false
	}
	return true

}

func IsCompressedPubKey(bytes []byte) bool {
	if len(bytes) != 33 {
		return false
	}
	if bytes[0] != 0x02 && bytes[0] != 0x03 {
		return false
	}
	return true
}

func CheckPubKeyEncoding(vchPubKey []byte, flags uint32) (bool, error) {
	if flags&SCRIPT_VERIFY_STRICTENC != 0 && !IsCompressedOrUncompressedPubKey(vchPubKey) {
		return false, ScriptErr(SCRIPT_ERR_PUBKEYTYPE)

	}
	// Only compressed keys are accepted when
	// SCRIPT_VERIFY_COMPRESSED_PUBKEYTYPE is enabled.
	if flags&SCRIPT_VERIFY_COMPRESSED_PUBKEYTYPE == 1 && !IsCompressedPubKey(vchPubKey) {
		return false, ScriptErr(SCRIPT_ERR_NONCOMPRESSED_PUBKEY)
	}
	return true, nil
}
