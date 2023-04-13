git-sm:
	git submodule update --init lang

ide:
	cd back; codium . &
	cd ui2; codium . &

dbs:
	podman start pg cass

dev: dbs ide

merge-dev:
	git fetch . `git branch --show-current`:dev
	git push origin dev
