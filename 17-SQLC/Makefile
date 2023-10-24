# Install migrate

# $ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
# $ sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
# $ sudo apt-get update
# $ sudo apt-get install -y migrate
# $ sudo snap install sqlc



createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose up

migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose down

.PHONY: migrate migratedown createmigration