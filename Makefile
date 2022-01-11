rsa-gen-key:
	mkdir -p cert/private
	mkdir -p cert/public
	openssl genrsa -out cert/private/private.pem
	openssl rsa -in cert/private/private.pem -pubout > cert/public/public.pem

run-server:
	go run cmd/api/main.go -env=local
