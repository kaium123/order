// Package repository provides the database operations for the Order endpoint.
package repository

import (
	"context"
	"fmt"
	"github.com/kaium123/order/internal/db"
	"github.com/kaium123/order/internal/log"
	"github.com/kaium123/order/internal/model"
	"time"
)

// IOrder Order is the repository for the Order endpoint.
type IOrder interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	FindAllOrders(ctx context.Context, req *model.FindAllRequest) (*model.FindAllResponse, error)
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
	// Check if the database connection is initialized
	if o.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	// Attempt a simple query to ensure the connection is still valid
	if err := o.db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Use the database connection to insert the order
	_, err := o.db.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		o.log.Error(ctx, err.Error())
		return nil, err
	}

	return order, nil
}

func (o *OrderReceiver) FindAllOrders(ctx context.Context, req *model.FindAllRequest) (*model.FindAllResponse, error) {
	// Step 1: Query for the total count of matching records
	total := 0
	countQuery := o.db.NewSelect().
		Model((*model.Order)(nil)).
		ColumnExpr("COUNT(*)").
		Where("user_id = ?", req.UserId)

	if req.TransferStatus != "" {
		countQuery.Where("transfer_status = ?", req.TransferStatus)
	}
	if req.Archive != 0 {
		countQuery.Where("archive = ?", req.Archive)
	}

	err := countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, err
	}

	// Step 2: Query for paginated data
	orders := []*model.Order{}
	query := o.db.NewSelect().
		Model((*model.Order)(nil)).
		Limit(req.Limit).
		Offset(req.Offset)

	if req.TransferStatus != "" {
		query.Where("transfer_status = ?", req.TransferStatus)
	}
	if req.Archive != 0 {
		query.Where("archive = ?", req.Archive)
	}

	err = query.Scan(ctx, &orders)
	if err != nil {
		return nil, err
	}

	// Step 3: Calculate pagination metadata
	currentPage := req.Offset/req.Limit + 1
	lastPage := (total + req.Limit - 1) / req.Limit
	totalInPage := len(orders)

	// Step 4: Prepare the utils
	response := &model.FindAllResponse{
		Orders:      orders,
		Total:       total,
		CurrentPage: currentPage,
		PerPage:     req.Limit,
		TotalInPage: totalInPage,
		LastPage:    lastPage,
	}

	return response, nil
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
