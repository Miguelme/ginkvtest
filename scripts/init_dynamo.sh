##!/bin/bash
#
#echo "Initializing DynamoDB Table..."
#
#aws dynamodb create-table \
#    --table-name KeyValueTable \
#    --attribute-definitions AttributeName=key,AttributeType=S \
#    --key-schema AttributeName=key,KeyType=HASH \
#    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
#    --endpoint-url http://localhost:8000
#
#echo "Inserting initial data into DynamoDB..."
#
#aws dynamodb put-item \
#    --table-name KeyValueTable \
#    --item '{"key": {"S": "test_key"}, "value": {"S": "test_value"}}' \
#    --endpoint-url http://localhost:8000
#
#echo "DynamoDB setup completed."
