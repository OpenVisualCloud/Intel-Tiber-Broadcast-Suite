@Library(['pact-shared-library', 'abi', 'vsval-gta-asset-library']) _

def setStatusCheck(String status_name, String state){
  def statusMap = [
    Build: "NEX: SED Visual Solutions/Ci Project Build",
    Gta: "daily testing",
  ]
  withCredentials([
    usernamePassword(
      credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
      usernameVariable: 'USERNAME',
      passwordVariable: 'PASSWORD'
    )]){
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
        kubernetes {
            cloud 'harvester-k8s'
            inheritFrom 'k8s-sdl-abi'
        }
    }
    environment {
      IMAGE_REGISTRY = "ger-is-registry.caas.intel.com/nex-vs-cicd-automation"
      IMAGE_CACHE_REGISTRY = "ger-is-registry.caas.intel.com/cache"
      IMAGE_TAG = '${BUILD_ID}'
      imageRegistry = "${IMAGE_REGISTRY}/sdb/tiber-broadcast-suite"
     
      GIT_SHA = '${GITHUB_PR_HEAD_SHA}'
      
      ARTIFACTORY_URL = 'https://af01p-igk.devtools.intel.com/artifactory'
      JOB_FILE_ARTIFACTORY_PREFIX = 'video_production_igk-igk-local/builds'
      
      PROTEX_SERVER  = "gerprotex012.devtools.intel.com"
      PROTEX_FOLDER  = "ProtexEvidence"
      PROTEX_PROJECT = "c_broadcastsuiteformediaentertainment_33366"      
      COVERITY_SERVER  = "https://coverityent.devtools.intel.com/prod1/"
      COVERITY_PROJECT = "Software Defined Broadcast main"
      COVERITY_FOLDER  = "CoverityEvidence"
      COVERITY_STREAM  = "sdb-ffmpeg-patch"
    }
    triggers {
      GenericTrigger(
          genericVariables: [
              /* params for GITHUB PR */
              [key: 'GITHUB_PR_TRIGGER_ACTION', value: '$.action'],
              [key: 'GITHUB_PR_TRIGGER_SENDER_AUTHOR', value: '$.sender.login'],
              [key: 'GITHUB_PR_TRIGGER_SENDER_ID', value: '$.sender.id'],
              [key: 'GITHUB_PR_COMMIT_STATUS', value: '$.pull_request._links.statuses.href'],
              [key: 'GITHUB_COMMENT_BODY', value: '$.comment.body'],
              [key: 'GITHUB_PR_TARGET_BRANCH', value: '$.pull_request.base.ref'],
              [key: 'GITHUB_PR_SOURCE_BRANCH', value: '$.pull_request.head.ref'],
              [key: 'GITHUB_PR_BODY', value: '$.pull_request'],
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
              [key: 'GITHUB_PR_API_URL', value: '$.pull_request.url'],
              /* params for GITHUB PR comment retrigger */
              [key: 'GITHUB_COMMENT_PR_API_URL', value: '$.issue.pull_request.url'],
              [key: 'GITHUB_COMMENT_PR_URL', value: '$.issue.html_url'],
              [key: 'GITHUB_COMMENT_PR_NUMBER', value: '$.issue.number']

          ],
          genericHeaderVariables: [
              [key: 'X-GitHub-Event', regexpFilter: '']
          ],
          token: 'app-srv-cloud-vsval-vcdp-13925d3',
          causeString: 'Triggered on GitHub webhook',
          printContributedVariables: false,
          printPostContent: true,
          silentResponse: false,
          regexpFilterText: '$GITHUB_PR_TRIGGER_ACTION $GITHUB_COMMENT_BODY',
          regexpFilterExpression: '^(opened|reopened|synchronize).*|.*REVERIFY.*$'
      )
    }
    options {
        timestamps()
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('Fetch Clean Repository'){
            steps {
                container('jnlp') {
                    script {
                        cleanWs()
                        def mainRepoBranch = 'main'
                        def mainRepoUrl = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline.git'
                        def mainRepoRemoteConfigs = [[url: mainRepoUrl, credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d']]
                        
                        def testRepoUrl = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.validation.itbs.git'
                        def testRepoRemoteConfigs = [[url: testRepoUrl, credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d']]
                        
                        /* PR open/reopen/synchronize delivers different payload than comment addition
                        This is why we need to distinguish which variables to use */
                        if (env.GITHUB_PR_URL){
                            PR_URL = env.GITHUB_PR_URL
                            println "assignin env.GITHUB_PR_URL to PR_URL"
                            println "PR_URL: '${PR_URL}'"
                        } else if (env.GITHUB_COMMENT_PR_URL){
                            PR_URL = env.GITHUB_COMMENT_PR_URL
                            println "assignin env.GITHUB_COMMENT_PR_URL to PR_URL"
                            println "PR_URL: '${PR_URL}'"
                        }
                        if (env.GITHUB_PR_API_URL){
                            PR_API_URL = env.GITHUB_PR_API_URL
                        } else if (env.GITHUB_COMMENT_PR_API_URL){
                            PR_API_URL = env.GITHUB_COMMENT_PR_API_URL
                        }
                        if (env.GITHUB_PR_NUMBER){
                            PR_NUMBER = env.GITHUB_PR_NUMBER
                        } else if (env.GITHUB_COMMENT_PR_NUMBER){
                            PR_NUMBER = env.GITHUB_COMMENT_PR_NUMBER
                        }
                        if(env.GITHUB_PR_URL){
                            currentBuild.description = 'Trigged by <a href="' + GITHUB_PR_URL + '" target="_blank">' + GITHUB_REPO_FULL_NAME + '</a>.'
                            println 'Trigged by ' + GITHUB_PR_URL + '.'
                            branches = [[name: "origin/pull/${GITHUB_PR_NUMBER}/merge"]]
                            mainRepoRemoteConfigs[0]['refspec'] = "+refs/pull/${GITHUB_PR_NUMBER}/merge:refs/remotes/origin/pull/${GITHUB_PR_NUMBER}/merge"
                        }     
                        def mainRepoScmVars = checkout([
                            $class: 'GitSCM',
                            branches: [[name:mainRepoBranch]],
                            extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'sdb']],
                            userRemoteConfigs: mainRepoRemoteConfigs
                        ])
                        GIT_SHA   = mainRepoScmVars.GIT_COMMIT
                        IMAGE_TAG = mainRepoScmVars.GIT_COMMIT

                        def testRepoScmVars = checkout([
                            $class: 'GitSCM',
                            branches: [[name:'main']],
                            extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'test-repo']],
                            userRemoteConfigs: testRepoRemoteConfigs
                        ])
                        stash includes: '**', name: 'code'
                        stash includes: '**', name: 'repo', useDefaultExcludes: false
                    }
                }
            }
        }
        stage('Build Image'){
            steps {
                container('buildx') {
                    script {
                        unstash 'code'
                        sh 'git config --global --add safe.directory $WORKSPACE'
                        withCredentials([usernamePassword(credentialsId: 'ger-is-registry-docker-robot-account-token-based', usernameVariable: 'username', passwordVariable: 'password')]) {
                            sh """echo '${password}' | docker login --username='${username}' --password-stdin '${imageRegistry}'"""
                        }
                        dir('sdb'){
                            sh """
                                docker buildx build \
                                    --build-arg IMAGE_CACHE_REGISTRY=${IMAGE_CACHE_REGISTRY} --build-arg http_proxy --build-arg https_proxy \
                                    --cache-to=type=registry,mode=max,image-manifest=true,oci-mediatypes=true,ref=${imageRegistry}:cache \
                                    --cache-from=type=registry,ref=${imageRegistry}:cache \
                                    --tag ${imageRegistry}:${IMAGE_TAG} \
                                    --output type=image,push=true \
                                    --target final-stage \
                                    --progress=plain \
                                    -f Dockerfile .
                            """
                        }
                    }
                }
            }
        }
        stage('Scans') {
          parallel {
            stage("Hadolint"){
              steps{
                container('hadolint'){
                  script{
                    dir('hadolint'){
                      unstash 'code'
                      sh """git config --global --add safe.directory '${WORKSPACE}/hadolint'"""
                      dir('sdb'){
                        sh 'jenkins/scripts/hadolint.sh'
                        archiveArtifacts allowEmptyArchive: true, artifacts: "Hadolint/hadolint-Dockerfile*"
                      }
                    }
                  }
                } 
              }
            }
            stage("Shellcheck"){
              steps{
                container('hadolint'){
                  script{
                    dir('shellcheck'){
                      unstash 'code'
                      sh """git config --global --add safe.directory '${WORKSPACE}/shellcheck'"""
                      dir('sdb'){
                        sh 'jenkins/scripts/shellcheck.sh'
                        archiveArtifacts allowEmptyArchive: true, artifacts: "shellcheck_logs/**"
                      }
                    }
                  }
                }
              } 
            }
            stage("Antivirus"){
              steps{
                container('abi'){
                  script{
                    dir('antivirus'){
                      unstash 'code'
                      sh """git config --global --add safe.directory '${WORKSPACE}/antivirus'"""
                      dir('sdb'){
                        sh 'jenkins/scripts/mcafee_scan.sh'
                        archiveArtifacts allowEmptyArchive: true, artifacts: "Malware/*"
                      }
                    }
                  }
                }
              }
            }
            stage("Protex"){
              steps{
                container('abi'){
                  script{
                    dir('protex'){
                      unstash 'code'
                      sh """git config --global --add safe.directory '${WORKSPACE}/protex'"""
                      withCredentials([usernamePassword(
                        credentialsId: 'SDB-faceless-account',
                        usernameVariable: 'USERNAME',
                        passwordVariable: 'PASSWORD'
                      )]){
                        dir('sdb'){
                          sh"""
                            abi ip_scan scan \
                                --scan_server '${env.PROTEX_SERVER}'   \
                                --scan_project '${env.PROTEX_PROJECT}' \
                                --username '${USERNAME}' \
                                --password '${PASSWORD}' \
                                --ing_path '.'           \
                                --report_type xlsx       \
                                --report_config cos      \
                                --report_config obl      \
                                --scan_output '${env.PROTEX_FOLDER}'
                          """
                          archiveArtifacts allowEmptyArchive: true, artifacts: "OWRBuild/${env.PROTEX_FOLDER}/*"
                        }
                      }
                    }
                  }
                }
              }
            }
            stage("Tarball Scripts"){
              stages{
                stage('Fetch Tarball'){
                  steps {
                    container('buildx') {
                      script {
                        dir('tarball'){
                          unstash 'code'
                          sh """git config --global --add safe.directory '${WORKSPACE}/tarball'"""
                          withCredentials([usernamePassword(credentialsId: 'ger-is-registry-docker-robot-account-token-based', usernameVariable: 'username', passwordVariable: 'password')]) {
                            sh """echo '${password}' | docker login --username='${username}' --password-stdin '${imageRegistry}'"""
                          }
                          dir('sdb'){
                            sh """
                              docker buildx build \
                                --build-arg IMAGE_CACHE_REGISTRY=${IMAGE_CACHE_REGISTRY} --build-arg http_proxy --build-arg https_proxy \
                                --cache-from=type=registry,ref=${imageRegistry}:cache \
                                --tag ${imageRegistry}:build-${IMAGE_TAG} \
                                --output "type=docker,dest=${WORKSPACE}/tarball/tiber-broadcast-suite.tar" \
                                --target build-stage \
                                --progress=plain \
                                -f Dockerfile .
                            """
                          }
                        }
                      }
                    }
                  }
                }
                stage("Trivy"){
                  steps{
                    container('trivy'){
                      script{
                        dir('trivy'){
                          unstash 'code'
                          sh """git config --global --add safe.directory '${WORKSPACE}/trivy'"""
                          dir('sdb'){
                            sh 'chmod a+x jenkins/scripts/trivy_image_scan.sh'
                            sh '${WORKSPACE}/trivy/sdb/jenkins/scripts/trivy_image_scan.sh "${WORKSPACE}/tarball/tiber-broadcast-suite.tar"'
                          }
                          archiveArtifacts allowEmptyArchive: true, artifacts: ("Trivy/**")
                        }
                      }
                    } 
                  }
                }
                stage("Docker CIS"){
                  steps{
                    container('abi'){
                      script{
                        dir('tarball/sdb'){
                        //   sh 'jenkins/scripts/docker_cis_benchmark.sh ${WORKSPACE}/tarball/tiber-broadcast-suite.tar'
                          sh 'echo ${WORKSPACE}/tarball/tiber-broadcast-suite.tar'
                          archiveArtifacts allowEmptyArchive: true, artifacts: "cisdockerbench_results/*"
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
        stage('Upload'){
          steps {
            script {
              dir('gtax'){
                unstash 'code'
                sh """git config --global --add safe.directory '${WORKSPACE}/gtax'"""
                dir('sdb'){
                  writeFile file: 'properties.json', text: jsonStringFromMap([
                    "video_production_sdb": [GIT_SHA],
                    "build.component_revision": [GIT_SHA],
                    "build.component_project": ['Tiber-Broadcast-Suite']
                  ])
                //   gtaAssetUpload([
                  gtaAssetUpload=[
                    propertiesFile:'properties.json',
                    creds:[user:"sys_vsval",password:credentials('ed6200e9-7a86-4cf8-b2f3-2cea840a23f5')],
                    gtaAsset:[
                      repo: 'video_production_igk-igk-local/builds',
                      rootUrl: 'https://af01p-igk.devtools.intel.com/artifactory',
                      assetName: 'ci',
                      assetVersion: "${env.JOB_BASE_NAME}_${env.BUILD_ID}",
                      localPath: '/tmp/gtax-upload',
                  ]]
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

