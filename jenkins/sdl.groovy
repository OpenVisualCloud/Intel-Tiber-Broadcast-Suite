@Library(['jenkins-common-pipelines', 'pact-shared-library']) _

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

    def abi_image="amr-registry.caas.intel.com/owr/abi_lnx:3.0.0"
    def _cmd =[
        hadolint: "jenkins/scripts/hadolint.sh" ,
        trivy: "jenkins/scripts/trivy.sh ${tar_image}",
        schellcheck: "jenkins/scripts/shellcheck.sh",
        mcAffee: "docker run --rm -v \$(pwd):/opt/ ${abi_image} /bin/bash -c \"cd /opt/; jenkins/scripts/mcafee_scan.sh\""
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
        GITHUB_PR_HEAD_SHA = "${params.github_pr_head_sha}"
        GITHUB_REPO_STATUS_URL = "${params.github_repo_status_url}"
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
