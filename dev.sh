function exampleRestart() {
echo "=================>"
 
#  export DEV_PRESETS=1
source dev_env &&  go run main.go  
}

export -f exampleRestart

find . -name "*.go" | entr -r bash -c "exampleRestart"
