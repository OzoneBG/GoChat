cd ..
echo "Building main executable.."
go build -ldflags="-s -w" -o build/chat.exe
echo "Done"
# IN ORDER TO SKIP COMPRESSION COMMENT
cd buildtools/
echo "Compressing built executable.."
./upx -9 -o ../build/chatCompressed.exe ../build/chat.exe
echo "Done"
cd ../build
echo "Renaming executable.."
rm chat.exe
mv chatCompressed.exe chat.exe
echo "Done"
cd ..
# UNTIL HERE
echo "Copying dependancies.."
cp -a ./templates ./build
echo "Done"