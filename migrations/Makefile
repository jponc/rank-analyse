deploy-staging:
	sls deploy --stage=staging --verbose
	sls invoke --stage=staging -f UpMigration

deploy-production:
	sls deploy --stage=production --verbose
	sls invoke --stage=production -f UpMigration

