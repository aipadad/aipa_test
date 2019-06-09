﻿/*
@Time : 2018/8/13 10:15 
@Author : 
@File : main
@Software: GoLand
*/
package main

import (
	cli "gopkg.in/urfave/cli.v1"
	"encoding/json"
	"os"
	"log"
	"fmt"
	"regexp"
	"strings"
	"math/big"
	//"encoding/hex"
	//pack "github.com/aipadad/msgpack-go"
	//user_proto "github.com/aipadad/magiccube/service/user/proto"
	//push_sign "github.com/aipadad/magiccube/service/common/signature/push"
	//"github.com/protobuf/proto"
	//"github.com/aipadad/magiccube/service/common/bean"
	//"bytes"
	//"github.com/aipadad/crypto-go/crypto"
	//"github.com/aipadad/magiccube/config"
	//"github.com/aipadad/magiccube/service/common/util"
	"github.com/micro/go-micro"
	chain "github.com/aipadad/aipa/api"
)

var ChainAddr string = "127.0.0.1:8689"
var ChainId   string = ""

// CLI responsible for processing command line arguments
type CLI struct {
	client chain.ChainClient
}

//NewCLI new console client
func NewCLI() *CLI {
	cli := &CLI{}
	service := micro.NewService()
	//avoid parameters conflict with those of bcli
	//service.Init()
	cli.client = chain.NewChainClient("aipa", service.Client())

	return cli
}

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		if ctx.GlobalIsSet("servaddr") {
			ChainAddr = ctx.GlobalString("servaddr")
		}
		if ctx.GlobalIsSet("chainid") {
			ChainId = ctx.GlobalString("chainid")
		}
		return action(ctx)
	}
}

func (cli *CLI) BcliGetChainInfo(ctx *cli.Context) error {
	chainInfo, err := cli.GetChainInfoOverHttp("http://"+ChainAddr+"/v1/block/height")
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return nil
	}
	fmt.Printf("\n==Chain Info==\n\n")
	
	b, _ := json.Marshal(chainInfo)
	cli.jsonPrint(b)
	
	return nil
}

func (cli *CLI) BcliGetBlockInfo(ctx *cli.Context) error {

	num := ctx.Uint64("num")
	hash := ctx.String("hash")

	blockInfo, err := cli.getBlockInfoOverHttp("http://"+ChainAddr+"/v1/block/detail", num, hash)
	if err != nil {
		return nil
	}
	fmt.Printf("\n==Block Info==\n\n")
	b, _ := json.Marshal(blockInfo)
	cli.jsonPrint(b)
	return nil
}

func (cli *CLI) BcliNewAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	pubkey := ctx.String("pubkey")
	
	cli.newaccount(username, pubkey)
	
	return nil
}

func (cli *CLI) BcliGetAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	
	cli.getaccount(username)
	
	return nil
}

func (cli *CLI) BcliTransfer(ctx *cli.Context) error {

	from := ctx.String("from")
	to   := ctx.String("to")
	amounttmp := ctx.String("amount")
	
	amount := big.NewInt(0)

	amount, _ = amount.SetString(amounttmp, 10)

	cli.transfer(from, to, *amount)
	
	return nil
}

func (cli *CLI) BCLIGetTransaction(ctx *cli.Context) error {

	trxhash := ctx.String("trxhash")

	cli.BcliGetTransaction(trxhash)
	
	return nil
}

func (cli *CLI) BCLIPushTransaction(ctx *cli.Context) error {
	
	var pushtrx BcliPushTrxInfo
        
        pushtrx.sender   = ctx.String("sender")
        pushtrx.contract = ctx.String("contract")
        pushtrx.method   = ctx.String("method")
	pushtrx.ParamMap = make(map[string]interface{})

	param1 := ctx.String("param")
	param1 = strings.Replace(param1, " ", "", -1)
	param2 := strings.Split(param1, ",")
	for _, item := range(param2) {
		param3 := strings.Split(item, ":")
		pushtrx.ParamMap[param3[0]] = param3[1]
	}

	cli.BcliPushTransaction(&pushtrx)
	
	return nil
}

func (cli *CLI) BCLIDeployCode(ctx *cli.Context) error {
	name := ctx.String("name")
	code := ctx.String("code")
	fileType := ctx.String("filetype")

	cli.deploycode(name, code, fileType)

	return nil
}

