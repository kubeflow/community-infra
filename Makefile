
# Make sure there are no nomos errors
.PHONY: 
acm-test:
	nomos vet --no-api-server-check --path=prod
