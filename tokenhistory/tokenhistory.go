package tokenhistory

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/event"
	klog "github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
	"log"
)

var (
	logger = klog.NewModuleLogger(klog.TokenHistory)
)

func NewDatabase(dbUser, dbPasswd, dbAddr, dbName string) *sql.DB {
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
		log.Fatal(err) // TODO 다른 방식으로 처리
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return db
}

func New(stack *node.Node) error {
	db := NewDatabase("root", "test1234!", "localhost:3306", "token_history")

	em := NewEmitterMysql(db, "token_history")

	srv := &TokenHistory{
		chainEventCh:     make(chan blockchain.ChainEvent, 1),
		chainEventChStop: make(chan struct{}),
		database:         db,
		emitter:          em,
	}

	// Register to LifeCycle
	stack.RegisterSubService(func(ctx *node.ServiceContext) (node.Service, error) {
		return srv, nil
	})
	return nil
}

type TokenHistory struct {
	blockChain       *blockchain.BlockChain
	chainEventCh     chan blockchain.ChainEvent
	chainEventChStop chan struct{}
	chainSub         event.Subscription
	emitter          *EmitterMysql
	database         *sql.DB
}

func (th *TokenHistory) Protocols() []p2p.Protocol {
	return nil
}

func (th *TokenHistory) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "tokenhistory",
			Version:   "1.0",
			Service:   NewPublicTokenHistoryAPI(th),
			Public:    true,
		},
	}
}

func (th *TokenHistory) Start(server p2p.Server) error {
	go th.handleEvent()
	th.chainSub = th.blockChain.SubscribeChainEvent(th.chainEventCh)
	return nil
}

func (th *TokenHistory) Stop() error {
	th.chainSub.Unsubscribe()
	close(th.chainEventChStop)
	return nil
}

func (th *TokenHistory) Components() []interface{} {
	return nil
}

func (th *TokenHistory) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			th.blockChain = v
		}
	}
}

func (th *TokenHistory) handleEvent() {
	for {
		select {
		case <-th.chainEventChStop:
			return
		case ev := <-th.chainEventCh:
			//logger.Info("got message", "message", ev.Block.Number(), "txCount", len(ev.Block.Transactions()))

			klayTransferMap := parseBlock2(ev)
			stateDB, err := th.blockChain.StateAt(ev.Block.Root())
			if err != nil {
				logger.Error("failed to get state", "error", err)
				return
			}
			for addr, transfers := range klayTransferMap {
				transfers[len(transfers)-1].Balance = stateDB.GetBalance(addr)
			}
			klayTransferMap.FillBalance()

			for addr, transfers := range klayTransferMap {
				logger.Info("Transfer", "addr", addr.Hex())
				for _, t := range transfers {
					logger.Info("[" + t.Account.Hex() + "] " + string(t.Direction) + " : " + t.Opposite.Hex() + " : " + t.Value.String() + " : Balance : " + t.Balance.String())
				}
			}

			th.emitter.EmitKlayTransfers(klayTransferMap)

			//tokenTransactions := parseBlock(ev)
			//for _, t := range tokenTransactions {
			//	stateDB, err := b.blockChain.StateAt(ev.Block.Root())
			//	if err != nil {
			//		logger.Error("failed to get state", "error", err)
			//		return
			//	}
			//	fromBalance := stateDB.GetBalance(*t.From)
			//	var toBalance *big.Int
			//	if t.To != nil {
			//		toBalance = stateDB.GetBalance(*t.To)
			//	}
			//
			//	toAddr := "Deploy New Contract"
			//	if t.To != nil {
			//		toAddr = t.To.Hex()
			//	}
			//	logger.Info("TOP-TX Balance", "fromAddr", t.From.Hex(),
			//		"fromBalance", fromBalance, "toAddr", toAddr, "toBalance", toBalance)
			//	b.emitter.EmitMessage(t)
			//
			//}

			//if len(ev.Logs) > 0 {
			//	for _, l := range ev.Logs {
			//		logger.Info("=======Log", "log", l)
			//		logger.Info("-- Topics", "topics", l.Topics)
			//		logger.Info("-- Data", "data", l.Data)
			//	}
			//}
		}
	}
}
