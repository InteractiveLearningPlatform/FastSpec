pipeline {
  agent { label 'self-hosted-linux && docker' }

  options {
    timestamps()
    ansiColor('xterm')
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '20'))
  }

  parameters {
    booleanParam(name: 'RUN_DEPLOY', defaultValue: false, description: 'Deploy the Helm chart after packaging.')
    string(name: 'DEPLOY_NAMESPACE', defaultValue: 'speclist', description: 'Kubernetes namespace for Helm deployment.')
    string(name: 'SPECLIST_DOMAIN', defaultValue: 'speclist.example.test', description: 'Ingress domain used by Helm.')
  }

  environment {
    CI = 'true'
    CARGO_TERM_COLOR = 'always'
    REGISTRY = 'ghcr.io/interactivelearningplatform/fastspec'
    API_IMAGE = "${env.REGISTRY}/speclist-api:${env.BUILD_NUMBER}"
    WEB_IMAGE = "${env.REGISTRY}/speclist-web:${env.BUILD_NUMBER}"
  }

  stages {
    stage('Prep') {
      steps {
        sh 'rustc --version'
        sh 'cargo --version'
        sh 'go version'
        sh 'node --version'
        sh 'npm --version'
        sh 'python3 --version'
        sh 'docker --version'
        sh 'helm version --short'
      }
    }

    stage('Rust') {
      steps {
        sh 'cargo fmt --all -- --check'
        sh 'cargo clippy --workspace --all-targets --all-features -- -D warnings'
        sh 'cargo test --workspace'
      }
    }

    stage('Go API') {
      steps {
        dir('apps/speclist-api') {
          sh 'go test ./...'
        }
      }
    }

    stage('Web App') {
      steps {
        dir('apps/speclist-web') {
          sh 'npm ci'
          sh 'npm run build'
        }
      }
    }

    stage('Platform Ops Validation') {
      steps {
        sh '''
          cat > deploy/compose/.env.ci <<EOF
TRAEFIK_DOMAIN=speclist.example.test
TRAEFIK_ACME_EMAIL=ops@example.test
TRAEFIK_TLS_ENABLED=true
SPECLIST_API_IMAGE=${API_IMAGE}
SPECLIST_WEB_IMAGE=${WEB_IMAGE}
POSTGRES_DB=speclist
POSTGRES_USER=speclist
POSTGRES_PASSWORD=Postgres_Str0ng!123
CLICKHOUSE_DB=speclist
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=Clickhouse_Str0ng!123
VALKEY_PASSWORD=Valkey_Str0ng!123
QDRANT_API_KEY=Qdrant_Str0ng!123
CROWDSEC_ENROLL_KEY=CrowdsecEnroll_Str0ng!123
CROWDSEC_BOUNCER_KEY=CrowdsecBouncer_Str0ng!123
VAULT_ADDR=https://vault.example.test
VAULT_ROLE=speclist-platform
SPECLIST_REPO_ROOT=/workspace
EOF
          cp deploy/compose/.env.ci deploy/compose/.env
          python3 deploy/compose/preflight/validate_config.py --env-file deploy/compose/.env.ci --compose-file deploy/compose/compose.platform.yml
          docker compose -f deploy/compose/compose.platform.yml --env-file deploy/compose/.env.ci config > /tmp/speclist-compose.yaml
          helm lint deploy/helm/speclist-platform
        '''
      }
      post {
        always {
          sh 'rm -f deploy/compose/.env deploy/compose/.env.ci'
        }
      }
    }

    stage('Security') {
      steps {
        sh 'cargo install cargo-audit --locked'
        sh 'cargo audit'
        dir('apps/speclist-api') {
          sh 'go install golang.org/x/vuln/cmd/govulncheck@latest'
          sh '$(go env GOPATH)/bin/govulncheck ./...'
        }
        dir('apps/speclist-web') {
          sh 'npm audit --omit=dev --audit-level=high'
        }
        sh 'docker run --rm -v "$PWD:/workspace:ro" aquasec/trivy:0.57.1 config --exit-code 1 --severity HIGH,CRITICAL /workspace/deploy'
      }
    }

    stage('Package') {
      when {
        anyOf {
          buildingTag()
          branch 'main'
          expression { return params.RUN_DEPLOY }
        }
      }
      steps {
        sh '''
          rm -rf dist
          mkdir -p dist
          cargo build --release -p fastspec-cli
          cp target/release/fastspec dist/fastspec
          tar -C dist -czf dist/fastspec-cli-linux-x86_64.tar.gz fastspec
        '''
        dir('apps/speclist-api') {
          sh '''
            mkdir -p dist
            CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/speclist-api ./cmd/speclist-api
            tar -C dist -czf dist/speclist-api-linux-x86_64.tar.gz speclist-api
          '''
        }
        dir('apps/speclist-web') {
          sh '''
            mkdir -p release
            tar -C dist -czf release/speclist-web-dist.tar.gz .
          '''
        }
        sh 'helm package deploy/helm/speclist-platform --destination dist'
        sh 'docker build -t "$API_IMAGE" apps/speclist-api'
        sh 'docker build -t "$WEB_IMAGE" apps/speclist-web'
      }
    }

    stage('Publish Images') {
      when {
        anyOf {
          buildingTag()
          expression { return params.RUN_DEPLOY }
        }
      }
      steps {
        withCredentials([usernamePassword(credentialsId: 'ghcr-push', usernameVariable: 'GHCR_USER', passwordVariable: 'GHCR_TOKEN')]) {
          sh '''
            echo "$GHCR_TOKEN" | docker login ghcr.io -u "$GHCR_USER" --password-stdin
            docker push "$API_IMAGE"
            docker push "$WEB_IMAGE"
          '''
        }
      }
    }

    stage('Deploy') {
      when {
        expression { return params.RUN_DEPLOY }
      }
      steps {
        withCredentials([file(credentialsId: 'kubeconfig-speclist', variable: 'KUBECONFIG_FILE')]) {
          sh '''
            export KUBECONFIG="$KUBECONFIG_FILE"
            helm upgrade --install speclist-platform deploy/helm/speclist-platform \
              --namespace "${DEPLOY_NAMESPACE}" \
              --create-namespace \
              --set global.domain="${SPECLIST_DOMAIN}" \
              --set api.image.repository="${REGISTRY}/speclist-api" \
              --set api.image.tag="${BUILD_NUMBER}" \
              --set web.image.repository="${REGISTRY}/speclist-web" \
              --set web.image.tag="${BUILD_NUMBER}"
          '''
        }
      }
    }
  }

  post {
    always {
      archiveArtifacts artifacts: 'dist/*.tar.gz,dist/*.tgz,apps/speclist-api/dist/*.tar.gz,apps/speclist-web/release/*.tar.gz', allowEmptyArchive: true
      junit testResults: '**/junit*.xml', allowEmptyResults: true
      cleanWs(deleteDirs: true, notFailBuild: true)
    }
  }
}
