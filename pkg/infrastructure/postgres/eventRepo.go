package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"photofinish/pkg/domain/dto"
	"photofinish/pkg/domain/event"
)

type EventRepositoryImpl struct {
	connPool *pgx.ConnPool
}

func NewEventRepository(connPool *pgx.ConnPool) *EventRepositoryImpl {
	u := new(EventRepositoryImpl)
	u.connPool = connPool
	return u
}

func (r *EventRepositoryImpl) CheckExist(pictureId int) error {
	const sql = "SELECT id FROM events WHERE id=$1"
	rows, err := r.connPool.Query(sql, pictureId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	if !rows.Next() {
		return errors.New("event not exist")
	}

	return nil
}

func (r *EventRepositoryImpl) FindAll(page *dto.Page) ([]event.Event, error) {
	const sql = "SELECT id, name, location FROM events LIMIT $1 OFFSET $2;"
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

	events, err := r.scanEvents(rows)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepositoryImpl) scanEvents(rows *pgx.Rows) ([]event.Event, error) {
	var events []event.Event
	var eventItem event.Event
	for rows.Next() {
		err := rows.Scan(
			&eventItem.Id,
			&eventItem.Name,
			&eventItem.Location,
		)
		if err != nil {
			return events, err
		}
		events = append(events, eventItem)
	}
	return events, nil
}

func (r *EventRepositoryImpl) Store(eventDto *event.CreateEventInputDto) (int, error) {

	var data []interface{}
	data = append(data, eventDto.Name, eventDto.Date, eventDto.Location)

	fmt.Println(data)
	lastInsertId := -1
	const sql = "INSERT INTO events (name, date, location) VALUES ($1, $2, $3) RETURNING id"
	row := r.connPool.QueryRow(sql, data...)
	err := row.Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (r *EventRepositoryImpl) Delete(eventId int) error {
	const sql = "DELETE FROM events WHERE id=$1"
	rows, err := r.connPool.Query(sql, eventId)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}
