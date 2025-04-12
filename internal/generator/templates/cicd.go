// internal/generator/templates/cicd.go - Templates for CI/CD files
package templates

import "github.com/username/goprojectgen/internal/config"

// GitHubWorkflowTemplate returns the content of the GitHub Actions workflow file
func GitHubWorkflowTemplate(cfg config.ProjectConfig) string {
	return `name: Build and Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.21"

      - name: Install dependencies
        run: go mod download

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

      - name: Run tests
        run: go test -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt
          token: ${{ secrets.CODECOV_TOKEN }}
          fail_ci_if_error: false

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ` + cfg.Username + `/` + cfg.ProjectName + `:latest
            ` + cfg.Username + `/` + cfg.ProjectName + `:${{ github.sha }}
          cache-from: type=registry,ref=` + cfg.Username + `/` + cfg.ProjectName + `:latest
          cache-to: type=inline

  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: build
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up kubeconfig
        uses: azure/k8s-set-context@v1
        with:
          kubeconfig: ${{ secrets.KUBECONFIG }}

      - name: Deploy to Kubernetes
        run: |
          # Update image tag in deployment.yaml
          sed -i 's|` + cfg.Username + `/` + cfg.ProjectName + `:latest|` + cfg.Username + `/` + cfg.ProjectName + `:${{ github.sha }}|' deployments/kubernetes/deployment.yaml
          
          # Apply Kubernetes manifests
          kubectl apply -f deployments/kubernetes/
          
          # Wait for deployment to complete
          kubectl rollout status deployment/` + cfg.ProjectName + ` --timeout=2m
`
}
