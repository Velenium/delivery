package courier

import (
	"errors"
	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

var (
	ErrInvalidStoragePlaceName   = errors.New("courier.storage_pace.name_invalid")
	ErrInvalidStoragePlaceVolume = errors.New("courier.storage_place.volume_invalid")
)

func NewStoragePlace(id uuid.UUID, name string, volume int, orderID *uuid.UUID) (_ *StoragePlace, err error) {
	if len(name) == 0 {
		err = errors.Join(err, ErrInvalidStoragePlaceName)
	}

	if volume <= 0 {
		err = errors.Join(err, ErrInvalidStoragePlaceVolume)
	}

	if err != nil {
		return nil, err
	}

	return &StoragePlace{
		id:          id,
		name:        name,
		totalVolume: volume,
		orderID:     orderID,
	}, nil
}

func (sp *StoragePlace) ID() uuid.UUID {
	return sp.id
}

func (sp *StoragePlace) Name() string {
	return sp.name
}

func (sp *StoragePlace) TotalVolume() int {
	return sp.totalVolume
}

func (sp *StoragePlace) OrderID() *uuid.UUID {
	return sp.orderID
}

func (sp *StoragePlace) Equals(sp2 *StoragePlace) bool {
	return sp.id == sp2.id
}

func (sp *StoragePlace) CanStore(volume int) bool {
	return sp.orderID == nil && volume <= sp.totalVolume
}

var (
	ErrStoragePlaceIsUnavailableForOrder = errors.New("courier.storage_pace.unavailable_for_order")
)

func (sp *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if !sp.CanStore(volume) {
		return ErrStoragePlaceIsUnavailableForOrder
	}

	sp.orderID = &orderID

	return nil
}

var (
	ErrStoragePlaceOrderNotFound = errors.New("courier.storage_place.order_not_found")
)

func (sp *StoragePlace) Clear(orderID uuid.UUID) error {
	if !sp.isOccupied() {
		return ErrStoragePlaceOrderNotFound
	}

	if *sp.orderID != orderID {
		return ErrStoragePlaceOrderNotFound
	}

	sp.orderID = nil

	return nil
}

func (sp *StoragePlace) isOccupied() bool {
	return sp.orderID == nil
}
