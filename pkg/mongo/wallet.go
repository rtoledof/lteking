package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
	"cubawheeler.io/pkg/derrors"
)

const WalletCollection Collections = "wallets"

var _ cubawheeler.WalletService = (*WalletService)(nil)

type WalletService struct {
	db *DB
}

func NewWalletService(db *DB) *WalletService {
	index := mongo.IndexModel{
		Keys: bson.D{{Key: "owner", Value: 1}},
	}
	_, err := db.Collection(WalletCollection).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create user index")
	}
	return &WalletService{db: db}
}

// Transactions implements cubawheeler.WalletService.
func (s *WalletService) Transactions(ctx context.Context, owner string) (_ []cubawheeler.TransferEvent, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Transactions")
	user := cubawheeler.UserFromContext(ctx)
	if user.ID != owner {
		return nil, fmt.Errorf("you are not allowed to do this: %w", cubawheeler.ErrForbidden)
	}
	w, err := findWalletByOwner(ctx, s.db, owner)
	if err != nil {
		return nil, err
	}
	return w.TransferEvent, nil
}

func (s *WalletService) Create(ctx context.Context, owner string) (_ *cubawheeler.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Create")
	w := cubawheeler.NewWallet()
	w.Owner = owner
	return w, storeWallet(ctx, s.db, w)
}

func (s *WalletService) FindByOwner(ctx context.Context, owner string) (_ *cubawheeler.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.FindByOwner")
	collection := s.db.Collection(WalletCollection)
	var w cubawheeler.Wallet
	f := bson.D{
		{Key: "$or", Value: []any{
			bson.D{{Key: "owner", Value: owner}},
			bson.D{{Key: "referer", Value: owner}},
		}},
	}
	err = collection.FindOne(ctx, f).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %w", err)
	}
	return &w, nil
}

func (s *WalletService) Deposit(ctx context.Context, owner string, amount int64, currency string) (_ *cubawheeler.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Deposit")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	if user.Role != cubawheeler.RoleAdmin {
		return nil, fmt.Errorf("you are not allowed to do this: %w", cubawheeler.ErrForbidden)
	}
	w, err := s.FindByOwner(ctx, owner)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %v: %w", err, cubawheeler.ErrNotFound)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %w", cubawheeler.ErrInvalidInput)
	}
	w.Deposit(amount, currency)
	return w, updateWallet(ctx, s.db, w)
}

func (s *WalletService) Withdraw(ctx context.Context, owner string, amount int64, currency string) (_ *cubawheeler.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Withdraw")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	w, err := s.FindByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}
	if w.Balance.Amount[currency]-amount < 0 {
		return nil, fmt.Errorf("insufficient funds: %w", cubawheeler.ErrInsufficientFunds)
	}
	w.Withdraw(amount, currency)
	// TODO: execute payout to the customer account in case of a driver
	return nil, updateWallet(ctx, s.db, w)
}

// Transfer implements cubawheeler.WalletService.
func (s *WalletService) Transfer(ctx context.Context, to string, amount int64, currency string) (_ *cubawheeler.TransferEvent, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Transfer")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return nil, cubawheeler.ErrAccessDenied
	}
	fromW, err := s.FindByOwner(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if !fromW.CanTransfer(amount, currency) {
		return nil, fmt.Errorf("insufficient funds: %w", cubawheeler.ErrInsufficientFunds)
	}

	transferEvent := cubawheeler.TransferEvent{
		ID:        cubawheeler.NewID().String(),
		From:      user.ID,
		To:        to,
		Type:      cubawheeler.TransferTypeTransfer,
		Status:    cubawheeler.TransferStatusPending,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(cubawheeler.Now().UTC().Unix()),
	}

	fromW.PendingTransfers = append(fromW.PendingTransfers, transferEvent)
	return &transferEvent, updateWallet(ctx, s.db, fromW)
}

func (s *WalletService) ConfirmTransfer(ctx context.Context, id, pin string) (err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.ConfirmTransfer")
	user := cubawheeler.UserFromContext(ctx)
	if user == nil {
		return cubawheeler.ErrAccessDenied
	}
	if err := user.ComparePin(pin); err != nil {
		return fmt.Errorf("invalid pin: %w", cubawheeler.ErrInvalidInput)
	}
	fromW, err := s.FindByOwner(ctx, user.ID)
	if err != nil {
		return err
	}
	pendingTransfer, index := fromW.FindPendingTransfer(id)
	if pendingTransfer == nil {
		return fmt.Errorf("transfer not found: %w", cubawheeler.ErrNotFound)
	}
	if !fromW.CanTransfer(pendingTransfer.Amount, pendingTransfer.Currency) {
		return fmt.Errorf("insufficient funds: %w", cubawheeler.ErrInsufficientFunds)
	}
	toW, err := s.FindByOwner(ctx, pendingTransfer.To)
	if err != nil {
		return err
	}

	if fromW.ID == toW.ID {
		return fmt.Errorf("invalid transfer: %w", cubawheeler.ErrInvalidInput)
	}

	fromW.Withdraw(pendingTransfer.Amount, pendingTransfer.Currency)
	toW.Deposit(pendingTransfer.Amount, pendingTransfer.Currency)
	fromW.TransferEvent = append(fromW.TransferEvent[:index], fromW.TransferEvent[index+1:]...)

	tx, err := s.db.client.StartSession()
	if err != nil {
		return err
	}
	err = tx.StartTransaction()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v: %w", err, cubawheeler.ErrInternal)
	}
	err = updateWallet(ctx, s.db, fromW)
	if err != nil {
		tx.AbortTransaction(ctx)
		return err
	}
	err = updateWallet(ctx, s.db, toW)
	if err != nil {
		tx.AbortTransaction(ctx)
		return err
	}
	return tx.CommitTransaction(ctx)
}

// Balance implements cubawheeler.WalletService.
func (s *WalletService) Balance(ctx context.Context, owner string) (_ cubawheeler.Balance, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Balance")
	w, err := findWalletByOwner(ctx, s.db, owner)
	if err != nil {
		return cubawheeler.Balance{}, err
	}
	return w.Balance, nil
}

func storeWallet(context context.Context, db *DB, w *cubawheeler.Wallet) error {
	collection := db.Collection(WalletCollection)
	_, err := collection.InsertOne(context, w)
	if err != nil {
		return fmt.Errorf("error inserting wallet: %w", err)
	}
	return nil
}

func updateWallet(context context.Context, db *DB, w *cubawheeler.Wallet) error {
	collection := db.Collection(WalletCollection)
	_, err := collection.UpdateOne(context, bson.M{"_id": w.ID}, bson.M{"$set": w})
	if err != nil {
		return fmt.Errorf("error updating wallet: %w", err)
	}
	return nil
}

func findWalletByOwner(context context.Context, db *DB, owner string) (*cubawheeler.Wallet, error) {
	collection := db.Collection(WalletCollection)
	var w cubawheeler.Wallet
	f := bson.D{
		{Key: "$or", Value: []any{
			bson.D{{Key: "owner", Value: owner}},
			bson.D{{Key: "referer", Value: owner}},
		}},
	}
	err := collection.FindOne(context, f).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %w", err)
	}
	return &w, nil
}
