CREATE TABLE users(
                      id serial not null unique,
                      name varchar(255) not null,
                      username varchar(255) not null unique,
                      password_hash varchar(255) not null
);

CREATE TABLE balance_list (
                              id serial not null unique,
                              user_id integer not null,
                              balance decimal(10, 2) not null,
                              foreign key (user_id) references users(id)
);

CREATE TABLE transaction_pattern (
                              id serial not null unique,
                              user_id integer not null,
                              amount float,
                              description  varchar(255),
                              date timestamp
);