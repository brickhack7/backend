build:
	DOCKER_BUILDKIT=1 docker build -t htn-backend .

pull:
	docker pull alphakilo07/htn-backend

push:
	docker tag htn-backend alphakilo07/htn-backend
	docker push alphakilo07/htn-backend

cloud:
	docker push gcr.io/team-dn-htn/htn-backend

run:
	docker run  --rm -d -p 8081:8081 -e PORT='8081' \
		--name htn-backend htn-backend

kill:
	docker kill htn-backend