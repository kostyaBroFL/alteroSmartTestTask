package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	logcontext "alteroSmartTestTask/common/log/context"
)

var (
	getDeviceQuery = `insert into device(name) 
		values (:device_name) 
		on conflict(name) do update 
		set name = :device_name returning id;`

	insertDeviceDataQuery = `insert into 
		device_data(device_id, data, timestamp) 
		VALUES (:device_id, :data, :timestamp);`

	getDeviceDataByDeviceName = `select id, device_id, data, timestamp 
		from device_data 
		where device_id = (select id from device where name = :device_name) 
		order by timestamp desc 
		limit :limit;`
)

type client struct {
	db *sqlx.DB
}

// Close is the method for turning off database client.
func (c *client) Close() error {
	return c.db.Close()
}

// NewClient is the constructor for the database client.
func NewClient(db *sqlx.DB) *client {
	return &client{db: db}
}

// GetDeviceIdReq it the request for GetDeviceId method.
type GetDeviceIdReq struct {
	Name string `db:"device_name"`
}

// GetDeviceId is the method for find identifier if the device by its name.
func (c *client) GetDeviceId(
	ctx context.Context, req *GetDeviceIdReq,
) (int64, error) {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("db_method", "GetDeviceId").
		WithField("device_name", req.Name))
	rows, err := c.db.NamedQuery(getDeviceQuery, req)
	if err != nil {
		logcontext.FromContext(ctx).WithError(err).
			Error("can not to make get device query")
		return 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logcontext.FromContext(ctx).WithError(err).
				Error("can not to close db request")
		}
	}()
	rows.Next()
	var id int64
	if err = rows.Scan(&id); err != nil {
		logcontext.FromContext(ctx).WithError(err).
			Error("can not to scan device id")
		return 0, err
	}
	logcontext.FromContext(ctx).Info("db success")
	return id, nil
}

// InsertDeviceDataReq it the request for InsertDeviceData method.
type InsertDeviceDataReq struct {
	DeviceId  int64   `db:"device_id"`
	Data      float64 `db:"data"`
	Timestamp string  `db:"timestamp"`
}

// InsertDeviceData is the method for inserting the chunk of the data.
func (c *client) InsertDeviceData(
	ctx context.Context, req *InsertDeviceDataReq,
) error {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("db_method", "GetDeviceId").
		WithField("device_id", req.DeviceId).
		WithField("data", req.Data).
		WithField("timestamp", req.Timestamp))
	rows, err := c.db.NamedQueryContext(ctx, insertDeviceDataQuery, req)
	if err != nil {
		logcontext.FromContext(ctx).WithError(err).
			Error("can not to insert device data")
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logcontext.FromContext(ctx).WithError(err).
				Error("can not to close db query")
		}
	}()
	logcontext.FromContext(ctx).Info("db success")
	return nil
}

// GetDataByDeviceNameRequest it the request for GetDataByDeviceName method.
type GetDataByDeviceNameRequest struct {
	DeviceName string `db:"device_name"`
	Limit      int32  `db:"limit"`
}

// DeviceData it the response of the GetDataByDeviceName method.
type DeviceData struct {
	Id        int32   `db:"id"`
	DeviceId  int32   `db:"device_id"`
	Data      float64 `db:"data"`
	Timestamp string  `db:"timestamp"`
}

// GetDataByDeviceName ix the method for for getting data of specific device.
func (c *client) GetDataByDeviceName(
	ctx context.Context,
	req *GetDataByDeviceNameRequest,
) ([]*DeviceData, error) {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("db_method", "GetDataByDeviceName").
		WithField("device_name", req.DeviceName).
		WithField("limit", req.Limit))
	rows, err := c.db.NamedQueryContext(ctx, getDeviceDataByDeviceName, req)
	if err != nil {
		logcontext.FromContext(ctx).WithError(err).
			Error("get device data error")
		return nil, err
	}
	var output []*DeviceData
	defer func() {
		if err = rows.Close(); err != nil {
			logcontext.FromContext(ctx).WithError(err).
				Error("can not to close db query")
		}
	}()
	for rows.Next() {
		data := &DeviceData{}
		if err := rows.StructScan(data); err != nil {
			logcontext.FromContext(ctx).WithError(err).
				Error("can not scan device data")
			return nil, err
		}
		output = append(output, data)
	}
	logcontext.FromContext(ctx).Info("db success")
	return output, nil
}
