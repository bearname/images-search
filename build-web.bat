@echo off
rmdir /s web
mkdir web
cd frontend
yarn build
cd ..
xcopy D:\Users\mikha\go\src\photofinish\frontend\dist D:\Users\mikha\go\src\photofinish\web /E /H /C
