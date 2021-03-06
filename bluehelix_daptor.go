package bluehelix_bridge

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/bluehelix-chain/chainnode/proto"
	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"

	"github.com/pkg/errors"

	stc "github.com/starcoinorg/starcoin-go/client"
	"github.com/starcoinorg/starcoin-go/types"
)

const PKLEN = 32
const ADDRESSLEN = 16

var NETURLMAP = make(map[string]string)

func init() {
	NETURLMAP["dev"] = "http://localhost:9850"
	NETURLMAP["banard"] = "https://barnard-seed.starcoin.org"
	NETURLMAP["proxima"] = "https://proxima-seed.starcoin.org"
	NETURLMAP["halley"] = "https://halley-seed.starcoin.org"
	NETURLMAP["main"] = "https://main-seed.starcoin.org/"
}

func findNetwork(name string) (string, error) {
	name = strings.ToLower(name)
	if url, found := NETURLMAP[name]; found {
		return url, nil
	} else {
		return "", fmt.Errorf("cant't found url by name %s", name)
	}
}

// StarcoinAdaptor use account model,only need to implements account related function,utxo function is not necessary
type StarcoinAdaptor struct {
}

func (adaptor *StarcoinAdaptor) ConvertAddress(req *proto.ConvertAddressRequest) (*proto.ConvertAddressReply, error) {
	if len(req.PublicKey) != PKLEN {
		return nil, fmt.Errorf("pk length should be 32")
	}

	var pk [PKLEN]byte
	copy(pk[:], req.PublicKey)

	address := stc.PublicKeyToAddress(pk)

	reply := proto.ConvertAddressReply{
		Code:    proto.ReturnCode_SUCCESS,
		Address: address,
	}
	return &reply, nil
}

func (adaptor *StarcoinAdaptor) ValidAddress(req *proto.ValidAddressRequest) (*proto.ValidAddressReply, error) {
	address := strings.Replace(req.Address, "0x", "", 1)

	addressBytes, err := hex.DecodeString(address)

	if err != nil {
		return &proto.ValidAddressReply{
			Code:  proto.ReturnCode_ERROR,
			Msg:   "address should be in hex",
			Valid: false,
		}, errors.Wrap(err, "address should be in hex")
	}

	if len(addressBytes) != ADDRESSLEN {
		return &proto.ValidAddressReply{
			Code:  proto.ReturnCode_ERROR,
			Msg:   "address length should be 16",
			Valid: false,
		}, errors.Wrap(err, "address length should be 16")
	}

	return &proto.ValidAddressReply{
		Code:  proto.ReturnCode_SUCCESS,
		Valid: true,
	}, nil
}

func (adaptor *StarcoinAdaptor) QueryBalance(req *proto.QueryBalanceRequest) (*proto.QueryBalanceReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.QueryBalanceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	client := stc.NewStarcoinClient(url)

	result, err := client.GetResource(context.Background(), req.Address)
	if err != nil {
		return &proto.QueryBalanceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find resource",
		}, errors.WithStack(err)
	}

	balance, err := result.GetBalanceOfStc()
	if err != nil {
		return &proto.QueryBalanceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find resource",
		}, errors.WithStack(err)
	}

	return &proto.QueryBalanceReply{
		Code:    proto.ReturnCode_SUCCESS,
		Balance: balance.String(),
	}, nil
}

func (adaptor *StarcoinAdaptor) QueryNonce(req *proto.QueryNonceRequest) (*proto.QueryNonceReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.QueryNonceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	client := stc.NewStarcoinClient(url)
	state, err := client.GetState(context.Background(), req.Address)
	if err != nil {
		return &proto.QueryNonceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find state",
		}, errors.WithStack(err)
	}

	return &proto.QueryNonceReply{
		Code:  proto.ReturnCode_SUCCESS,
		Nonce: state.SequenceNumber,
	}, nil
}

func (adaptor *StarcoinAdaptor) QueryGasPrice(req *proto.QueryGasPriceRequest) (*proto.QueryGasPriceReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.QueryGasPriceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	client := stc.NewStarcoinClient(url)
	price, err := client.GetGasUnitPrice(context.Background())
	if err != nil {
		return &proto.QueryGasPriceReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find gas price",
		}, errors.WithStack(err)
	}

	return &proto.QueryGasPriceReply{
		Code:     proto.ReturnCode_SUCCESS,
		GasPrice: fmt.Sprint(price),
	}, nil
}

