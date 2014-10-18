use symdb;

drop table if exists user;
drop table if exists token;
drop table if exists aniversary;

create table user (
  id BIGINT(20) not null primary key auto_increment,
  name VARCHAR(64) unique NOT NULL,
  password VARCHAR(128) NOT NULL,
  token VARCHAR(128) unique NOT NULL
) ENGINE=InnoDB;

create table token (
  token VARCHAR(255) NOT NULL PRIMARY KEY,
  secret VARCHAR(256) NOT NULL
) ENGINE=InnoDB;

create table aniversary (
  id varchar(64) not null primary key,
  prefix varchar(256) not null,
  message varchar(256) not null,
  users varchar(1024) not null,
  embed varchar(2048)
) ENGINE=InnoDB;