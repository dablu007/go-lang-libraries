-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE employee (
    id             SERIAL PRIMARY KEY NOT NULL,
    employee_id     TEXT,
    name            TEXT,
    address         TEXT
);