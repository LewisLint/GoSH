# Set up the prompt
#autoload -Uz promptinit
#promptinit
#prompt adam1

#setopt histignorealldups sharehistory

# Use emacs keybindings even if our EDITOR is set to vi
# Keep 1000 lines of history within the shell and save it to ~/.zsh_history:
HISTSIZE=1000
SAVEHIST=1000
HISTFILE=~/.zsh_history

setopt PROMPT_SUBST
# Use modern completion system
autoload -Uz compinit
compinit

zstyle ':completion:*' auto-description 'specify: %d'
zstyle ':completion:*' completer _expand _complete _correct _approximate
zstyle ':completion:*' format 'Completing %d'
zstyle ':completion:*' group-name ''
zstyle ':completion:*' menu select=2
eval "$(dircolors -b)"
#export LS_COLORS='di=38;5;110:ln=38;5;109:ex=38;5;108:pi=38;5;180:so=38;5;139:bd=38;5;180:cd=38;5;180:or=38;5;131:mi=38;5;131'
#export LS_COLORS="${LS_COLORS}:di=1;38;5;75:ow=1;38;5;75;48;5;236:tw=1;38;5;75;48;5;236:st=1;38;5;75;48;5;236"
export LS_COLORS="${LS_COLORS}:di=1;38;5;75:ow=1;38;5;139:tw=1;38;5;139:st=1;38;5;139:*.go=1;38;5;75"
zstyle ':completion:*:default' list-colors ${(s.:.)LS_COLORS}
zstyle ':completion:*' list-colors ''
zstyle ':completion:*' list-prompt %SAt %p: Hit TAB for more, or the character to insert%s
zstyle ':completion:*' matcher-list '' 'm:{a-z}={A-Z}' 'm:{a-zA-Z}={A-Za-z}' 'r:|[._-]=* r:|=* l:|=*'
zstyle ':completion:*' menu select=long
zstyle ':completion:*' select-prompt %SScrolling active: current selection at %p%s
zstyle ':completion:*' use-compctl false
zstyle ':completion:*' verbose true

zstyle ':completion:*:*:kill:*:processes' list-colors '=(#b) #([0-9]#)*=0=01;31'
zstyle ':completion:*:kill:*' command 'ps -u $USER -o pid,%cpu,tty,cputime,cmd'

alias ll="ls -la"
alias v="vi"
alias '^l'="clear"
alias '~'="cd ~"
alias 'vzsh'="v ~/.zshrc && source ~/.zshrc"
alias cp='cp -v'
alias mv='mv -v'
alias img='ristretto'
alias ig='ristretto *'

#cd() {
#    builtin cd $@
#    ls
#}
#bindkey -M viins 'jj' vicmd

#bindkey -s '^A' 'sudo '
#Enable Vi Mode
bindkey -v

# Enable Zsh color support
autoload -U colors && colors


#PROMPT="%{$fg[cyan]%}%n%{$reset_color%} %{$fg[yellow]%}%~%{$reset_color%} %# "
# Fix backspace/delete refusing to clear old text in vi-insert mode
bindkey -M viins '^?' backward-delete-char
bindkey -M viins '^H' backward-delete-char
echo -ne "\e[6 q"
# Map jj to escape to normal mode
bindkey -M viins 'jj' vi-cmd-mode


# my bullshit:
# Enable dynamic prompt updating
setopt prompt_subst

# Set the default starting state (Blue I)
# Force the cursor to be a vertical bar on startup
printf '\e[5 q'

# Reset cursor back to a bar whenever a new prompt line draws
_set_cursor_bar() {
    printf '\e[5 q'
}
autoload -Uz add-zsh-hook

add-zsh-hook precmd _set_cursor_bar
# Detect mode changes and swap the RPROMPT character/color


RPROMPT="%F{blue}\x1b[1mI%f"


function zle-keymap-select() {
  case $KEYMAP in
    vicmd)      RPROMPT="%F{green}N%f" ;; # Green Normal mode
    viins|main) RPROMPT="%F{blue}I%f"  ;; # Blue Insert mode
    visual)     RPROMPT="%F{244}V%f"  ;; # Gray Visual mode
  esac
  zle reset-prompt
}
zle -N zle-keymap-select

# Ensure the prompt resets back to Insert mode for a clean new line
function zle-line-init() {
  RPROMPT="%F{blue}I%f"
  zle reset-prompt
}
zle -N zle-line-init
KEYTIMEOUT=15

bindkey -M vicmd '?' history-incremental-search-backward
alias ytdlp='yt-dlp -f "best[height<=720]" --impersonate chrome --cookies-from-browser chromium'
alias py="python"
alias c=" clear"








alias clear="clear"
# Enable prompt expansion so functions run inside the prompt

# Enable vim keybindings

# Function to build the conditional arrow and exit code
build_arrow() {
    local exit_code=$?
    if [ $exit_code -eq 0 ]; then
        # Success: Green arrow (Replace '❯' with any symbol you prefer)
        echo "%F{green}❯%f"
    else
        # Failure: Red exit code followed by a Red arrow
        echo "%F{red}$exit_code ❯%f"
    fi
}

# Construct the final prompt


autoload -Uz vcs_info
precmd() { vcs_info }
zstyle ':vcs_info:git:*' formats '(%b)'

