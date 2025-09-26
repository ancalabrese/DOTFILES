if status --is-interactive
    gcertstatus --check_remaining=60m >/dev/null 2>&1
    or gcert -s
end

set -gx GOPATH /Users/ancalabrese/go

function forwardAdb
    ssh -R 5037:127.0.0.1:5037 avengers.c.googlers.com
end

