.PHONY: all

REPO=pfandzelter

all: adapt aggregate camera cfd centraldashboard generatedashboard pkgcntrl predict cntrl prognosis sensor

adapt aggregate camera cfd centraldashboard generatedashboard pkgcntrl predict cntrl prognosis sensor:
	cd $@ ; docker build . --no-cache -t $(REPO)/$@:latest
	docker push $(REPO)/$@ 
	cd $@ ; go build .