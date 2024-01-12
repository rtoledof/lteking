package mongo

import (
	"context"
	"testing"

	"cubawheeler.io/pkg/cubawheeler"
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
	if w.CreatedAt == 0 {
		t.Fatal("expected wallet CreatedAt to be set")
	}
	if w.UpdatedAt == 0 {
		t.Fatal("expected wallet UpdatedAt to be set")
	}
}

func TestWalletServiceFindByOwner(t *testing.T) {

	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "test",
		Role: cubawheeler.RoleRider,
	})

	database = "test"
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	w, err := s.Create(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	w2, err := s.FindByOwner(ctx, w.Owner)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(w, w2); diff != "" {
		t.Fatalf("WalletService.FindByOwner() mismatch (-want +got):\n%s", diff)
	}
}

func TestWalletServiceDeposit(t *testing.T) {
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "test",
		Role: cubawheeler.RoleAdmin,
	})
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	var test = []struct {
		owner      string
		amount     int64
		currency   string
		wantAmount cubawheeler.Balance
		wantErr    bool
	}{
		{"test", 100, "CUP", cubawheeler.Balance{Amount: map[string]int64{"CUP": 100}}, false},
		{"test", 100, "CUP", cubawheeler.Balance{Amount: map[string]int64{"CUP": 200}}, false},
		{"test", 100, "CUP", cubawheeler.Balance{Amount: map[string]int64{"CUP": 300}}, false},
	}

	for _, tt := range test {
		w, err := s.Deposit(ctx, tt.owner, tt.amount, tt.currency)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}
		if diff := cmp.Diff(w.Balance, tt.wantAmount); diff != "" {
			t.Fatalf("WalletService.Deposit() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestWalletServiceWithdraw(t *testing.T) {
	ctx := context.Background()
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "test",
		Role: cubawheeler.RoleAdmin,
	})
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(ctx, "test", 200, "CUP")
	ctx = cubawheeler.NewContextWithUser(ctx, &cubawheeler.User{
		ID:   "test",
		Role: cubawheeler.RoleRider,
	})
	var test = []struct {
		owner      string
		amount     int64
		currency   string
		wantAmount int64
		wantErr    bool
	}{
		{"test", 100, "CUP", 100, false},
		{"test", 100, "CUP", 0, false},
		{"test", 100, "CUP", 0, true},
	}

	for _, tt := range test {
		w, err := s.Withdraw(ctx, tt.owner, tt.amount, tt.currency)
		if err != nil && !tt.wantErr {
			t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
		}
		if w != nil {
			if diff := cmp.Diff(w.Balance, tt.wantAmount); diff != "" {
				t.Fatalf("WalletService.Withdraw() mismatch (-want +got):\n%s", diff)
			}
		}
	}
}

func TestWalletServiceTransfer(t *testing.T) {
	ctx := context.Background()
	user := &cubawheeler.User{
		ID: "test",
	}
	user.EncryptPin("1234")
	ctx = cubawheeler.NewContextWithUser(ctx, user)

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

	s.Deposit(ctx, "test", 200, "CUP")

	var test = []struct {
		from     string
		to       string
		amount   int64
		currency string
		want     *cubawheeler.TransferEvent
		wantErr  bool
	}{
		{"test", "test2", 100, "CUP", &cubawheeler.TransferEvent{
			From:     "test",
			To:       "test2",
			Type:     cubawheeler.TransferTypeTransfer,
			Amount:   100,
			Currency: "CUP",
			Status:   cubawheeler.TransferStatusPending,
		}, false},
		{"test", "test2", 200, "CUP", &cubawheeler.TransferEvent{
			From:     "test",
			To:       "test2",
			Type:     cubawheeler.TransferTypeTransfer,
			Amount:   200,
			Currency: "CUP",
			Status:   cubawheeler.TransferStatusPending,
		}, false},
		{"test", "test2", 400, "CUP", nil, true},
	}

	for _, tt := range test {
		event, err := s.Transfer(ctx, tt.to, tt.amount, tt.currency)
		if err != nil && !tt.wantErr {
			t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
		}
		if tt.want != nil {
			tt.want.ID = event.ID
			tt.want.CreatedAt = event.CreatedAt
		}
		if diff := cmp.Diff(event, tt.want); diff != "" {
			t.Fatalf("WalletService.Transfer() mismatch (-want +got):\n%s", diff)
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

	s.Deposit(context.Background(), "test", 200, "CUP")

	balance, err := s.Balance(context.Background(), w.Owner)
	if err != nil {
		t.Fatal(err)
	}
	if balance.Amount["CUP"] != 200 {
		t.Fatalf("expected wallet balance to be %d, got %d", 200, balance.Amount["CUP"])
	}
}

func TestWalletServiceConfirmTransfer(t *testing.T) {
	ctx := context.Background()
	user := &cubawheeler.User{
		ID: "test",
	}
	user.EncryptPin("1234")
	ctx = cubawheeler.NewContextWithUser(ctx, user)
	database = "test"
	db := NewTestDB()
	defer func() {
		db.client.Database(database).Collection("wallets").Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx, "test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Create(ctx, "test2")
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(ctx, "test", 200, "CUP")

	var tests = []struct {
		name     string
		to       string
		amount   int64
		currency string
		wantErr  bool
	}{
		{"ok", "test2", 100, "CUP", false},
		{"invalid amount", "test2", 400, "CUP", true},
		{"valid", "test2", 100, "CUP", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := s.Transfer(ctx, tt.to, tt.amount, tt.currency)
			if err != nil && !tt.wantErr {
				t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
			}
			if event != nil {
				err = s.ConfirmTransfer(ctx, event.ID, "1234")
				if err != nil && !tt.wantErr {
					t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
				}
			}
		})
	}
}
