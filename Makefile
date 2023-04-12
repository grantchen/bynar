apps := service/payments
changedFiles := $(shell git show --pretty="" --name-only)
rootDir := $(shell pwd)

common:
	git fetch origin

build: common
	@for app in $(apps) ; do \
    		# if [ $$(git show --pretty="" --name-only | grep $$app | wc -l) -eq 0 ] ; then \
    		# 	echo "No change in $$app. Skipped $@~~"; \
    		# 	continue; \
    		# fi; \
    		# echo "$$app is changed. $@~~"; \
    		$(MAKE) -C $(rootDir)/$$app build; \
    	done

build-local: common
	@for app in $(apps) ; do \
			echo "$$app build-local. $@~~"; \
			$(MAKE) -C $(rootDir)/$$app build-local; \
		done

build-lambda: common
	@for app in $(apps) ; do \
			echo "$$app build-lambda. $@~~"; \
			$(MAKE) -C $(rootDir)/$$app build-lambda; \
		done
		