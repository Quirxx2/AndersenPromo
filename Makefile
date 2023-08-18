MAKEFLAGS += --silent

.PHONY: clean
clean:
	rm -rf ./tmp

.PHONY: build.docker
build.docker:
	docker buildx build -t=promo --target=release .

.PHONY: build.mocks
build.mocks: build.mocks.requires build.mock.DBConnexion

.PHONY: build.mocks.requires
build.mocks.requires:
	if ! type mockery > /dev/null; then \
  		echo "mockery required: go install github.com/vektra/mockery/v2@latest"; \
		exit 1;\
	fi;


.PHONY: build.mock.DBConnexion
build.mock.DBConnexion: build.mocks.requires
	mockery --name=DBConnexion --inpackage --testonly --case underscore --with-expecter;

COMPOSER=export HOST_UID=$$(id -u):$$(id -g); docker compose -f docker-compose.yml
.PHONY: up
up:
	$(COMPOSER) up -d -V

.PHONY: down
down:
	docker compose down

.PHONY: restart
restart: down up

.PHONY: db.up
db.up:
	$(COMPOSER) up -d -V db

.PHONY: db.down
db.down:
	docker compose down

.PHONY: db.wait
db.wait: db.up
	echo "waiting for postgres";
	docker compose exec db sh -c 'until pg_isready -q; do \
								  	{ printf .; sleep 0.1; }; \
								  done;'
	echo "\\npostgres is ready";
	sleep 1;
