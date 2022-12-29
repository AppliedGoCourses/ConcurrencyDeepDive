FROM gitpod/workspace-go:2022-12-28-23-50-51

# Graphviz provides the dot tool for pprof
RUN brew install graphviz