#!/bin/bash

openssl genrsa -out ./keys/private_access.pem 2048
openssl rsa -in ./keys/private_access.pem -pubout -out ./keys/public_access.pem

openssl genrsa -out ./keys/private_refresh.pem 2048
openssl rsa -in ./keys/private_refresh.pem -pubout -out ./keys/public_refresh.pem