func (adaptor *StarcoinAdaptor) CreateUtxoTransaction(req *proto.CreateUtxoTransactionRequest) (*proto.CreateUtxoTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) CreateAccountTransaction(req *proto.CreateAccountTransactionRequest) (*proto.CreateAccountTransactionReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	client := stc.NewStarcoinClient(url)
	gasPrice, err := strconv.Atoi(req.GasPrice)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "gas price should int",
		}, errors.WithStack(err)
	}

	gaslimit, err := strconv.Atoi(req.GasLimit)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "gas limit should int",
		}, errors.WithStack(err)
	}

	amountBigInt := new(big.Int)
	amountBigInt, ok := amountBigInt.SetString(req.Amount, 10)
	if !ok {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "amount should int",
		}, errors.WithStack(err)
	}

	amount, err := stc.BigIntToU128(amountBigInt)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "gas limit should int",
		}, errors.WithStack(err)
	}

	from, err := stc.ToAccountAddress(req.From)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "parse address failed",
		}, errors.WithStack(err)
	}

	to, err := stc.ToAccountAddress(req.To)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "parse address failed",
		}, errors.WithStack(err)
	}
	txn, err := client.BuildTransferStcTxn(context.Background(), *from, *to, *amount, gasPrice, uint64(gaslimit), req.Nonce)
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "create txn failed",
		}, errors.WithStack(err)
	}

	txnData, err := txn.BcsSerialize()
	if err != nil {
		return &proto.CreateAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "txn serialize failed",
		}, errors.WithStack(err)
	}

	return &proto.CreateAccountTransactionReply{
		Code:   proto.ReturnCode_SUCCESS,
		TxData: txnData,
	}, nil
}

func (adaptor *StarcoinAdaptor) CreateUtxoSignedTransaction(req *proto.CreateUtxoSignedTransactionRequest) (*proto.CreateSignedTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) CreateAccountSignedTransaction(req *proto.CreateAccountSignedTransactionRequest) (*proto.CreateSignedTransactionReply, error) {
	result := stc.Verify(req.PublicKey, req.TxData, req.Signature)
	if !result {
		return &proto.CreateSignedTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "verify sign failed",
		}, fmt.Errorf("verify sign failed")
	}

	rawUserTransaction, err := types.BcsDeserializeRawUserTransaction(req.TxData)
	if err != nil {
		return &proto.CreateSignedTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "Deserialize RawUserTransaction failed",
		}, errors.WithStack(err)
	}

	transactionAuthenticator := types.TransactionAuthenticator__Ed25519{
		PublicKey: types.Ed25519PublicKey(req.PublicKey),
		Signature: req.Signature,
	}

	signedUserTxn := types.SignedUserTransaction{
		RawTxn:        rawUserTransaction,
		Authenticator: &transactionAuthenticator,
	}
	signedTxn, err := signedUserTxn.BcsSerialize()
	if err != nil {
		return &proto.CreateSignedTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "serialize signedtxn failed",
		}, errors.WithStack(err)
	}

	hash, err := stc.GetSignedUserTransactionHash(signedUserTxn)
	if err != nil {
		return &proto.CreateSignedTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "get signedtxn hash failed",
		}, errors.WithStack(err)
	}

	return &proto.CreateSignedTransactionReply{
		Code:         proto.ReturnCode_SUCCESS,
		SignedTxData: signedTxn,
		Hash:         hash,
	}, nil
}

func (adaptor *StarcoinAdaptor) QueryAccountTransactionFromData(req *proto.QueryTransactionFromDataRequest) (*proto.QueryAccountTransactionReply, error) {
	_, err := types.BcsDeserializeRawUserTransaction(req.RawData)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "Deserialize RawUserTransaction failed",
		}, errors.WithStack(err)
	}

	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryAccountTransactionFromSignedData(req *proto.QueryTransactionFromSignedDataRequest) (*proto.QueryAccountTransactionReply, error) {
	signedTxn, err := types.BcsDeserializeSignedUserTransaction(req.SignedTxData)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "Deserialize RawUserTransaction failed",
		}, errors.WithStack(err)
	}

	txnHash, err := stc.GetSignedUserTransactionHash(signedTxn)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "get txn hash from signed user txn failed",
		}, errors.WithStack(err)
	}

	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	return adaptor.QueryAccountTransactionByHash(url, hex.EncodeToString(txnHash))
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransactionFromData(req *proto.QueryTransactionFromDataRequest) (*proto.QueryUtxoTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransactionFromSignedData(req *proto.QueryTransactionFromSignedDataRequest) (*proto.QueryUtxoTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) BroadcastTransaction(req *proto.BroadcastTransactionRequest) (*proto.BroadcastTransactionReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.BroadcastTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	client := stc.NewStarcoinClient(url)
	state, err := client.SubmitSignedTransactionBytes(context.Background(), req.SignedTxData)
	if err != nil {
		return &proto.BroadcastTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't broadcast txn",
		}, errors.WithStack(err)
	}

	return &proto.BroadcastTransactionReply{
		Code:   proto.ReturnCode_SUCCESS,
		TxHash: state,
	}, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxo(req *proto.QueryUtxoRequest) (*proto.QueryUtxoReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) QueryUtxoInsFromData(req *proto.QueryUtxoInsFromDataRequest) (*proto.QueryUtxoInsReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransaction(req *proto.QueryTransactionRequest) (*proto.QueryUtxoTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) QueryAccountTransaction(req *proto.QueryTransactionRequest) (*proto.QueryAccountTransactionReply, error) {
	url, err := findNetwork(req.Chain)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't find network",
		}, errors.WithStack(err)
	}

	return adaptor.QueryAccountTransactionByHash(url, req.TxHash)
}

