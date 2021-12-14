package tokenhistory

import (
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
)

var (
	logger = log.NewModuleLogger(log.TokenHistory)
)

func New(stack *node.Node) error {
	em := NewEmitterMysql("root", "test1234!", "localhost:3306", "token_history")

	srv := &TokenHistory{
		chainEventCh:     make(chan blockchain.ChainEvent, 1),
		chainEventChStop: make(chan struct{}),
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
}

func (b *TokenHistory) Protocols() []p2p.Protocol {
	return nil
}

func (b *TokenHistory) APIs() []rpc.API {
	return nil
}

func (b *TokenHistory) Start(server p2p.Server) error {
	go b.handleEvent()
	b.chainSub = b.blockChain.SubscribeChainEvent(b.chainEventCh)
	return nil
}

func (b *TokenHistory) Stop() error {
	b.chainSub.Unsubscribe()
	close(b.chainEventChStop)
	return nil
}

func (b *TokenHistory) Components() []interface{} {
	return nil
}

func (b *TokenHistory) SetComponents(components []interface{}) {
	for _, component := range components {
		switch v := component.(type) {
		case *blockchain.BlockChain:
			b.blockChain = v
		}
	}
}

func (b *TokenHistory) handleEvent() {
	for {
		select {
		case <-b.chainEventChStop:
			return
		case ev := <-b.chainEventCh:
			logger.Info("got message", "message", ev.Block.Number(), "txCount", len(ev.Block.Transactions()))
			tokenTransactions := parseBlock(ev)
			for _, t := range tokenTransactions {
				b.emitter.EmitMessage(t)
			}
			if len(ev.Logs) > 0 {
				for _, l := range ev.Logs {
					logger.Info("=======Log", "log", l)
					logger.Info("-- Topics", "topics", l.Topics)
					logger.Info("-- Data", "data", l.Data)
				}
			}
		}
	}
}