func (cli *CLI) BCLIDeployAbi(ctx *cli.Context) error {
	name := ctx.String("name")
	Abi := ctx.String("abi")

	cli.deployabi(name, Abi)

	return nil
}

func (cli *CLI) BCLIDeployBoth(ctx *cli.Context) error {
	name := ctx.String("name")
	Abi := ctx.String("abi")
	code := ctx.String("code")
	fileType := ctx.String("filetype")
	user := ctx.String("user")

	if len(name) <= 0 {
		fmt.Println("Error: Please input the contract name.")
		return nil
	}

	if len(Abi) <= 0 {
		fmt.Println("Error: Please input the contract abi file's path/name.")
		return nil
	}

	if len(code) <= 0 {
		fmt.Println("Error: Please input the contract code file's path/name.")
		return nil
	}

	if len(fileType) <= 0 {
		fmt.Println("Error: Please input file type: wasm or js.")
		return nil
	}

	if len(user) <= 0 {
		fmt.Println("Error: Please input the user name whom deploy the abi and contract.")
		return nil
	}

	cli.deploycode(name, code, user, fileType)

	cli.deployabi(name, Abi, user)

	return nil
}

func (cli *CLI) BCLIGetContractCode(ctx *cli.Context) error {
	name := ctx.String("name")
	SaveToAbiPath  := ctx.String("abi")
	SaveTocodePath := ctx.String("code")
	
	cli.BcliGetContractCode(name, SaveTocodePath, SaveToAbiPath)

	return nil
}

func (cli *CLI) BCLIGetTableInfo(ctx *cli.Context) error {
	contract := ctx.String("contract")
	table := ctx.String("table")
	key  := ctx.String("key")
	
	cli.BCliGetTableInfo(contract, table, key)

	return nil
}

func (cli *CLI) BcliAccountStake(ctx *cli.Context) error {
	account := ctx.String("account")
	amounttmp  := ctx.String("amount")
	
	amount := big.NewInt(0)

	amount, _ = amount.SetString(amounttmp, 10)
	
	cli.BCliAccountStakeInfo(account, *amount)

	return nil
}

func (cli *CLI) BcliAccountUnStake(ctx *cli.Context) error {
	account := ctx.String("account")
	amounttmp  := ctx.String("amount")
	
	amount := big.NewInt(0)

	amount, _ = amount.SetString(amounttmp, 10)
	
	cli.BCliAccountUnStakeInfo(account, *amount)

	return nil
}

func (cli *CLI) BcliAccountClaim(ctx *cli.Context) error {
	account := ctx.String("account")
	amounttmp  := ctx.String("amount")
	
	amount := big.NewInt(0)

	amount, _ = amount.SetString(amounttmp, 10)
	
	cli.BCliAccountClaimInfo(account, *amount)

	return nil
}

func (cli *CLI) BcliVote(ctx *cli.Context) error {
	vouter := ctx.String("vouter")
	delegate  := ctx.String("delegate")
	
	cli.BCliVoteInfo(vouter, delegate)

	return nil
}

func (cli *CLI) BcliCancelVote(ctx *cli.Context) error {
	vouter := ctx.String("vouter")
	delegate  := ctx.String("delegate")
	
	cli.BCliCancelVoteInfo(vouter, delegate)

	return nil
}

func (cli *CLI) BcliDelegateReg(ctx *cli.Context) error {
	account := ctx.String("account")
	signkey  := ctx.String("signkey")
	location := ctx.String("location")
	description := ctx.String("description")

	cli.BCliDelegateRegInfo(account, signkey, location, description)

	return nil
}

func (cli *CLI) BcliDelegateUnReg(ctx *cli.Context) error {
	account := ctx.String("account")
	
	cli.BCliDelegateUnRegInfo(account)

	return nil
}

