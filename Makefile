all:
	@echo "make deploy - Deploys services to local cluster"
	@echo "make fractalk8s - Deploys app to local cluster"

VERSION := $(shell date +%s)

build-docker:
	docker build --tag fractalk8s:$(VERSION) . 
	docker tag fractalk8s:$(VERSION) fractalk8s:latest

docker: build-docker
	docker tag fractalk8s:$(VERSION) k3d-myregistry.localhost:12345/fractalk8s:latest
	docker push k3d-myregistry.localhost:12345/fractalk8s:latest

docker-public: build-docker
	docker tag fractalk8s:$(VERSION) jstrohm/fractalk8s:$(VERSION)
	docker tag fractalk8s:$(VERSION) jstrohm/fractalk8s:latest
	docker push jstrohm/fractalk8s:$(VERSION)
	docker push jstrohm/fractalk8s:latest

build-fractalk8s:
	docker build --tag k3d-myregistry.localhost:12345/fractalk8s:latest -f Dockerfile .

build: build-fractalk8s 
	@echo Build done

deploy-fractalk8s: 
	docker push k3d-myregistry.localhost:12345/fractalk8s:latest

#deploy: build deploy-auth deploy-account deploy-sociallist 
	#-killall kubectl
	#-kubectl delete -f grapevine.yaml
	#kubectl create -f grapevine.yaml
	#kubectl port-forward --namespace default svc/auth 8080:8080 &
	#kubectl port-forward --namespace default svc/account 8081:8080 &
	#@echo Deploy done

fractalk8s: docker-public
	-kubectl delete -f fractalk8s.yaml
	kubectl create -f fractalk8s.yaml
	@echo $(VERSION)

