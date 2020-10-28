GIT_VERSION:=$(shell git describe --dirty --always)

PROJECT:=kf-infra-gitops

IMAGE:=gcr.io/$(PROJECT)/groups-sync:$(GIT_VERSION)

# Echo commands
SHELL = sh -xv

LATESTIMAGES=latest_image.yaml

CONTEXT:=kf-autobot

# Build it locally
build:
	mkdir -p .build
	go build -o .build/groups ./cmd

# Synchronize the groups using personal credentials.
sync: build
	./.build/groups run --input=./groups/*.yaml \
		--continuous=false \
		--credentials-file=gs://kf-infra-gitops_secrets/autobot-at-kubeflow_client_secret.json

build-submit:
	gcloud --project=$(PROJECT) builds submit --machine-type=n1-highcpu-32 --tag=$(IMAGE) ./

build-output:
	yq w -i latest_image.yaml $(PROJECT).groups.image \
	  $(IMAGE)@$(shell gcloud container images describe $(IMAGE) | yq r - image_summary.digest)

set-image:
	cd  manifests && kustomize edit set image groups=$(shell yq r $(LATESTIMAGES) $(PROJECT).groups.image)

update-image: build-submit build-output set-image


deploy:
	kustomize build manifests | kubectl --context=$(CONTEXT) apply -f -

# Build the image using GCB
# We can't use skaffold due to: https://github.com/GoogleContainerTools/skaffold/issues/3468
build: build-submit build-output
