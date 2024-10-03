@Library(['jenkins-common-pipelines', 'pact-shared-library', 'vsval-gta-asset-library']) _

def GtaStatusCheck = "daily testing"
def relativeDir = 'vp-build'

def setStatusCheck(String status_name, String state){
    def statusMap = [
        Build: "NEX: SED Visual Solutions/Ci Project Build",
        Gta: "daily testing",
    ]
    withCredentials([usernamePassword(
        credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
        usernameVariable: 'USERNAME',
        passwordVariable: 'PASSWORD')]){
            status_check(
                gh_token : PASSWORD,
                commit : GITHUB_PR_HEAD_SHA,
                pull_request : GITHUB_REPO_STATUS_URL,
                name : statusMap[status_name],
                state : state
            )
    }
}

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
            token: '344050040c8011ef82a100155d3925d3',
            causeString: 'Triggered on GitHub webhook',
            printContributedVariables: false,
            printPostContent: true,
            silentResponse: false,
            regexpFilterText: '${GITHUB_PR_TRIGGER_ACTION}',
            regexpFilterExpression: '^(opened|reopened|synchronize).*$'
        )
    }
    parameters {
        string(name: 'BRANCH',    defaultValue: 'main', description: 'select branch to build')

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
                    def repo = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline.git'
                    def branches = [[name: params.BRANCH ]]
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
                    def dockerBuildArgs = ["--build-arg http_proxy",
                                           "--build-arg https_proxy",
                                           "-t ${env.IMAGE_TAG_NAME}_build_stage",
                                           "--target build-stage",
                                           "-f Dockerfile ."]
                    dir(relativeDir){
                        sh """
                          git config --global --add safe.directory \$(pwd)
                          docker build ${dockerBuildArgs.join(" ")}
                        """

                        dockerBuildArgs.remove("--target buildstage")
                        dockerBuildArgs = dockerBuildArgs.collect { 
                            it.contains('_build_stage') ? it.replaceAll( /_build_stage/, '' ) : it }
                        sh """ docker build ${dockerBuildArgs.join(" ")}  """
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
                    dir('/tmp/test-repo'){
                        def repo = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.validation.itbs'
                        def branches = [[name: 'main']]
                        def userRemoteConfigs = [[
                                credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
                                url: repo ]]     
                        checkout([
                        $class: 'GitSCM',
                        branches: branches,
                        userRemoteConfigs: userRemoteConfigs
                    ])
                    sh """ tar czf \${UPLOAD_DIR}/itbs_tests_repo_${JOB_BUILD_ID}.tar.gz . """
                  }    
                }
            }
        }
        stage("trigger sdl"){
            steps{
                script{
                    def pr_head_sha = (env.GITHUB_PR_HEAD_SHA ? env.GITHUB_PR_HEAD_SHA : "manual_trigger")
                    def github_repo_status_url = (env.GITHUB_REPO_STATUS_URL ? env.GITHUB_REPO_STATUS_URL : "manual_trigger")
                    build job: 'sdl-scans' , parameters: [
                        string(name: 'tar_docker_image', value: "${UPLOAD_DIR}/${IMAGE_TAG_NAME}.tar.gz"),
                        string(name: 'github_repo_status_url', value:"${github_repo_status_url}" ),
                        string(name: 'github_pr_head_sha', value: "${pr_head_sha}"),
                        string(name: 'relative_dir', value: "${WORKSPACE}/${relativeDir}"),
                        string(name: 'parent_job_id', value: "${JOB_BUILD_ID}"),
                    ], wait: false
                }
            }
        }
        stage('Upload'){
            steps {
                script {
                    dir(relativeDir){
                        // https://repo/name -> repo/name
                        def repoName = (env.GITHUB_REPO_FULL_NAME ? env.GITHUB_REPO_FULL_NAME.split('/')[1] : "none") 
                        def pr_head_sha = (env.GITHUB_PR_HEAD_SHA ? env.GITHUB_PR_HEAD_SHA : "manual_trigger")
                        dir(relativeDir){
                            def propertiesFile = 'properties.json'
                            writeFile file: propertiesFile, text: jsonStringFromMap([
                                    "video_production_sdb": [pr_head_sha],
                                    "build.component_revision": [pr_head_sha],
                                    "build.component_project": [repoName]])
                            def assetConfig = [ 
                                creds: [
                                    user: "sys_vsval",
                                    password: env.ARTIFACTORY_SECRET,
                                ],
                                gtaAsset: [
                                    rootUrl: env.ARTIFACTORY_URL,
                                    repo: env.JOB_FILE_ARTIFACTORY_PREFIX,
                                    assetName: env.GTA_ASSET_NAME,
                                    assetVersion: "${env.JOB_BASE_NAME}_${env.JOB_BUILD_ID}",
                                    localPath: env.UPLOAD_DIR,
                                ],
                                propertiesFile: propertiesFile
                            ]
                            gtaAssetUpload(assetConfig)
                    }
                }
            }
         }
       }  
    }
    post{
        success{
            script {
                if(env.GITHUB_PR_URL){
                    setStatusCheck("Build", "success")
                }
            }
        }
        failure{
            script {
                if(env.GITHUB_PR_URL){
                    setStatusCheck("Build", "fail")
                    setStatusCheck("Gta", "fail")
                }
            }
        }
        aborted{
            script {
                if(env.GITHUB_PR_URL){
                    setStatusCheck("Build", "aborted")
                    setStatusCheck("Gta", "aborted")
                }
            }
        }
    }
}
