package core

import (
	cmn "github.com/DSiSc/apigateway/common"
	apitypes "github.com/DSiSc/apigateway/core/types"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
)

type blockdata struct {
	Number           cmn.Uint64           `json:"number"`
	Hash             cmn.Hash             `json:"hash"`
	ParentHash       cmn.Hash             `json:"parentHash"`
	MixHash          cmn.Hash             `json:"mixHash"`
	StateRoot        cmn.Hash             `json:"stateRoot"`
	Miner            apitypes.Address     `json:"miner"`
	Timestamp        cmn.Uint64           `json:"timestamp"`
	TransactionsRoot cmn.Hash             `json:"transactionsRoot"`
	ReceiptsRoot     cmn.Hash             `json:"receiptsRoot"`
	Transactions     []*types.Transaction `json:"transactions"`
}

func GetBlockByHash(blockHash cmn.Hash, fullTx bool) (*blockdata, error) {
	blockchain, err := blockchain.NewLatestStateBlockChain()
	if err == nil {
		block, err := blockchain.GetBlockByHash(TypeConvert(&blockHash))
		if block != nil {
			return rpcOutputBlock(block, true, fullTx)
		}
		return nil, err
	}
	return nil, err
}

func GetBlockByNumber(blockNr apitypes.BlockNumber, fullTx bool) (*blockdata, error) {
	blockchain, err := blockchain.NewLatestStateBlockChain()
	if err == nil {
		height := blockNr.Touint64()
		block, err := blockchain.GetBlockByHeight(height)
		if block != nil {
			return rpcOutputBlock(block, true, fullTx)
		}
		return nil, err
	}
	return nil, err
}

func TypeConvert(a *cmn.Hash) types.Hash {
	var hash types.Hash
	if a != nil {
		copy(hash[:], a[:])
	}
	return hash
}

func rpcOutputBlock(b *types.Block, inclTx bool, fullTx bool) (*blockdata, error) {
	fields, err := RPCMarshalBlock(b, inclTx, fullTx)
	if err != nil {
		return nil, err
	}
	//fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(b.Hash()))
	return fields, err
}

func RPCMarshalBlock(b *types.Block, inclTx bool, fullTx bool) (*blockdata, error) {
	head := b.Header // copies the header once
	fields := blockdata {
		Number:           (cmn.Uint64)(head.Height),
		Hash:             (cmn.Hash)(b.HeaderHash),
		ParentHash:       (cmn.Hash)(head.PrevBlockHash),
		MixHash:          (cmn.Hash)(head.MixDigest),
		StateRoot:        (cmn.Hash)(head.StateRoot),
		Miner:            (apitypes.Address)(head.Coinbase),
		Timestamp:        (cmn.Uint64)(head.Timestamp),
		TransactionsRoot: (cmn.Hash)(head.TxRoot),
		ReceiptsRoot:     (cmn.Hash)(head.ReceiptsRoot),
	}

	if inclTx {
		txs := b.Transactions
		if fullTx {
			fields.Transactions = txs
		}
	}

	return &fields, nil
}