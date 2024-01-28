package mongo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/jwt"

	"wallet.io/pkg/wallet"
)

func TestWalletServiceCreate(t *testing.T) {
	ctx := prepateContext(t)
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	w, err := s.Create(ctx)
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
	ctx1 := prepateContext(t)

	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx1)
		db.client.Disconnect(ctx1)
	}()
	s := NewWalletService(db)

	w, err := s.Create(ctx1)
	if err != nil {
		t.Fatal(err)
	}

	w2, err := s.Wallet(ctx1)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(w, w2); diff != "" {
		t.Fatalf("WalletService.FindByOwner() mismatch (-want +got):\n%s", diff)
	}
}

func TestWalletServiceDeposit(t *testing.T) {
	ctx := prepateContext(t)

	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}
	user := wallet.UserFromContext(ctx)
	ctx = prepateContext(t, wallet.RoleAdmin)

	var test = []struct {
		owner      string
		amount     int64
		currency   string
		wantAmount wallet.Balance
		wantErr    bool
	}{
		{user.ID, 100, "CUP", wallet.Balance{Amount: map[string]int64{"CUP": 100}}, false},
		{user.ID, 100, "CUP", wallet.Balance{Amount: map[string]int64{"CUP": 200}}, false},
		{user.ID, 100, "CUP", wallet.Balance{Amount: map[string]int64{"CUP": 300}}, false},
	}

	for _, tt := range test {
		err := s.Deposit(ctx, tt.owner, tt.amount, tt.currency)
		if err != nil && !tt.wantErr {
			t.Fatal(err)
		}
	}
}

func TestWalletServiceWithdraw(t *testing.T) {
	ctx := prepateContext(t)
	// TODO: add claims on context
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}
	user := wallet.UserFromContext(ctx)
	ctx = prepateContext(t, wallet.RoleAdmin)
	s.Deposit(ctx, user.ID, 200, "CUP")
	// TODO: add claims on context
	var test = []struct {
		owner      string
		amount     int64
		currency   string
		wantAmount int64
		wantErr    bool
	}{
		{user.ID, 100, "CUP", 100, false},
		{user.ID, 100, "CUP", 0, false},
		{user.ID, 100, "CUP", 0, true},
	}

	for _, tt := range test {
		err := s.Withdraw(ctx, tt.amount, tt.currency)
		if err != nil && !tt.wantErr {
			t.Fatalf("expected no error, got %v, want: %v", err, tt.wantErr)
		}
	}
}

func TestWalletServiceTransfer(t *testing.T) {
	ctx1 := prepateContext(t)
	ctx2 := prepateContext(t)
	// TODO: add claims on context

	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx1)
		db.client.Disconnect(ctx1)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Create(ctx2)
	if err != nil {
		t.Fatal(err)
	}
	sender := wallet.UserFromContext(ctx1)
	receipt := wallet.UserFromContext(ctx2)
	adminContext := prepateContext(t, wallet.RoleAdmin)

	s.Deposit(adminContext, sender.ID, 200, "CUP")

	var test = []struct {
		to       string
		amount   int64
		currency string
		want     *wallet.TransferEvent
		wantErr  bool
	}{
		{receipt.ID, 100, "CUP", &wallet.TransferEvent{
			From:     sender.ID,
			To:       receipt.ID,
			Type:     wallet.TransferTypeTransfer,
			Amount:   100,
			Currency: "CUP",
			Status:   wallet.TransferStatusPending,
		}, false},
		{receipt.ID, 200, "CUP", &wallet.TransferEvent{
			From:     sender.ID,
			To:       receipt.ID,
			Type:     wallet.TransferTypeTransfer,
			Amount:   200,
			Currency: "CUP",
			Status:   wallet.TransferStatusPending,
		}, false},
		{receipt.ID, 400, "CUP", nil, true},
	}

	for _, tt := range test {
		event, err := s.Transfer(ctx1, tt.to, tt.amount, tt.currency)
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
	ctx := prepateContext(t)
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s.Deposit(ctx, "test", 200, "CUP")

	balance, err := s.Balance(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if balance.Amount["CUP"] != 200 {
		t.Fatalf("expected wallet balance to be %d, got %d", 200, balance.Amount["CUP"])
	}
}

func TestWalletServiceConfirmTransfer(t *testing.T) {
	ctx := context.Background()
	// TODO: add claims on context
	db := NewTestDB()
	defer func() {
		db.Collection(WalletCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()
	s := NewWalletService(db)

	_, err := s.Create(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Create(ctx)
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

func prepateContext(t *testing.T, roles ...wallet.Role) context.Context {
	t.Helper()
	ctx := context.Background()

	token := jwt.New()
	token.Set("id", wallet.NewID().String())
	user := wallet.User{
		ID:    wallet.NewID().String(),
		Name:  "test",
		Email: "test",
		Role:  "rider",
	}
	if roles != nil {
		user.Role = roles[0]
	}
	userData, _ := json.Marshal(user)
	token.Set("user", userData)
	return jwtauth.NewContext(ctx, token, nil)
}
