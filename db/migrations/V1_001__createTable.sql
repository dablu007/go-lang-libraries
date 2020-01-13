CREATE TABLE Flows(
   	Id             SERIAL PRIMARY KEY NOT NULL,
	Name           VARCHAR NOT NULL,
	Version        VARCHAR NOT NULL,
	Type           VARCHAR NOT NULL,
	Status         VARCHAR NOT NULL,
	MerchantId     UUID NOT NULL,
	ModuleVersions JSON,
	CreatedOn      TIMESTAMP NOT NULL,
	DeletedOn      TIMESTAMP
);

CREATE TABLE Fields(
 	Id        SERIAL PRIMARY KEY NOT NULL,
	Name      VARCHAR NOT NULL,
	SectionId SMALLINT NOT NULL,
	IsVisible BOOLEAN NOT NULL,
	Version   VARCHAR NOT NULL,
	CreatedOn TIMESTAMP NOT NULL,
	DeletedOn TIMESTAMP
);

CREATE TABLE FieldVersions(
   	Id        SERIAL PRIMARY KEY NOT NULL,
	Name VARCHAR NOT NULL,
	ExternalId VARCHAR NOT NULL,
	FieldId SMALLINT NOT NULL,
	IsVisible BOOLEAN NOT NULL,
	Version VARCHAR NOT NULL,
	CreatedOn TIMESTAMP NOT NULL,
	DeletedOn TIMESTAMP
);

CREATE TABLE Sections(
   	Id        SERIAL PRIMARY KEY NOT NULL,
	Name      VARCHAR NOT NULL,
	ModuleId  SMALLINT NOT NULL,
	IsVisible Boolean NOT NULL,
	Version   VARCHAR NOT NULL, 
	CreatedOn TIMESTAMP NOT NULL,
	DeletedOn TIMESTAMP
);

CREATE TABLE SectionVersions(
 	Id        SERIAL PRIMARY KEY NOT NULL,
	Name VARCHAR NOT NULL,
	ExternalId VARCHAR NOT NULL,
	SectionId VARCHAR NOT NULL,
	IsVisible BOOLEAN NOT NULL,
	Version VARCHAR NOT NULL,
	CreatedOn TIMESTAMP NOT NULL,
	DeletedOn  TIMESTAMP
);

CREATE TABLE Module(
 	Id        SERIAL PRIMARY KEY NOT NULL,
	Name      VARCHAR NOT NULL,
	Status    VARCHAR NOT NULL,
	IsVisible BOOLEAN NOT NULL,
	TenantId  VARCHAR NOT NULL,
	CreatedOn TIMESTAMP NOT NULL,
	DeletedOn TIMESTAMP
);

CREATE TABLE ModuleVersions(
 	Id        SERIAL PRIMARY KEY NOT NULL,
	ModuleId        SMALLINT NOT NULL,
	ExternalId      VARCHAR NOT NULL,
	Version         VARCHAR NOT NULL,
	CreatedOn       TIMESTAMP NOT NULL,
	DeletedOn       TIMESTAMP,
	Properties      JSON,
	SectionVersions JSON
);


ALTER TABLE Sections 
ADD CONSTRAINT FK_Sections_ModuleId FOREIGN KEY (ModuleId) REFERENCES Modules (Id);

ALTER TABLE ModuleVersions 
ADD CONSTRAINT FK_ModuleVersions_ModuleId FOREIGN KEY (ModuleId) REFERENCES Modules (Id);

ALTER TABLE SectionVersions 
ADD CONSTRAINT FK_SectionVersions_SectionId FOREIGN KEY (SectionId) REFERENCES Sections (Id);

ALTER TABLE Fields 
ADD CONSTRAINT FK_Fields_SectionId FOREIGN KEY (SectionId) REFERENCES Sections (Id);

ALTER TABLE FieldVersions 
ADD CONSTRAINT FK_FieldVersions_FieldId FOREIGN KEY (FieldId) REFERENCES Fields (Id);