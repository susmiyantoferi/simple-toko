package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type orderRepositoryImpl struct {
	Db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) *orderRepositoryImpl {
	return &orderRepositoryImpl{
		Db: db,
	}
}

const (
	Waiting   string = "waiting"
	Confirmed string = "confirmed"
	Canceled  string = "canceled"
	OnProcess string = "on process"
	Delivered string = "delivered"
)

var (
	ErrEmptyItems      = errors.New("order has no items")
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrAddressNotFound = errors.New("address not found")
)

func (o *orderRepositoryImpl) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if order == nil {
		return nil, fmt.Errorf("nil order")
	}

	if len(order.OrderProducts) == 0 {
		return nil, ErrEmptyItems
	}

	err := o.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		order.StatusOrder = Waiting
		order.StatusDelivery = Waiting

		var address entity.Address
		if err := tx.WithContext(ctx).First(&address, order.AddressID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAddressNotFound
			}
		}

		if err := tx.Omit("OrderProducts").Create(order).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		//create data on table pivot
		for i := range order.OrderProducts {
			item := &order.OrderProducts[i]
			item.OrderID = order.ID

			var p entity.Product
			if err := tx.Select("id, price").First(&p, item.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrProductNotFound
				}
				return fmt.Errorf("find product: %w", err)
			}
			item.UnitPrice = p.Price

			//reduce stock
			stock := tx.Model(&entity.Product{}).Where("id = ? AND stock >= ?", item.ProductID, item.Qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Qty))
			if stock.Error != nil {
				return fmt.Errorf("reduce stock: %w", stock.Error)
			}

			if stock.RowsAffected == 0 {
				return ErrNotEnoughStock
			}
		}

		if err := tx.Create(&order.OrderProducts).Error; err != nil {
			return fmt.Errorf("crate order item: %w", err)
		}

		amountPay := 0.0
		for _, v := range order.OrderProducts {
			amountPay += v.UnitPrice * float64(v.Qty)
		}

		if err := tx.Model(order).Update("amount_pay", amountPay).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := o.Db.WithContext(ctx).
		Preload("OrderProducts").Preload("OrderProducts.Product").Preload("User").
		Preload("Address").First(order, order.ID).Error; err != nil {
		return nil, fmt.Errorf("order repo: preload order: %w", err)
	}

	return order, nil
}

func (o *orderRepositoryImpl) UpdateAddress(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	newAddress := order.AddressID

	if err := o.Db.WithContext(ctx).Where("status_order = ? AND id = ?", Waiting, order.ID).First(order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order repo: find id update: %w", err)
	}

	var adrs entity.Address
	if err := o.Db.WithContext(ctx).First(&adrs, newAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("order repo: find id address update: %w", err)
	}

	updateDta := map[string]interface{}{
		"address_id": newAddress,
	}

	if err := o.Db.WithContext(ctx).Model(order).Select("address_id").Updates(updateDta).Error; err != nil {
		return nil, fmt.Errorf("order repo: update: %w", err)
	}

	if err := o.Db.WithContext(ctx).
		Preload("OrderProducts").Preload("OrderProducts.Product").Preload("User").
		Preload("Address").First(order, order.ID).Error; err != nil {
		return nil, fmt.Errorf("order repo: preload order: %w", err)
	}

	return order, nil
}

func (o *orderRepositoryImpl) Delete(ctx context.Context, id uint) error {
	order := entity.Order{}

	result := o.Db.WithContext(ctx).Delete(&order, id)
	if result.Error != nil {
		return fmt.Errorf("order repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil

}

func (o *orderRepositoryImpl) FindById(ctx context.Context, id uint) (*entity.Order, error) {
	order := entity.Order{}

	if err := o.Db.WithContext(ctx).Preload("OrderProducts").Preload("User").
		Preload("Address").Preload("OrderProducts.Product").
		First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order repo: find by id: %w", err)
	}

	return &order, nil
}

func (o *orderRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*entity.Order, int64, error) {
	var order []*entity.Order
	var totalItems int64

	if err := o.Db.WithContext(ctx).Model(&entity.Order{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := o.Db.WithContext(ctx).Limit(pageSize).Offset(offset).Preload("OrderProducts").
		Preload("OrderProducts.Product").Preload("User").
		Preload("Address").Find(&order).Error; err != nil {
		return nil, 0, err
	}

	return order, totalItems, nil
}

func (o *orderRepositoryImpl) FindByOrderId(ctx context.Context, orderId uint) ([]*entity.OrderProduct, error) {
	var order []*entity.OrderProduct

	if err := o.Db.WithContext(ctx).Preload("Order").Preload("Product").
		First(&order, orderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("order repo: find by user id: %w", err)
	}

	return order, nil
}

func (o *orderRepositoryImpl) ConfirmOrder(ctx context.Context, orderId uint, order *entity.Order) (*entity.Order, error) {
	data := map[string]interface{}{}

	if order.StatusOrder != "" {
		data["status_order"] = order.StatusOrder
	}

	if order.StatusDelivery != "" {
		data["status_delivery"] = order.StatusDelivery
	}

	if err := o.Db.WithContext(ctx).Where("id = ? AND (status_order = ? OR status_order = ?)", orderId, Waiting, Confirmed).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order repo: confirm order find id: %w", err)
	}

	if err := o.Db.WithContext(ctx).Model(&order).Select("status_order", "status_delivery").
		Updates(data).Error; err != nil {
		return nil, fmt.Errorf("order repo: confirm order: %w", err)
	}

	if err := o.Db.WithContext(ctx).
		Preload("OrderProducts").Preload("OrderProducts.Product").
		Preload("User").Preload("Address").First(&order, orderId).Error; err != nil {
		return nil, fmt.Errorf("order repo: preload order confirm: %w", err)
	}

	return order, nil
}

//func (o *orderRepositoryImpl) AddOrderItem(ctx context.Context, orderId uint, item *entity.OrderProduct) (*entity.Order, error) {
// 	var order entity.Order

// 	err := o.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

// 		if err := tx.First(&order, orderId).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return ErrOrderNotFound
// 			}
// 			return fmt.Errorf("add item repo: find order: %w", err)
// 		}

// 		var p entity.Product
// 		if err := tx.First(&p, item.ProductID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return ErrProductNotFound
// 			}
// 			return fmt.Errorf("add item repo: find product: %w", err)
// 		}

// 		item.UnitPrice = p.Price
// 		item.OrderID = orderId

// 		if err := tx.Create(item).Error; err != nil {
// 			return err
// 		}

// 		newAmount := order.AmountPay + (item.UnitPrice * float64(item.Qty))
// 		if err := tx.Model(&order).Update("amount_pay", newAmount).Error; err != nil {
// 			return err
// 		}

// 		if err := tx.Preload("OrderProduct").First(&order, orderId).Error; err != nil {
// 			return err
// 		}

// 		return nil

// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &order, nil
// }

// func (o *orderRepositoryImpl) RemoveOrderItem(ctx context.Context, orderId, productId uint) (*entity.Order, error) {

// }

// func (o *orderRepositoryImpl) UpdateOrderQty(ctx context.Context, orderId, productId uint, qty int) (*entity.Order, error) {

// }
