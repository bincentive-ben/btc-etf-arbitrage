# Run Client Portal API Gateway
runPortal:
	cd third_party/clientportal.gw/ && ./bin/run.sh root/conf.yaml
	cd -


checkAuth:
	go run main.go arbitrage checkAuthStatus --config ./configs/env.dev.yaml

