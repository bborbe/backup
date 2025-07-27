
precommit: lint build test

deps:
	sudo npm install -g npm-check-updates

test:
	npm run test -- --run

build:
	npm install
	npm run build

watch:
	npm install
	npm run build-watch

update:
	npm audit fix
	npm update
	ncu -u

lint:
	npm run lint:analyse
