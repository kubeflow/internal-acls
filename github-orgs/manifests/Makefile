PROJECT := kubeflow-admin
IMG = gcr.io/$(PROJECT)/peribolos


TAG := $(shell date +v%Y%m%d)-$(shell git describe --always --dirty)
TAG := $(shell date +v%Y%m%d)-$(shell git describe --tags --always --dirty)-$(shell git diff | shasum -a256 | cut -c -6)
DIR := ${CURDIR}

# Apply all the updates
apply: apply-config-map apply-cron-job

# Create a config map containing the python code to valide the image
apply-config-map:
	kubectl -n github-admin delete configmap github-sync-code || echo "configmap does not exit"
	kubectl -n github-admin create configmap github-sync-code --from-file=validate_config.py

apply-cron-job:
	kubectl -n github-admin apply -f github-sync-cron-job.yaml

run-github-sync-once:
	kubectl -n github-admin create -f github-sync-job.yaml
