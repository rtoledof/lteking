package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"cubawheeler.io/pkg/cubawheeler"
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
	_, err := db.client.Database(database).Collection(WalletCollection.String()).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		panic("unable to create user index")
	}
	return &WalletService{db: db}
}

func (s *WalletService) Create(ctx context.Context, owner string) (*cubawheeler.Wallet, error) {
	w := cubawheeler.NewWallet()
	w.Owner = owner
	return w, storeWallet(ctx, s.db, w)
}

func (s *WalletService) FindByOwner(ctx context.Context, owner string) (*cubawheeler.Wallet, error) {
	collection := s.db.Collection(WalletCollection)
	var w cubawheeler.Wallet
	err := collection.FindOne(ctx, bson.M{"owner": owner}).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %w", err)
	}
	return &w, nil
}

func (s *WalletService) Deposit(ctx context.Context, owner string, amount int64) (*cubawheeler.Wallet, error) {
	w, err := s.FindByOwner(ctx, owner)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %v: %w", err, cubawheeler.ErrNotFound)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %w", cubawheeler.ErrInvalidInput)
	}
	w.Deposit(amount)
	return w, updateWallet(ctx, s.db, w)
}

func (s *WalletService) Withdraw(ctx context.Context, owner string, amount int64) (*cubawheeler.Wallet, error) {
	w, err := s.FindByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}
	if w.Balance-amount < 0 {
		return nil, fmt.Errorf("insufficient funds: %w", cubawheeler.ErrInsufficientFunds)
	}
	w.Withdraw(amount)
	return nil, updateWallet(ctx, s.db, w)
}

// Transfer implements cubawheeler.WalletService.
func (s *WalletService) Transfer(ctx context.Context, from, to string, amount int64) (*cubawheeler.Wallet, *cubawheeler.Wallet, error) {
	fromW, err := s.FindByOwner(ctx, from)
	if err != nil {
		return nil, nil, err
	}
	toW, err := s.FindByOwner(ctx, to)
	if err != nil {
		return nil, nil, err
	}
	if amount <= 0 {
		return nil, nil, fmt.Errorf("invalid amount: %w", cubawheeler.ErrInvalidInput)
	}
	if fromW.Balance-amount < 0 {
		return nil, nil, fmt.Errorf("insufficient funds: %w", cubawheeler.ErrInsufficientFunds)
	}
	fromW.Withdraw(amount)
	toW.Deposit(amount)
	tx, err := s.db.client.StartSession()
	if err != nil {
		return nil, nil, err
	}
	err = tx.StartTransaction()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting transaction: %v: %w", err, cubawheeler.ErrInternal)
	}
	err = updateWallet(ctx, s.db, fromW)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, nil, err
	}
	err = updateWallet(ctx, s.db, toW)
	if err != nil {
		tx.AbortTransaction(ctx)
		return nil, nil, err
	}

	return fromW, toW, tx.CommitTransaction(ctx)
}

// Balance implements cubawheeler.WalletService.
func (s *WalletService) Balance(ctx context.Context, owner string) (int64, error) {
	w, err := s.FindByOwner(ctx, owner)
	if err != nil {
		return 0, err
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
	err := collection.FindOne(context, bson.M{"owner": owner}).Decode(&w)
	if err != nil {
		return nil, fmt.Errorf("error finding wallet: %w", err)
	}
	return &w, nil
}
