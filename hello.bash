#!/bin/bash
set -e


hi() {
    inside_files_array="$1"
    for file in "${inside_files_array[@]}"; do
        if [[ $file == *.json ]]; then
            echo "file=$file"
        fi
    done
}


readarray -t inside_files_array < <(ls)
hi "$inside_files_array"

# for file in "${inside_files_array[@]}"; do
#     if [[ $file == *.json ]]; then
#         echo "file=$file"
#     fi
# done

db.runCommand (
   {
      listIndexes: "bb.one",
   }
)