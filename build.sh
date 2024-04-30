# This script is used to build SpaceJunk3000 project

#!/bin/bash

# Create build directory
if [ ! -d "bin" ]; then
    mkdir bin
fi

# Compile source files
go build .

# Move executable to bin directory
mv spacejunk3000 bin/

# Delete test user file
rm data/u-*.json

echo "Build completed successfully!"
echo "Executable is located in bin/ directory"

