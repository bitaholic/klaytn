package tokenhistory

import (
	"fmt"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"math/big"
)

var (
	/* ERC-20, KIP-7 Transfer Topic Hash */
	tokenTransferEventHash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

type TokenTransaction struct {
	BlockNumber   uint64
	TxIdx         uint64
	logIdx        uint64
	From          *common.Address
	To            *common.Address
	Value         *big.Int
	Hash          *common.Hash
	Unit          uint64
	TokenContract *common.Address
}

func parseBlock(msg blockchain.ChainEvent) []TokenTransaction {
	var ret []TokenTransaction

	blockNumber := msg.Block.NumberU64() // uint64 최고값을 넘어가면?
	for txIdx, tx := range msg.Block.Transactions() {
		var from common.Address
		if tx.IsLegacyTransaction() {
			signer := types.NewEIP155Signer(tx.ChainId())
			from, _ = types.Sender(signer, tx)
		} else {
			from, _ = tx.From()
		}
		to := tx.To()
		value := tx.Value()
		hash := tx.Hash()
		// 제외해야 할 Transactions
		// * Status 가 OK 가 아닐 경우
		// * Deploy Contract
		// * 0 Value 전달
		logger.Info("transaction", "blockNumber", blockNumber, "txIdx", uint64(txIdx),
			"from", from, "to", to, "value", value, "hash", hash)
		ret = append(ret, TokenTransaction{
			BlockNumber: blockNumber,
			TxIdx:       uint64(txIdx),
			logIdx:      0,
			From:        &from,
			To:          to,
			Value:       value,
			Hash:        &hash,
		})
	}

	printItx(msg.InternalTxTraces)
	//for _, itx := range msg.InternalTxTraces {
	//	logger.Info("---found internal tx", "itx", itx)
	//	if itx.Reverted != nil {
	//		logger.Info("---reverted", "msg", itx.Reverted.Message)
	//	}
	//logger.Info("---found internal tx", "i-from", itx.From.Hex(), "i-to", itx.To.Hex(), "i-value", itx.Value)
	//itx.Calls TODO Iteration 을 돌려서 해야 함 테스트 케이스 (Contract) 필요
	//ret = append(ret, TokenTransaction{
	//	BlockNumber: blockNumber,
	//	TxIdx:       uint64(txIdx),
	//	logIdx:      0,
	//	From:        &from,
	//	To:          to,
	//	Value:       value,
	//	Hash:        &hash,
	//})
	//}

	for logIdx, l := range msg.Logs {
		if len(l.Topics) > 0 && l.Topics[0].String() == tokenTransferEventHash {
			words, err := splitToWords(l.Data)
			if err != nil {
				logger.Error("failed to split data", "error", err)
				continue
			}
			data := append(l.Topics, words...)
			from := wordToAddress(data[1])
			to := wordToAddress(data[2])
			value := new(big.Int).SetBytes(data[3].Bytes())
			logger.Info("transfer", "from", from, "to", to,
				"value", value, "txIdx", l.TxIndex, "txHash", l.TxHash.Hex())
			ret = append(ret, TokenTransaction{
				BlockNumber:   blockNumber,
				TxIdx:         uint64(l.TxIndex),
				logIdx:        uint64(logIdx),
				From:          &from,
				To:            &to,
				Value:         value,
				Hash:          &l.TxHash,
				TokenContract: &l.Address,
			})
		}
	}
	return ret
}

func printItx(itx []*vm.InternalTxTrace) {
	for _, i := range itx {
		var from string
		if i.From != nil {
			from = i.From.Hex()
		}
		var to string
		if i.To != nil {
			to = i.To.Hex()
		}
		logger.Info("---found internal tx", "itx", i, "from", from, "to", to)
		printItx(i.Calls)
	}
}

// splitToWords divides log data to the words.
func splitToWords(data []byte) ([]common.Hash, error) {
	if len(data)%common.HashLength != 0 {
		return nil, fmt.Errorf("data length is not valid. want: %v, actual: %v", common.HashLength, len(data))
	}
	var words []common.Hash
	for i := 0; i < len(data); i += common.HashLength {
		words = append(words, common.BytesToHash(data[i:i+common.HashLength]))
	}
	return words, nil
}

// wordToAddress trims input word to get address field only.
func wordToAddress(word common.Hash) common.Address {
	return common.BytesToAddress(word[common.HashLength-common.AddressLength:])
}
