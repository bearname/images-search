rm -rf frontend
rm -rf dist
rm -rf web
git clone https://github.com/col3name/images-search-frontend frontend
cd frontend
npm install
npm run build
cd ..
cp -r frontend/dist dist

