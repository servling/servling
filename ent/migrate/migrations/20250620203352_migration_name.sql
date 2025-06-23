-- Create "templates" table
CREATE TABLE "templates" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "password" character varying NOT NULL,
  "token_version" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_name_key" to table: "users"
CREATE UNIQUE INDEX "users_name_key" ON "users" ("name");
-- Create "applications" table
CREATE TABLE "applications" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NOT NULL,
  "image_url" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "template_applications" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "applications_templates_applications" FOREIGN KEY ("template_applications") REFERENCES "templates" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create "services" table
CREATE TABLE "services" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "service_name" character varying NOT NULL,
  "image" character varying NOT NULL,
  "ports" jsonb NULL,
  "environment" jsonb NULL,
  "entrypoint" character varying NULL,
  "labels" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "application_services" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "services_applications_services" FOREIGN KEY ("application_services") REFERENCES "applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
