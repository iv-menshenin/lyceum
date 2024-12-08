#!/bin/bash

GH_TOKEN=$1
GH_REPO=$2
GH_PULL=$3

# Получаем список файлов, которые были изменены последним комитом
FILES=$(git diff --name-only HEAD^ HEAD)

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

wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 > /dev/null || exit 1
chmod +x ./jq > /dev/null

ACC=""
FAILED=""
PASSED=""
echo "" | tee ./results/tests.log
# Проходим по списку каталогов и запускаем в них тесты, если есть
for path in "${UNIQ[@]}"; do
    go test --json ./"$path" | tee ./results/test.out
    while IFS= read -r line
    do
      RES=$(echo "$line" | ./jq '.Action' -r)
      TST=$(echo "$line" | ./jq '.Test' -r)
      if [[ "$RES" == "output" ]]; then
        LN=$(echo "$line" | ./jq '.Output' -r) > /dev/null
        if [[ "$LN" == "null" ]]; then
          ACC=$(printf '%s\n' "$ACC")
        else
          ACC=$(printf '%s\n%s' "$ACC" "${LN}")
        fi
      fi
      if [[ "$RES" == "pass" ]]; then
        PASSED=$(printf '%s\n%s' "$PASSED" "$TST")
        ACC=""
      fi
      if [[ "$RES" == "fail" ]]; then
        FAILED=$(printf '%s\n\n%s' "$FAILED" "$ACC")
        ACC=""
      fi
      if [[ "$RES" == "run" ]]; then
        ACC=""
      fi
    done < ./results/test.out
done

if [[ -n "$FAILED" ]]; then
  printf 'FAILED\n```\n%s\n```\n' "$FAILED" | tee ./results/test_stat.log > /dev/null
else
  printf 'PASSED\n```\n%s\n```\n' "$PASSED" | tee ./results/test_stat.log > /dev/null
fi

./jq -n --rawfile a ./results/test_stat.log '{body: $a}' | tee ./results/test_stat.json > /dev/null || exit 1

if [[ -n "$GH_PULL" ]]; then
  curl --http1.1 -v -L \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: token $GH_TOKEN" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "https://api.github.com/repos/${GH_REPO}/issues/${GH_PULL}/comments" \
    -d @results/test_stat.json \
    -o ./results/curl.out.txt
fi
