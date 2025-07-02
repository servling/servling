-- Create index "applications_name_key" to table: "applications"
CREATE UNIQUE INDEX "applications_name_key" ON "applications" ("name");
-- Create index "templates_name_key" to table: "templates"
CREATE UNIQUE INDEX "templates_name_key" ON "templates" ("name");
-- Create "domains" table
CREATE TABLE "domains" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "certificate" character varying NULL,
  "key" character varying NULL,
  "cloudflare_email" character varying NULL,
  "cloudflare_api_key" character varying NULL,
  PRIMARY KEY ("id")
);
-- Create index "domains_name_key" to table: "domains"
CREATE UNIQUE INDEX "domains_name_key" ON "domains" ("name");
-- Create index "services_name_key" to table: "services"
CREATE UNIQUE INDEX "services_name_key" ON "services" ("name");
-- Create index "services_service_name_key" to table: "services"
CREATE UNIQUE INDEX "services_service_name_key" ON "services" ("service_name");
-- Create "ingresses" table
CREATE TABLE "ingresses" (
  "id" character varying NOT NULL,
  "name" character varying NOT NULL,
  "target_port" integer NOT NULL,
  "domain_ingresses" character varying NULL,
  "service_ingresses" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "ingresses_domains_ingresses" FOREIGN KEY ("domain_ingresses") REFERENCES "domains" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "ingresses_services_ingresses" FOREIGN KEY ("service_ingresses") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "ingresses_name_key" to table: "ingresses"
CREATE UNIQUE INDEX "ingresses_name_key" ON "ingresses" ("name");
