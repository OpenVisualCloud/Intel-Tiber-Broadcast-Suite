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

def fetchArtifacts(String artifactsPath){
    archiveArtifacts allowEmptyArchive: true, artifacts: "${artifactsPath}*"
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
        string(name: 'parent_job_id', defaultValue: '')
    }
    environment {
        ABI_IMAGE="amr-registry.caas.intel.com/owr/abi_lnx:3.0.0"
        GITHUB_PR_HEAD_SHA = "${params.github_pr_head_sha}"
        GITHUB_REPO_STATUS_URL = "${params.github_repo_status_url}"
        PROTEX_SERVER = "gerprotex012.devtools.intel.com"
        PROTEX_FOLDER = "ProtexEvidence"
        PROTEX_PROJECT = "c_broadcastsuiteformediaentertainment_33366"
        COVERITY_SERVER = "https://coverityent.devtools.intel.com/prod1/"
        COVERITY_PROJECT = "Software Defined Broadcast main"
        COVERITY_FOLDER = "CoverityEvidence"
        COVERITY_STREAM = "sdb-ffmpeg-patch"
        DOCKER_ABI="docker run --rm -v \$(pwd):/opt/ -v /tmp:/tmp ${env.ABI_IMAGE}  /bin/bash -c "
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
        stage("Hadolint"){
            steps{
                script{
                    dir(params.relative_dir){
                        sh """ jenkins/scripts/hadolint.sh """
                        fetchArtifacts("Hadolint/hadolint-Dockerfile")
                    }
                }
            } 
        }
        stage("Trivy"){
            steps{
                script{
                    dir(params.relative_dir){
                        sh """ 
                            rm -rf Trivy
                            mkdir Trivy
                            jenkins/scripts/trivy_image_scan.sh ${params.tar_docker_image}
                            jenkins/scripts/trivy_dockerfile_scan.sh
                        """
                        fetchArtifacts("Trivy/")
                        fetchArtifacts("Trivy/dockerfile/")
                        fetchArtifacts("Trivy/image/")
                    }
                }
            } 
        }
        stage("Shellcheck"){
            steps{
                script{
                    dir(params.relative_dir){
                        sh """ jenkins/scripts/shellcheck.sh """
                        fetchArtifacts("shellcheck_logs/")
                    }
                }
            } 
        }
        stage("McAfee"){
            steps{
                script{
                    dir(params.relative_dir){
                        sh """ ${env.DOCKER_ABI} \"cd /opt/; jenkins/scripts/mcafee_scan.sh\" """
                        fetchArtifacts("Malware/")
                    }
                }
            } 
        }
        stage("Protex"){
            steps{
                script{
                    withCredentials([usernamePassword(
                        credentialsId: 'bbff5d12-094b-4009-9dce-b464d51f96d4',
                        usernameVariable: 'USERNAME',
                        passwordVariable: 'PASSWORD')]){
                        dir(params.relative_dir){                                    
                            sh"""
                                ${env.DOCKER_ABI} \"cd /opt/; abi ip_scan scan \
                                    --scan_server ${env.PROTEX_SERVER} \
                                    --scan_project ${env.PROTEX_PROJECT} \
                                    --username ${USERNAME} \
                                    --password ${PASSWORD} \
                                    --ing_path \".\" \
                                    --report_type xlsx \
                                    --report_config cos \
                                    --report_config obl \
                                    --scan_output ${env.PROTEX_FOLDER}\"
                                 sudo chown -R \${USER}:\${USER} Logs
                                 sudo chown -R \${USER}:\${USER} OWRBuild
                                    
                            """
                            fetchArtifacts("OWRBuild/${env.PROTEX_FOLDER}/")
                        }
                    }
                }
            }
        }
        stage("Docker CIS"){
            steps{
                script{
                    dir(params.relative_dir){
                        sh """ jenkins/scripts/docker_cis_benchmark.sh ${params.tar_docker_image} """
                        fetchArtifacts("cisdockerbench_results/")
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
