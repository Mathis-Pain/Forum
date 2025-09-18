
docker build -t forum .
docker run -p 5090:5080 -v $(pwd)/data:/data forum