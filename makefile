NAMESPACE ?= shortly-production

.PHONY: deploy undeploy template lint minikube-start minikube-build minikube-deploy minikube-show

## Deploy the entire stack with Helm
deploy:
	helm upgrade --install shortly ./charts/shortly \
		--namespace $(NAMESPACE) --create-namespace

## Remove the Helm release
undeploy:
	helm uninstall shortly --namespace $(NAMESPACE)

## Render manifests without applying
template:
	helm template shortly ./charts/shortly --namespace $(NAMESPACE)

## Lint the Helm chart
lint:
	helm lint ./charts/shortly

## Local development with Minikube
minikube-start:
	minikube start

minikube-build:
	eval $$(minikube docker-env) && \
	docker build -t shortener_svc:latest ./src/shortener_svc && \
	docker build -t redirect_svc:latest ./src/redirect_svc

minikube-deploy: minikube-build
	kubectl apply -f k8s/

minikube-show:
	minikube service krakend --url
