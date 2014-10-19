use symdb;

drop table if exists user;

create table user (
  id BIGINT(20) not null primary key auto_increment,
  name VARCHAR(64) unique NOT NULL,
  password VARCHAR(128) NOT NULL,
  token VARCHAR(128) unique NOT NULL
) ENGINE=InnoDB;
