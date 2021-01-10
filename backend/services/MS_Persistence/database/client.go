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
    device_data(device_id, data, timestamp) 
    VALUES (:device_id, :data, :timestamp);`

var getDeviceDataByDeviceName = `select * from device_data 
	where device_id = (select id from device where name = :device_name) 
	order by timestamp desc 
	limit :limit;`

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

type GetDataByDeviceNameRequest struct {
	DeviceName string `db:"device_name"`
	Limit      int32  `db:"limit"`
}

type DeviceData struct {
	id        int32   `db:"id"`
	deviceId  int32   `db:"device_id"`
	Data      float64 `db:"data"`
	Timestamp string  `db:"timestamp"`
}

func (c *client) GetDataByDeviceName(
	ctx context.Context,
	req *GetDataByDeviceNameRequest,
) ([]*DeviceData, error) {
	rows, err := c.db.NamedQueryContext(ctx, getDeviceDataByDeviceName, req)
	if err != nil {
		return nil, err
	}
	var output []*DeviceData
	defer rows.Close() // todo err & search olso
	for rows.Next() {
		var data *DeviceData
		err := rows.StructScan(data)
		if err != nil {
			return nil, err
		}
		output = append(output, data)
	}
	return output, nil
}
