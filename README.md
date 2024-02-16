# ensure-quay-repo

uses quay's api to ensure a given repo exists.

usage: `AUTH_TOKEN=<secret-token> ensure-quay-repo <org> <repo name>`

example log output:
```
2024/02/16 14:39:38 INFO Ensuring Repo org=erdii-private name=ensure-quay-repo
2024/02/16 14:39:38 INFO Checking if Repo exists namespacedName=erdii-private/ensure-quay-repo
2024/02/16 14:39:39 INFO Repo does not exist. Creating Repo org=erdii-private name=ensure-quay-repo
```

This repo uses [ko.build](https://ko.build/) to build and publish its own container image.
The container image is published at [quay.io/erdii-private/ensure-quay-repo](https://quay.io/repository/erdii-private/ensure-quay-repo?tab=tags)
