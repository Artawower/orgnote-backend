#!/bin/bash

mkdir -p app/infrastructure/generated
oapi-codegen -package subscription -generate client,types http://65.108.90.95:8000/api/schema/ > app/infrastructure/generated/subscription.go

echo "Subscumber API layer was successfully generated."
