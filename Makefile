update_codecov:
	echo "Updating Code Coverage Report"
	cd Processors;\
	go test -coverprofile=coverage.out -covermode=atomic
	cd Bot;\
	go test -coverprofile=coverage.out -covermode=atomic
	echo "Done Updating Code Coverage Report"