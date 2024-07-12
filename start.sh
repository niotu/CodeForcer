#!/bin/sh

# Debug: List files to verify server binary is present
ls -al

# Start the React frontend
npm run start
#npm run build
#npm install -g serve
#serve -s build -l 80

# Start the Go server in the background
./server &