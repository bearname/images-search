package postgres

import (
	"fmt"
	"github.com/jackc/pgx"
	"photofinish/pkg/domain/order"
	"strconv"
)

type OrderRepositoryImpl struct {
	connPool *pgx.ConnPool
}

func NewOrderRepository(connPool *pgx.ConnPool) *OrderRepositoryImpl {
	u := new(OrderRepositoryImpl)
	u.connPool = connPool
	return u
}

func (r *OrderRepositoryImpl) Store(order *order.CreateOrderDTO) error {
	sql := `INSERT INTO orders (id, userid, totalprice) VALUES ($1, $2, $3);
INSERT INTO picture_in_order (orderId, pictureId) VALUES `
	var data []interface{}
	order.UserId = 7
	orderId := order.OrderId.String()
	data = append(data, orderId, order.UserId, order.TotalPrice)
	pictures := len(order.Data)
	j := 4
	for i, price := range order.Data {
		sql += "($" + strconv.Itoa(j) + ", $" + strconv.Itoa(j+1) + ")"
		if i != pictures-1 {
			sql += ","
		} else {
			sql += ";"
		}
		j += 2
		data = append(data, orderId, price.PictureID)
	}

	//picturesSql += `($` + strconv.Itoa(pictureI) + `, $` + strconv.Itoa(pictureI+1) + `, $` + strconv.Itoa(pictureI+2) + `, $` + strconv.Itoa(pictureI+3) + `) RETURNING ID;`
	//id := uuid.Generate().String()
	//data = append(data, id, image.DropboxPath, image.EventId)
	fmt.Println(sql)
	//return db.WithTransactionSQL(r.connPool, sql, data)
	tx, err := r.connPool.Begin()
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}
	_, err = tx.Exec(sql, data...)
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}

	return tx.Commit()
}

func (r *OrderRepositoryImpl) UpdateStatus(order *order.UpdateOrderStatusDTO) error {
	const sql = "UPDATE orders SET status = $1 WHERE id = $2;"
	var data []interface{}
	data = append(data, order.Status, order.OrderId)
	tx, err := r.connPool.Begin()
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}
	_, err = tx.Exec(sql, data...)
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}

	return tx.Commit()
	//return db.WithTransactionSQL(r.connPool, sql, data)
}

func (r *OrderRepositoryImpl) GetOrder(dto *order.GetOrderDTO) (*order.ReturnOrderDTO, error) {
	const sql = `SELECT orders.id, userid, status, totalprice, pay_id, receipt_url, receipt_number, count(pio.orderId) AS countPictures
FROM orders
         LEFT JOIN picture_in_order pio on orders.id = pio.orderId
WHERE orders.id = $1
  AND userid = $2
GROUP BY orders.id;`
	var data []interface{}
	data = append(data, dto.OrderId, dto.UserId)
	rows, err := r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	var returnDto order.ReturnOrderDTO
	if rows.Next() {
		err = rows.Scan(&returnDto.OrderId,
			&returnDto.UserId,
			&returnDto.Status,
			&returnDto.TotalPrice,
			&returnDto.PayId,
			&returnDto.ReceiptUrl,
			&returnDto.ReceiptNumber,
			&returnDto.CountPictures)
		if err != nil {
			return nil, err
		}
	}
	return &returnDto, nil
}

func (r *OrderRepositoryImpl) SavePayResult(payResult *order.PayResultDTO) error {
	const sql = "UPDATE orders SET pay_id = $1, receipt_url = $2, receipt_number = $3 WHERE id = $4"
	var data []interface{}
	data = append(data, payResult.ID, payResult.ReceiptURL, payResult.ReceiptNumber, payResult.OrderId)
	tx, err := r.connPool.Begin()
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}
	_, err = tx.Exec(sql, data...)
	if err != nil {
		if tx != nil {
			return tx.Rollback()
		}
		return err
	}

	return tx.Commit()
	//return db.WithTransactionSQL(r.connPool, sql, data)
}
