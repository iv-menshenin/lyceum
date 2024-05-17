#!/bin/bash

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

# Проходим по списку каталогов и запускаем в них тесты, если есть
for path in "${UNIQ[@]}"; do
    go test --json ./"$path"
done