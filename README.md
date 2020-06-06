# billing-slack-notifier

go 1.13 

## Edit Makefile

Edit your param

```
PROJECT_ID=
PROJECT_NO=
ENDPOINT=
CONTAINER_NAME=
SERVICE_NAME=
SPLIT_TABLE_NAME=
```

## first　Deploy

```
make first-deploy
```

## Deploy

Build and deploy the container image

```
make deploy
```

## Set Scheduler

Set cloud scheduler

```
make setup-scheduler
```