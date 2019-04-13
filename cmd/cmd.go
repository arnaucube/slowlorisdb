package cmd

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/arnaucube/slowlorisdb/config"
	"github.com/arnaucube/slowlorisdb/core"
	"github.com/arnaucube/slowlorisdb/db"
	"github.com/arnaucube/slowlorisdb/node"
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

func cmdCreate(c *cli.Context) error {
	log.Info("creating new keys of the node")
	privK, err := core.NewKey()
	if err != nil {
		return err
	}
	fmt.Println(privK)
	return nil
}

func cmdStart(c *cli.Context) error {
	conf, err := config.MustRead(c)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", conf.DbPath)
	if err != nil {
		return err
	}
	db, err := db.New(dir)
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

	// TODO parse privK from path in the config file
	privK, err := core.NewKey()
	node, err := node.NewNode(privK, bc, true)
	if err != nil {
		return err
	}
	err = node.Start()
	if err != nil {
		return err
	}

	return nil
}