func (adaptor *StarcoinAdaptor) QueryAccountTransactionByHash(url, hash string) (*proto.QueryAccountTransactionReply, error) {
	client := stc.NewStarcoinClient(url)
	txn, err := client.GetTransactionByHash(context.Background(), hash)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "can't get txn by hash",
		}, errors.WithStack(err)
	}

	nonce, err := strconv.Atoi(txn.UserTransaction.RawTransaction.SequenceNumber)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "nonce should be int",
		}, errors.WithStack(err)
	}

	blockHeight, err := strconv.Atoi(txn.BlockMetadata.Number)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "block height should be int",
		}, errors.WithStack(err)
	}

	sign, err := stc.HexStringToBytes(txn.UserTransaction.Authenticator.Ed25519.Signature)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "signature should be hexed string",
		}, errors.WithStack(err)
	}

	payload, err := decodePayload(txn.UserTransaction.RawTransaction.Payload)
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "decode payload failed",
		}, errors.WithStack(err)
	}

	payloadScriptFunction, ok := payload.(*types.TransactionPayload__ScriptFunction)
	if !ok {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "payload should be scriptfunction",
		}, fmt.Errorf("payload should be scriptfunction")
	}

	address, err := stc.ToAccountAddress("0x00000000000000000000000000000001")
	if err != nil {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "parse address failed",
		}, errors.WithStack(err)
	}

	if *address != payloadScriptFunction.Value.Module.Address {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "module address should be 0x00000000000000000000000000000001",
		}, fmt.Errorf("module address should be 0x00000000000000000000000000000001")
	}

	if payloadScriptFunction.Value.Module.Name != "TransferScripts" {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "module name should be TransferScripts",
		}, fmt.Errorf("module name should be TransferScripts")
	}

	if payloadScriptFunction.Value.Function != "peer_to_peer" {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "function name should be peer_to_peer",
		}, fmt.Errorf("function name should be peer_to_pee")
	}
	// TODO ?????????????????????type args???stc

	if len(payloadScriptFunction.Value.Args) != 2 {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "transfer script should have 2 args",
		}, fmt.Errorf("transfer script should have 2 args")
	}

	toArgs := payloadScriptFunction.Value.Args[0]

	toAddress, ok := toArgs.(*types.TransactionArgument__Address)
	if !ok {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "args0 should be address",
		}, fmt.Errorf("args0 should be address")
	}

	amountArgs := payloadScriptFunction.Value.Args[1]

	amount, ok := amountArgs.(*types.TransactionArgument__U128)
	if !ok {
		return &proto.QueryAccountTransactionReply{
			Code: proto.ReturnCode_ERROR,
			Msg:  "args1 should be u128",
		}, fmt.Errorf("args1 should be u128")
	}

	return &proto.QueryAccountTransactionReply{
		Code:        proto.ReturnCode_SUCCESS,
		TxHash:      txn.TransactionHash,
		From:        txn.BlockMetadata.Author,
		To:          stc.BytesToHexString(toAddress.Value[:]),
		Amount:      stc.U128ToBigInt((*serde.Uint128)(amount)).String(),
		Nonce:       uint64(nonce),
		GasPrice:    txn.UserTransaction.RawTransaction.GasUnitPrice,
		GasLimit:    txn.UserTransaction.RawTransaction.MaxGasAmount,
		BlockHeight: uint64(blockHeight),
		BlockTime:   uint64(txn.BlockMetadata.Timestamp),
		SignHash:    sign,
	}, nil

}

func (adaptor *StarcoinAdaptor) VerifyAccountSignedTransaction(req *proto.VerifySignedTransactionRequest) (*proto.VerifySignedTransactionReply, error) {
	//
	return nil, nil
}

func (adaptor *StarcoinAdaptor) VerifyUtxoSignedTransaction(req *proto.VerifySignedTransactionRequest) (*proto.VerifySignedTransactionReply, error) {
	panic("utxo txn is not nesscery function")
}

func (adaptor *StarcoinAdaptor) GetLatestBlockHeight() (int64, error) {
	return 0, nil
}

func (adaptor *StarcoinAdaptor) GetAccountTransactionByHeight(height int64, replyCh chan *proto.QueryAccountTransactionReply, errCh chan error) {
}

func (adaptor *StarcoinAdaptor) GetUtxoTransactionByHeight(height int64, replyCh chan *proto.QueryUtxoTransactionReply, errCh chan error) {
	panic("utxo txn is not nesscery function")
}

func decodePayload(payloadString string) (types.TransactionPayload, error) {
	payloadBytes, err := stc.HexStringToBytes(payloadString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	payload, err := types.BcsDeserializeTransactionPayload(payloadBytes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return payload, nil
}