PROMPT='%F{#9CCFD8}%B%~%b%f $(build_arrow) '
setopt SHARE_HISTORY
setopt HIST_IGNORE_DUPS
setopt HIST_REDUCE_BLANKS
setopt INC_APPEND_HISTORY
setopt AUTO_CD
setopt AUTO_PUSHD
setopt PUSHD_IGNORE_DUPS


# aliases:
alias ls='ls -aFBh --color=auto'
alias ll='ls -l'
if [ "$USER" = "root" ]; then
    # --- Root Warning ---
    print -P ""
    print -P "  %F{196}WARNING: You are now logged in as root%f"
    print -P "  %F{244}All commands you run will affect the root FS%f"
    print -P "  %F{244}Current Node: %n on %m%f"
    print -P ""

    # --- Root-Specific Prompt ---
    build_arrow() {
        local exit_code=$?
        if [ $exit_code -eq 0 ]; then
            echo "%F{196}❯%f"
        else
            echo "%F{124} $exit_code ❯%f"
        fi
    }

    PROMPT='%F{196}%n%f%F{244}@%m:%F{250}%~%f $(build_arrow) '

    # --- Root Vi Mode ---
    bindkey -v
    KEYTIMEOUT=15

    bindkey -M viins 'jj' vi-cmd-mode
    bindkey -M viins '^?' backward-delete-char
    bindkey -M viins '^H' backward-delete-char
    bindkey -M vicmd '?' history-incremental-search-backward

    # --- Root Cursor ---
    _set_cursor_bar() { printf '\e[5 q' }

    autoload -Uz add-zsh-hook
    add-zsh-hook precmd _set_cursor_bar

    # --- Root Vi Mode Indicator ---
    function zle-keymap-select() {
        case $KEYMAP in
            vicmd)      RPROMPT="%F{226}N%f" ;;
            viins|main) RPROMPT="%F{196}I%f" ;;
            visual)     RPROMPT="%F{244}V%f" ;;
        esac
        zle reset-prompt
    }
    zle -N zle-keymap-select

    function zle-line-init() {
        RPROMPT="%F{196}I%f"
        _set_cursor_bar
        zle reset-prompt
    }
    zle -N zle-line-init

    # --- Root Safety Aliases ---
    alias rm='rm -iv'
    alias cp='cp -v'
    alias mv='mv -v'
    alias mkdir='mkdir -pv'
    alias chown='chown --preserve-root -v'
    alias chmod='chmod --preserve-root -v'
    alias chgrp='chgrp --preserve-root -v'
    alias ln='ln -iv'
    alias rmdir='rmdir -v'
    alias shred='shred -v'
    alias dd='dd status=progress'

    # --- Root Power Command Warnings ---
    alias reboot='echo "WARNING: Rebooting system as root"; command reboot'
    alias poweroff='echo "WARNING: Powering off system as root"; command poweroff'
    alias shutdown='echo "WARNING: Shutting down system as root"; command shutdown'

    # --- Root Sudo Prompt ---
    export SUDO_PROMPT="[ROOT ACCESS] Password for %u: "
fi



source /usr/share/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
# Rosé Pine color mapping for zsh-syntax-highlighting
# type sudo pacman -S zsh-history-substring-search zsh-syntax-highlighting
typeset -A ZSH_HIGHLIGHT_STYLES

# Base text colors (Foam / Iris / Pine)
ZSH_HIGHLIGHT_STYLES[command]='fg=116,bold'  
ZSH_HIGHLIGHT_STYLES[alias]='fg=116'         
ZSH_HIGHLIGHT_STYLES[builtin]='fg=182,italic'
ZSH_HIGHLIGHT_STYLES[function]='fg=182,bold' 
ZSH_HIGHLIGHT_STYLES[precommand]='fg=31,bold'

# Misspelled or invalid commands (Love)
ZSH_HIGHLIGHT_STYLES[unknown-token]='fg=210'

# Arguments and Options (Gold / Rose / Muted)
ZSH_HIGHLIGHT_STYLES[single-hyphen-option]='fg=216,bold'
ZSH_HIGHLIGHT_STYLES[double-hyphen-option]='fg=216,bold'
ZSH_HIGHLIGHT_STYLES[string]='fg=181,bold'              
ZSH_HIGHLIGHT_STYLES[path]='fg=116,underline'

# Comments (Muted)
#ZSH_HIGHLIGHT_STYLES[comment]='fg=243,italic'



source /usr/share/zsh/plugins/zsh-history-substring-search/zsh-history-substring-search.zsh

bindkey '^[k' history-substring-search-up
bindkey '^[j' history-substring-search-down
bindkey '˚' history-substring-search-up
bindkey '∆' history-substring-search-down
# ZLE BULLSHIT -- 
autoload -Uz edit-command-line
zle -N edit-command-line
bindkey '^[;' edit-command-line
bindkey '…' edit-command-line


export EDITOR=vim
export VISUAL=vim

last-arg() {
    LBUFFER+="${(Q)${(z)history[$((HISTCMD-1))]}[-1]}"
}

zle -N last-arg
bindkey '\e.' last-arg
bindkey '≥' last-arg

setopt HIST_IGNORE_SPACE
zshaddhistory() {
    [[ "$1" == "clear" ]] && return 1
    return 0
}
