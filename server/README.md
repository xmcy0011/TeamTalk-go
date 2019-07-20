# server

## build

### mac

1.protobuf3.7
```bash
brew install protobuf
protoc --version
cd ../pb
chmod 777 *.sh
./create.sh && ./sync.sh
```