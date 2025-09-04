CREATE TABLE IF NOT EXISTS delivery
(
	delivery_id SERIAL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	phone VARCHAR(32) NOT NULL,
	zip VARCHAR(255) NOT NULL,
	city VARCHAR(255) NOT NULL,
	address VARCHAR(255) NOT NULL,
	region VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL
);

BEGIN TRANSACTION;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS payment
(
	payment_id SERIAL PRIMARY KEY,
	"transaction" UUID NOT NULL,
	request_id TEXT NOT NULL,
	currency VARCHAR(32) NOT NULL,
	provider VARCHAR(255) NOT NULL,
	amount INT NOT NULL,
	payment_dt BIGINT NOT NULL,
	bank VARCHAR(255) NOT NULL,
	delivery_cost INT NOT NULL,
	goods_total INT NOT NULL,
	custom_fee INT NOT NULL
);

COMMIT;

CREATE TABLE IF NOT EXISTS "order" 
(
	order_uid UUID PRIMARY KEY,
	track_number TEXT UNIQUE NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_id INT NOT NULL,
    payment_id INT NOT NULL,
    locale VARCHAR(32) NOT NULL,
    internal_signature TEXT,
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(255) NOT NULL,
	FOREIGN KEY (delivery_id) REFERENCES delivery (delivery_id),
    FOREIGN KEY (payment_id) REFERENCES payment (payment_id)  
);

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS item
(
	item_id SERIAL PRIMARY KEY,
	chrtID BIGINT NOT NULL,
	track_number TEXT NOT NULL,
    price int NOT NULL,
    rid TEXT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    sale INT NOT NULL,
    "size" VARCHAR(32) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INT NOT NULL,
	FOREIGN KEY (track_number) REFERENCES "order"(track_number) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_item_track_number ON item(track_number);

COMMIT;
