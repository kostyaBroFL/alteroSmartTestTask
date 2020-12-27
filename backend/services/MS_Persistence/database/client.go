package main

import (
	"context"
	"fmt"
	// _ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	id, err := NewClient().GetDeviceId(context.Background(), &GetDeviceIdReq{
		Name: "dsf",
	})
	fmt.Printf("id %d err %+v\n", id, err)

	id, err = NewClient().GetDeviceId(context.Background(), &GetDeviceIdReq{
		Name: "dsf",
	})
	fmt.Printf("id %d err %+v\n", id, err)

	err = NewClient().InsertDeviceData(context.Background(), &InsertDeviceDataReq{
		DeviceId:         1,
		Data:             23423423,
		TimestampSeconds: 234,
		TimestampNanos:   32423,
	})
	fmt.Printf("err %+v\n", err)
	err = NewClient().InsertDeviceData(context.Background(), &InsertDeviceDataReq{
		DeviceId:         1,
		Data:             3423657,
		TimestampSeconds: 234,
		TimestampNanos:   32423,
	})
	fmt.Printf("err %+v\n", err)
}

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

func NewClient() *client {
	db, err := sqlx.Connect("postgres", "host=localhost user=postgres password=qwerty dbname=ms_persistence sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
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
