-- +migrate Up

CREATE TABLE package_sizes (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    size INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE (product_id, size)
);
    
-- +migrate Down

DROP TABLE package_sizes;
