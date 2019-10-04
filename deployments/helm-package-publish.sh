helm package helm-cabin/
version="$(grep -Po "^version:\s\K(.*)$" helm-cabin/Chart.yaml)"
name="$(grep -Po "^name:\s\K(.*)$" helm-cabin/Chart.yaml)"
chartname="$name-$version.tgz"
mv "$chartname" ../../helmcabin-helm-repo/
cd ../../helmcabin-helm-repo/
. ./make-repo-index.sh
git add "$chartname" index.yaml
git commit -m "Update chart"
git push
