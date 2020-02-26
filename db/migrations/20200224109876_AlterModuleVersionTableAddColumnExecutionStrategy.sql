-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE module_version ADD COLUMN execution_strategy SMALLINT NOT NULL DEFAULT 1;



-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back.

ALTER TABLE module_version DROP COLUMN execution_strategy;