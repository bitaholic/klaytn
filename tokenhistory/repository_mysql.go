package tokenhistory

import (
	"database/sql"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
)

func (th *TokenHistory) ListLatestKlayTransfer() []*KlayTransfer {
	q := `SELECT account_addr,
			   block_num,
			   tx_idx,
			   itx_idx,
			   opposite_addr,
			   tx_value,
			   balance,
			   tx_hash,
			   direction,
       		   block_time
		FROM klay_transfer_history
		ORDER BY account_addr, block_num desc, tx_idx desc, itx_idx desc
		LIMIT 100`
	rows, err := th.database.Query(q)
	if err != nil {
		logger.Error("failed to query db", "error", err)
		return nil
		// TODO
	}
	defer rows.Close()

	var result []*KlayTransfer
	var colBlockNum, colTxIdx, colItxIdx, colBlockTime uint64
	var colAccountRaw, colOppositeRaw, colTxHashRaw sql.RawBytes
	var colValueRaw, colBalanceRaw string
	var colDirection Direction

	for rows.Next() {
		err := rows.Scan(&colAccountRaw, &colBlockNum, &colTxIdx, &colItxIdx, &colOppositeRaw, &colValueRaw, &colBalanceRaw,
			&colTxHashRaw, &colDirection, &colBlockTime)
		if err != nil {
			logger.Error("error", "err", err)
			return nil
			// TODO
		}
		colAccount := new(common.Address)
		colAccount.SetBytes(colAccountRaw)
		colOpposite := new(common.Address)
		colOpposite.SetBytes(colOppositeRaw)
		colTxHash := new(common.Hash)
		colTxHash.SetBytes(colTxHashRaw)

		result = append(result, &KlayTransfer{
			BlockNumber:   colBlockNum,
			TxIdx:         colTxIdx,
			InternalTxIdx: colItxIdx,
			Account:       colAccount,
			Opposite:      colOpposite,
			Value:         hexutil.MustDecodeBig(colValueRaw), // TODO  Must 로 해도 되는건지?
			Direction:     colDirection,
			TxHash:        colTxHash,
			Balance:       hexutil.MustDecodeBig(colBalanceRaw), // TODO  Must 로 해도 되는건지?
		})
	}
	return result
}

func (th *TokenHistory) ListAccountKlayTransferBetweenTime(addr *common.Address, fromTime, toTime uint64) []*KlayTransfer {
	q := `SELECT account_addr,
			   block_num,
			   tx_idx,
			   itx_idx,
			   opposite_addr,
			   tx_value,
			   balance,
			   tx_hash,
			   direction,
       		   block_time
		FROM klay_transfer_history
		WHERE account_addr = ? AND (block_time between ? AND ?)
		ORDER BY account_addr, block_num desc, tx_idx desc, itx_idx desc
		LIMIT 100`
	rows, err := th.database.Query(q, addr.Bytes(), fromTime, toTime)
	if err != nil {
		logger.Error("failed to query db", "error", err)
		return nil
		// TODO
	}
	defer rows.Close()

	var result []*KlayTransfer
	var colBlockNum, colTxIdx, colItxIdx, colBlockTime uint64
	var colAccountRaw, colOppositeRaw, colTxHashRaw sql.RawBytes
	var colValueRaw, colBalanceRaw string
	var colDirection Direction

	for rows.Next() {
		err := rows.Scan(&colAccountRaw, &colBlockNum, &colTxIdx, &colItxIdx, &colOppositeRaw, &colValueRaw, &colBalanceRaw,
			&colTxHashRaw, &colDirection, &colBlockTime)
		if err != nil {
			logger.Error("error", "err", err)
			return nil
			// TODO
		}
		colAccount := new(common.Address)
		colAccount.SetBytes(colAccountRaw)
		colOpposite := new(common.Address)
		colOpposite.SetBytes(colOppositeRaw)
		colTxHash := new(common.Hash)
		colTxHash.SetBytes(colTxHashRaw)

		result = append(result, &KlayTransfer{
			BlockNumber:   colBlockNum,
			TxIdx:         colTxIdx,
			InternalTxIdx: colItxIdx,
			Account:       colAccount,
			Opposite:      colOpposite,
			Value:         hexutil.MustDecodeBig(colValueRaw), // TODO  Must 로 해도 되는건지?
			Direction:     colDirection,
			TxHash:        colTxHash,
			Balance:       hexutil.MustDecodeBig(colBalanceRaw), // TODO  Must 로 해도 되는건지?
			BlockTime:     colBlockTime,
		})
	}
	return result
}

