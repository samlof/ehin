-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS price_history(
    delivery_start timestamptz  NOT NULL PRIMARY KEY,
    delivery_end timestamptz  NOT NULL,
    price NUMERIC(8, 2) NOT NULL,
    created timestamptz  NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS price_history;
-- +goose StatementEnd
