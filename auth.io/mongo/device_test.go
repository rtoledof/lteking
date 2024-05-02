package mongo

import (
	"auth.io/models"
	"context"
	"reflect"
	"testing"
)

func TestDeviceService_AddDevice(t *testing.T) {
	type fields struct {
		DB *DB
	}
	type args struct {
		ctx    context.Context
		device models.Device
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeviceService{
				DB: tt.fields.DB,
			}
			if err := d.AddDevice(tt.args.ctx, tt.args.device); (err != nil) != tt.wantErr {
				t.Errorf("AddDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeviceService_FindAll(t *testing.T) {
	type fields struct {
		DB *DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Device
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeviceService{
				DB: tt.fields.DB,
			}
			got, err := d.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceService_Remove(t *testing.T) {
	type fields struct {
		DB *DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeviceService{
				DB: tt.fields.DB,
			}
			if err := d.Remove(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDeviceService(t *testing.T) {
	type args struct {
		db *DB
	}
	tests := []struct {
		name string
		args args
		want *DeviceService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeviceService(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeviceService() = %v, want %v", got, tt.want)
			}
		})
	}
}
