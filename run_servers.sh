# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- tm 0 0;exec bash"'
# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- tm 0 1;exec bash"'
# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- tm 0 2;exec bash"'

gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 0 0;exec bash"'
gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 0 1;exec bash"'
gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 0 2;exec bash"'

# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 1 0;exec bash"'
# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 1 1;exec bash"'
# gnome-terminal --window -e '/bin/bash -c "go run overdb.go -- kv 1 2;exec bash"'
