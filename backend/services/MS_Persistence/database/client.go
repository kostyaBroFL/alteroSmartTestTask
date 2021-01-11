package database

import (
	"context"

	// _ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	log_context "alteroSmartTestTask/common/log/context"
)

var getDeviceQuery = `insert into device(name) 
	values (:device_name) 
	on conflict(name) do update 
	set name = :device_name returning id;`

var insertDeviceDataQuery = `insert into 
    device_data(device_id, data, timestamp) 
    VALUES (:device_id, :data, :timestamp);`

var getDeviceDataByDeviceName = `select id, device_id, data, timestamp 
	from device_data 
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
	DeviceId  int64   `db:"device_id"`
	Data      float64 `db:"data"`
	Timestamp string  `db:"timestamp"`
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
	Id        int32   `db:"id"`
	DeviceId  int32   `db:"device_id"`
	Data      float64 `db:"data"`
	Timestamp string  `db:"timestamp"`
}

func (c *client) GetDataByDeviceName(
	ctx context.Context,
	req *GetDataByDeviceNameRequest,
) ([]*DeviceData, error) {
	ctx = log_context.WithLogger(ctx, log_context.FromContext(ctx).
		WithField("db_method", "GetDataByDeviceName").
		WithField("device_name", req.DeviceName).
		WithField("limit", req.Limit))
	rows, err := c.db.NamedQueryContext(ctx, getDeviceDataByDeviceName, req)
	if err != nil {
		log_context.FromContext(ctx).WithError(err).
			Error("get device data error")
		return nil, err
	}
	var output []*DeviceData
	defer func() {
		if err = rows.Close(); err != nil {
			log_context.FromContext(ctx).WithError(err).
				Error("can not to close row")
		}
	}()
	for rows.Next() {
		data := &DeviceData{}
		if err := rows.StructScan(data); err != nil {
			log_context.FromContext(ctx).WithError(err).
				Error("can not scan device data")
			return nil, err
		}
		output = append(output, data)
	}
	log_context.FromContext(ctx).Info("db success")
	return output, nil
}
