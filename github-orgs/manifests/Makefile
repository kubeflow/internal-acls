PROJECT := kubeflow-admin
IMG = gcr.io/$(PROJECT)/peribolos


TAG := $(shell date +v%Y%m%d)-$(shell git describe --always --dirty)
TAG := $(shell date +v%Y%m%d)-$(shell git describe --tags --always --dirty)-$(shell git diff | shasum -a256 | cut -c -6)
DIR := ${CURDIR}

# Create a config map containing the python code to valide the image
create-config-map:
	kubectl -n github-admin create configmap github-sync-code --from-file=validate_config.py