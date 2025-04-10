MOCK_DEST := ./mocks
MOCKGEN := $(GOPATH)/bin/mockgen

.PHONY: mocks
mocks:
	$(MOCKGEN) -source=./uc/tasks.go -destination=$(MOCK_DEST)/mock_uc/tasks.go -package=mock
	$(MOCKGEN) -source=./domain/tasks.go -destination=$(MOCK_DEST)/mock_domain/tasks.go -package=mock
