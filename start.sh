#!/bin/sh

# Debug: List files to verify server binary is present
ls -al

# Start the Go server in the background
./server &

# Start the React frontend
#npm run start --prefix ./web/front
cd ./web/front
npm run build
npm install -g serve
serve -s build -l 80