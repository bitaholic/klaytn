package tokenhistory

import (
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
)

type PublicTokenHistoryAPI struct {
	th *TokenHistory
}

func NewPublicTokenHistoryAPI(th *TokenHistory) *PublicTokenHistoryAPI {
	return &PublicTokenHistoryAPI{th: th}
}

func (api *PublicTokenHistoryAPI) ListKlayTransfer(addr *common.Address, fromBlock, toBlock rpc.BlockNumberOrHash) []*KlayTransfer {
	logger.Info("listKlayTransfer", "addr", addr.Hex(), "from", fromBlock, "to", toBlock)
	return api.th.ListKlayTransfer(addr, fromBlock, toBlock)
}

func (api *PublicTokenHistoryAPI) ListAccountKlayTransferBetweenTime(addr *common.Address, fromTimestamp, toTimestamp uint64) []*KlayTransfer {
	logger.Info("ListAccountKlayTransferBetweenTime", "addr", addr.Hex(), "fromTimestamp",
		fromTimestamp, "toTimeStamp", toTimestamp)
	return api.th.ListAccountKlayTransferBetweenTime(addr, fromTimestamp, toTimestamp)
}

func (api *PublicTokenHistoryAPI) LatestKlayTransfers() []*KlayTransfer {
	logger.Info("latestKlayTransfers")
	return api.th.ListLatestKlayTransfer()
}
