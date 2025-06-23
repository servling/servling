-- Modify "applications" table
ALTER TABLE "applications" ADD COLUMN "status" character varying NOT NULL DEFAULT 'stopped', ADD COLUMN "error" character varying NULL;
-- Modify "services" table
ALTER TABLE "services" DROP CONSTRAINT "services_applications_services", ALTER COLUMN "application_services" DROP NOT NULL, ADD COLUMN "status" character varying NOT NULL DEFAULT 'stopped', ADD COLUMN "error" character varying NULL, ADD CONSTRAINT "services_applications_services" FOREIGN KEY ("application_services") REFERENCES "applications" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