func (th *TokenHistory) ListKlayTransfer(addr *common.Address, fromBlock, toBlock rpc.BlockNumberOrHash) []*KlayTransfer {
	var fromBlockNum, toBlockNum uint64
	if bn, ok := fromBlock.Number(); ok {
		switch *fromBlock.BlockNumber {
		case rpc.PendingBlockNumber:
			return nil
		case rpc.EarliestBlockNumber:
			fromBlockNum = 1
		case rpc.LatestBlockNumber:
			fromBlockNum = th.blockChain.CurrentBlock().Number().Uint64()
		default:
			fromBlockNum = bn.Uint64()
		}
	} else if bn, ok := fromBlock.Hash(); ok {
		b := th.blockChain.GetBlockByHash(bn)
		if b == nil {
			return nil // TODO error
		}
		fromBlockNum = b.NumberU64()
	} else {
		// TODO
	}

	if bn, ok := toBlock.Number(); ok {
		switch *toBlock.BlockNumber {
		case rpc.PendingBlockNumber:
			return nil
		case rpc.EarliestBlockNumber:
			toBlockNum = 1
		case rpc.LatestBlockNumber:
			toBlockNum = th.blockChain.CurrentBlock().Number().Uint64()
		default:
			toBlockNum = bn.Uint64()
		}
	} else if bn, ok := toBlock.Hash(); ok {
		b := th.blockChain.GetBlockByHash(bn)
		if b == nil {
			return nil // TODO error
		}
		toBlockNum = b.NumberU64()
	} else {
		// TODO
	}

	logger.Info("from_to", "from", fromBlockNum, "to", toBlockNum)

	q := `SELECT account_addr,
			   block_num,
			   tx_idx,
			   itx_idx,
			   opposite_addr,
			   tx_value,
			   balance,
			   tx_hash,
			   direction,
               block_time
		FROM klay_transfer_history
		WHERE account_addr = ? AND (block_num between ? AND ?)
		ORDER BY account_addr, block_num desc, tx_idx desc, itx_idx desc
		LIMIT 100`
	rows, err := th.database.Query(q, addr.Bytes(), fromBlockNum, toBlockNum)
	if err != nil {
		logger.Error("failed to query db", "error", err)
		return nil
		// TODO
	}
	defer rows.Close()

	var result []*KlayTransfer
	var colBlockNum, colTxIdx, colItxIdx, colBlockTime uint64
	var colAccountRaw, colOppositeRaw, colTxHashRaw sql.RawBytes
	var colValueRaw, colBalanceRaw string
	var colDirection Direction

	for rows.Next() {
		err := rows.Scan(&colAccountRaw, &colBlockNum, &colTxIdx, &colItxIdx, &colOppositeRaw, &colValueRaw, &colBalanceRaw,
			&colTxHashRaw, &colDirection, &colBlockTime)
		if err != nil {
			logger.Error("error", "err", err)
			return nil
			// TODO
		}
		colAccount := new(common.Address)
		colAccount.SetBytes(colAccountRaw)
		colOpposite := new(common.Address)
		colOpposite.SetBytes(colOppositeRaw)
		colTxHash := new(common.Hash)
		colTxHash.SetBytes(colTxHashRaw)

		result = append(result, &KlayTransfer{
			BlockNumber:   colBlockNum,
			TxIdx:         colTxIdx,
			InternalTxIdx: colItxIdx,
			Account:       colAccount,
			Opposite:      colOpposite,
			Value:         hexutil.MustDecodeBig(colValueRaw), // TODO  Must 로 해도 되는건지?
			Direction:     colDirection,
			TxHash:        colTxHash,
			Balance:       hexutil.MustDecodeBig(colBalanceRaw), // TODO  Must 로 해도 되는건지?
			BlockTime:     colBlockTime,
		})
	}
	return result
}
