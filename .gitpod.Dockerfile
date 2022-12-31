FROM gitpod/workspace-go:2022-12-28-23-50-51

# Install Homebrew (taken from https://github.com/gitpod-io/workspace-images/blob/main/chunks/tool-brew/Dockerfile)
ENV TRIGGER_REBUILD=4

RUN mkdir ~/.cache && /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
ENV PATH=/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin/:$PATH
ENV MANPATH="$MANPATH:/home/linuxbrew/.linuxbrew/share/man"
ENV INFOPATH="$INFOPATH:/home/linuxbrew/.linuxbrew/share/info"
ENV HOMEBREW_NO_AUTO_UPDATE=1

RUN sudo apt remove -y cmake \
    && brew install cmake

# Graphviz provides the dot tool for pprof
RUN brew install graphviz