#!/bin/bash
set -euxo pipefail

export KO_DOCKER_REPO="quay.io/erdii-private/ensure-quay-repo"
ko build --bare ./cmd/ensure-quay-repo
