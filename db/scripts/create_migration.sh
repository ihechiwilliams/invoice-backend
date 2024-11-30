migrate \
		-database "postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" \
		create -ext sql -dir db/migrations "$1"
