.PHONY: build

build:
	sam build

docker-build:
	docker network create lambda-local
	docker compose up -d --build


# docker-compose up -d
# go run main.go

# aws dynamodb \
#   --region ap-northeast-1 \
#   --endpoint-url http://dynamodb:8000 \
#     create-table \
#   --table-name SampleTable \
#   --attribute-definitions \
#     AttributeName=id,AttributeType=S \
#   --key-schema \
#     AttributeName=id,KeyType=HASH \
#   --billing-mode PAY_PER_REQUEST


# テストデータ
# aws dynamodb put-item \
#     --region ap-northeast-1 \
#     --endpoint-url http://dynamodb:8000 \
#     --table-name SampleTable \
#     --item '{
#         "id": {"S": "123"},
#         "name": {"S": "Test"}
#       }'

# docker network create lambda-local


# AWS profile(初回のみ)
# localのDynamoDB
# $ aws configure set aws_access_key_id dummy     --profile local
# $ aws configure set aws_secret_access_key dummy --profile local
# $ aws configure set region ap-northeast-1       --profile local