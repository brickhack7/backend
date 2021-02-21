build:
	DOCKER_BUILDKIT=1 docker build -t brickhack-backend .

pull:
	docker pull alphakilo07/brickhack-backend

push:
	docker tag brickhack-backend alphakilo07/brickhack-backend
	docker push alphakilo07/brickhack-backend

cloud:
	docker push gcr.io/team-dn-htn/htn-backend

run:
	docker run  --rm -d -p 8081:8081 -e PORT='8081' \
		--name brickhack-backend brickhack-backend

kill:
	docker kill brickhack-backend
	