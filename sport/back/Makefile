
.PHONY: search trigger-gen

trigger-gen:
	cd trigger_gen; go build -o ../bin/trigger_gen *.go
	cd bin; ./trigger_gen
	
build:
	go build -o bin/run_api main.go

build-dbg:
	go build -o bin/run_dbg_api -tags DBG main.go

rsv-manip:
	cd rsv_manip; go build -o ../bin/rsv_manip *.go
	cd bin; ./rsv_manip

test-adyen-flow:
	cd adyen_flow; go build -o ../bin/adyen.flow *.go
	cd bin; ./adyen.flow

test-mail:
	cd mail_test; go build -o ../bin/mail.test *.go
	cd bin; ./mail.test

gen-train:
	cd train_gen; go build -o ../bin/train_gen *.go
	cd bin; ./train_gen

search-gen:
	cd search_gen; go build -o ../bin/search_gen *.go
	cd bin; ./search_gen

search:
	cd search_perf; go build -o ../bin/search_perf *.go
	cd bin; ./search_perf

ttr=

# test-run
define tr
	cd $(1)_test; go test -c -o ../bin/$(1).test;
	cd bin; ./$(1).test $(2) $(ttr) -test.failfast -test.v;
endef

test-user:
	$(call tr,user)

test-static:
	$(call tr,static)

test-train:
	$(call tr,train)

test-api:
	$(call tr,api)

test-search:
	$(call tr,search)

test-search-coord:
	$(call tr,search,-test.run Coord)

test-search-sug:
	$(call tr,search,-test.run Sug)

test-schedule:
	$(call tr,schedule)

test-rsv:
	$(call tr,rsv)

test-instr:
	$(call tr,instr)
test-instr-profile:
	$(call tr,instr,-test.run Profile*)

test-review:
	$(call tr,review)

test-rsv-qr:
	$(call tr,rsv_qr)

test-sub-qr:
	$(call tr,sub_qr)

test-search-lock:
	$(call tr,search,-test.run Lock)

test-sub:
	$(call tr,sub)

test-rsv-contact:
	$(call tr,rsv,-test.run Contact)
	
test-rsv-grp:
	$(call tr,rsv,-test.run GroupTrainings)

test-helpers:
	$(call tr,helpers)

test-dc:
	$(call tr,dc)

test-invoicing:
	$(call tr,invoicing)

test-adyen:
	$(call tr,adyen)

test-lang:
	$(call tr,lang)

test-chat:
	$(call tr,chat)

test-chat-integrator:
	$(call tr,chat_integrator)

TESTS := test-api
TESTS += test-adyen
TESTS += test-static
TESTS += test-user
TESTS += test-instr
TESTS += test-train
TESTS += test-rsv
TESTS += test-search
TESTS += test-review
TESTS += test-rsv-qr
TESTS += test-dc
TESTS += test-sub
TESTS += test-sub-qr
TESTS += test-schedule
TESTS += test-lang
TESTS += test-chat
TESTS += test-chat-integrator
TESTS += test-invoicing

# there is no point [for now] to run helper tests in bundle
# uncomment this if situation changes
# TESTS += test-helpers


test: build $(TESTS)

api: build
	cd bin; ./run_api --config ../config.yml $(FLAGS)

api-dbg: build-dbg
	cd bin; ./run_dbg_api

container=docker

psql:
	$(container) exec -it pg psql -U postgres -d sportdb

cqlsh:
	$(container) exec -it cass cqlsh

db:
	$(container) container rm --force pg || true
	$(container) container rm --force cass || true 
	$(container) run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name pg docker.io/postgres:13.2
	$(container) run -d -p 9042:9042 -e JVM_OPTS="-Xms512M -Xmx512M" --name cass docker.io/library/cassandra:latest

resetcass:
	podman exec -it cass cqlsh -e 'drop keyspace sportdb' || :
	podman exec -it cass cqlsh -e 'drop keyspace sportdb_test' || :

ip-deploy:
	# download data file if not exists
	cd ip; if [ ! -f data.csv ]; then \
		scp -P 27634 anon@yoda.$(domain):/static/dagobah/data.csv .; \
	fi
	# delete existing ips
	dp -c db.yml --demand dropip -e
	# exec ips
	cd ip; go run main.go

ui:
	cd ui2 && npm start

dbname=sportdb
	
# rpm : install golang-godoc if you dont have it
docs:
	godoc -http :6060
	chromium-browser http://localhost:6060/pkg/sport/

coverage-cleanup: ./*_test
	for x in $^; do \
		echo "rm $${x}/api.yml" ; \
		rm -f $${x}/api.yml; \
	done

coverage: ./*_test
	for x in $^; do \
		cp bin/api.yml $${x}; \
		cd $${x}; go test -coverpkg ../$${x//_test/} | grep coverage; \
		cd ..; \
	done
