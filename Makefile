update_codecov:
	echo "Updating Code Coverage Report"
	cd processors;\
	go test -coverprofile=coverage.out -covermode=atomic
	cd handlers;\
	go test -coverprofile=coverage.out -covermode=atomic
	echo "Done Updating Code Coverage Report"