PROJECT_ID=
PROJECT_NO=
ENDPOINT=
CONTAINER_NAME=
SERVICE_NAME=
SPLIT_TABLE_NAME=

# Iamの設定
first-deploy:
	deploy
	gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
        --member=serviceAccount:cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
	    --role=roles/run.invoker --platform managed
	gcloud secrets add-iam-policy-binding SLACK_TOKEN \
        --role roles/secretmanager.secretAccessor \
        --member serviceAccount:${PROJECT_NO}-compute@developer.gserviceaccount.com
	gcloud secrets add-iam-policy-binding CHANNEL_ID \
        --role roles/secretmanager.secretAccessor \
        --member serviceAccount:${PROJECT_NO}-compute@developer.gserviceaccount.com


# CloudRun　デプロイ
deploy:
	gcloud builds submit --tag gcr.io/${PROJECT_ID}/${CONTAINER_NAME}
	gcloud run deploy ${SERVICE_NAME} --image gcr.io/${PROJECT_ID}/${CONTAINER_NAME} \
	       --platform managed --update-env-vars SPLIT_TABLE_NAME=${SPLIT_TABLE_NAME},TABLE_NAME=biling,MODE=release

# Schedulerの設定
setup-scheduler:
	gcloud scheduler jobs create http task-job --schedule="0 0 10 * *" \
       --http-method=POST \
       --uri=${ENDPOINT}/notification \
       --oidc-service-account-email=cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com   \
       --oidc-token-audience=${ENDPOINT}