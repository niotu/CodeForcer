#!/bin/sh

# Debug: List files to verify server binary is present
ls -al

# Start the Go server in the background
./server &

# Start the React frontend
npm start
