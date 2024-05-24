@Library(['jenkins-common-pipelines', 'pact-shared-library']) _

def GitHubStatus = "NEX: SED Visual Solutions/Ci Project Build"
def GtaStatusCheck = "daily testing"
def SdlStatusCheck = "SDL scans"
def relativeDir = 'vp-build'
pipeline {
    agent {
        label 'video-production-build'
    }
    environment {
        ARTIFACTORY_URL = 'https://af01p-igk.devtools.intel.com/artifactory'
        GTA_ASSET_NAME = "ci"
        JOB_BUILD_ID  = "${BUILD_ID}"
        JOB_FILE_ARTIFACTORY_PREFIX = 'video_production_igk-igk-local/builds'
        ARTIFACTORY_SECRET = credentials('ed6200e9-7a86-4cf8-b2f3-2cea840a23f5')
        UPLOAD_DIR = '/tmp/gtax-upload'
        IMAGE_TAG_NAME = "video_production_image_${JOB_BUILD_ID}"
        JENKINS_INVOKE_TOKEN = credentials('c783a259-c17c-4a89-90e8-61bdfee916a4')
    }
    triggers {
        GenericTrigger(
            genericVariables: [
                /* params for GITHUB PR */
                [key: 'GITHUB_PR_TRIGGER_ACTION', value: '$.action'],
                [key: 'GITHUB_PR_TRIGGER_SENDER_AUTHOR', value: '$.sender.login'],
                [key: 'GITHUB_PR_TRIGGER_SENDER_ID', value: '$.sender.id'],
                [key: 'GITHUB_PR_COMMIT_STATUS', value: '$.pull_request._links.statuses.href'],
                [key: 'GITHUB_COMMENT_PR_URL', value: '$.issue.pull_request.url'],
                [key: 'GITHUB_COMMENT_BODY', value: '$.comment.body'],
                [key: 'GITHUB_PR_TARGET_BRANCH', value: '$.pull_request.base.ref'],
                [key: 'GITHUB_PR_SOURCE_BRANCH', value: '$.pull_request.head.ref'],
                [key: 'GITHUB_PR_BODY', value: '$.pull_request'],
                [key: 'GITHUB_PR_BODY_MERGED', value: '$.pull_request.merged'],
                [key: 'GITHUB_PR_TITLE', value: '$.pull_request.title'],
                [key: 'GITHUB_PR_URL', value: '$.pull_request.html_url'],
                [key: 'GITHUB_PR_SOURCE_REPO_OWNER', value: '$.pull_request.base.user.login'],
                [key: 'GITHUB_PR_HEAD_SHA', value: '$.pull_request.head.sha'],
                [key: 'GITHUB_PR_NUMBER', value: '$.pull_request.number'],
                [key: 'GITHUB_PR_STATE', value: '$.pull_request.state'],
                [key: 'GITHUB_PR_COMMENT_BODY', value: '$.pull_request.body'],
                [key: 'GITHUB_PR_LABELS', value: '$.pull_request.labels'],
                /* params for GITHUB BRANCH */
                [key: 'GITHUB_BRANCH_NAME', value: '$.ref'],
                [key: 'GITHUB_BRANCH_TYPE', value: '$.ref_type'],
                /* params for GITHUB REPO */
                [key: 'GITHUB_REPO_GIT_URL', value: '$.repository.git_url'],
                [key: 'GITHUB_REPO_SSH_URL', value: '$.repository.ssh_url'],
                [key: 'GITHUB_REPO_URL', value: '$.repository.clone_url'],
                [key: 'GITHUB_REPO_FULL_NAME', value: '$.repository.full_name'],
                [key: 'GITHUB_REPO_STATUS_URL', value: '$.pull_request.url']

            ],
            genericHeaderVariables: [
                [key: 'X-GitHub-Event', regexpFilter: '']
            ],
            token: '${JENKINS_INVOKE_TOKEN}',
            causeString: 'Triggered on GitHub webhook',
            printContributedVariables: false,
            printPostContent: true,
            silentResponse: false,
            regexpFilterText: '${GITHUB_PR_TRIGGER_ACTION}',
            regexpFilterExpression: '^(opened|reopened|synchronize).*$'
        )
    }
    options {
        timestamps()
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('Build'){
            steps {
                    script {
                        cleanWs()
                        if(env.GITHUB_PR_URL){
                                currentBuild.description = 'Trigged by <a href="' + GITHUB_PR_URL + '" target="_blank">' + GITHUB_REPO_FULL_NAME + '</a>.'
                                println 'Trigged by ' + GITHUB_PR_URL + '.'
                                withCredentials([usernamePassword(credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                                    status_check(
                                        gh_token : PASSWORD,
                                        commit : GITHUB_PR_HEAD_SHA,
                                        pull_request : GITHUB_REPO_STATUS_URL,
                                        name : GitHubStatus,
                                        state : "pending"
                                    )
                                    status_check(
                                        gh_token : PASSWORD,
                                        commit : GITHUB_PR_HEAD_SHA,
                                        pull_request : GITHUB_REPO_STATUS_URL,
                                        name : GtaStatusCheck,
                                        message: "GTA tests",
                                        description: "GTA tests",
                                        state : "pending"
                                    )
                                    status_check(
                                        gh_token : PASSWORD,
                                        commit : GITHUB_PR_HEAD_SHA,
                                        pull_request : GITHUB_REPO_STATUS_URL,
                                        name : SdlStatusCheck,
                                        message: "SDL scans",
                                        description: "SDL scans",
                                        state : "pending"
                                    )
                                }
                        }
                        
                        def repo = 'https://github.com/intel-innersource/libraries.media.encoding.svt-jpeg-xs.git'
                        def jpegRelativeDir = 'libraries.media.encoding.svt-jpeg-xs'
                        def branches = [[name: 'main']]
                        def userRemoteConfigs = [[
                            credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
                            url: repo
                        ]]
                        checkout([
                            $class: 'GitSCM',
                            branches: branches,
                            extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: jpegRelativeDir]],
                            userRemoteConfigs: userRemoteConfigs
                        ])
                    }
                    script {
                            def repo = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline.git'
                            def branches = [[name: 'main']]
                            def userRemoteConfigs = [[
                                credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
                                url: repo
                            ]]
                            if(env.GITHUB_PR_URL){
                                currentBuild.description = 'Trigged by <a href="' + GITHUB_PR_URL + '" target="_blank">' + GITHUB_REPO_FULL_NAME + '</a>.'
                                println 'Trigged by ' + GITHUB_PR_URL + '.'
                                branches = [[name: "origin/pull/${GITHUB_PR_NUMBER}/merge"]]
                                userRemoteConfigs[0]['refspec'] = "+refs/pull/${GITHUB_PR_NUMBER}/merge:refs/remotes/origin/pull/${GITHUB_PR_NUMBER}/merge"
                            }
                                 
                            checkout([
                                $class: 'GitSCM',
                                branches: branches,
                                extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: relativeDir]],
                                userRemoteConfigs: userRemoteConfigs
                            ])
                            dir(relativeDir){
                                sh """
                                    git config --global --add safe.directory \$(pwd)
                                    mv ${WORKSPACE}/libraries.media.encoding.svt-jpeg-xs libraries.media.encoding.svt-jpeg-xs
                                    git config --global --add safe.directory libraries.media.encoding.svt-jpeg-xs
                                    docker build \
                                        --build-arg http_proxy=http://proxy-dmz.intel.com:912 \
                                        --build-arg https_proxy=http://proxy-dmz.intel.com:912 \
                                        -t \"${IMAGE_TAG_NAME}\" \
                                        -f Dockerfile .
                                """
                            }
                    }
                }
        }
        stage("Pack"){
            steps{
                script {
                    dir(relativeDir){
                        sh """
                            rm -rf \${UPLOAD_DIR}/
                            mkdir -p \${UPLOAD_DIR}
                            docker save -o \${UPLOAD_DIR}/\"${IMAGE_TAG_NAME}\".tar.gz \"${IMAGE_TAG_NAME}\"
                            tar czf \${UPLOAD_DIR}/tests_repo_${JOB_BUILD_ID}.tar.gz .
                        """
                    }
                }
            }
        }
        stage("scans"){
            parallel {
                stage("Hadolint"){
                    steps{
                        script{
                            dir(relativeDir){
                                sh """ 
                                    jenkins/scripts/hadolint.sh
                                """
                                archiveArtifacts allowEmptyArchive: true, artifacts: "Hadolint/hadolint-Dockerfile*"
                            }
                        }
                    } 
                }
                stage("Trivy"){
                    steps{
                        script{
                            dir(relativeDir){
                                sh """ 
                                    jenkins/scripts/trivy.sh \${UPLOAD_DIR}/\"${IMAGE_TAG_NAME}\".tar.gz
                                """
                                archiveArtifacts allowEmptyArchive: true, artifacts: "Trivy/*"
                            }
                        }
                    } 
                }
                stage("Schellcheck"){
                    steps{
                        script{
                            dir(relativeDir){
                                sh """ 
                                    jenkins/scripts/shellcheck.sh
                                """
                                archiveArtifacts allowEmptyArchive: true, artifacts: "shellcheck_logs/*"
                            }
                        }
                    } 
                }
                stage("McAfee"){
                    steps{
                        script{
                            dir(relativeDir){
                                sh """ 
                                    DOCKER_IMAGE_NAME="amr-registry.caas.intel.com/owr/abi_lnx:3.0.0"
                                    docker run --rm -v \$(pwd):/opt/ \${DOCKER_IMAGE_NAME} /bin/bash -c "cd /opt/; jenkins/scripts/mcafee_scan.sh"
                                """
                                archiveArtifacts allowEmptyArchive: true, artifacts: "Malware/*"
                            }
                        }
                    } 
                }
            }
        }
        stage('Upload'){
            steps {
                script {
                    def pr_head_sha = (env.GITHUB_PR_HEAD_SHA ? env.GITHUB_PR_HEAD_SHA : "manual_trigger")
                    dir(relativeDir){
                        sh """
                            set -x
                            GIT_COMMIT=\$(git log --pretty=%H -n 1)
                            GIT_SHA=${pr_head_sha}
                            GIT_REPO_NAME=\$(echo ${env.GITHUB_REPO_FULL_NAME} | cut -d \"/\" -f2)
                            echo '{\"video_production_sdb\":[\"'\${GIT_COMMIT}'\"], \
                                       \"build.component_revision\":[\"'\${GIT_SHA}'\"], \
                                       \"build.component_project\":[\"'\${GIT_REPO_NAME}'\"]}' > recipe.json
                            gta-asset push --no-archive --properties-file recipe.json \
                                       -u sys_vsval \
                                       -p \"${ARTIFACTORY_SECRET}\" \
                                       --root-url \"${env.ARTIFACTORY_URL}\" \
                                       \"${env.JOB_FILE_ARTIFACTORY_PREFIX}\" \
                                       \"${env.GTA_ASSET_NAME}\" \
                                       \"${env.JOB_BASE_NAME}_${env.JOB_BUILD_ID}\" \
                                       \"${UPLOAD_DIR}/\"
                        """
                    }
                }
            }
       }
    }
    post{
        success{
            withCredentials([usernamePassword(credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                script {
                    if(env.GITHUB_PR_URL){
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GitHubStatus,
                            state : "success"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GtaStatusCheck,
                            message: "GTA tests",
                            description: "GTA tests",
                            state : "pending"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : SdlStatusCheck,
                            message: "SDL scans",
                            description: "SDL scans",
                            state : "success"
                        )
                    }
                   }
            }
        }
        failure{
               withCredentials([usernamePassword(credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                script {
                    if(env.GITHUB_PR_URL){
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GitHubStatus,
                            state : "fail"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GtaStatusCheck,
                            message: "GTA tests",
                            description: "GTA tests",
                            state : "fail"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : SdlStatusCheck,
                            message: "SDL scans",
                            description: "SDL scans",
                            state : "fail"
                        )
                     }
                 }
                }
        }
        aborted{
            withCredentials([usernamePassword(credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                script {
                    if(env.GITHUB_PR_URL){
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GitHubStatus,
                            state : "aborted"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : GtaStatusCheck,
                            message: "GTA tests",
                            description: "GTA tests",
                            state : "aborted"
                        )
                        status_check(
                            gh_token : PASSWORD,
                            commit : GITHUB_PR_HEAD_SHA,
                            pull_request : GITHUB_REPO_STATUS_URL,
                            name : SdlStatusCheck,
                            message: "SDL scans",
                            description: "SDL scans",
                            state : "aborted"
                        )
                    }
                }
            }
        }
    }
}
