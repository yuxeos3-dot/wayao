.PHONY: build run dev frontend clean deploy

# Build Go server
build:
	CGO_ENABLED=1 go build -o brandsite-server ./cmd/server/

# Run server locally
run: build
	./brandsite-server

# Dev mode (auto-restart not included, use air or similar)
dev: build
	DATA_DIR=./data TEMPLATE_DIR=./templates ADMIN_DIR=./frontend/dist HUGO_PATH=hugo ./brandsite-server

# Build frontend
frontend:
	cd frontend && npm install && npm run build

# Clean
clean:
	rm -f brandsite-server
	rm -rf frontend/dist
	rm -rf data/builds

# Deploy to remote server
# Usage: make deploy SERVER=1.2.3.4
deploy: build frontend
ifndef SERVER
	$(error SERVER is required. Usage: make deploy SERVER=1.2.3.4)
endif
	ssh root@$(SERVER) 'mkdir -p /opt/brandsite/frontend/dist /opt/brandsite/templates /opt/brandsite/scripts /opt/brandsite/nginx /opt/brandsite/deploy'
	rsync -avz --exclude='.git' --exclude='node_modules' --exclude='data/*.db' \
		brandsite-server root@$(SERVER):/opt/brandsite/
	rsync -avz templates/ root@$(SERVER):/opt/brandsite/templates/
	rsync -avz frontend/dist/ root@$(SERVER):/opt/brandsite/frontend/dist/
	rsync -avz scripts/ root@$(SERVER):/opt/brandsite/scripts/
	rsync -avz nginx/ root@$(SERVER):/opt/brandsite/nginx/
	rsync -avz deploy/ root@$(SERVER):/opt/brandsite/deploy/
	ssh root@$(SERVER) 'systemctl restart brandsite'
