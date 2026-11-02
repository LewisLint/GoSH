FROM archlinux:base

RUN pacman -Syu --noconfirm \
    && pacman -S --noconfirm --needed git go shadow sudo zsh \
    && groupadd --force wheel \
    && useradd --create-home --shell /bin/zsh --groups wheel codespace \
    && printf 'codespace ALL=(ALL) NOPASSWD: ALL\n' > /etc/sudoers.d/codespace \
    && chmod 0440 /etc/sudoers.d/codespace \
    && pacman -Scc --noconfirm
