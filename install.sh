#!/usr/bin/env bash
set -Eeuo pipefail

repo_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

if [[ ! -f "$repo_dir/.devcontainer/devcontainer.json" ]]; then
    printf 'Expected .devcontainer/devcontainer.json at the repository root.\n' >&2
    exit 1
fi

if [[ ! -f /etc/arch-release ]] || ! command -v pacman >/dev/null 2>&1; then
    printf '%s\n' \
        'This installer requires an Arch Linux environment.' \
        'A default Codespaces dotfiles repository cannot select the container image.' \
        'Use this repository as the workspace repository, or put its .devcontainer config in the workspace repository.' >&2
    exit 1
fi

if (( EUID == 0 )); then
    printf 'Run this installer as a regular user; it uses sudo for pacman.\n' >&2
    exit 1
fi

for config in .vimrc .zshrc; do
    if [[ -f "$repo_dir/$config" ]]; then
        cp --backup=numbered -- "$repo_dir/$config" "$HOME/$config"
    fi
done

sudo pacman -Syu --needed --noconfirm base-devel git go gopls go-tools vim

if ! command -v yay >/dev/null 2>&1; then
    build_dir="$(mktemp -d)"
    trap 'rm -rf -- "$build_dir"' EXIT

    git clone --depth=1 https://aur.archlinux.org/yay.git "$build_dir/yay"
    (
        cd "$build_dir/yay"
        makepkg -si --noconfirm
    )
fi

yay -S --needed --noconfirm vim-plug
go install golang.org/x/tools/gopls@latest
go install golang.org/x/tools/cmd/goimports@latest
vim +PlugInstall +qall
