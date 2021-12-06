MANIFESTS := https://raw.githubusercontent.com/reactive-tech/kubegres/v1.13/kubegres.yaml
DEST := /tmp/kubegres.yaml

.DEFAULT_GOAL := project

download_manifests:
	curl -s $(MANIFESTS) -o $(DEST)

project: download_manifests
	PULUMI_SKIP_UPDATE_CHECK=true pulumi stack init test
	PULUMI_SKIP_UPDATE_CHECK=true pulumi config set manifestPath $(DEST)
	PULUMI_SKIP_UPDATE_CHECK=true pulumi up -yf

clean:
	#kubectl delete pvc --all  -A
	PULUMI_SKIP_UPDATE_CHECK=true pulumi destroy -yf
	PULUMI_SKIP_UPDATE_CHECK=true pulumi stack rm test -y
	rm -rfv $(DEST)

crd: download_manifests
	crd2pulumi --goPath ./kubegres/crd/generated $(DEST)
