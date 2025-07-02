set windows-shell := ["nu", "-c"]

generate:
  go generate "./..."

js-dev:
  pnpm dev

go-dev:
  air

dev:
  @echo "Starting development server"
  @go run cmd/run-parallel/main.go "just go-dev" "just js-dev"

new-migration name:
  atlas migrate diff {{name}} --dir "file://ent/migrate/migrations" --to "ent://ent/schema" --dev-url "docker://postgres/15/test?search_path=public"

migrate url:
  atlas migrate apply --dir "file://ent/migrate/migrations" --url "{{url}}" --revisions-schema=public

migrate-dev: (migrate "postgres://servling:servling@localhost:5432/servling?search_path=public&sslmode=disable")