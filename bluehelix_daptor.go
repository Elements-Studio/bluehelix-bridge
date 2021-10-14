package bluehelix_bridge

import (
	"github.com/bluehelix-chain/chainnode/proto"
)

type StarcoinAdaptor struct {
}

func (adaptor *StarcoinAdaptor) ConvertAddress(req *proto.ConvertAddressRequest) (*proto.ConvertAddressReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) ValidAddress(req *proto.ValidAddressRequest) (*proto.ValidAddressReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryBalance(req *proto.QueryBalanceRequest) (*proto.QueryBalanceReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryNonce(req *proto.QueryNonceRequest) (*proto.QueryNonceReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryGasPrice(req *proto.QueryGasPriceRequest) (*proto.QueryGasPriceReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) CreateUtxoTransaction(req *proto.CreateUtxoTransactionRequest) (*proto.CreateUtxoTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) CreateAccountTransaction(req *proto.CreateAccountTransactionRequest) (*proto.CreateAccountTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) CreateUtxoSignedTransaction(req *proto.CreateUtxoSignedTransactionRequest) (*proto.CreateSignedTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) CreateAccountSignedTransaction(req *proto.CreateAccountSignedTransactionRequest) (*proto.CreateSignedTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryAccountTransactionFromData(req *proto.QueryTransactionFromDataRequest) (*proto.QueryAccountTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryAccountTransactionFromSignedData(req *proto.QueryTransactionFromSignedDataRequest) (*proto.QueryAccountTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransactionFromData(req *proto.QueryTransactionFromDataRequest) (*proto.QueryUtxoTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransactionFromSignedData(req *proto.QueryTransactionFromSignedDataRequest) (*proto.QueryUtxoTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) BroadcastTransaction(req *proto.BroadcastTransactionRequest) (*proto.BroadcastTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxo(req *proto.QueryUtxoRequest) (*proto.QueryUtxoReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxoInsFromData(req *proto.QueryUtxoInsFromDataRequest) (*proto.QueryUtxoInsReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryUtxoTransaction(req *proto.QueryTransactionRequest) (*proto.QueryUtxoTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) QueryAccountTransaction(req *proto.QueryTransactionRequest) (*proto.QueryAccountTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) VerifyAccountSignedTransaction(req *proto.VerifySignedTransactionRequest) (*proto.VerifySignedTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) VerifyUtxoSignedTransaction(req *proto.VerifySignedTransactionRequest) (*proto.VerifySignedTransactionReply, error) {
	return nil, nil
}

func (adaptor *StarcoinAdaptor) GetLatestBlockHeight() (int64, error) {
	return 0, nil
}

func (adaptor *StarcoinAdaptor) GetAccountTransactionByHeight(height int64, replyCh chan *proto.QueryAccountTransactionReply, errCh chan error) {
}

func (adaptor *StarcoinAdaptor) GetUtxoTransactionByHeight(height int64, replyCh chan *proto.QueryUtxoTransactionReply, errCh chan error) {

}
