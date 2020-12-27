package database

import (
	"context"
	// _ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var getDeviceQuery = `insert into device(name) 
	values (:device_name) 
	on conflict(name) do update 
	set name = :device_name returning id;`

var insertDeviceDataQuery = `insert into 
    device_data(device_id, data, timestamp_seconds, timestamp_nanos) 
    VALUES (:device_id, :data, :timestamp_seconds, :timestamp_nanos);`

type client struct {
	db *sqlx.DB
}

func (c *client) Close() error {
	return c.db.Close()
}

func NewClient(db *sqlx.DB) *client {
	return &client{db: db}
}

type GetDeviceIdReq struct {
	Name string `db:"device_name"`
}

func (c *client) GetDeviceId(ctx context.Context, req *GetDeviceIdReq) (int64, error) {
	rows, err := c.db.NamedQuery(getDeviceQuery, req)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

type InsertDeviceDataReq struct {
	DeviceId         int64   `db:"device_id"`
	Data             float64 `db:"data"`
	TimestampSeconds int64   `db:"timestamp_seconds"`
	TimestampNanos   int32   `db:"timestamp_nanos"`
}

func (c *client) InsertDeviceData(ctx context.Context, req *InsertDeviceDataReq) error {
	rows, err := c.db.NamedQueryContext(ctx, insertDeviceDataQuery, req)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
