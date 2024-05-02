package mongo

import (
	"auth.io/models"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestTenantServiceCreate(t *testing.T) {
	ctx := prepareContext(t)

	var tests = []struct {
		name string
		in   struct {
			name   string
			domain string
		}
		want    *models.Tenant
		wantErr bool
	}{
		{
			name: "success",
			in: struct {
				name   string
				domain string
			}{
				name:   "Test Tenant",
				domain: "test.com",
			},
			want: &models.Tenant{
				Name:   "Test Tenant",
				Domain: "test.com",
				Status: models.TenantStatusActive,
			},
		},
	}

	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection(tenantCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewTenantService(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tenant := models.Tenant{
				Name:   tt.in.name,
				Domain: tt.in.domain,
			}
			got, err := srv.Create(ctx, tenant)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TenantService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != nil {
				tt.want.ID = got.ID
				tt.want.CreatedAt = got.CreatedAt
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("TenantService.Create() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTenantServiceFindAll(t *testing.T) {
	ctx := prepareContext(t)

	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection(tenantCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewTenantService(db)

	tenant, err := srv.Create(ctx, models.Tenant{
		Name:   "Test Tenant",
		Domain: "test.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	tenant1, err := srv.Create(ctx, models.Tenant{
		Name:   "Test Tenant",
		Domain: "test.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		name    string
		filter  models.TenantFilter
		want    []models.Tenant
		token   string
		wantErr bool
	}{
		{
			name:    "success find all",
			filter:  models.TenantFilter{},
			want:    []models.Tenant{*tenant, *tenant1},
			token:   "",
			wantErr: false,
		},
		{
			name: "success find all",
			filter: models.TenantFilter{
				Limit: 1,
			},
			want:    []models.Tenant{*tenant},
			token:   tenant1.ID.String(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, token, err := srv.FindAll(ctx, tt.filter)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TenantService.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}

			if token != tt.token {
				t.Fatalf("TenantService.FindAll() token = %v, want %v", token, tt.token)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("TenantService.FindAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTenantServiceFindByID(t *testing.T) {
	ctx := prepareContext(t)

	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection(tenantCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewTenantService(db)

	tenant, err := srv.Create(ctx, models.Tenant{
		Name:   "Test Tenant",
		Domain: "test.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name    string
		id      models.ID
		want    models.Tenant
		wantErr bool
	}{
		{
			name:    "success find by id",
			id:      tenant.ID,
			want:    *tenant,
			wantErr: false,
		},
		{
			name:    "failure invalid id",
			id:      models.ID("invalid"),
			want:    models.Tenant{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.FindByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TenantService.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("TenantService.FindByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTenantServiceUpdate(t *testing.T) {
	ctx := prepareContext(t)

	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection(tenantCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewTenantService(db)

	tenant, err := srv.Create(ctx, models.Tenant{
		Name:   "Test Tenant",
		Domain: "test.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		name    string
		id      models.ID
		want    models.Tenant
		wantErr bool
	}{
		{
			name: "success update",
			id:   tenant.ID,
			want: *tenant,
		},
		{
			name:    "failure invalid id",
			id:      models.ID("invalid"),
			want:    models.Tenant{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.Update(ctx, tt.want)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TenantService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				tt.want.UpdatedAt = got.UpdatedAt
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("TenantService.Update() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTenantServiceDelete(t *testing.T) {
	ctx := prepareContext(t)

	database = "test"
	db := NewTestDB(t)
	defer func() {
		db.client.Database(database).Collection(tenantCollection).Drop(ctx)
		db.client.Disconnect(ctx)
	}()

	srv := NewTenantService(db)

	tenant, err := srv.Create(ctx, models.Tenant{
		Name:   "Test Tenant",
		Domain: "test.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	var tests = []struct {
		name    string
		id      models.ID
		wantErr bool
	}{
		{
			name:    "success delete",
			id:      tenant.ID,
			wantErr: false,
		},
		{
			name:    "failure invalid id",
			id:      models.ID("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.Delete(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TenantService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
