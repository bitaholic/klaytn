package tokenhistory

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

type EmitterMysql struct {
	db           *sql.DB
	databaseName string
}

func NewEmitterMysql(dbUser, dbPasswd, dbAddr, dbName string) *EmitterMysql {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPasswd,
		Net:    "tcp",
		Addr:   dbAddr,
		DBName: dbName,
	}
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
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
