.SILENT:

ifneq ($(shell id -u), 0)
	@echo "please run make install as root"
	exit 1
endif

install:
	sudo ./install.sh