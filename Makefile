PROJECT_NO=[YOUR_PROJECT_NO]
PROJECT_ID=[YOUR_PROJECT_ID]
ENDPOINT=[YOU_ENDPOINT]
CONTAINER_NAME=[CONTAINER_NAME]
SERVICE_NAME=[YOUR_SERVICE_NAME]

# Iamの設定
initialize:
	gcloud run services add-iam-policy-binding ${SERVICE_NAME} \
        --member=serviceAccount:cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com \
	    --role=roles/run.invoker --platform managed

# CloudRun　デプロイ
deploy:
	gcloud builds submit --tag gcr.io/${PROJECT_ID}/${CONTAINER_NAME}
	gcloud run deploy ${SERVICE_NAME} --image gcr.io/${PROJECT_ID}/${CONTAINER_NAME} \
	       --platform managed

# Schedulerの設定
setup-scheduler:
	gcloud scheduler jobs create http task-job --schedule="0 0 10 * *" \
       --http-method=POST \
       --uri=${ENDPOINT}/task \
       --oidc-service-account-email=cloud-run-scheduler-invoker@${PROJECT_ID}.iam.gserviceaccount.com   \
       --oidc-token-audience=${ENDPOINT}