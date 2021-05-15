init-workspaces:
	terraform workspace new staging-rank-analyse
	terraform workspace new production-rank-analyse

plan-staging:
	terraform workspace select staging-rank-analyse
	terraform plan -var-file="staging.tfvars"

plan-production:
	terraform workspace select production-rank-analyse
	terraform plan -var-file="production.tfvars"

apply-staging:
	terraform workspace select staging-rank-analyse
	terraform apply -var-file="staging.tfvars"

apply-production:
	terraform workspace select production-rank-analyse
	terraform apply -var-file="production.tfvars"
