ifndef UNIQUE_BUILD_ID
	UNIQUE_BUILD_ID=latest
endif

PROJECT_NAME=ottct-poller-service
DOCKER_IMAGE_APP=$(PROJECT_NAME):$(UNIQUE_BUILD_ID)

include scripts/Makefile.dev
include scripts/Makefile.build
include scripts/Makefile.help
