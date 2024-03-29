pipeline {
    agent any

    tools { go 'Go 1.18' }

    environment {
        PATH = "${env.WORKSPACE}/bin:${env.HOME}/go/bin:${env.PATH}"
        GIT_CREDENTIAL_ID = 'wf-jenkins-github'
        GITHUB_TOKEN = credentials('GITHUB_TOKEN')
        REPO_NAME = 'wavefront-kubernetes-adapter'
    }

    parameters {
        choice(name: 'BUMP_COMPONENT', choices: ['patch', 'minor', 'major'], description: 'Specify a semver component to bump.')
    }

    stages {
        stage('Build and Run Tests') {
            steps {
                sh 'make fmt lint build test'
            }
        }

        stage('Create Bump Version PR') {
            environment {
                BUMP_COMPONENT = "${params.BUMP_COMPONENT}"
            }

            steps {
                sh 'git config --global user.email "svc.wf-jenkins@vmware.com"'
                sh 'git config --global user.name "svc.wf-jenkins"'
                sh 'git remote set-url origin https://${GITHUB_TOKEN}@github.com/wavefronthq/${REPO_NAME}.git'
                sh 'CGO_ENABLED=0 go install github.com/davidrjonas/semver-cli@latest'
                sh './scripts/update_release_version.sh -v $(cat release/VERSION) -s ${BUMP_COMPONENT}'
                sh './scripts/create_bump_version_pr.sh -v $(cat release/VERSION) -t ${GITHUB_TOKEN}'
            }
        }
    }

    // Notify only on null->failure or success->failure or failure->success
    post {
        failure {
            script {
                if(currentBuild.previousBuild == null) {
                    slackSend (channel: '#tobs-k8po-team', color: '#FF0000', message: "Bump version failed: <${env.BUILD_URL}|${env.JOB_NAME} [${env.BUILD_NUMBER}]>")
                }
            }
        }
        regression {
            slackSend (channel: '#tobs-k8po-team', color: '#FF0000', message: "Bump version failed: <${env.BUILD_URL}|${env.JOB_NAME} [${env.BUILD_NUMBER}]>")
        }
        fixed {
            slackSend (channel: '#tobs-k8po-team', color: '#008000', message: "Bump version fixed: <${env.BUILD_URL}|${env.JOB_NAME} [${env.BUILD_NUMBER}]>")
        }
        always {
            cleanWs()
        }
    }
}
