#!/bin/bash

cd ./kubernetes

if [ ! -z "$DEPLOY_DB_SERVICE" ]; then
    kubectl apply -f mongodbcommunity.mongodb.com_mongodbcommunity.yaml
    kubectl apply -f mongodb-deployment.yaml
fi

if [ ! -z "$DEPLOY_REDIS_SERVICE" ]; then
    kubectl apply -f redis-deployment.yaml
fi

if [ ! -z "$DEPLOY_IDENTITY_SERVICE" ]; then
    kubectl apply -f identity-deployment.yaml
fi

if [ ! -z "$DEPLOY_ORDER_SERVICE" ]; then
    kubectl apply -f order-deployment.yaml
fi

if [ ! -z "$DEPLOY_WALLET_SERVICE" ]; then
    kubectl apply -f wallet-deployment.yaml
fi

if [ -z "$DEPLOY_ORDER_SERVICE" ] &&
    [ -z "$DEPLOY_WALLET_SERVICE" ] &&
    [ -z "$DEPLOY_AUTH_SERVICE" ]; then
    kubectl apply -f mongodbcommunity.mongodb.com_mongodbcommunity.yaml
    kubectl apply -f mongodb-deployment.yaml
    kubectl apply -f redis-deployment.yaml
    kubectl apply -f order-deployment.yaml
    kubectl apply -f wallet-deployment.yaml
    kubectl apply -f identity-deployment.yaml
fi

kubectl apply -f gateway-deployment.yaml

echo "Deployment completed"