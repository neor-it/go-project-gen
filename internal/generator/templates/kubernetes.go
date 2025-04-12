// internal/generator/templates/kubernetes.go - Templates for Kubernetes files
package templates

import "github.com/username/goprojectgen/internal/config"

// KubernetesDeploymentTemplate returns the content of the deployment.yaml file
func KubernetesDeploymentTemplate(cfg config.ProjectConfig) string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + cfg.ProjectName + `
  labels:
    app: ` + cfg.ProjectName + `
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ` + cfg.ProjectName + `
  template:
    metadata:
      labels:
        app: ` + cfg.ProjectName + `
    spec:
      containers:
        - name: ` + cfg.ProjectName + `
          image: ` + cfg.Username + `/` + cfg.ProjectName + `:latest
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              name: http
          resources:
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 100m
              memory: 128Mi
          livenessProbe:
            httpGet:
              path: /api/v1/health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /api/v1/health
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
          env:
            - name: SERVER_PORT
              value: "8080"
          envFrom:
            - configMapRef:
                name: ` + cfg.ProjectName + `-config
`
}

// KubernetesServiceTemplate returns the content of the service.yaml file
func KubernetesServiceTemplate(cfg config.ProjectConfig) string {
	return `apiVersion: v1
kind: Service
metadata:
  name: ` + cfg.ProjectName + `
  labels:
    app: ` + cfg.ProjectName + `
spec:
  selector:
    app: ` + cfg.ProjectName + `
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  type: ClusterIP
`
}

// KubernetesConfigMapTemplate returns the content of the configmap.yaml file
func KubernetesConfigMapTemplate(cfg config.ProjectConfig) string {
	return `apiVersion: v1
kind: ConfigMap
metadata:
  name: ` + cfg.ProjectName + `-config
data:
  SERVER_PORT: "8080"
  SERVER_READ_TIMEOUT: "10s"
  SERVER_WRITE_TIMEOUT: "10s"
  LOGGING_LEVEL: "info"
  LOGGING_FORMAT: "json"
  SHUTDOWN_TIMEOUT: "10s"
`
}