func (Cli *CLI) RunNewCLI() {
	app := cli.NewApp()
	app.Name = "aipa Cmd"
	app.Usage = "block chain bcli"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag {
			cli.StringFlag{
				Name:  "servaddr",
				Value: "127.0.0.1:8689",
			},
			cli.StringFlag{
				Name:  "chainid",
				Value: "000000000000000000",
			},
	}

	app.Commands = []cli.Command {
		{
			Name: "getinfo",
			Usage: "Get chian info",
			Category: "general",
			Action: MigrateFlags(Cli.BcliGetChainInfo),
		},
		{
			Name: "getblock",
			Usage: "Get block info",
			Category: "general",
			Flags: []cli.Flag {
				cli.Uint64Flag{
					Name: "num",
					Value: 100,
					Usage: "get block by number",
				},
				cli.StringFlag{
					Name: "hash",
					Value: "",
					Usage: "get block by hash",
				},
			},
			Action: MigrateFlags(Cli.BcliGetBlockInfo),
		},
		{
			Name: "gettable",
			Usage: "get table info",
			Category: "general",
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "contract",
					Value:"usermng",
					Usage: "contract name",
				},
				cli.StringFlag{
					Name: "table",
					Usage: "table name",
				},
				cli.StringFlag{
					Name: "key",
					Usage: "key value",
				},
			},
			Action: MigrateFlags(Cli.BCLIGetTableInfo),
		},
		{
			Name: "account",
			Usage: "Create or Get account",
			Category: "account",
			Subcommands: []cli.Command {
				{
					Name: "create",
					Usage: "Create account",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "username",
							Value:"",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "pubkey",
							Value:"",
							Usage: "account public key",
						},
					},
					Action: MigrateFlags(Cli.BcliNewAccount),
				},
				{
					Name: "get",
					Usage: "Getter account info",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "username",
							Value:"",
							Usage: "acocunt name",
						},
					},
					Action: MigrateFlags(Cli.BcliGetAccount),
				},
				{
					Name: "stake",
					Usage: "stake of account",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "account",
							Value:"",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "ammount",
							Value:"",
							Usage: "amount",
						},
					},
					Action: MigrateFlags(Cli.BcliAccountStake),
				},
				{
					Name: "unstake",
					Usage: "unstake of account",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "account",
							Value:"",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "ammount",
							Value:"",
							Usage: "amount",
						},
					},
					Action: MigrateFlags(Cli.BcliAccountUnStake),
				},
				{
					Name: "claim",
					Usage: "claim of stake",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "account",
							Value:"",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "ammount",
							Value:"",
							Usage: "amount",
						},
					},
					Action: MigrateFlags(Cli.BcliAccountClaim),
				},
			},
		},
		{
			Name: "transfer",
			Usage: "transfer",
			Category: "transfer",
			Flags:[]cli.Flag {
				cli.StringFlag{
					Name: "from",
					Usage: "",
					Value: "aipa",
				},
				cli.StringFlag{
					Name: "to",
					Usage: "",
					Value: "aipa",
				},
				cli.StringFlag{
					Name: "amount",
					Usage: "",
					Value: "0",
				},
				cli.StringFlag{
					Name: "sign",
					Usage: "",
				},
			},
			Action: MigrateFlags(Cli.BcliTransfer),
		},
		{
			Name: "transaction",
			Usage: "transaction lists",
			Category: "transaction",
			Subcommands: []cli.Command {
				{
					Name: "get",
					Usage: "Getter tx details",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "trxhash",
						},
					},
					Action: MigrateFlags(Cli.BCLIGetTransaction),
				},
				{
					Name: "push",
					Usage: "push transaction",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "sender",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "contract",
							Usage: "contract name",
						},
						cli.StringFlag{
							Name: "method",
							Usage: "method name",
						},
						cli.StringFlag{
							Name: "param",
							Usage: "param value",
						},
						cli.StringFlag{
							Name: "sign",
							Usage: "sign value",
						},
					},
					Action: MigrateFlags(Cli.BCLIPushTransaction),
				},
			},
		},
		{
			Name: "contract",
			Usage: "contract info",
			Category: "contract",
			Subcommands: []cli.Command {
				{
					Name:  "deploy",
					Usage: "contract deploy",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "the contract name",
						},
						cli.StringFlag{
							Name:  "code",
							Usage: "the contract's wasm file path ( includes wasm file name )",
						},
						cli.StringFlag{
							Name:  "filetype",
							Value: "wasm",
							Usage: "the contract's file type: wasm or js",
						},
						cli.StringFlag{
							Name:  "abi",
							Usage: "the contract's abi file path ( includes abi file name )",
						},
						cli.StringFlag{
							Name:  "user",
							Usage: "the user account",
						},
					},
					Action: MigrateFlags(Cli.BCLIDeployBoth),
				},
				{
					Name:  "deploycode",
					Usage: "contract  deploycode",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "the contract name",
						},
						cli.StringFlag{
							Name:  "code",
							Usage: "the contract's wasm file path ( includes wasm file name )",
						},
						cli.StringFlag{
							Name:  "filetype",
							Value: "wasm",
							Usage: "the contract's file type: wasm or js",
						},
						cli.StringFlag{
							Name:  "user",
							Usage: "the user account",
						},
					},
					Action: MigrateFlags(Cli.BCLIDeployCode),
				},
				{
					Name: "deployabi",
					Usage: "contract  deployabi",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "abi",
							Usage:"",
						},
						cli.StringFlag{
							Name: "sign",
							Usage:"",
						},
					},
					Action: MigrateFlags(Cli.BCLIDeployAbi),
				},
				{
					Name: "get",
					Usage: "Getter contract",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "code",
							Usage:"",
						},
						cli.StringFlag{
							Name: "abi",
							Usage:"",
						},
					},
					Action: MigrateFlags(Cli.BCLIGetContractCode),
				},
			},
		},
		{
			Name:     "p2p",
			Category: "p2p",
			Subcommands: []cli.Command{
				{
					Name:  "connect",
					Usage: "connect address or port",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "peer",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("peer"))
						return nil
					},
				},
				{
					Name:  "disconnect",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "peer",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("peer"))
						return nil
					},
				},
				{
					Name:  "status",
					Usage: "p2p status",
					Action: func(c *cli.Context) error {
						// TODO

						return nil
					},
				},
				{
					Name:  "peers",
					Usage: "peers info",
					Action: func(c *cli.Context) error {
						// TODO

						return nil
					},
				},
			},
		},
		{
			Name: "delegate",
			Category: "delegate",
			Subcommands: []cli.Command{
				{
					Name:  "reg",
					Usage: "connect address or port",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Usage:"account name",
						},
						cli.StringFlag{
							Name: "signkey",
							Usage:"sign key",
						},
						cli.StringFlag{
							Name: "location",
						},
						cli.StringFlag{
							Name: "description",
						},
					},
					Action: MigrateFlags(Cli.BcliDelegateReg),
				},
				{
					Name:  "unreg",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
						},
					},
					Action: MigrateFlags(Cli.BcliDelegateUnReg),
				},
				{
					Name:  "list",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name: "limit",
							Value:100,
						},
						cli.Int64Flag{
							Name: "start",
							Value:0,
						},

					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("limit"))
						fmt.Println(c.String("start"))
						return nil
					},
				},
				{
					Name: "vote",
					Usage: "vote for producers",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "voter",
							Value:"",
							Usage: "vouter",
						},
						cli.StringFlag{
							Name: "delegate",
							Value:"",
							Usage: "delegate",
						},
					},
					Action: MigrateFlags(Cli.BcliVote),
				},
				{
					Name: "cancelvote",
					Usage: "cancel vote for producers",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "voter",
							Value:"",
							Usage: "vouter",
						},
						cli.StringFlag{
							Name: "delegate",
							Value:"",
							Usage: "delegate",
						},
					},
					Action: MigrateFlags(Cli.BcliCancelVote),
				},
			},
		},
		{
			Name: "wallet",
			Category: "wallet",
			Subcommands: []cli.Command{
				{
					Name:  "generatekey",
					Usage: "generate key",
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},
				{
					Name:  "create",
					Usage: "create wallet",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Value:"",
						},
						cli.StringFlag{
							Name: "privatekey",
							Value:"",
						},

					},
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},
				{
					Name:  "lock",
					Usage: "lock wallet",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Value:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},
				{
					Name:  "unlock",
					Usage: "unlock wallet",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Value:"",
						},
						cli.StringFlag{
							Name: "password",
							Value:"",
						},
						cli.StringFlag{
							Name: "timeout",
							Value:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list wallet",
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},
				{
					Name:  "listkey",
					Usage: "listkey of wallet",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Value:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						return nil
					},
				},

			},
		},
	}


	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}

func isNotEmpty(str string) bool {
	if len(str) > 0 {
		return true
	}
	return false
}

func validatorUsername(str string) (bool,error) {
	if !isNotEmpty(str) {
		return false,  fmt.Errorf("Parameter anomaly！")
	}

	match, err := regexp.MatchString("^[a-z][a-z1-9]{2,15}$", str);
	if err != nil {
		return false, err
	}

	if !match {
		return false, fmt.Errorf("Error parameter!")
	}

	return true, nil
}

