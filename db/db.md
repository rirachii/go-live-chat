# DB SETUP
Set up postgres server on linux

Make sure docker is installed

Run the following in order:
- `sudo make postgresinit`
- `sudo make createdb`

Now posgres db server is running in detached mode on port :5432

## Install Golang-migrate

- https://www.geeksforgeeks.org/how-to-install-golang-migrate-on-ubuntu/

- `make migrateup`

Check docker postgres db tables

- `sudo make postgres`
- `\l`
- `\c go-chat`
- `\d`

# Create New Migration SQL File

In project root run

`migrate create -ext sql -dir service/db/migrations name_of_table`
