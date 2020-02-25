-- +goose Up
-- SQL in this section is executed when the migration is applied.

alter table module_version add column execution_strategy smallint not null default 1;