gcloud run revisions list --region=europe-west3 --filter="status.conditions.type:Active AND status.conditions.status:'False'" --format='value(metadata.name)' | xargs -r -L1 gcloud run revisions delete --quiet --region=europe-west3
gcloud run revisions list --region=europe-west2 --filter="status.conditions.type:Active AND status.conditions.status:'False'" --format='value(metadata.name)' | xargs -r -L1 gcloud run revisions delete --quiet --region=europe-west2
gcloud run revisions list --region=europe-west1 --filter="status.conditions.type:Active AND status.conditions.status:'False'" --format='value(metadata.name)' | xargs -r -L1 gcloud run revisions delete --quiet --region=europe-west1
gcloud run revisions list --region=europe-west4 --filter="status.conditions.type:Active AND status.conditions.status:'False'" --format='value(metadata.name)' | xargs -r -L1 gcloud run revisions delete --quiet --region=europe-west4