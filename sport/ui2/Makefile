
define cpconf
	@printf '// This file was auto-generated DONT MODIFY IT\n\n%s' "`cat src/conf.$(1).js`" > src/conf.js
endef



conf-local:
	$(call cpconf,local)

conf-dev:
	$(call cpconf,dev)

conf-prod:
	$(call cpconf,prod)

inject-monitoring-prod:
	sed -ie '4r newrelic-prod.html' public/index.html

npm-ci:
	npm install --force

install-dev: conf-dev npm-ci

install-prod: conf-prod npm-ci inject-monitoring-prod
	

build:
	npm run build

start:
	npm start

install:
	npm install
