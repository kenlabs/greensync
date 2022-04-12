package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func Init(out io.Writer) (*Config, error) {
	identity, err := CreateIdentity(out)
	if err != nil {
		return nil, err
	}
	return InitWithIdentity(identity)
}

func InitWithIdentity(identity Identity) (*Config, error) {
	return &Config{
		Identity:  identity,
		Datastore: NewDatastore(),
	}, nil
}

// CreateIdentity initializes a new identity.
func CreateIdentity(out io.Writer) (Identity, error) {
	ident := Identity{}

	var sk crypto.PrivKey
	var pk crypto.PubKey

	fmt.Fprint(out, "generating ED25519 keypair...")
	priv, pub, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return ident, err
	}
	fmt.Fprintln(out, "done")

	sk = priv
	pk = pub

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := crypto.MarshalPrivateKey(sk)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	fmt.Fprintf(out, "peer identity: %s\n", ident.PeerID)
	return ident, nil
}