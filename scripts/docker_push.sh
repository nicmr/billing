#!/bin/bash
docker tag $DOCKER_NS/$DOCKER_REPO:$COMMIT $DOCKER_REPO:latest
docker tag $DOCKER_NS/$DOCKER_REPO:$COMMIT $DOCKER_REPO:travis-$TRAVIS_BUILD_NUMBER
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push $DOCKER_NS/$DOCKER_REPO:latest
docker push $DOCKER_NS/$DOCKER_REPO:travis-$TRAVIS_BUILD_NUMBER
docker push $DOCKER_NS/$DOCKER_REPO:$COMMIT