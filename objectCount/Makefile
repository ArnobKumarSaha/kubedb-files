all:
	export CGO_ENABLED=0
	go build
	docker build -t arnobkumarsaha/object-count .
	docker push arnobkumarsaha/object-count
	kubectl apply -f pod.yaml

clean:
	kubectl delete -f pod.yaml