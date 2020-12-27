create table device (
 id serial4 not null primary key,
 name varchar(128) not null unique
);

create table device_data (
 id serial8 not null primary key,
 device_id int4 not null,
 data float8 not null,
 timestamp_seconds int8 not null,
 timestamp_nanos int4 not null,
 foreign key (device_id) references device(id)
);
