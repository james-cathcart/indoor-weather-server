node('workers') {

    def weatherEnv
    def logLevel

    switch (env.BRANCH_NAME) {
    case 'dev':
        weatherEnv = 'dev'
        logLevel = 'info'
        break
    case 'uat':
        weatherEnv = 'uat'
        logLevel = 'warn'
        break
    case 'master':
        weatherEnv = 'prd'
        logLevel = 'error'
        break
    default:
        weatherEnv = 'dev'
        logLevel = 'info'
    }

    withEnv([
            "IMAGE_NAME=apollosamp/weather-server",
            "REGISTRY_URL=https://index.docker.io/v1/",
            "DOCKER_BUILDKIT=1",
        ]) {

        stage('Checkout') {
            cleanWs()
            scmInfo = checkout scm
            env.GIT_COMMIT = scmInfo.GIT_COMMIT
        }

        def imageTest = docker.build("${env.IMAGE_NAME}-test", "-f test.Dockerfile .")

        stage('Compile & Unit Test') {
            imageTest.inside{
                sh 'make'
            }
        }

        stage('Build') {

            if (env.BRANCH_NAME == 'master') {
                docker.build(env.IMAGE_NAME, "-f Dockerfile .")
            } else {
                sh 'echo "skipping container image build for feature branch"'
            }
        }

        stage('Push') {
            docker.withRegistry(env.REGISTRY_URL, 'Docker Hub RegCreds') {

                if (env.BRANCH_NAME == 'master') {
                    env.IMAGE_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT}-${weatherEnv}"
                    docker.image(env.IMAGE_NAME).push(env.IMAGE_TAG)
                } else {
                    sh 'echo "skipping image push for feature branch"'
                }
            }
        }

        stage('Deploy') {
            if (env.BRANCH_NAME == 'master') {
                sh "sed -i 's#image:#image: ${env.IMAGE_NAME}:${env.IMAGE_TAG}#' kube/${weatherEnv}/deployment.yml"
                sh "kubectl apply -f kube/${weatherEnv}"
            } else {
                sh 'echo "skipping deployment for feature branch"'
            }
        }
    }
}