#docker build \
#  --tag 'ms_generation:latest' \
#  --file backend/services/MS_Generation/Dockerfile .
cd ./backend/services/MS_Generation/common/api
./build.sh
cd ../../../../..