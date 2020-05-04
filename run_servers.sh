gnome-terminal --window -e '/bin/bash -c "go run kgv.go -- 8000 8001 8002;exec bash"'
gnome-terminal --window -e '/bin/bash -c "go run kgv.go -- 8001 8000 8002;exec bash"'
gnome-terminal --window -e '/bin/bash -c "go run kgv.go -- 8002 8000 8001;exec bash"'
