#!/bin/bash
docker tag $TAG_COMMIT $TAG_BUILD
docker tag $TAG_COMMIT $TAG_LATEST
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push $TAG_COMMIT
docker push $TAG_BUILD
docker push $TAG_LATEST