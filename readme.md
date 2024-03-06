Run `make run`

Server on port 8080

## Versions
- Ubuntu 22.04
- Go 1.22
- Echo - why, bc returns error and uses the orignal handler

## Setup
- `sudo snap install docker`
- install make ( i forgot cmd )

## To run app
- `make` to run server on post :8080
- `sudo make postgresup` to start postgres server, make sure that you have set up postgres successfully.

Set up postgres server on linux
Run the following in order, might need to copy the command in the make file and append sudo before idk why can't bypass permission error
- `make postgresinit`
- `make postgres`
- `make createdb` 
Postgres server is in detached mode

Install Golang-migrate
- https://www.geeksforgeeks.org/how-to-install-golang-migrate-on-ubuntu/

Create postgres tables, kinda like git for db tables
- `make createmigration`
- `make migrateup`


Check docker db tables
- `make postgres`
- `\l`
- 
Connect to go-chat db: `\c go-chat`
TO get items `\d`
