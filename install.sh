#! /usr/bin/env bash

set -eo pipefail

case $(uname) in
	Linux)
		OS_IDENTIFIED="Linux"
    ;;
  Darwin)
		OS_IDENTIFIED="Mac"
		;;
	Mac)
		OS_IDENTIFIED="Mac"
		;;
	*)
		echo "Could not determine the Operating System"
		exit 1
		;;
esac

ARCH_IDENTIFIED=$(uname -m)

OS=${OS:-${OS_IDENTIFIED}}
ARCH=${ARCH:-${ARCH_IDENTIFIED}}

echo ">>> OS: ${OS}"
echo ">>> ARCH: ${ARCH}"

echo ">>> Downloading"
tag=$(curl -s https://api.github.com/repos/arunvelsriram/sodexwoe/releases/latest | jq -r '.tag_name')
filename=sodexwoe_${tag}_${OS}_${ARCH}.tar.gz
curl -L https://github.com/arunvelsriram/sodexwoe/releases/download/${tag}/${filename} -o /tmp/${filename}
tar -C /tmp -xzf /tmp/${filename}
sudo mv /tmp/sodexwoe /usr/local/bin/sodexwoe
echo ">>> Installation completed"
