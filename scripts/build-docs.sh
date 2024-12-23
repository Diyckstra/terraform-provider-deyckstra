#!/bin/bash

# Path of the directory with the documentation
SOURCE_DIR="site/"

help () {
    echo "
    --help  Show this message and exit.
    --tools Install all required dependencies.
    --push  Build and push internal documentation to the bucket.
    --local Build and run internal documentation locally."
}

tools () {
  pip3 install mkdocs mkdocs-material mkdocs-awesome-pages-plugin | pip install mkdocs mkdocs-material mkdocs-awesome-pages-plugin
}

# Copying of configuration for the documentatiion to the website folder
copy () {
  cp -r mkdocs/images docs/mkdocs_images
  cp -r mkdocs/assets docs/mkdocs_assets
  cp mkdocs/mkdocs.yml ./
}

# Generation of the documentation
build () {
  echo "Generation of the documentation"
  mkdocs build -f ./mkdocs.yml --clean
}

# Run the documentation locally
run_local () {
  echo "Run the documentation locally"
  mkdocs serve -f ./mkdocs.yml
}

# Funtion to upload the documentation files to the bucket
upload_other_files () {
    s3cmd sync "$SOURCE_DIR" "s3://$S3_DOCS_BUCKET_NAME" --acl-public
}

# Funtion to upload .css files to the bucket with specified Content-Type
upload_css_files () {
    find "$SOURCE_DIR" -type f -name "*.css" | while read -r file; do
    s3cmd modify "s3://$S3_DOCS_BUCKET_NAME/${file#$SOURCE_DIR}" --add-header='Content-Type:text/css'
done
}

# Funtion to upload .js files to the bucket with specified Content-Type
upload_js_files () {
    find "$SOURCE_DIR" -type f -name "*.js" | while read -r file; do
    s3cmd modify "s3://$S3_DOCS_BUCKET_NAME/${file#$SOURCE_DIR}" --add-header='Content-Type:application/javascript'
done
}

# Removing of the temp files
cleanup () {
  echo "Removing of the temp files"
  rm -fr docs/mkdocs_images docs/mkdocs_assets mkdocs.yml site
  echo Complete!
}

# Uploading of the documentation to the S3 bucket

if [[ "$1" == "--push" ]]; then
    if [[ -n $S3_DOCS_BUCKET_NAME ]]; then
        copy
        build
        upload_other_files
        modify_content_type css text/css
        modify_content_type js application/javascript
        cleanup
    else
        echo "Define S3_DOCS_BUCKET_NAME environment variable."
        exit 1
    fi
elif [[ "$1" == "--local" ]]; then
    copy
    run_local
    cleanup
elif [[ "$1" == "--tools" ]]; then
    tools
elif [[ "$1" == "--help" ]]; then
    help
else
    echo "Choose one of the options: --push, --local, --tools or --help"
    exit 1
fi