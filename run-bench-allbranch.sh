#!/bin/bash

GH_TOKEN=$1
GH_REPO=$2
GH_PULL=$3

TESTR="$(pwd)/results"

# Получаем список всех измененных файлов за все время существования ветки
BRANCH=$(git rev-parse --abbrev-ref HEAD)
STARTDATE=$(git log --until="$(git show -s --format=%ct $BRANCH)" --format="%ct" -1)
FILES=$(git log --since="$STARTDATE" --name-only --pretty=format: --diff-filter=ACMRTUXB | sort -u)

UNIQ=()

# Проходим по списку этих файлов и собираем адреса директорий в которых лежал эти файлы
for file in $FILES; do
    FPATH=$(dirname "$file")
    # Проверяем, есть ли в этой папке файлы, имя которых заканчивается на "test.go"
    TESTFILES=$(find "$FPATH" -maxdepth 1 -type f -name "*test.go")
    if [[ -n "$TESTFILES" ]]; then
      if ! [[ " ${UNIQ[*]} " == *" $FPATH "* ]]; then
        UNIQ+=("$FPATH")
      fi
    fi
done

BENCHED=""
# Проходим по списку каталогов и запускаем в них тесты, если есть
for path in "${UNIQ[@]}"; do
    pushd ./"$path" || exit 1
    if [[ -n "$(go test --list . | grep Benchmark)" ]]; then
      go test --run=^$ --bench=. --count=10 --timeout=60m | tee -a "${TESTR}/benchmarks.log"
      BENCHED="OK"
    fi
    popd || exit 1
done

if [[ -z "$BENCHED" ]]; then
  echo "there is no benchmarks"
  exit 0
fi

mkdir -p ./bin
GOBIN="$(pwd)/bin" go install golang.org/x/perf/cmd/...@latest

echo '```' | tee "${TESTR}/stat.log"
./bin/benchstat -filter ".unit:(ns/op OR allocs/op)" -table .config -col .name ./results/benchmarks.log | tee -a "${TESTR}/stat.log"
echo '```' | tee -a "${TESTR}/stat.log"

wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64  || exit 1
chmod +x ./jq

./jq -n --rawfile a ./results/stat.log '{body: $a}' | tee ./results/stat.json || exit 1

curl --http1.1 -v -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: token $GH_TOKEN" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "https://api.github.com/repos/${GH_REPO}/issues/${GH_PULL}/comments" \
  -d @results/stat.json \
  -o ./results/curl.out.txt
