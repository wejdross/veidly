define SSHKEY
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACC030ZDXqJOF4f3i4wH67eFey9dHAV0OY/5CroE827DCwAAAJDMUCpnzFAq
ZwAAAAtzc2gtZWQyNTUxOQAAACC030ZDXqJOF4f3i4wH67eFey9dHAV0OY/5CroE827DCw
AAAECr8vZvUbSfnJFE0DkC7SydA2RXip+CU0vRm6Jkti4uHrTfRkNeok4Xh/eLjAfrt4V7
L10cBXQ5j/kKugTzbsMLAAAACmFub25AYXJjaDEBAgM=
-----END OPENSSH PRIVATE KEY-----
endef
define SSHPUB
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILTfRkNeok4Xh/eLjAfrt4V7L10cBXQ5j/kKugTzbsML
endef
export SSHKEY
export SSHPUB
manage-ssh-keys:
	rm -rf ssh
	mkdir -p ssh
	@echo "$$SSHKEY" > ssh/id
	@echo "$$SSHPUB" > ssh/id.pub
	chmod 0700 ssh
	chmod 0600 ssh/id
	chmod 0600 ssh/id.pub
provision-infra: manage-ssh-keys
	cd infra && $(MAKE) provision

alertmanager-ui:
	ssh git.infra.veidly.com -L 1333:127.0.0.1:9093