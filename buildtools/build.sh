cd ..
echo "Building main executable"
go build -ldflags="-s -w" -o build/chat.exe
cd buildtools/
echo "Compressing built executable"
./upx -9 -o ../build/chatCompressed.exe ../build/chat.exe
cd ../build
rm chat.exe
mv chatCompressed.exe chat.exe
read