@Library(['jenkins-common-pipelines', 'pact-shared-library', 'abi']) _

def setSdlStatusCheck(String state){
    withCredentials([usernamePassword(
        credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
        usernameVariable: 'USERNAME',
        passwordVariable: 'PASSWORD')]){
            status_check(
                gh_token : PASSWORD,
                commit : GITHUB_PR_HEAD_SHA,
                pull_request : GITHUB_REPO_STATUS_URL,
                name : "SDL scans",
                state : state
            )
    }
}

def scan(String type, String tar_image){

    def _cmd =[
        hadolint: "jenkins/scripts/hadolint.sh" ,
        trivy: "jenkins/scripts/trivy.sh ${tar_image}",
        schellcheck: "jenkins/scripts/shellcheck.sh",
        mcAffee: "${env.DOCKER_ABI} \"cd /opt/; jenkins/scripts/mcafee_scan.sh\""
    ]
    def _artifacts_path =[
        hadolint: "Hadolint/hadolint-Dockerfile*" ,
        trivy: "Trivy/*",
        schellcheck: "shellcheck_logs/*",
        mcAffee: "Malware/*"
    ]
    sh """
        ${_cmd[type]}
    """
    archiveArtifacts allowEmptyArchive: true, artifacts: "${_artifacts_path[type]}"
}



pipeline {
    agent {
        label 'video-production-build'
    }
    parameters {
        string(name: 'tar_docker_image', defaultValue: '', description: 'full patch to built and packed image')
        string(name: 'github_repo_status_url', defaultValue: '', description: 'api endpoint to return sdl status')
        string(name: 'github_pr_head_sha', defaultValue: '', description: 'github pr commit id')
        string(name: 'relative_dir', defaultValue: '', description: 'jenkins workspace with codespace ')
    }
    environment {
        ABI_IMAGE="amr-registry.caas.intel.com/owr/abi_lnx:3.0.0"
        GITHUB_PR_HEAD_SHA = "${params.github_pr_head_sha}"
        GITHUB_REPO_STATUS_URL = "${params.github_repo_status_url}"
        PROTEX_SERVER = "gerprotex012.devtools.intel.com"
        EVIDENCEFOLDER = "InternalEvidence"
        PROTEX_PROJECT = "c_broadcastsuiteformediaentertainment_33366"
        CRED_DEFAULT = "build_sie"
        CRED_DEFAULT_EMAIL = "build_sie-email"
        DOCKER_ABI="docker run --rm -v \$(pwd):/opt/ ${env.ABI_IMAGE} /bin/bash -c "
        EVIDENCE="protex.log"
    }
    stages {
        stage("set status"){
            steps{
                script{
                    if(params.github_repo_status_url){
                        setSdlStatusCheck("pending")
                    }
                }
            }
        }
        stage("scans"){
            parallel {
                stage("Hadolint"){
                    steps{
                        script{
                            dir(params.relative_dir){
                                scan("hadolint", "")
                            }
                        }
                    } 
                }
                stage("Trivy"){
                    steps{
                        script{
                            dir(params.relative_dir){
                                scan("trivy", params.tar_docker_image)
                            }
                        }
                    } 
                }
                stage("Schellcheck"){
                    steps{
                        script{
                            dir(params.relative_dir){
                                scan("schellcheck", "")
                            }
                        }
                    } 
                }
                stage("McAfee"){
                    steps{
                        script{
                            dir(params.relative_dir){
                                scan("mcAffee", "")
                            }
                        }
                    } 
                }
                stage("Protex"){
                    steps{
                        script{
                            withCredentials([usernamePassword(
                                credentialsId: 'de62b8cb-2c0a-4a67-83ea-cc51c7486c05',
                                usernameVariable: 'USERNAME',
                                passwordVariable: 'PASSWORD')]){
                                dir(params.relative_dir){
                                    sh"""
                                    ${env.DOCKER_ABI} \"cd /opt/; abi ip_scan scan \
                                            --scan_server ${env.PROTEX_SERVER} \
                                            --scan_project ${env.PROTEX_PROJECT} \
                                            --username ${USERNAME} \
                                            --password ${PASSWORD} \
                                            --root_dir ${WORKSPACE} \
                                            --ing_path \".\" \
                                            --report_type xlsx \
                                            --report_config cos \
                                            --report_config obl \
                                            --scan_output ${env.WORKSPACE}/${env.EVIDENCE}\"
                                    """
                                    archiveArtifacts allowEmptyArchive: true, artifacts: "${env.WORKSPACE}/${env.EVIDENCE}"
                                }
                            }
                        }
                    }
                }
        
            }
        }
    }
    post{
        success{
            script {
                    if(params.github_repo_status_url){
                        setSdlStatusCheck("success")
                    }
                }
        }
        failure{
            script {
                    if(params.github_repo_status_url){
                        setSdlStatusCheck("fail")
                    }
                }
        }
        aborted{
            script {
                    if(params.github_repo_status_url){
                        setSdlStatusCheck("aborted")
                    }
                }
        }
    }
}
