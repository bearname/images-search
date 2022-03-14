package postgres

import (
	"github.com/col3name/images-search/pkg/common/infrarstructure/db"
	"github.com/col3name/images-search/pkg/common/util/uuid"
	"github.com/col3name/images-search/pkg/domain/dto"
	"github.com/col3name/images-search/pkg/domain/tasks"
	"github.com/jackc/pgx"
)

type TasksRepoImpl struct {
	connPool *pgx.ConnPool
}

func NewTasksRepo(connPool *pgx.ConnPool) *TasksRepoImpl {
	u := new(TasksRepoImpl)
	u.connPool = connPool
	return u
}

func (r *TasksRepoImpl) Store(e *tasks.AddImageDto) error {
	sql := `INSERT INTO tasks (id, dropbox_path, eventid) VALUES ($1,$2,$3)
  INSERT INTO outbox (id, broker_topic, broker_key, broker_value) VALUES ($1,$2,$3,$4)`
	var data []interface{}
	outboxId := uuid.Generate().String()
	t := e.Task
	data = append(data, t.Id, t.DropboxPath, t.EventId, outboxId, e.BrokerTopic, e.BrokerTopic, e.TaskData)
	return db.WithTransactionSQL(r.connPool, sql, data)
}

func (r *TasksRepoImpl) GetStatsByTask(taskId string) (*tasks.TaskStats, error) {
	sql := r.getPictureProcessingStatusSQL()
	var data []interface{}
	data = append(data, taskId)
	i, err := db.Query(r.connPool, sql, data, r.scanTaskStatistic())
	if err != nil {
		return nil, err
	}
	taskStatus := i.(tasks.TaskStats)

	sql = r.getStatisticTimeSQL()

	data = make([]interface{}, 0)
	data = append(data, taskId)

	rows, err := r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	stats, err := r.scanTaskStatTime(rows, taskStatus)
	if err != nil {
		return stats, err
	}
	return &taskStatus, nil
}

func (r *TasksRepoImpl) scanTaskStatTime(rows *pgx.Rows, taskStatus tasks.TaskStats) (*tasks.TaskStats, error) {
	if rows.Next() {
		err := rows.Scan(
			&taskStatus.StartedAt,
			&taskStatus.LastUpdatedAt,
		)
		if err != nil {
			return &taskStatus, err
		}
	}
	return nil, nil
}

func (r *TasksRepoImpl) scanTaskStatistic() func(rows *pgx.Rows) (interface{}, error) {
	return func(rows *pgx.Rows) (interface{}, error) {
		var taskStatus tasks.TaskStats
		var item tasks.TaskStatsItem
		for rows.Next() {
			err := rows.Scan(
				&item.Status,
				&item.Count,
			)
			if err != nil {
				return &taskStatus, err
			}
			taskStatus.Stats = append(taskStatus.Stats, item)
		}
		return taskStatus, nil
	}
}

func (r *TasksRepoImpl) GetTaskList(page *dto.Page) (*[]tasks.TaskReturnDTO, error) {
	sql := r.getTaskListSQL()
	var data []interface{}
	data = append(data, page.Limit, page.Offset)
	rows, err := r.connPool.Query(sql, data...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return r.scanTaskList(rows)
}

func (r *TasksRepoImpl) scanTaskList(rows *pgx.Rows) (*[]tasks.TaskReturnDTO, error) {
	var result []tasks.TaskReturnDTO
	var task tasks.TaskReturnDTO
	var avgRating float64
	var err error
	for rows.Next() {
		err = rows.Scan(
			&task.Id,
			&task.CountImages,
			&avgRating,
			&task.StartedAt,
			&task.LastUpdateAt,
		)
		if err != nil {
			return &result, err
		}
		task.IsCompleted = !(avgRating > 0)

		result = append(result, task)
	}
	return &result, err
}

func (r *TasksRepoImpl) getTaskListSQL() string {
	return `SELECT task_id, p.count_images, avg_rating, started_at, last_update
                      FROM (
                               SELECT AVG(processing_status) as avg_rating, MAX(update_at) as last_update, COUNT(id) AS count_images, task_id
                               FROM pictures
                               GROUP BY task_id
                           ) AS p
                               INNER JOIN tasks ON task_id = tasks.id
ORDER BY started_at
LIMIT $1 OFFSET $2;`

}

func (r *TasksRepoImpl) getPictureProcessingStatusSQL() string {
	return "SELECT processing_status, count(id) FROM pictures WHERE task_id = $1 GROUP BY processing_status;"
}

func (r *TasksRepoImpl) getStatisticTimeSQL() string {
	return `SELECT t.started_at, MAX(update_at) AS last_updated_at
			FROM pictures p
				LEFT JOIN tasks t ON p.task_id = t.id
			WHERE task_id = $1 AND processing_status = 0
			GROUP BY t.count_images, t.started_at;`
}
