package graph

import (
	"time"

	"wallet.io/graph/model"
	"wallet.io/pkg/wallet"
)

func assembleModelTransfer(transfer *wallet.TransferEvent) *model.Transfer {
	t := time.Unix(int64(transfer.CreatedAt), 0)
	return &model.Transfer{
		ID:        transfer.ID,
		Amount:    int(transfer.Amount),
		Currency:  transfer.Currency,
		From:      transfer.From,
		To:        transfer.To,
		CreatedAt: t.Format("2006-01-02 15:04:05"),
	}
}
