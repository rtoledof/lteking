package mongo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"wallet.io/pkg/derrors"
	"wallet.io/pkg/wallet"
)

const WalletCollection Collections = "wallets"

var _ wallet.WalletService = (*WalletService)(nil)

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

// Transactions implements wallet.WalletService.
func (s *WalletService) Transactions(ctx context.Context) (_ []wallet.TransferEvent, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Transactions")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return nil, wallet.ErrAccessDenied
	}
	w, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return nil, err
	}
	return w.TransferEvent, nil
}

func (s *WalletService) Create(ctx context.Context) (_ *wallet.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Create")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return nil, wallet.ErrAccessDenied
	}

	w := wallet.NewWallet()
	w.Owner = *user
	return w, storeWallet(ctx, s.db, w)
}

func (s *WalletService) Wallet(ctx context.Context) (_ *wallet.Wallet, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.FindByOwner")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return nil, wallet.ErrAccessDenied
	}

	w, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WalletService) Deposit(ctx context.Context, owner string, amount int64, currency string) (err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Deposit")
	user := wallet.UserFromContext(ctx)
	if user == nil || user.Role != wallet.RoleAdmin {
		return wallet.ErrAccessDenied
	}
	w, err := findWallet(ctx, s.db, owner)
	if err != nil {
		return fmt.Errorf("error finding wallet: %v: %w", err, wallet.ErrNotFound)
	}
	if amount <= 0 {
		return fmt.Errorf("invalid amount: %w", wallet.ErrInvalidInput)
	}
	w.Deposit(amount, currency)
	return updateWallet(ctx, s.db, w)
}

func (s *WalletService) Withdraw(ctx context.Context, amount int64, currency string) (err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Withdraw")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return wallet.ErrAccessDenied
	}
	if user.Role != wallet.RoleDriver {
		return wallet.ErrAccessDenied
	}
	w, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return err
	}
	if w.Balance.Amount[currency]-amount < 0 {
		return fmt.Errorf("insufficient funds: %w", wallet.ErrInsufficientFunds)
	}
	w.Withdraw(amount, currency)
	// TODO: execute payout to the customer account in case of a driver
	return updateWallet(ctx, s.db, w)
}

// Transfer implements wallet.WalletService.
func (s *WalletService) Transfer(ctx context.Context, to string, amount int64, currency string) (_ *wallet.TransferEvent, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Transfer")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return nil, wallet.ErrAccessDenied
	}

	fromW, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return nil, err
	}
	if !fromW.CanTransfer(amount, currency) {
		return nil, fmt.Errorf("insufficient funds: %w", wallet.ErrInsufficientFunds)
	}

	transferEvent := wallet.TransferEvent{
		ID:        wallet.NewID().String(),
		From:      user.ID,
		To:        to,
		Type:      wallet.TransferTypeTransfer,
		Status:    wallet.TransferStatusPending,
		Amount:    amount,
		Currency:  currency,
		CreatedAt: uint(wallet.Now().UTC().Unix()),
	}

	fromW.PendingTransfers = append(fromW.PendingTransfers, transferEvent)
	return &transferEvent, updateWallet(ctx, s.db, fromW)
}

func (s *WalletService) ConfirmTransfer(ctx context.Context, id, pin string) (err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.ConfirmTransfer")
	claim, err := claimFromContext(ctx)
	if err != nil {
		return err
	}
	userData := claim.String("user")
	var user wallet.User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return wallet.ErrAccessDenied
	}
	fromW, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return err
	}
	pendingTransfer, index := fromW.FindPendingTransfer(id)
	if pendingTransfer == nil {
		return fmt.Errorf("transfer not found: %w", wallet.ErrNotFound)
	}
	if !fromW.CanTransfer(pendingTransfer.Amount, pendingTransfer.Currency) {
		return fmt.Errorf("insufficient funds: %w", wallet.ErrInsufficientFunds)
	}
	toW, err := findWallet(ctx, s.db, pendingTransfer.To)
	if err != nil {
		return err
	}

	if fromW.ID == toW.ID {
		return fmt.Errorf("invalid transfer: %w", wallet.ErrInvalidInput)
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
		return fmt.Errorf("error starting transaction: %v: %w", err, wallet.ErrInternal)
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

// Balance implements wallet.WalletService.
func (s *WalletService) Balance(ctx context.Context) (_ wallet.Balance, err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.Balance")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return wallet.Balance{}, wallet.ErrAccessDenied
	}

	w, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return wallet.Balance{}, err
	}
	return w.Balance, nil
}

// SetPin implements wallet.WalletService.
func (s *WalletService) SetPin(ctx context.Context, old string, new string) (err error) {
	defer derrors.Wrap(&err, "mongo.WalletService.SetPin")
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return wallet.ErrAccessDenied
	}
	w, err := findWallet(ctx, s.db, user.ID)
	if err != nil {
		return err
	}
	if err := w.ComparePin(old); err != nil && w.PIN != nil {
		return fmt.Errorf("invalid pin: %w", wallet.ErrInvalidInput)
	}
	if err := w.SetPin(new); err != nil {
		return err
	}
	return updateWallet(ctx, s.db, w)
}

func storeWallet(context context.Context, db *DB, w *wallet.Wallet) error {
	collection := db.Collection(WalletCollection)
	_, err := collection.InsertOne(context, w)
	if err != nil {
		return fmt.Errorf("error inserting wallet: %w", err)
	}
	return nil
}

func updateWallet(context context.Context, db *DB, w *wallet.Wallet) error {
	collection := db.Collection(WalletCollection)
	_, err := collection.UpdateOne(context, bson.M{"_id": w.ID}, bson.M{"$set": w})
	if err != nil {
		return fmt.Errorf("error updating wallet: %w", err)
	}
	return nil
}

func findWallet(context context.Context, db *DB, owner string) (*wallet.Wallet, error) {
	collection := db.Collection(WalletCollection)
	var w *wallet.Wallet
	f := bson.D{
		{Key: "$or", Value: []any{
			bson.D{{Key: "owner._id", Value: owner}},
			bson.D{{Key: "referer", Value: owner}},
		}},
	}
	err := collection.FindOne(context, f).Decode(&w)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w = wallet.NewWallet()
			w.Owner.ID = owner
			err = storeWallet(context, db, w)
			if err != nil {
				return nil, err
			}
		}
		return nil, wallet.ErrNotFound
	}
	return w, nil
}

func claimFromContext(ctx context.Context) (wallet.Claim, error) {
	_, claim, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return claim, nil
}

func create(ctx context.Context, db *DB) (*wallet.Wallet, error) {
	user := wallet.UserFromContext(ctx)
	if user == nil {
		return nil, wallet.ErrAccessDenied
	}

	w := wallet.NewWallet()
	w.Owner = *user
	return w, storeWallet(ctx, db, w)
}
