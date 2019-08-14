# Altemista Billing setup

To charge customers according to their use of cloud resources through Altemista, we've developed the program [altemista-billing](www.github.com/altemista/altemista-billing)

This document will explain how to get the program running in your cluster.

Setup in AWS
------------

1. Create an S3 bucket `altemista-billing`. This is where bills will be placed by the tool later.
1. Create an IAM user with the following permissions:
    - Full Cost Explorer access
    - Access to the S3 bucket `altemista-billing
1. Create an aws access key pair for the IAM user.


Setup in Kubernetes
-------------------

### 1. Create a billing namespace in your Kubernetes cluster


`namespace-billing.yaml`
```yaml
kind: Namespace
apiVersion: v1
metadata:
    name: billing
    labels: {name: billing}

```

```shell
kubectl create -f ./namespace-billing.yaml
```

### 2. Create and switch to the billing context

Run the following command to get the current cluster and current user

```shell
kubectl config view
```

Then create the context with the following command, adding in your current user and current cluster from the command above:
```shell
kubectl config set-context billing --namespace=billing \
  --cluster={{cluster}} \
  --user={{user}}
```

Now make it the active context

```shell
kubectl config use-context billing
```
Confirm with
```shell
kubectl config current-context
# billing
```

### 3. Deploy your AWS secrets to the namespace:

Replace {{aws_access_key_id}} and {{aws_secret_access_key}} with the keys you create during the [setup in AWS](#setup-in-aws)

`secrets-billing.yaml`
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secrets-billing
type: Opaque
stringData:
  credentials: |-
    [default]
    aws_access_key_id = {{aws_access_key_id}}
    aws_secret_access_key = {{aws_secret_access_key}}
    [default]
    region = eu-central-1
```


```shell
kubectl create -f ./secrets-billing.yaml
```

### 4. Deploy an altemista/altemista-billing container to the cluster as a cronjob

Replace the `4` in the cronjob definition with the time of the day the cronjob should run (24h format)

Replace the `1` in the cronjob definition with the day of the month the cronjob should run

`cronjob-billing.yaml`
```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: altemista-billing
spec:
  schedule: "* 4 1 * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: altemista-billing-container
            image : altemista/altemista-billing
            args:
            - /app/billing
            - invoice 
            - --month
            - last
            - --provider 
            - aws
            - --bucket
            - altemista-billing
            volumeMounts:
            - name: credentials
              mountPath: "/home/runner/.aws"
              readOnly: true
          restartPolicy: OnFailure
          volumes:
          - name: credentials
            secret:
              secretName: secrets-billing
              items:
              - key: credentials
                path: credentials
              - key: config
                path: config

```

```shell
kubectl create -f ./cronjob-billing.yaml
```