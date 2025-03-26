--liquibase formatted sql

--changeset samlof:price_history_table
CREATE TABLE price_history(
delivery_start timestamptz  NOT NULL PRIMARY KEY,
delivery_end timestamptz  NOT NULL,
price NUMERIC(8, 2) NOT NULL,
created timestamptz  NOT NULL DEFAULT now()
);

--rollback DROP TABLE price_history
