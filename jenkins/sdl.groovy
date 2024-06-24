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
        mcAffee: "${env.DOCKER_ABI} \"cd /opt/; jenkins/scripts/mcafee_scan.sh\"",
        docker_benchmark: "jenkins/scripts/docker_cis_benchmark.sh"
    ]
    def _artifacts_path =[
        hadolint: "Hadolint/hadolint-Dockerfile*" ,
        trivy: "Trivy/*",
        schellcheck: "shellcheck_logs/*",
        mcAffee: "Malware/*",
        docker_benchmark: "docker-bench-security/cisdockerbench_results/*" 
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
        DOCKER_ABI="docker run --rm -v \$(pwd):/opt/ ${env.ABI_IMAGE} /bin/bash -c "
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
                            dir("Hadolint-scan"){
                                sh """ cp -r ${params.relative_dir}/* . """
                                scan("hadolint", "")
                            }
                        }
                    } 
                }
                stage("Trivy"){
                    steps{
                        script{
                            dir("Trivy-scan"){
                                sh """ cp -r ${params.relative_dir}/* . """
                                scan("trivy", params.tar_docker_image)
                            }
                        }
                    } 
                }
                stage("Schellcheck"){
                    steps{
                        script{
                            dir("Shellcheck-scan"){
                                sh """ cp -r ${params.relative_dir}/* . """
                                scan("schellcheck", "")
                            }
                        }
                    } 
                }
                stage("McAfee"){
                    steps{
                        script{
                            dir("Mcafee_scan"){
                                sh """ cp -r ${params.relative_dir}/* . """
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
                                dir("Protex-scan"){                                    
                                    sh"""
                                      cp -r ${params.relative_dir}/* .
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
                                        sudo chown -R $USER:$USER ./*
                                    """
                                    archiveArtifacts allowEmptyArchive: true, artifacts: "OWRBuild/${env.PROTEX_FOLDER}/*"
                                // Protex run in docker as root, it creates files as root on host machine.
                                // jenkins cleanup cannot remove those files( permision denied 13), 
                                // so they needs to be removed manually
                                sh """ sudo rm -rf OWRBuild/ """
                                }
                            }
                        }
                    }
                }
                stage("coverity"){
                    steps{
                        script{
                        withCredentials([usernamePassword(
                            credentialsId: 'bbff5d12-094b-4009-9dce-b464d51f96d4',
                            usernameVariable: 'USERNAME',
                            passwordVariable: 'PASSWORD')]){
                                dir("Coverity-scan"){
                                    sh """
                                        cp -r ${params.relative_dir}/* .
                                        ${env.DOCKER_ABI} \"cd /opt/; abi coverity analyze \
                                                --debug \
                                                --aggressiveness-level high \
                                                --server "${COVERITY_SERVER}" \
                                                --username "${COVERITY_USR}" \
                                                --password "${COVERITY_PSW}" \
                                                --report_output_dir "${WORKSPACE}" \
                                                --build_command "jenkins/scripts/coverity_build_script.sh" \
                                                --stream "${COVERITY_STREAM}"

                                    """
                                    archiveArtifacts allowEmptyArchive: true, artifacts: "OWRBuild/${COVERITY_FOLDER}/*"
                                }
                            }
                        }
                    }
                }
                stage("Docker CIS"){
                    steps{
                        script{
                            dir("Docker-CIS"){
                                sh """ cp -r ${params.relative_dir}/* . """
                                scan("docker_benchmark", "")
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
