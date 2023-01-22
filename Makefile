.PHONY: client server

client:
	@cd client && ENV=local go run . && cd -

server:
	@cd server && ENV=local go run . && cd -