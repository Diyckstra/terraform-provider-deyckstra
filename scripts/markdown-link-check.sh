#!/usr/bin/env bash

# Local script runner for recursive markdown-link-check

set -euo pipefail

link_check_container="${LINK_CHECK_CONTAINER:-ghcr.io/tcort/markdown-link-check:stable}"
link_check_config="${LINK_CHECK_CONFIG:-.markdownlinkcheck.json}"
link_check_paths="${LINK_CHECK_PATHS:-docs}"
link_check_platform="${LINK_CHECK_PLATFORM:-linux/amd64}"

if [ ! -r "${link_check_config}" ]; then
  echo "Markdown link check config is not accessible: ${link_check_config}"
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "This script requires Docker to run markdown-link-check."
  echo ""
  echo "More information: https://github.com/tcort/markdown-link-check"
  exit 1
fi

if ! docker image inspect "${link_check_container}" >/dev/null 2>&1; then
  echo "==> Pulling ${link_check_container}..."
  docker pull --platform "${link_check_platform}" "${link_check_container}"
fi

echo "==> Checking Markdown links..."

read -r -a paths <<< "${link_check_paths}"

for path in "${paths[@]}"; do
  if [ ! -e "${path}" ]; then
    echo "Markdown link check path does not exist: ${path}"
    exit 1
  fi
done

error_file="$(mktemp -t markdown-link-check-errors.XXXXXX)"
output_file="$(mktemp -t markdown-link-check-output.XXXXXX)"

trap 'rm -f "${error_file}" "${output_file}"' EXIT

set +e
docker run --rm -i \
  --platform "${link_check_platform}" \
  -v "$(pwd):/github/workspace:ro" \
  -w /github/workspace \
  --entrypoint /usr/bin/find \
  "${link_check_container}" \
  "${paths[@]}" -type f \( -name "*.md" -o -name "*.markdown" \) -exec /src/markdown-link-check --config "${link_check_config}" --quiet --verbose {} \; \
  2>&1 | tee -a "${output_file}"
link_check_status=${PIPESTATUS[0]}
set -e

touch "${error_file}"
PREVIOUS_LINE=""
while IFS= read -r LINE; do
  if [[ $LINE = *"FILE"* ]]; then
    PREVIOUS_LINE=$LINE
    if [[ $(tail -1 "${error_file}") != *FILE* ]]; then
        echo -e "\n" >> "${error_file}"
        echo "$LINE" >> "${error_file}"
    fi
  elif [[ $LINE = *"✖"* ]] && [[ $PREVIOUS_LINE = *"FILE"* ]]; then
    echo "$LINE" >> "${error_file}"
  else 
    PREVIOUS_LINE=""
  fi
done < "${output_file}"

if [ "${link_check_status}" -ne 0 ] || grep -q "ERROR:" "${output_file}" || grep -q "✖" "${output_file}"; then
  echo -e "==================> MARKDOWN LINK CHECK FAILED <=================="
  if [[ $(tail -1 "${error_file}") = *FILE* ]]; then
    sed '$d' "${error_file}"
  else
    cat "${error_file}"
  fi
  printf "\n"
  echo -e "=================================================================="
  exit 1
else
  echo -e "==================> MARKDOWN LINK CHECK SUCCESS <=================="
  printf "\n"
  echo -e "[✔] All links are good!"
  printf "\n"
  echo -e "==================================================================="
fi
