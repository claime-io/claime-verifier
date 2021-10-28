.PHONY: setup test build deploy
TARGET = dev

setup:
	export GO111MODULE=on
	go mod vendor
	go get github.com/aws/aws-lambda-go/cmd/build-lambda-zip
build: test
	export GO111MODULE=on
	for module_dir in $$(ls lib/functions | grep -v lib); do\
	  echo  "building start... $${module_dir}";\
	  	cd lib/functions/$${module_dir};\
		pwd;\
		mkdir -p bin;\
		env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main main.go || exit 1;\
		zip bin/main.zip bin/main;\
		cd ../../..;\
	  echo  "building finished. $${module_dir}";\
	done
test: 
	for module_dir in $$(find lib/functions -type f -name "*_test.go" | sed -e "s/[^/]*_test\.go//g" | uniq); do\
	  echo  "testing start... $${module_dir}";\
		cd $${module_dir} && go test -v;\
		if [ $$? != 0 ]; then\
		  exit 1;\
		fi;\
		cd -;\
	  echo  "testing finished. $${module_dir}";\
	done
deploy:
	cdk deploy -c target=dev --all --require-approval never
abi:
	abigen --abi claime-registry/abi/contracts/ClaimRegistry.sol/ClaimRegistry.json --pkg contracts --out lib/functions/lib/contracts/claimregistry.go	
