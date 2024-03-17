#!/bin/sh
openssl req -new -x509 -days 365 -key mytls.key -out mytls.crt -subj "/CN=sealed-secret/O=sealed-secret"