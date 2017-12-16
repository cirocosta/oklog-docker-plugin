PLUGIN_NAME := oklog-docker-plugin
ROOTFS_IMAGE := cirocosta/$(PLUGIN_NAME)-rootfs
ROOTFS_CONTAINER := rootfs
PLUGIN_FULL_NAME := cirocosta/$(PLUGIN_NAME)


install:
	go install -v


fmt:
	go fmt
	cd ./http && go fmt
	cd ./driver && go fmt
	cd ./docker && go fmt


rootfs-image:
	docker build -t $(ROOTFS_IMAGE) .


rootfs: rootfs-image
	docker rm -vf $(ROOTFS_CONTAINER) || true
	docker create --name $(ROOTFS_CONTAINER) $(ROOTFS_IMAGE) || true
	mkdir -p plugin/rootfs
	rm -rf plugin/rootfs/*
	docker export $(ROOTFS_CONTAINER) | tar -x -C plugin/rootfs
	docker rm -vf $(ROOTFS_CONTAINER)


plugin: rootfs
	docker plugin disable $(PLUGIN_NAME) || true
	docker plugin rm --force $(PLUGIN_NAME) || true
	docker plugin create $(PLUGIN_NAME) ./plugin
	docker plugin enable $(PLUGIN_NAME)


plugin-push: rootfs
	docker plugin rm --force $(PLUGIN_FULL_NAME) || true
	docker plugin create $(PLUGIN_FULL_NAME) ./plugin
	docker plugin push $(PLUGIN_FULL_NAME)


.PHONY: fmt rootfs-image rootfs plugin plugin-push
