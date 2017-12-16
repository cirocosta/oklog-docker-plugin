#!/bin/bash

set -o errexit

main() {
	setup_dependencies

	echo "INFO:
  Done! Finished setting up travis machine.
  "
}

setup_dependencies() {
	echo "INFO:
  Setting up dependencies.
  "

	sudo apt update -y
	sudo apt install realpath -y
	sudo apt install --only-upgrade docker-ce -y

	docker info
}

main

