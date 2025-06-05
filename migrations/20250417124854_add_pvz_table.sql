-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    registration_date TIMESTAMP WITH TIME ZONE NOT NULL,
    city VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    pvz_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    type VARCHAR(50) NOT NULL,
    reception_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS pvzuser (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    pass_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL
);

ALTER TABLE reception DROP CONSTRAINT IF EXISTS fk_reception_pvz_id;
ALTER TABLE product DROP CONSTRAINT IF EXISTS fk_product_reception_id;

ALTER TABLE reception ADD CONSTRAINT fk_reception_pvz_id FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE;
ALTER TABLE product ADD CONSTRAINT fk_product_reception_id FOREIGN KEY (reception_id) REFERENCES reception(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
