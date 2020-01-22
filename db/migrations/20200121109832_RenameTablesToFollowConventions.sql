-- +goose Up
-- SQL in this section is executed when the migration is applied.

alter table flows rename to flow;

alter table modules rename to module;

alter table module_versions rename to module_version;

alter table sections rename to section;

alter table section_versions rename to section_version;

alter table fields rename to field;

alter table field_versions rename to field_version;
