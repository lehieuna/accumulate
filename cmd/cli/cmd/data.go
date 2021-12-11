package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/AccumulateNetwork/accumulate/internal/api/v2"
	"github.com/AccumulateNetwork/accumulate/internal/url"
	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/spf13/cobra"
)

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Create, add, and query adi data accounts",
	Run: func(cmd *cobra.Command, args []string) {
		var out string
		var err error
		switch args[0] {
		case "get":
			switch len(args) {
			case 3:
				fallthrough
			case 4:
				out, err = GetDataEntry(args[1], args[2:])
			case 5:
				out, err = GetDataEntrySet(args[1], args[2:])
			default:
				PrintDataGet()
			}
		case "create":
			if len(args) > 2 {
				out, err = CreateDataAccount(args[1], args[2:])
			} else {
				PrintDataAccountCreate()
			}
		case "write":
			if len(args) > 2 {
				out, err = CreateDataAccount(args[1], args[2:])
			} else {
				PrintDataWrite()
			}
		default:
			PrintData()
		}
		printOutput(cmd, out, err)
	},
}

func PrintDataGet() {
	fmt.Println("  accumulate data get [DataAccountURL]			  Get existing Key Page by URL")
	fmt.Println("  accumulate data get [DataAccountURL] [EntryHash]  Get data entry by entryHash in hex")
	fmt.Println("  accumulate data get [DataAccountURL] [start index] [count]  Get a set of data entries starting from start and going to start+count")
	//./cli data get acc://actor/dataAccount
	//./cli data get acc://actor/dataAccount entryHash
	//./cli data get acc://actor/dataAccount start limit
}

func PrintDataAccountCreate() {
	//./cli data create acc://actor key idx height acc://actor/dataAccount acc://actor/keyBook (optional)
	fmt.Println("  accumulate data create [actor adi url] [signing key name] [key index (optional)] [key height (optional)] [adi data account url] [key book (optional)] Create new data account")
	fmt.Println("\t\t example usage: accumulate data create acc://actor signingKeyName acc://actor/dataAccount acc://actor/ssg0")
}

func PrintDataWrite() {
	fmt.Println("./cli data write [data account url] [signingKey] [extid_0 optional)] ... [extid_n (optional)] [data] Write entry to your data account. Note: extid's and data needs to be a quoted string or hex")
}

func PrintData() {
	PrintDataAccountCreate()
	PrintDataGet()
	PrintDataWrite()
}

func GetDataEntry(accountUrl string, args []string) (string, error) {
	u, err := url.Parse(accountUrl)
	if err != nil {
		return "", err
	}

	params := api.DataEntryQuery{}
	params.Url = u.String()
	if len(args) > 0 {
		n, err := hex.Decode(params.EntryHash[:], []byte(args[0]))
		if err != nil {
			return "", err
		}
		if n != 32 {
			return "", fmt.Errorf("entry hash must be 64 hex characters in length")
		}
	}

	var res api.QueryResponse

	data, err := json.Marshal(params)
	jsondata := json.RawMessage(data)
	if err != nil {
		return "", err
	}

	err = Client.RequestV2(context.Background(), "query-data", jsondata, &res)
	if err != nil {
		return "", err
	}

	return PrintQueryResponseV2(&res)
}

func GetDataEntrySet(accountUrl string, args []string) (string, error) {
	u, err := url.Parse(accountUrl)
	if err != nil {
		return "", err
	}

	if len(args) != 2 {
		return "", fmt.Errorf("expecting the start index and count parameters")
	}

	params := api.DataEntrySetQuery{}
	params.Url = u.String()

	v, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid start argument %s, %v", args[1], err)
	}
	params.Start = uint64(v)

	v, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid count argument %s, %v", args[1], err)
	}
	params.Count = uint64(v)

	var res api.QueryResponse
	data, err := json.Marshal(params)
	jsondata := json.RawMessage(data)
	if err != nil {
		return "", err
	}

	err = Client.RequestV2(context.Background(), "query-data-set", jsondata, &res)
	if err != nil {
		return "", err
	}

	return PrintQueryResponseV2(&res)
}

func CreateDataAccount(actorUrl string, args []string) (string, error) {
	actor, err := url.Parse(actorUrl)
	if err != nil {
		return "", err
	}

	args, si, privkey, err := prepareSigner(actor, args)
	if err != nil {
		return "", fmt.Errorf("insufficient number of command line arguments")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("expecting account url")
	}

	accountUrl, err := url.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("invalid account url %s", args[0])
	}
	if actor.Authority != accountUrl.Authority {
		return "", fmt.Errorf("account url to create (%s) doesn't match the authority adi (%s)", accountUrl.Authority, actor.Authority)
	}

	var keybook string
	if len(args) > 1 {
		kbu, err := url.Parse(args[1])
		if err != nil {
			return "", fmt.Errorf("invalid key book url")
		}
		keybook = kbu.String()
	}

	cda := protocol.CreateDataAccount{}
	cda.Url = accountUrl.String()
	cda.KeyBookUrl = keybook

	data, err := json.Marshal(cda)
	if err != nil {
		return "", err
	}

	dataBinary, err := cda.MarshalBinary()
	if err != nil {
		return "", err
	}

	nonce := nonceFromTimeNow()
	params, err := prepareGenTxV2(data, dataBinary, actor, si, privkey, nonce)
	if err != nil {
		return "", err
	}

	var res api.TxResponse

	if err := Client.RequestV2(context.Background(), "create-data-account", params, &res); err != nil {
		return PrintJsonRpcError(err)
	}

	return ActionResponseFrom(&res).Print()
}

func WriteData(accountUrl string, args []string) (string, error) {
	actor, err := url.Parse(accountUrl)
	if err != nil {
		return "", err
	}

	args, si, privkey, err := prepareSigner(actor, args)
	if err != nil {
		return "", fmt.Errorf("insufficient number of command line arguments")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("expecting account url")
	}

	wd := protocol.WriteData{}
	for i := 0; i < len(args); i++ {
		data := make([]byte, len(args[i]))
		if args[i][0:1] != "\"" {
			//attempt to hex decode it
			_, err := hex.Decode(data, []byte(args[i]))
			if err != nil {
				return "", fmt.Errorf("extid is neither hex nor quoted string")
			}
		} else {
			copy(data, args[i])
		}
		if i == len(args)-1 {
			wd.Entry.ExtIds = append(wd.Entry.ExtIds, data)
		} else {
			wd.Entry.Data = data
		}
	}

	data, err := json.Marshal(wd)
	if err != nil {
		return "", err
	}

	dataBinary, err := wd.MarshalBinary()
	if err != nil {
		return "", err
	}

	nonce := nonceFromTimeNow()
	params, err := prepareGenTxV2(data, dataBinary, actor, si, privkey, nonce)
	if err != nil {
		return "", err
	}

	var res api.TxResponse

	if err := Client.RequestV2(context.Background(), "write-data", params, &res); err != nil {
		return PrintJsonRpcError(err)
	}

	return ActionResponseFrom(&res).Print()
}
