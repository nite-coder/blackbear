count = 10


log.bench:
	go test -benchmem -run NONE -bench . -count 1

log.bench/out:
	go test -benchmem -run NONE -bench . -count $(count) > $@