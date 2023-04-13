.DEFAULT_GOAL := all
MARKETING_DIR='./email_templates/marketing/'
BACKEND_TEMPLATES_DIR='../back/lang/email_templates/'
UNAME := $(shell uname)
sed=""
ifeq ($(UNAME), Linux)
	sed=/usr/bin/sed
else
	sed=gsed
endif


.PHONY: generate-emails
generate-emails:
	node mailgen.js

.PHONY: rsync
rsync:
	find $(BACKEND_TEMPLATES_DIR)/* -type d -execdir  rm -rf $(basename {}) \;
	rsync -avz email_templates/* $(BACKEND_TEMPLATES_DIR)
	rm -rf $(BACKEND_TEMPLATES_DIR)/marketing

clean-marketing:
	grep '{{ define "content" }}' -rl $(MARKETING_DIR) | xargs $(sed) -i '/{{ define "content" }}/d'
	grep '{{ end }}' -rl $(MARKETING_DIR) | xargs $(sed) -i '/{{ end }}/d'
.PHONY: all

all: generate-emails rsync clean-marketing