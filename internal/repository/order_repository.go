package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/orders_api/internal/models"
)

var (
	ErrOrderAlreadyExistsUUID  = errors.New("order with this uuid already exists")
	ErrOrderAlreadyExistsTrack = errors.New("order with this track_number already exists")
	ErrOrderNotFoundByUUID     = errors.New("orders with this ID not found")
)

type OrderPostgresRepository struct {
	Db *pgx.Conn
}

func NewOrderPostgresRepository(db *pgx.Conn) *OrderPostgresRepository {
	return &OrderPostgresRepository{
		Db: db,
	}
}

func (r *OrderPostgresRepository) GetOrderByUID(ctx context.Context, uid uuid.UUID) (*models.Order, error) {

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| begin transaction]: , %w", err)
	}
	// откат транзакции при ошибке в ней
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// найдем сам заказ ( если он есть )
	respOrder, err := r.getOrder(ctx, uid, tx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| get order]: , %w", err)

	}

	// подтянем данные о Delivery ( по delivery_id )
	del, err := r.getDelivery(ctx, respOrder.Delivery.ID, tx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| get delivery]: , %w", err)
	}
	respOrder.Delivery = *del

	// подтянем данные о Payment ( по payment_id )
	pay, err := r.getPayment(ctx, respOrder.Payment.ID, tx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| get payment]: , %w", err)
	}
	respOrder.Payment = *pay

	// подтянем данные о items ( по track_number )
	items, err := r.getItems(ctx, respOrder.TrackNumber, tx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| get items]: , %w", err)
	}
	respOrder.Items = items

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetOrderByUID| commit transaction]: , %w", err)
	}
	return respOrder, nil
}

func (r *OrderPostgresRepository) InsertOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| begin transaction]: , %w", err)
	}
	// откат транзакции при ошибке в ней
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// вставляем данные о Delivery
	delivery_id, err := r.insertDelivery(ctx, tx, &order.Delivery)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| insert delivery in transaction]: , %w", err)

	}
	order.Delivery.ID = delivery_id

	// вставляем данные о Payment
	payment_id, err := r.insertPayment(ctx, tx, &order.Payment)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| insert payment in transaction]: , %w", err)
	}
	order.Payment.ID = payment_id

	// вставляем данные самого заказа
	err = r.insertOrder(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| insert order in transaction]: , %w", err)
	}

	// вставляем данные о Items
	err = r.insertItems(ctx, tx, &order.Items)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| insert items in transaction]: , %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("[InsertOrder| commit transaction]: , %w", err)
	}
	return order, nil
}

func (r *OrderPostgresRepository) insertOrder(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `INSERT INTO "order" ( order_uid,track_number,entry,delivery_id,payment_id,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created, oof_shard) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

	_, err := tx.Exec(ctx, query, order.OrderUID, order.TrackNumber, order.Entry, order.Delivery.ID, order.Payment.ID, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		// проверим, связана ли ошибка с тем что заказ с уникальным полем уже есть
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// проверим по какому полю ошибка уникальности
			// проверять будем по именам constraint, для поля track_number будет присвоего дефолтное имя order_track_number_key
			switch {
			case pgErr.ConstraintName == "order_track_number_key":
				return fmt.Errorf("[insertOrder| exec insert order]: , %w", ErrOrderAlreadyExistsTrack)

			default:
				return fmt.Errorf("[insertOrder| exec insert order]: , %w", ErrOrderAlreadyExistsUUID)
			}
		}

		// если ошибка не вызвана дублированием уникального поля
		return fmt.Errorf("[insertOrder| exec insert order]: , %w", err)
	}

	return nil
}

func (r *OrderPostgresRepository) insertDelivery(ctx context.Context, tx pgx.Tx, del *models.Delivery) (int, error) {
	query := `INSERT INTO delivery(
	name,phone,zip,city,address,region,email) 
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING delivery_id`

	row := tx.QueryRow(ctx, query, del.Name, del.Phone, del.Zip, del.City, del.Address, del.Region, del.Email)

	var delivery_id int

	err := row.Scan(&delivery_id)

	if err != nil {
		return 0, fmt.Errorf("[insertDelivery| exec insert delivery]: , %w", err)
	}

	return delivery_id, nil
}

func (r *OrderPostgresRepository) insertPayment(ctx context.Context, tx pgx.Tx, p *models.Payment) (int, error) {
	query := `INSERT INTO payment(
	transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	RETURNING payment_id`

	row := tx.QueryRow(ctx, query, p.Transaction, p.RequestID, p.Currency, p.Provider, p.Amount, p.PaymentDt, p.Bank, p.DeliveryCost, p.GoodsTotal, p.CustomFee)

	var payment_id int

	err := row.Scan(&payment_id)
	if err != nil {
		return 0, fmt.Errorf("[insertPayment| exec insert payment]: , %w", err)
	}
	return payment_id, nil
}

func (r *OrderPostgresRepository) insertItems(ctx context.Context, tx pgx.Tx, items *[]models.Item) error {
	// для атомарности вставки будем использовать batch
	batch := &pgx.Batch{}

	query := `INSERT INTO item(
	chrtID,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	for _, item := range *items {
		batch.Queue(query, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	}

	res := tx.SendBatch(ctx, batch)
	defer res.Close()

	// проверим вставки на ошибки
	for i := 0; i < batch.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return fmt.Errorf("[insertItems| exec insert %d item]: , %w", i+1, err)
		}
	}
	return nil
}

func (r *OrderPostgresRepository) getOrder(ctx context.Context, id uuid.UUID, tx pgx.Tx) (*models.Order, error) {
	var responceOrder models.Order

	query := `SELECT *
	FROM "order" 
	WHERE order_uid = $1`

	row := tx.QueryRow(ctx, query, id)

	err := row.Scan(&responceOrder.OrderUID, &responceOrder.TrackNumber, &responceOrder.Entry, &responceOrder.Delivery.ID, &responceOrder.Payment.ID, &responceOrder.Locale, &responceOrder.InternalSignature, &responceOrder.CustomerID, &responceOrder.DeliveryService, &responceOrder.Shardkey, &responceOrder.SmID, &responceOrder.DateCreated, &responceOrder.OofShard)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("[getOrder|row scan]: , %w", ErrOrderNotFoundByUUID)
		}
		return nil, fmt.Errorf("[getOrder|row scan]: , %w", err)
	}
	return &responceOrder, nil
}

