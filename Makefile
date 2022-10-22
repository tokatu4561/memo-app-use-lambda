.PHONY: build

build:
	sam build

start-api:
	sam local start-api --env-vars env.json
	
invoke:
	sam local invoke --env-vars env.json


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
