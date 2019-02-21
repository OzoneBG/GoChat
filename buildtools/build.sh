cd ..
echo "Building main executable.."
go build -ldflags="-s -w" -o build/chat.exe
cd buildtools/
echo "Compressing built executable.."
./upx -9 -o ../build/chatCompressed.exe ../build/chat.exe
cd ../build
echo "Renaming executable.."
rm chat.exe
mv chatCompressed.exe chat.exe
cd ..
echo "Copying dependancies.."
cp -a ./templates ./build
read