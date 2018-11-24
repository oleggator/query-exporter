create table people (
  id serial not null primary key,
  name varchar(256) not null,
  lastname varchar(256) not null,
  birthday date not null,
  some_flag integer not null,
  created timestamp(0) not null default current_timestamp
);

create table cities (
  id serial not null primary key,
  name varchar(256) not null,
  country_id integer not null,
  created timestamp(0) not null default current_timestamp
);

create table countries (
  id serial not null primary key,
  name varchar(256) not null,
  created timestamp(0) not null default current_timestamp
);
