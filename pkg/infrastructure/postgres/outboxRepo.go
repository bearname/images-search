package postgres

import (
	"github.com/col3name/images-search/pkg/common/infrarstructure/db"
	"github.com/col3name/images-search/pkg/domain/broker"
	"github.com/col3name/images-search/pkg/domain/domainerror"
	"github.com/jackc/pgx"
)

type OutboxRepoImpl struct {
	connPool *pgx.ConnPool
}

func NewOutboxRepo(connPool *pgx.ConnPool) *OutboxRepoImpl {
	u := new(OutboxRepoImpl)
	u.connPool = connPool
	return u
}

func (r *OutboxRepoImpl) CheckExist(outboxId string) error {
	const sql = "SELECT id FROM outbox WHERE id=$1"
	rows, err := r.connPool.Query(sql, outboxId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	if !rows.Next() {
		return domainerror.ErrNotExists
	}

	return nil
}

func (r *OutboxRepoImpl) UpdateStatus(outboxId string, status broker.ProcessingStatus) error {
	sql := `UPDATE status = $1, updated_at = now() WHERE id = $2;`
	var data []interface{}
	data = append(data, outboxId, status)
	err := db.WithTransactionSQL(r.connPool, sql, data)
	return err
}

func (r *OutboxRepoImpl) FindNotCompletedOutboxList(limit int) (*[]broker.Outbox, error) {
	sql := "SELECT id FROM outbox WHERE  (status = $1 AND updated_at +  INTERVAL '10 min' <  now()) OR status = $2 LIMIT $3;"
	var data []interface{}
	data = append(data, broker.OutboxProcessing, broker.OutboxNotProcessing, limit)
	rows, err := r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	var res []broker.Outbox
	var item broker.Outbox
	for rows.Next() {
		err = rows.Scan(
			&item.Id,
			&item.BrokerTopic,
			&item.BrokerKey,
			&item.BrokerTopic,
			&item.UpdatedAt,
			&item.Status,
		)
		if err != nil {
			return &res, err
		}
	}
	return &res, nil
}
