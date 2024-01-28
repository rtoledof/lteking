package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/jwt"

	"identity.io/pkg/identity"
)

func TestClientServiceCreate(t *testing.T) {
	// Setup db connection
	ctx := prepateContext(t)

	// Setup test client
	var tests = []struct {
		name    string
		domain  string
		want    string
		wantErr bool
	}{
		{
			name:   "success create",
			domain: "example.com",
			want:   "client id",
		},
		{
			name:    "failure invalid email",
			domain:  "invalid",
			wantErr: true,
		},
	}

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection("clients").Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call method being tested
			client := &identity.Client{
				Name:   tt.name,
				Domain: tt.domain,
			}
			err := srv.Create(ctx, client)

			// Assertions
			if err != nil && !tt.wantErr {
				t.Fatalf("got error: %v, want error? %v", err, tt.wantErr)
			}
			if err == nil {
				if tt.domain != client.Domain {
					t.Fatalf("got: %s, want: %s", client.Domain, tt.domain)
				}

				if client.ID == nil {
					t.Fatalf("client id should not be nil")
				}
			}
		})
	}
}

func TestClientServiceUpdate(t *testing.T) {
	ctx := prepateContext(t)

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection("clients").Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	client := &identity.Client{
		Name:   "Test Client",
		Domain: "test.com",
	}
	err := srv.Create(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name      string
		client    func(identity.Client) *identity.Client
		want      *identity.Client
		wantError bool
	}{
		{
			name: "success update",
			client: func(c identity.Client) *identity.Client {
				c.Name = "Updated Name"
				return &c
			},
			want: func() *identity.Client {
				c := client
				c.Name = "Updated Name"
				return c
			}(),
		},
		{
			name: "failure invalid id",
			client: func(c identity.Client) *identity.Client {
				c.ID = identity.ID("invalid id")
				return &c
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := tt.client(*client)

			err = srv.Update(ctx, client)

			// Assertions
			if err != nil && !tt.wantError {
				t.Fatalf("got error: %v, want error? %v", err, tt.wantError)
			}
			if tt.want != nil {
				if diff := cmp.Diff(client, tt.want); diff != "" {
					t.Errorf("client.Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestClientServerUpdateKey(t *testing.T) {
	ctx := prepateContext(t)

	// Create a PlanService instance with the mock collection
	db := NewTestDB(t)
	defer func() {
		db.Collection(ClientCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	client := &identity.Client{
		Name:   "Test Client",
		Domain: "test.com",
	}

	err := srv.Create(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	err = srv.UpdateKey(ctx, client.ID, true)
	if err != nil {
		t.Fatalf("failed to update keys: %v", err)
	}

}

func TestClientServiceFindClients(t *testing.T) {
	ctx := prepateContext(t)

	// Create a PlanService instance with the mock collection
	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.Collection(ClientCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	// Create test clients
	database = "test"
	clients := make([]*identity.Client, 3)
	for i := 0; i < 3; i++ {
		client := &identity.Client{
			Name:   fmt.Sprintf("Test Client %d", i),
			Domain: "test.com",
		}
		err := srv.Create(ctx, client)
		if err != nil {
			t.Fatal(err)
		}
		clients[i] = client
	}

	var tests = []struct {
		name      string
		filter    identity.ClientFilter
		want      []*identity.Client
		wantToken string
		wantErr   bool
	}{
		{
			name: "success find all",
			want: clients,
		},
		{
			name: "pagination",
			filter: identity.ClientFilter{
				Limit: 2,
			},
			want:      clients[:2],
			wantToken: clients[2].ID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, token, err := srv.FindClients(ctx, tt.filter)

			// Assertions
			if err != nil && !tt.wantErr {
				t.Fatalf("got error: %v, want error? %v", err, tt.wantErr)
			}
			if err == nil {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("Find() mismatch (-want +got):\n%s", diff)
				}
				if token != tt.wantToken {
					t.Errorf("token mismatch, got %s want %s", token, tt.wantToken)
				}
			}
		})
	}
}

func TestClientServiceFindByID(t *testing.T) {
	ctx := prepateContext(t)

	// Create a PlanService instance with the mock collection
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection("clients").Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	client := &identity.Client{
		Name:   "Test Client",
		Domain: "test.com",
	}

	err := srv.Create(ctx, client)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name      string
		id        identity.ID
		want      *identity.Client
		wantError bool
	}{
		{
			name: "success find by id",
			id:   client.ID,
			want: client,
		},
		{
			name:      "failure invalid id",
			id:        identity.ID("invalid"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.FindByID(ctx, tt.id)

			// Assertions
			if err != nil && !tt.wantError {
				t.Fatalf("got error: %v, want error? %v", err, tt.wantError)
			}
			if err == nil {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("FindByID() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestClientServiceFindByKey(t *testing.T) {
	ctx := prepateContext(t)

	// Create a PlanService instance with the mock collection
	db := NewTestDB(t)
	defer func() {
		db.Collection(ClientCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewClientService(db)

	client := &identity.Client{
		Name:   "Test Client",
		Domain: "test.com",
	}

	err := srv.Create(ctx, client)
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		name    string
		in      string
		want    *identity.Client
		public  bool
		wantErr bool
	}{
		{
			name: "success find by private key",
			want: client,
		},
		{
			name:   "success find by public key",
			public: true,
			want:   client,
		},
		{
			name:    "failure invalid client id",
			in:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func prepateContext(t *testing.T, roles ...identity.Role) context.Context {
	t.Helper()
	ctx := context.Background()

	token := jwt.New()
	token.Set("id", identity.NewID().String())
	user := identity.User{
		ID:    identity.NewID().String(),
		Name:  "test",
		Email: "test",
		Role:  identity.RoleAdmin,
	}
	if roles != nil {
		user.Role = roles[0]
	}
	userData, _ := json.Marshal(user)
	token.Set("user", userData)
	return jwtauth.NewContext(ctx, token, nil)
}
