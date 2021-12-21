package tokenhistory

import (
	"database/sql"
	"fmt"
	"github.com/klaytn/klaytn/common/hexutil"
)

type EmitterMysql struct {
	db           *sql.DB
	databaseName string
}

func NewEmitterMysql(db *sql.DB, dbName string) *EmitterMysql {
	return &EmitterMysql{db: db, databaseName: dbName}
}

func (e *EmitterMysql) EmitMessage(msg TokenTransaction) {
	// TODO InternalTX
	// TODO id (block_index, transaction_index, log_index)

	if msg.To == nil {
		logger.Info("Contract Deployed", "msg", msg)
		return
	}

	var tokenContract []byte
	if msg.TokenContract != nil {
		tokenContract = msg.TokenContract.Bytes()
	}

	q := fmt.Sprintf(
		`INSERT INTO %s.token_history
				(block_num, tx_idx, log_idx, from_addr, to_addr, tx_value, token_contract_addr)
				VALUES (?,?,?,?,?,?,?)
				`,
		e.databaseName)
	_, err := e.db.Exec(
		q,
		msg.BlockNumber,
		msg.TxIdx,
		msg.logIdx,
		msg.From.Bytes(),
		msg.To.Bytes(),
		msg.Value.String(),
		tokenContract,
	)
	if err != nil {
		logger.Error("failed to insert row", "error", err)
	}
}

func (e *EmitterMysql) EmitKlayTransfers(tm KlayTransferMap) {
	for _, transfers := range tm {
		q := fmt.Sprintf(
			`INSERT INTO %s.klay_transfer_history
				(account_addr, block_num, tx_idx, itx_idx, direction, opposite_addr, tx_value, balance, tx_hash, block_time)
                VALUES (?,?,?,?,?,?,?,?,?,?)`, e.databaseName)
		for _, t := range transfers {
			_, err := e.db.Exec(q,
				t.Account.Bytes(),
				t.BlockNumber,
				t.TxIdx,
				t.InternalTxIdx,
				t.Direction,
				t.Opposite.Bytes(),
				hexutil.EncodeBig(t.Value),
				hexutil.EncodeBig(t.Balance),
				t.TxHash.Bytes(),
				t.BlockTime)
			if err != nil {
				logger.Error("failed to insert klay transfer", "error", err)
			}
		}
	}
}
