-- +migrate Up

CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- +migrate Down

DROP TABLE products;
