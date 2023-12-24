package mongo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWalletServiceCreate(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallet").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	w, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	if w.ID == "" {
		t.Fatal("expected wallet ID to be set")
	}
	if w.Balance != 0 {
		t.Fatal("expected wallet balance to be 0")
	}
	if w.CreatedAt == 0 {
		t.Fatal("expected wallet CreatedAt to be set")
	}
	if w.UpdatedAt == 0 {
		t.Fatal("expected wallet UpdatedAt to be set")
	}
}

func TestWalletServiceFindByOwner(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	w, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	w2, err := s.FindByOwner(context.Background(), w.Owner)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(w, w2); diff != "" {
		t.Fatalf("WalletService.FindByOwner() mismatch (-want +got):\n%s", diff)
	}
}

func TestWalletServiceDeposit(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	_, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	var test = []struct {
		owner      string
		amount     int64
		wantAmount int64
		wantErr    bool
	}{
		{"test", 100, 100, false},
		{"test", 100, 200, false},
		{"test", 100, 300, false},
	}

	for _, tt := range test {
		w, err := s.Deposit(context.Background(), tt.owner, tt.amount)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}
		if w.Balance != tt.wantAmount {
			t.Fatalf("expected wallet balance to be %d, got %d", tt.wantAmount, w.Balance)
		}
	}
}

func TestWalletServiceWithdraw(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	_, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(context.Background(), "test", 200)

	var test = []struct {
		owner      string
		amount     int64
		wantAmount int64
		wantErr    bool
	}{
		{"test", 100, 100, false},
		{"test", 100, 0, false},
		{"test", 100, 0, true},
	}

	for _, tt := range test {
		w, err := s.Withdraw(context.Background(), tt.owner, tt.amount)
		if err != nil && !tt.wantErr {
			t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
		}
		if w != nil && w.Balance != tt.wantAmount {
			t.Fatalf("expected wallet balance to be %d, got %d", tt.wantAmount, w.Balance)
		}
	}
}

func TestWalletServiceTransfer(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	_, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Create(context.Background(), "test2")
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(context.Background(), "test", 200)

	var test = []struct {
		from           string
		to             string
		amount         int64
		wantFromAmount int64
		wantToAmount   int64
		wantErr        bool
	}{
		{"test", "test2", 100, 100, 100, false},
		{"test", "test2", 100, 0, 200, false},
		{"test", "test2", 100, 0, 200, true},
	}

	for _, tt := range test {
		from, to, err := s.Transfer(context.Background(), tt.from, tt.to, tt.amount)
		if err != nil && !tt.wantErr {
			t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
		}
		if from != nil && from.Balance != tt.wantFromAmount {
			t.Fatalf("expected wallet balance to be %d, got %d", tt.wantFromAmount, from.Balance)
		}
		if to != nil && to.Balance != tt.wantToAmount {
			t.Fatalf("expected wallet balance to be %d, got %d", tt.wantToAmount, to.Balance)
		}
	}
}

func TestWalletServiceBalance(t *testing.T) {

	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(context.Background())
		db.client.Disconnect(context.Background())
	}()
	s := NewWalletService(db)

	w, err := s.Create(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(context.Background(), "test", 200)

	balance, err := s.Balance(context.Background(), w.Owner)
	if err != nil {
		t.Fatal(err)
	}
	if balance != 200 {
		t.Fatalf("expected wallet balance to be %d, got %d", 200, balance)
	}
}