func (r *OrderPostgresRepository) getDelivery(ctx context.Context, id int, tx pgx.Tx) (*models.Delivery, error) {
	var del models.Delivery

	query := `SELECT name,phone,zip,city,address,region,email
	FROM delivery 
	WHERE delivery_id = $1`

	err := tx.QueryRow(ctx, query, id).Scan(&del.Name, &del.Phone, &del.Zip, &del.City, &del.Address, &del.Region, &del.Email)

	if err != nil {
		return nil, fmt.Errorf("[getDelivery|row scan]: , %w", err)
	}
	return &del, nil
}

func (r *OrderPostgresRepository) getPayment(ctx context.Context, id int, tx pgx.Tx) (*models.Payment, error) {

	var pay models.Payment

	query := `SELECT transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee 
	FROM payment 
	WHERE payment_id = $1`

	err := tx.QueryRow(ctx, query, id).Scan(&pay.Transaction, &pay.RequestID, &pay.Currency, &pay.Provider, &pay.Amount, &pay.PaymentDt, &pay.Bank, &pay.DeliveryCost, &pay.GoodsTotal, &pay.CustomFee)

	if err != nil {
		return nil, fmt.Errorf("[getPayment|row scan]: , %w", err)
	}

	return &pay, nil
}

func (r *OrderPostgresRepository) getItems(ctx context.Context, track string, tx pgx.Tx) ([]models.Item, error) {
	var items []models.Item

	query := `SELECT chrtID,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status FROM item 
	WHERE track_number = $1`

	rows, err := tx.Query(ctx, query, track)
	if err != nil {
		return nil, fmt.Errorf("[getItems|rows scan]: , %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var it models.Item
		err := rows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name, &it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status)
		if err != nil {
			return nil, fmt.Errorf("[getItems|rows scan item]: , %w", err)
		}

		items = append(items, it)
	}

	return items, nil
}

func (r *OrderPostgresRepository) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	orders := []*models.Order{}

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetAllOrders| begin transaction]: , %w", err)
	}
	// откат транзакции при ошибке в ней
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// получим массив заказов, заполним только данными самих заказов
	query := `SELECT * FROM "order"`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("[GetAllOrders|rows scan]: , %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var curOrder models.Order
		err = rows.Scan(&curOrder.OrderUID, &curOrder.TrackNumber, &curOrder.Entry, &curOrder.Delivery.ID, &curOrder.Payment.ID, &curOrder.Locale, &curOrder.InternalSignature, &curOrder.CustomerID, &curOrder.DeliveryService, &curOrder.Shardkey, &curOrder.SmID, &curOrder.DateCreated, &curOrder.OofShard)

		if err != nil {
			return nil, fmt.Errorf("[GetAllOrders|row scan order]: , %w", err)
		}

		orders = append(orders, &curOrder)
	}

	// пройдемся по массиву, для каждого заказа найдем соответствующие delivery,payment и items и подтянем информацию о них
	for _, curOrder := range orders {
		curDelivery, err := r.getDelivery(ctx, curOrder.Delivery.ID, tx)
		if err != nil {
			return nil, fmt.Errorf("[GetAllOrders| get delivery]: , %w", err)
		}
		curOrder.Delivery = *curDelivery

		curPayment, err := r.getPayment(ctx, curOrder.Payment.ID, tx)
		if err != nil {
			return nil, fmt.Errorf("[GetAllOrders| get payment]: , %w", err)
		}
		curOrder.Payment = *curPayment

		curItems, err := r.getItems(ctx, curOrder.TrackNumber, tx)
		if err != nil {
			return nil, fmt.Errorf("[GetAllOrders| get items]: , %w", err)
		}
		curOrder.Items = curItems
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("[GetAllOrders| commit transaction]: , %w", err)
	}
	return orders, nil
}
