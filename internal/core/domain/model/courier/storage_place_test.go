package courier

import (
	"delivery/pkg/constructor"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	validStoragePlaceName        = "storage name"
	validStoragePlaceTotalVolume = 1
)

var (
	id         = uuid.New()
	orderID    = uuid.New()
	orderIDPtr = &orderID
)

func TestNewStoragePlace(t *testing.T) {
	tests := []struct {
		id          uuid.UUID
		testName    string
		name        string
		totalVolume int
		orderID     *uuid.UUID
		want        *StoragePlace
		wantErrs    []error
	}{
		{
			testName:    "Короткое имя",
			id:          id,
			name:        "",
			totalVolume: validStoragePlaceTotalVolume,
			orderID:     nil,
			want:        nil,
			wantErrs:    []error{ErrInvalidStoragePlaceName},
		},
		{
			testName:    "Отрицательная вместительность",
			id:          id,
			name:        validStoragePlaceName,
			totalVolume: -1,
			orderID:     nil,
			want:        nil,
			wantErrs:    []error{ErrInvalidStoragePlaceVolume},
		},
		{
			testName:    "Нулевая вместительность",
			id:          id,
			name:        validStoragePlaceName,
			totalVolume: 0,
			orderID:     nil,
			want:        nil,
			wantErrs:    []error{ErrInvalidStoragePlaceVolume},
		},
		{
			testName:    "Не заполненный идентификатор заказа",
			id:          id,
			name:        validStoragePlaceName,
			totalVolume: validStoragePlaceTotalVolume,
			orderID:     nil,
			want: &StoragePlace{
				id:          id,
				name:        validStoragePlaceName,
				totalVolume: validStoragePlaceTotalVolume,
				orderID:     nil,
			},
			wantErrs: nil,
		},
		{
			testName:    "Заполненный идентификатор заказа",
			id:          id,
			name:        validStoragePlaceName,
			totalVolume: validStoragePlaceTotalVolume,
			orderID:     orderIDPtr,
			want: &StoragePlace{
				id:          id,
				name:        validStoragePlaceName,
				totalVolume: validStoragePlaceTotalVolume,
				orderID:     orderIDPtr,
			},
			wantErrs: nil,
		},
	}

	for _, tt := range tests {
		result, err := NewStoragePlace(tt.id, tt.name, tt.totalVolume, tt.orderID)

		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(t, result, tt.want)

			for _, wantErr := range tt.wantErrs {
				assert.ErrorIs(t, err, wantErr)
			}
		})
	}
}

func TestStoragePlace_Equals(t *testing.T) {
	tests := []struct {
		name       string
		sp1        *StoragePlace
		sp2        *StoragePlace
		wantEquals bool
	}{
		{
			name:       "Расчет равенства для разных валидных StoragePlace",
			sp1:        constructor.MustBeValid(NewStoragePlace(uuid.New(), validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			sp2:        constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			wantEquals: false,
		},
		{
			name:       "Расчет равенства для разных валидных StoragePlace с одинаковым идентификатором",
			sp1:        constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName+"test", validStoragePlaceTotalVolume+1, nil)),
			sp2:        constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			wantEquals: true,
		},
		{
			name:       "Расчет равенства для одинаковых валидных StoragePlace",
			sp1:        constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			sp2:        constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			wantEquals: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantEquals, tt.sp1.Equals(tt.sp2))
			assert.Equal(t, tt.wantEquals, tt.sp2.Equals(tt.sp1))
		})
	}
}

func TestStoragePlace_CanStore(t *testing.T) {
	tests := []struct {
		name   string
		sp     *StoragePlace
		volume int
		want   bool
	}{
		{
			name:   "StoragePlace занят другим заказом",
			sp:     constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			volume: validStoragePlaceTotalVolume,
			want:   false,
		},
		{
			name:   "StoragePlace не достаточно вместителен",
			sp:     constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, nil)),
			volume: validStoragePlaceTotalVolume + 1,
			want:   false,
		},
		{
			name:   "StoragePlace вместителен ровно настолько, насколько нужно",
			sp:     constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, nil)),
			volume: validStoragePlaceTotalVolume,
			want:   true,
		},
		{
			name:   "StoragePlace вместительнее, чем нужно",
			sp:     constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume+1, nil)),
			volume: validStoragePlaceTotalVolume,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sp.CanStore(tt.volume))
		})
	}
}

func TestStoragePlace_Store(t *testing.T) {
	tests := []struct {
		name        string
		sp          *StoragePlace
		orderID     uuid.UUID
		volume      int
		wantOrderID *uuid.UUID
		wantError   error
	}{
		{
			name:        "StoragePlace занят другим заказом",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			orderID:     orderID,
			volume:      validStoragePlaceTotalVolume,
			wantOrderID: orderIDPtr,
			wantError:   ErrStoragePlaceIsUnavailableForOrder,
		},
		{
			name:        "StoragePlace не достаточно вместителен",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, nil)),
			orderID:     orderID,
			volume:      validStoragePlaceTotalVolume + 1,
			wantOrderID: nil,
			wantError:   ErrStoragePlaceIsUnavailableForOrder,
		},
		{
			name:        "StoragePlace вместителен ровно настолько, насколько нужно",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, nil)),
			orderID:     orderID,
			volume:      validStoragePlaceTotalVolume,
			wantOrderID: orderIDPtr,
			wantError:   nil,
		},
		{
			name:        "StoragePlace вместительнее, чем нужно",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume+1, nil)),
			orderID:     orderID,
			volume:      validStoragePlaceTotalVolume,
			wantOrderID: orderIDPtr,
			wantError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sp.Store(tt.orderID, tt.volume)

			assert.Equal(t, tt.wantOrderID, tt.sp.OrderID())
			assert.ErrorIs(t, err, tt.wantError)
		})
	}
}

func TestStoragePlace_Clear(t *testing.T) {
	tests := []struct {
		name        string
		sp          *StoragePlace
		orderID     uuid.UUID
		wantOrderID *uuid.UUID
		wantError   error
	}{
		{
			name:        "StoragePlace занят другим заказом",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, orderIDPtr)),
			orderID:     uuid.New(),
			wantOrderID: orderIDPtr,
			wantError:   ErrStoragePlaceOrderNotFound,
		},
		{
			name:        "StoragePlace не занят каким-либо заказом",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume, nil)),
			orderID:     orderID,
			wantOrderID: nil,
			wantError:   ErrStoragePlaceOrderNotFound,
		},
		{
			name:        "StoragePlace занят искомым заказом",
			sp:          constructor.MustBeValid(NewStoragePlace(id, validStoragePlaceName, validStoragePlaceTotalVolume+1, orderIDPtr)),
			orderID:     orderID,
			wantOrderID: nil,
			wantError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sp.Clear(tt.orderID)

			assert.Equal(t, tt.wantOrderID, tt.sp.OrderID())
			assert.ErrorIs(t, err, tt.wantError)
		})
	}
}
