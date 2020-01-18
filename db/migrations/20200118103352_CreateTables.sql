-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE flows(
   	id             SERIAL PRIMARY KEY NOT NULL,
	external_id		UUID NOT NULL UNIQUE,
	name           VARCHAR NOT NULL,
	version        VARCHAR NOT NULL,
	type           VARCHAR NOT NULL,
	status         VARCHAR NOT NULL,
	flow_context   JSON NOT NULL,
	module_versions JSON NOT NULL,
	created_on      TIMESTAMP NOT NULL,
	deleted_on      TIMESTAMP DEFAULT NULL
);

CREATE TABLE fields(
 	id        SERIAL PRIMARY KEY NOT NULL,
	name      VARCHAR NOT NULL,
	tenant_id UUID UNIQUE DEFAULT NULL, 
	created_on TIMESTAMP NOT NULL,
	deleted_on TIMESTAMP DEFAULT NULL
);

CREATE TABLE field_versions(
   	id        SERIAL PRIMARY KEY NOT NULL,
	Name VARCHAR NOT NULL,
	external_id UUID NOT NULL UNIQUE,
	field_id INT NOT NULL,
	is_visible BOOLEAN NOT NULL,
	version VARCHAR NOT NULL,
	properties JSON NOT NULL,
	created_on TIMESTAMP NOT NULL,
	deleted_on TIMESTAMP DEFAULT NULL
);

CREATE TABLE sections(
   	id        SERIAL PRIMARY KEY NOT NULL,
	name      VARCHAR NOT NULL,
	tenant_id  UUID UNIQUE DEFAULT NULL,
	created_on TIMESTAMP NOT NULL,
	deleted_on TIMESTAMP DEFAULT NULL
);

CREATE TABLE section_versions(
 	id        SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR NOT NULL,
	external_id UUID NOT NULL UNIQUE,
	section_id INT NOT NULL,
	is_visible BOOLEAN NOT NULL,
	version VARCHAR NOT NULL,
	properties JSON NOT NULL,
	field_versions JSON NOT NULL,
	created_on TIMESTAMP NOT NULL,
	deleted_on  TIMESTAMP DEFAULT NULL
);

CREATE TABLE modules(
 	id        SERIAL PRIMARY KEY NOT NULL,
	name      VARCHAR NOT NULL,
	status    VARCHAR NOT NULL,
	tenant_id  UUID UNIQUE DEFAULT NULL,
	created_on TIMESTAMP NOT NULL,
	deleted_on TIMESTAMP DEFAULT NULL
);

CREATE TABLE module_versions(
 	id        SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR NOT NULL,
	module_id        INT NOT NULL,
	external_id      UUID NOT NULL UNIQUE,
	version         VARCHAR NOT NULL,
	created_on       TIMESTAMP NOT NULL,
	deleted_on       TIMESTAMP DEFAULT NULL,
	properties      JSON,
	section_versions JSON
);

ALTER TABLE module_versions 
ADD CONSTRAINT FK_moduleVersions_moduleId FOREIGN KEY (module_id) REFERENCES modules (id);

ALTER TABLE field_versions 
ADD CONSTRAINT FK_fieldVersions_fieldId FOREIGN KEY (field_id) REFERENCES fields (id);

ALTER TABLE section_versions 
ADD CONSTRAINT FK_sectionVersions_sectionId FOREIGN KEY (section_id) REFERENCES sections (id);