IMAGE = "locker"

.phony: docker-image
docker-image:
	@docker build -t $(IMAGE) .

.phony: clean
clean:
	docker builder prune -f
	docker image prune -f