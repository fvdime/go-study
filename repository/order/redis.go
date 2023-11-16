package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fvdime/go-study/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	Client *redis.Client
}

func orderByIdKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepository) Insert(ctx context.Context, order model.Order) error {
	// marshal means encoding
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	// generating unique key
	key := orderByIdKey(order.OrderId)

	// transaction client
	txn := r.Client.TxPipeline()

	// user not override any data that exists
	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add orders set: %w", err)
	}
	
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrorNotExist = errors.New("order does not exist")

func (r *RedisRepository) FindById(ctx context.Context, id uint64) (model.Order, error) {
	key := orderByIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrorNotExist		
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode json: %w", err)
	}

	return order, nil
}

func (r *RedisRepository) DeleteById(ctx context.Context, id uint64) error {
	key := orderByIdKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrorNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}


	return nil
}

func (r *RedisRepository) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderByIdKey(order.OrderId)

	// setXX method we can set method if it is only exists
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrorNotExist
	} else if err != nil {
		return fmt.Errorf("set order: %w", err)
	}

	return nil
}

// pagination
type FindAllPage struct {
	Size uint64
	OffSet uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepository) FindAll(ctx context.Context, page FindAllPage) (FindResult, error){
	res := r.Client.SScan(ctx, "orders", page.OffSet, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
		Orders: []model.Order{},
		}, nil
	}

	// mget = multi get
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders %w", err)
	}

	// order slice
	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x:= x.(string)
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}