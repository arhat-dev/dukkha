#!/bin/sh

export COSIGN_PASSWORD="testdata"

cosign generate-key-pair
cosign sign-blob --key cosign.key blob.txt --output blob.txt.sig
