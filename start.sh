# This script is used to start SpaceJunk3000 with a dummy door32.sys file

#!/bin/bash

# Create a dummy door32.sys file
echo "Creating dummy door32.sys file"
echo "0" > data/door32.sys
echo "0" >> data/door32.sys
echo "67600" >> data/door32.sys
echo "Talisman v0.51-dev" >> data/door32.sys
echo "1" >> data/door32.sys
echo "Robbie Whiting" >> data/door32.sys
echo "Alpha" >> data/door32.sys
echo "100" >> data/door32.sys
echo "90" >> data/door32.sys
echo "1" >> data/door32.sys
echo "1" >> data/door32.sys

# Start SpaceJunk3000
echo "Starting SpaceJunk3000 with dummy door32.sys file..."
bin/spacejunk3000 -door32 data/


