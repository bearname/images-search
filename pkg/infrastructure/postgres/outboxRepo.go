package postgres

import (
	"github.com/jackc/pgx"
	"photofinish/pkg/common/infrarstructure/db"
	"photofinish/pkg/domain/broker"
	"photofinish/pkg/domain/errors"
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
		return errors.ErrNotExists
	}

	return nil
}

//func (r *OutboxRepoImpl) Store(o *broker.Outbox) error {
//    sql := `INSERT INTO outbox (id, broker_topic, broker_key, broker_value, updated_at, status) VALUES ($1,$2,$3,$4, $5, $6)
//ON CONFLICT (id) DO UPDATE
//    SET broker_topic = $7, broker_key= $8, broker_value= $8, updated_at = now(), status = $9, value ;`
//    var data []interface{}
//    data = append(data, o.Id, o.BrokerTopic, o.BrokerKey, o.BrokerValue, o.UpdatedAt, o.Status,
//        o.BrokerTopic, o.BrokerKey, o.BrokerValue, o.Status)
//    err := db.WithTransactionSQL(r.connPool, sql, data)
//    return err
//}

func (r *OutboxRepoImpl) UpdateStatus(outboxId string, status broker.ProcessingStatus) error {
	sql := `UPDATE status = $1, updated_at = now() WHERE id = $2;`
	var data []interface{}
	data = append(data, outboxId, status)
	err := db.WithTransactionSQL(r.connPool, sql, data)
	return err
}

//func (r *OutboxRepoImpl)  Store(imageTextDetectionDto *Outbox) (int, error) {
//
//}
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

//
//func (r *OutboxRepoImpl) Delete(outboxId string) error {
//   const sql = "DELETE FROM outbox WHERE id=$1"
//   var s []interface{}
//   s = append(s, outboxId)
//   return db.WithTransactionSQL(r.connPool, sql, s)
//}

//func (r *TasksRepoImpl) Store(e *tasks.AddImageDto) error {
//    sql := `INSERT INTO tasks (id, dropbox_path, eventid) VALUES ($1,$2,$3)
//  INSERT INTO outbox (id, broker_topic, broker_key, broker_value) VALUES ($1,$2,$3,$4)`
//    var data []interface{}
//    outboxId := uuid.Generate().String()
//    t := e.Task
//    data = append(data, t.Id, t.DropboxPath, t.EventId, outboxId, e.BrokerTopic, e.BrokerTopic, e.TaskData)
//    return db.WithTransactionSQL(r.connPool, sql, data)
//}
//
//func (r *TasksRepoImpl) GetStatsByTask(taskId string) (*tasks.TaskStats, error) {
//    sql := r.getPictureProcessingStatusSQL()
//    var data []interface{}
//    data = append(data, taskId)
//    i, err := db.Query(r.connPool, sql, data, r.scanTaskStatistic())
//    if err != nil {
//        return nil, err
//    }
//    taskStatus := i.(tasks.TaskStats)
//
//    sql = r.getStatisticTimeSQL()
//
//    data = make([]interface{}, 0)
//    data = append(data, taskId)
//
//    rows, err := r.connPool.Query(sql, data...)
//    if err != nil {
//        return nil, err
//    }
//    defer rows.Close()
//    if rows.Err() != nil {
//        return nil, rows.Err()
//    }
//    stats, err := r.scanTaskStatTime(rows, taskStatus)
//    if err != nil {
//        return stats, err
//    }
//    return &taskStatus, nil
//}
//
//func (r *TasksRepoImpl) scanTaskStatTime(rows *pgx.Rows, taskStatus tasks.TaskStats) (*tasks.TaskStats, error) {
//    if rows.Next() {
//        err := rows.Scan(
//            &taskStatus.StartedAt,
//            &taskStatus.LastUpdatedAt,
//        )
//        if err != nil {
//            return &taskStatus, err
//        }
//    }
//    return nil, nil
//}
//
//func (r *TasksRepoImpl) scanTaskStatistic() func(rows *pgx.Rows) (interface{}, error) {
//    return func(rows *pgx.Rows) (interface{}, error) {
//        var taskStatus tasks.TaskStats
//        var item tasks.TaskStatsItem
//        for rows.Next() {
//            err := rows.Scan(
//                &item.Status,
//                &item.Count,
//            )
//            if err != nil {
//                return &taskStatus, err
//            }
//            taskStatus.Stats = append(taskStatus.Stats, item)
//        }
//        return taskStatus, nil
//    }
//}
//
//func (r *TasksRepoImpl) GetTaskList(page *dto.Page) (*[]tasks.TaskReturnDTO, error) {
//    sql := r.getTaskListSQL()
//    var data []interface{}
//    data = append(data, page.Limit, page.Offset)
//    rows, err := r.connPool.Query(sql, data...)
//    if err != nil {
//        return nil, err
//    }
//    defer rows.Close()
//    if rows.Err() != nil {
//        return nil, rows.Err()
//    }
//    return r.scanTaskList(rows, err)
//}
//
//func (r *TasksRepoImpl) scanTaskList(rows *pgx.Rows, err error) (*[]tasks.TaskReturnDTO, error) {
//    var result []tasks.TaskReturnDTO
//    var task tasks.TaskReturnDTO
//    var avgRating float64
//    for rows.Next() {
//        err = rows.Scan(
//            &task.Id,
//            &task.CountImages,
//            &avgRating,
//            &task.StartedAt,
//            &task.LastUpdateAt,
//        )
//        if err != nil {
//            return &result, err
//        }
//        task.IsCompleted = !(avgRating > 0)
//
//        result = append(result, task)
//    }
//    return &result, err
//}
//
//func (r *TasksRepoImpl) getTaskListSQL() string {
//    return `SELECT task_id, p.count_images, avg_rating, started_at, last_update
//                      FROM (
//                               SELECT AVG(processing_status) as avg_rating, MAX(update_at) as last_update, COUNT(id) AS count_images, task_id
//                               FROM pictures
//                               GROUP BY task_id
//                           ) AS p
//                               INNER JOIN tasks ON task_id = tasks.id
//ORDER BY started_at
//LIMIT $1 OFFSET $2;`
//
//}
//
//func (r *TasksRepoImpl) getPictureProcessingStatusSQL() string {
//    return "SELECT processing_status, count(id) FROM pictures WHERE task_id = $1 GROUP BY processing_status;"
//}
//
//func (r *TasksRepoImpl) getStatisticTimeSQL() string {
//    return `SELECT t.started_at, MAX(update_at) AS last_updated_at
//			FROM pictures p
//				LEFT JOIN tasks t ON p.task_id = t.id
//			WHERE task_id = $1 AND processing_status = 0
//			GROUP BY t.count_images, t.started_at;`
//}
//
