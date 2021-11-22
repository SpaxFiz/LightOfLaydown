pnpm --prefix ./webapp install
pnpm --prefix ./webapp build

go mod vendor;
go build -o bin/web cmd/web/main.go;

sudo supervisorctl restart lol
