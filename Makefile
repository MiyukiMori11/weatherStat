build:
	cd ./apiGateway/ && docker build -f ./Dockerfile -t apigw .
	cd ./explorer/ && docker build -f ./Dockerfile -t explorer --target explorer .
	cd ./explorer/ && docker build -f ./Dockerfile -t explorer-migration --target migration .

compose:
	docker compose up

test:
	cd ./explorer/ && go test ./...