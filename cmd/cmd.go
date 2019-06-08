package cmd

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arnaucube/slowlorisdb/config"
	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
	"github.com/arnaucube/slowlorisdb/node"
	"github.com/arnaucube/slowlorisdb/peer"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	{
		Name:    "create",
		Aliases: []string{},
		Usage:   "create the node",
		Action:  cmdCreate,
	},
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the node",
		Action:  cmdStart,
	},
}

func writePrivKToFile(privK *ecdsa.PrivateKey, path string) error {
	x509Encoded, err := x509.MarshalECPrivateKey(privK)
	if err != nil {
		return err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	// write privK to file
	err = ioutil.WriteFile(path, pemEncoded, 0777)
	return err
}

func readPrivKFromFile(path string) (*ecdsa.PrivateKey, error) {
	pemEncoded, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privK, err := x509.ParseECPrivateKey(x509Encoded)
	return privK, err
}

// creates the node, this needs to be executed for first time
func cmdCreate(c *cli.Context) error {
	conf, err := config.MustRead(c)
	if err != nil {
		return err
	}

	log.Info("creating new keys of the node")
	privK, err := core.NewKey()
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.StoragePath, 0777)
	if err != nil {
		return err
	}
	err = writePrivKToFile(privK, conf.StoragePath+"/privK.pem")
	if err != nil {
		return err
	}

	fmt.Println("pubK", hex.EncodeToString(core.PackPubK(&privK.PublicKey)))
	fmt.Println("addr", core.AddressFromPubK(&privK.PublicKey).String())

	return nil
}

func cmdStart(c *cli.Context) error {
	conf, err := config.MustRead(c)
	if err != nil {
		return err
	}

	db, err := db.New(conf.StoragePath + "/db")
	if err != nil {
		return err
	}

	// parse AuthNodes from the config file
	var authNodes []*ecdsa.PublicKey
	for _, authNode := range conf.AuthNodes {
		packedPubK, err := hex.DecodeString(authNode)
		if err != nil {
			return err
		}
		pubK := core.UnpackPubK(packedPubK)
		authNodes = append(authNodes, pubK)
	}

	bc := core.NewPoABlockchain(db, authNodes)

	// parse privK from path in the config file
	privK, err := readPrivKFromFile(conf.StoragePath + "/privK.pem")
	if err != nil {
		return err
	}
	fmt.Println("pubK", hex.EncodeToString(core.PackPubK(&privK.PublicKey)))
	fmt.Println("addr", core.AddressFromPubK(&privK.PublicKey).String())

	node, err := node.NewNode(privK, bc, true)
	if err != nil {
		return err
	}

	p := peer.NewPeer(node, conf)
	p.Start()

	return nil
}
