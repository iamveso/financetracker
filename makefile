MIGRATE = migrate
MIGRATION_DIR = db/migrations
DB_URL = "sqlite3://financetracker.db"

migrate_up:
	$(MIGRATE) -path $(MIGRATION_DIR) -database $(DB_URL) -verbose up

migrate_down:
	$(MIGRATE) -path $(MIGRATION_DIR) -database $(DB_URL) -verbose down

migrate-create:
	@i	 [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name: make migrate-create name=init_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migration/ -seq $(name)
