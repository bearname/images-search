package postgres

import (
	"github.com/jackc/pgx"
	"photofinish/pkg/domain/tasks"
)

type TasksRepositoryImpl struct {
	connPool *pgx.ConnPool
}

func NewTasksRepositoryImpl(connPool *pgx.ConnPool) *TasksRepositoryImpl {
	u := new(TasksRepositoryImpl)
	u.connPool = connPool
	return u
}

func (r *TasksRepositoryImpl) GetStatsByTask(taskId string) (*tasks.TaskStats, error) {
	sql := "SELECT processing_status, count(id) FROM pictures WHERE task_id = $1 GROUP BY processing_status;"
	var data []interface{}
	data = append(data, taskId)
	rows, err := r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var taskStatus tasks.TaskStats
	var item tasks.TaskStatsItem
	for rows.Next() {
		err = rows.Scan(
			&item.Status,
			&item.Count,
		)
		if err != nil {
			return &taskStatus, err
		}
		taskStatus.Stats = append(taskStatus.Stats, item)
	}

	sql = `select t.started_at, max(update_at) as last_updated_at
			from pictures p
					 left join tasks t on p.task_id = t.id
			where task_id = $1
			  and processing_status = 0
			group by t.count_images, t.started_at;`

	data = nil
	data = append(data, taskId)
	rows, err = r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if rows.Next() {
		err = rows.Scan(
			&taskStatus.StartedAt,
			&taskStatus.LastUpdatedAt,
		)
		if err != nil {
			return &taskStatus, err
		}
	}
	return &taskStatus, nil
}
