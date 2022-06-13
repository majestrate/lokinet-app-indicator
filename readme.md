lokinet app indicator for linux

uses systemd (for now)
written in go (requires golang >= 1.17)

to build run:

    $ go get -u github.com/majestrate/lokinet-app-indicator
    $ go install github.com/majestrate/lokinet-app-indicator

to install:

    $ go build github.com/majestrate/lokinet-app-indicator
    $ mkdir -p $HOME/.local/bin/
    $ mkdir -p $HOME/.local/share/applications/
    $ cp lokinet-app-indicator $HOME/.local/bin/
    $ cp lokinet-app-indicator.desktop $HOME/.local/share/applications/
