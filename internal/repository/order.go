// Package repository provides the database operations for the Order endpoint.
package repository

import (
	"context"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"time"
)

// IOrder Order is the repository for the Order endpoint.
type IOrder interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	FindAllOrders(ctx context.Context, req *model.FindAllRequest) ([]*model.Order, *model.PaginationResponse, error)
	CancelOrder(ctx context.Context, req *model.OrderCancelRequest) error
}

type InitOrderRepository struct {
	Db  *db.DB
	Log *log.Logger
}

type OrderReceiver struct {
	log *log.Logger
	db  *db.DB
}

// NewOrder returns a new instance of the Order repository.
func NewOrder(initOrderRepository *InitOrderRepository) IOrder {
	return &OrderReceiver{
		log: initOrderRepository.Log,
		db:  initOrderRepository.Db,
	}
}

// CreateOrder creates a new order in the database.
func (o *OrderReceiver) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {

	// Use the database connection to insert the order
	_, err := o.db.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		o.log.Error(ctx, err.Error())
		return nil, err
	}

	return order, nil
}

func (o *OrderReceiver) FindAllOrders(ctx context.Context, req *model.FindAllRequest) ([]*model.Order, *model.PaginationResponse, error) {
	// Step 1: Query for paginated data
	orders := []*model.Order{}
	query := o.db.NewSelect().
		Model((*model.Order)(nil)).
		Limit(req.Limit).
		Offset(req.Offset).
		Where("user_id = ? and deleted_at is null ", req.UserId)

	if req.TransferStatus != "" {
		query.Where("transfer_status = ?", req.TransferStatus)
	}
	if req.Archive != 0 {
		query.Where("archive = ?", req.Archive)
	}

	query.Order("created_at DESC")
	total, err := query.ScanAndCount(ctx, &orders)
	if err != nil {
		o.log.Error(ctx, err.Error())
		return nil, nil, err
	}

	// Step 3: Calculate pagination metadata
	currentPage := req.Offset/req.Limit + 1
	lastPage := (total + req.Limit - 1) / req.Limit
	totalInPage := len(orders)

	paginationResponse := &model.PaginationResponse{
		Total:       total,
		CurrentPage: currentPage,
		PerPage:     req.Limit,
		TotalInPage: totalInPage,
		LastPage:    lastPage,
	}

	return orders, paginationResponse, nil
}

func (o *OrderReceiver) CancelOrder(ctx context.Context, req *model.OrderCancelRequest) error {
	_, err := o.db.NewUpdate().Model((*model.Order)(nil)).
		Set("deleted_at = ?", time.Now().UTC()).
		Where("order_consignment_id = ? and user_id = ?", req.ConsignmentID, req.UserId).
		Exec(ctx)
	if err != nil {
		o.log.Error(ctx, err.Error())
		return err
	}

	return nil
}
