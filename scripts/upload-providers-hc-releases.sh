#!/bin/bash

PROVIDER_NAME=${PROVIDER_NAME:-"terraform-provider-local"} # Change to required provider
S3_BUCKET=${S3_BUCKET:-"ank-website"}
PROVIDER_BASE_URL="https://releases.hashicorp.com/${PROVIDER_NAME}"
LOCAL_TMP_DIR="/tmp/${PROVIDER_NAME}"
PROVIDER_SHORT_NAME=$(echo "$PROVIDER_NAME" | sed 's/^terraform-provider-//')
S3_CMD_CFG_LOCATION=${S3_CMD_CFG_LOCATION:-"$HOME/.s3cfg"}
S3_BACKUP_DIR=${S3_BACKUP_DIR:-"./s3_backup"}

if [[ -n "${S3_BACKUP_DIR}" ]]; then
  timestamp=$(date +%Y%m%d-%H%M%S)

  echo "Backup s3 bucket to ${S3_BACKUP_DIR}/${S3_BUCKET}-${timestamp}/"

  mkdir -p "${S3_BACKUP_DIR}"

  s3cmd sync --config=$S3_CMD_CFG_LOCATION --quiet --no-preserve "s3://${S3_BUCKET}/" "${S3_BACKUP_DIR}/${S3_BUCKET}-${timestamp}/"

  echo "Finish backup"
fi

while true; do
    read -p "All is OK [yes/no]? " answer
    if [[ $answer == [yY] || $answer == [yY][eE][sS] ]]; then
        echo "Go on"
        break
    elif [[ $answer == [nN] || $answer == [nN][oO] ]]; then
        echo "Canceled by user"
        exit 1
    else
        echo "Incorrect input. Please, enter 'yes' or 'no'."
    fi
done

if [ ${#PROVIDER_SHORT_NAME} -eq 3 ]; then
    PROVIDER_FORMATTED_NAME=$(echo "$PROVIDER_SHORT_NAME" | awk '{print toupper($0)}')
else
    PROVIDER_FORMATTED_NAME=$(echo "$PROVIDER_SHORT_NAME" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}')
fi

mkdir -p $LOCAL_TMP_DIR

REMOTE_VERSIONS=$(curl -s $PROVIDER_BASE_URL/ | grep -oP '(?<=<a href="/'${PROVIDER_NAME}'/)[^/"]*')

LOCAL_VERSIONS=$(s3cmd ls "s3://${S3_BUCKET}/${PROVIDER_NAME}/" | awk '/DIR/ {split($0, arr, "/"); print arr[length(arr)-1]}')

generate_index_html() {
    local target_directory=$1
    local title=$2
    local items=("$@")
    local num_items=$((${#items[@]} - 2))

    echo '<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>'"$title"'</title>
    <style type="text/css">
        BODY { font-family : monospace, sans-serif;  color: black;}
        P { font-family : monospace, sans-serif; color: black; margin:0px; padding: 0px;}
        A:visited { text-decoration : none; margin : 0px; padding : 0px;}
        A:link    { text-decoration : none; margin : 0px; padding : 0px;}
        A:hover   { text-decoration: underline; background-color : yellow; margin : 0px; padding : 0px;}
        A:active  { margin : 0px; padding : 0px;}
        .VERSION { font-size: small; font-family : arial, sans-serif; }
        .NORM  { color: black;  }
        .FIFO  { color: purple; }
        .CHAR  { color: yellow; }
        .DIR   { color: blue;   }
        .BLOCK { color: yellow; }
        .LINK  { color: aqua;   }
        .SOCK  { color: fuchsia;}
        .EXEC  { color: green;  }
    </style>
</head>
<body>
    <p>
    <a href=".">.</a><br>
    ├── <a href="../">..</a><br>' > "$target_directory/index.html"

    for (( i=0; i<num_items; i++ )); do
        item="${items[i+2]}"
        item_path=${item%/}

        if [[ $item_path != "index.html" ]]; then
            if [ $i -lt $((num_items - 1)) ]; then
                echo '    ├── <a href="./'"$item_path"'">'"$item_path"'</a><br>' >> "$target_directory/index.html"
            else
                echo '    └── <a href="./'"$item_path"'">'"$item_path"'</a><br>' >> "$target_directory/index.html"
            fi
        fi
    done

    echo '    </p>
</body>
</html>' >> "$target_directory/index.html"

echo "index.html file is created in $target_directory"
}

for VERSION in $REMOTE_VERSIONS; do
    if ! echo "$LOCAL_VERSIONS" | grep -q "^${VERSION}$"; then
        echo "New version found: $VERSION. Downloading..."
        
        TARGET_DIR="${LOCAL_TMP_DIR}/${PROVIDER_NAME}/${VERSION}/"
        
        mkdir -p $TARGET_DIR
        
        INDEX_PAGE=$(curl -s "${PROVIDER_BASE_URL}/${VERSION}/")

        FILES=$(echo "$INDEX_PAGE" | grep -oP '<a .*?href="\K[^"]*')

        for FILE in $FILES; do
            if [[ $(basename $FILE) == ${PROVIDER_NAME}_${VERSION}_* ]]; then
                DOWNLOAD_URL="${FILE}"
                echo "Downloading $FILE..."
                wget -q -O "${TARGET_DIR}/$(basename ${FILE})" $FILE
            fi
        done

        cd $TARGET_DIR

        files=(*)
        generate_index_html "$TARGET_DIR" "Terraform Provider: $PROVIDER_FORMATTED_NAME | v$VERSION" "${files[@]}"

        cd -
    else
        echo "Version $VERSION already exists in bucket."
    fi
done

cd "$LOCAL_TMP_DIR/$PROVIDER_NAME"
versions=($(ls -d */ | sed 's/\/$//'))
versions_with_slash=("${versions[@]/%//}")
EXISTING_VERSIONS=($(s3cmd ls "s3://${S3_BUCKET}/${PROVIDER_NAME}/" | grep DIR | awk -F'/' '{print $(NF-1)"/"}'))
all_versions=("${versions_with_slash[@]/%//}" "${EXISTING_VERSIONS[@]/%//}")
generate_index_html "$LOCAL_TMP_DIR/$PROVIDER_NAME" "Terraform Provider: $PROVIDER_FORMATTED_NAME" "${all_versions[@]}"

s3cmd sync $LOCAL_TMP_DIR/ "s3://$S3_BUCKET"
s3cmd modify --recursive "s3://$S3_BUCKET/$PROVIDER_NAME/" --acl-public

rm -rf $LOCAL_TMP_DIR

echo "Done!"
