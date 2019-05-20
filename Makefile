DOCKER_IMAGE ?= containersol/externalsecret-operator
DOCKER_TAG ?= $(shell grep -Po 'Version = "\K.*?(?=")' version/version.go)

# export these if you want to use AWS secrets manager
AWS_ACCESS_KEY_ID ?= AKIACONFIGUREME
AWS_SECRET_ACCESS_KEY ?= Secretsecretconfigureme 
AWS_REGION ?= eu-west-1

NAMESPACE ?= "default"

.PHONY: build
build:
	operator-sdk build $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: push
push:
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: deploy
.EXPORT_ALL_VARIABLES: deploy
deploy:
	kubectl apply -n $(NAMESPACE) -f ./deploy/service_account.yaml
	kubectl apply -n $(NAMESPACE) -f ./deploy/role.yaml
	envsubst < ./deploy/role_binding.yaml | kubectl apply -n $(NAMESPACE) -f  -
	kubectl apply -n $(NAMESPACE) -f ./deploy/crds/externalsecret-operator_v1alpha1_externalsecret_crd.yaml
	kubectl apply -n $(NAMESPACE) -f ./deploy/crds/externalsecret-operator_v1alpha1_externalsecretbackend_crd.yaml
	envsubst < deploy/operator.yaml | kubectl apply -n $(NAMESPACE) -f -

.PHONY: test
test:
	go test -v ./pkg/secrets/... ./pkg/controller/...
