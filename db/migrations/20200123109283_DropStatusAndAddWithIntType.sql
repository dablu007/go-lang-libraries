-- +goose Up
-- SQL in this section is executed when the migration is applied.

alter table flow drop column status ;

alter table flow add column status smallint not null default 1;

alter table flow drop column type ;

alter table flow add column type smallint not null default 1;

alter table module drop column status ;

alter table module add column status smallint not null default 1;

alter table section drop column status ;

alter table section add column status smallint not null default 1;

alter table field drop column status ;

alter table field add column status smallint not null default 1;

