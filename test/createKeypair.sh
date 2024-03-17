#!/bin/bash

openssl req -x509 -days 365 -nodes -newkey rsa:4096 -keyout "mytls.key" -out "mytls.crt" -subj "/CN=sealed-secret/O=sealed-secret"