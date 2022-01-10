rsa-gen-key:
	mkdir -p cert/private
	mkdir -p cert/public
	openssl genrsa -out cert/private/private.pem
	openssl rsa -in cert/private/private.pem -pubout > cert/public/public.pem

gen-psql:
	jet -source=PostgreSQL -host=0.0.0.0 -port=5432 -user=postgres -password=postgres -dbname=lenswear -schema=public -path=internal/rdbms/postgresql/gen

run-server:
	go run cmd/api/main.go -env=local
