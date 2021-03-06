// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (h Header) MarshalJSON() ([]byte, error) {
	type HeaderJSON struct {
		ParentHash  common.Hash    `json:"parentHash"`
		UncleHash   common.Hash    `json:"sha3Uncles"`
		Coinbase    common.Address `json:"miner"`
		Root        common.Hash    `json:"stateRoot"`
		TxHash      common.Hash    `json:"transactionsRoot"`
		ReceiptHash common.Hash    `json:"receiptsRoot"`
		Bloom       Bloom          `json:"logsBloom"`
		Difficulty  *hexutil.Big   `json:"difficulty"`
		Number      *hexutil.Big   `json:"number"`
		GasLimit    *hexutil.Big   `json:"gasLimit"`
		GasUsed     *hexutil.Big   `json:"gasUsed"`
		Time        *hexutil.Big   `json:"timestamp"`
		Extra       hexutil.Bytes  `json:"extraData"`
		MixDigest   common.Hash    `json:"mixHash"`
		Nonce       BlockNonce     `json:"nonce"`
	}
	var enc HeaderJSON
	enc.ParentHash = h.ParentHash
	enc.UncleHash = h.UncleHash
	enc.Coinbase = h.Coinbase
	enc.Root = h.Root
	enc.TxHash = h.TxHash
	enc.ReceiptHash = h.ReceiptHash
	enc.Bloom = h.Bloom
	enc.Difficulty = (*hexutil.Big)(h.Difficulty)
	enc.Number = (*hexutil.Big)(h.Number)
	enc.GasLimit = (*hexutil.Big)(h.GasLimit)
	enc.GasUsed = (*hexutil.Big)(h.GasUsed)
	enc.Time = (*hexutil.Big)(h.Time)
	enc.Extra = h.Extra
	enc.MixDigest = h.MixDigest
	enc.Nonce = h.Nonce
	return json.Marshal(&enc)
}

func (h *Header) UnmarshalJSON(input []byte) error {
	type HeaderJSON struct {
		ParentHash  *common.Hash    `json:"parentHash"`
		UncleHash   *common.Hash    `json:"sha3Uncles"`
		Coinbase    *common.Address `json:"miner"`
		Root        *common.Hash    `json:"stateRoot"`
		TxHash      *common.Hash    `json:"transactionsRoot"`
		ReceiptHash *common.Hash    `json:"receiptsRoot"`
		Bloom       *Bloom          `json:"logsBloom"`
		Difficulty  *hexutil.Big    `json:"difficulty"`
		Number      *hexutil.Big    `json:"number"`
		GasLimit    *hexutil.Big    `json:"gasLimit"`
		GasUsed     *hexutil.Big    `json:"gasUsed"`
		Time        *hexutil.Big    `json:"timestamp"`
		Extra       hexutil.Bytes   `json:"extraData"`
		MixDigest   *common.Hash    `json:"mixHash"`
		Nonce       *BlockNonce     `json:"nonce"`
	}
	var dec HeaderJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	var x Header
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for Header")
	}
	x.ParentHash = *dec.ParentHash
	if dec.UncleHash == nil {
		return errors.New("missing required field 'sha3Uncles' for Header")
	}
	x.UncleHash = *dec.UncleHash
	if dec.Coinbase == nil {
		return errors.New("missing required field 'miner' for Header")
	}
	x.Coinbase = *dec.Coinbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for Header")
	}
	x.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for Header")
	}
	x.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for Header")
	}
	x.ReceiptHash = *dec.ReceiptHash
	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Header")
	}
	x.Bloom = *dec.Bloom
	if dec.Difficulty == nil {
		return errors.New("missing required field 'difficulty' for Header")
	}
	x.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	x.Number = (*big.Int)(dec.Number)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Header")
	}
	x.GasLimit = (*big.Int)(dec.GasLimit)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	x.GasUsed = (*big.Int)(dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	x.Time = (*big.Int)(dec.Time)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for Header")
	}
	x.Extra = dec.Extra
	if dec.MixDigest == nil {
		return errors.New("missing required field 'mixHash' for Header")
	}
	x.MixDigest = *dec.MixDigest
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' for Header")
	}
	x.Nonce = *dec.Nonce
	*h = x
	return nil
}
