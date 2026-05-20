.DEFAULT_GOAL := help
.PHONY: help start stop restart status reset migrate migration-new migrations diff lint

# The Supabase CLI operates on the project root (the parent of this directory).
SUPABASE := supabase --workdir ..

help: ## List all available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

start: ## Start the local Supabase stack
	$(SUPABASE) start

stop: ## Stop the local Supabase stack
	$(SUPABASE) stop

restart: ## Restart the local Supabase stack
	-$(SUPABASE) stop
	$(SUPABASE) start

status: ## Show status, service URLs and keys
	$(SUPABASE) status

reset: ## Reset the local database (re-runs migrations + seed)
	$(SUPABASE) db reset

migrate: ## Apply pending migrations to the local database
	$(SUPABASE) migration up

migration-new: ## Create a new migration (usage: make migration-new name=add_items)
	@test -n "$(name)" || { echo "usage: make migration-new name=<name>"; exit 1; }
	$(SUPABASE) migration new $(name)

migrations: ## List migrations and their applied status
	$(SUPABASE) migration list

diff: ## Show schema changes not yet captured in a migration
	$(SUPABASE) db diff

lint: ## Lint the database schema
	$(SUPABASE) db lint
