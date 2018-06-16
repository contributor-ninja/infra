UP = up
APEX = apex

generate-config:
	cat api/index/up.json.tpl | envsubst > api/index/up.json
	cat project.json.tpl | envsubst > project.json

deploy-api-prod:
	cd api/index && $(UP) deploy production

deploy-crawler-prod:
	$(APEX) deploy
