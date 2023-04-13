apply:
	cd terraform && $(MAKE) apply

destroy:
	cd terraform && $(MAKE) destroy

show:
	cd terraform && $(MAKE) show
roles:
	ansible-galaxy collection install community.crypto --force
	ansible-galaxy collection install ansible.posix --force
	ansible-galaxy collection install community.general --force
	pip3 install jmespath passlib

provision: roles
	cd ansible && $(MAKE) provision

provision-fast:
	cd ansible && $(MAKE) provision
clean-roles:
	rm -rf ansible/roles
wait-for-vms:
	sleep 10

fullreset: destroy apply wait-for-vms provision
