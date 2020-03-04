-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE employee (
    id             SERIAL PRIMARY KEY NOT NULL,
    external_id    UUID NOT NULL,
    customer_id    UUID NOT NULL,
    created_on     TIMESTAMP NOT NULL,
    updated_on     TIMESTAMP DEFAULT NULL
);