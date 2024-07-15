@Library(['jenkins-common-pipelines', 'pact-shared-library']) _

pipeline {
    agent { 
        label 'video-production-build'
    }
    environment {
        COVERITY_SERVER = "https://coverityent.devtools.intel.com/prod1/"
        COVERITY_PROJECT = "Software Defined Broadcast main"
        COVERITY_FOLDER = "CoverityEvidence"
        COVERITY_STREAM = "sdb-ffmpeg-patch"
    }
    stages{
        stage("coverity build"){
            steps{
                script{
                    def repo = 'https://github.com/intel-innersource/applications.services.cloud.visualcloud.vcdp.video-production-pipeline.git'
                    def branches = [[name: 'lgrab/fix_coverity_docker']]
                    def userRemoteConfigs = [[
                        credentialsId: '0febae38-30c4-4243-88f1-b85eb771452d',
                        url: repo
                    ]]
                    checkout([
                        $class: 'GitSCM',
                        branches: branches,
                        userRemoteConfigs: userRemoteConfigs
                    ])
                    def latest_sdb_image = sh(script: "docker images --format '{{.Repository}}:{{.Tag}} {{.CreatedAt}}' |  grep 'build_stage' | sort -k2 -r | head -n1 | cut -d \" \" -f1", returnStdout: true).trim()
                    def docker_build_args = ["--build-arg http_proxy=http://proxy-dmz.intel.com:912", 
                                         "--build-arg https_proxy=http://proxy-dmz.intel.com:912",
                                         "--build-arg no_proxy=localhost,intel.com,192.168.0.0/16,172.16.0.0/12,127.0.0.0/8,10.0.0.0/8",
                                         "--build-arg IMAGE_TAG_NAME=${latest_sdb_image}",
                                         "-t coverity_image_${BUILD_ID} ",
                                         "-f Dockerfile ."].join(" ")

                sh """
                    git config --global --add safe.directory \$(pwd)
                    cd jenkins/docker/coverity
                    docker build ${docker_build_args}
                """
            }
          }
        }
        stage("Coverity Scan"){
            steps{
                script{
                    withCredentials([usernamePassword(
                        credentialsId: 'bbff5d12-094b-4009-9dce-b464d51f96d4',
                        usernameVariable: 'USERNAME',
                        passwordVariable: 'PASSWORD')]){
                        def docker_runtime_args = ["-v \"\$(pwd)\":/tmp/host",
                                         "-e http_proxy=http://proxy-dmz.intel.com:912", 
                                         "-e https_proxy=http://proxy-dmz.intel.com:912",
                                         "-e no_proxy=localhost,intel.com,192.168.0.0/16,172.16.0.0/12,127.0.0.0/8,10.0.0.0/8",
                                         "-e COVERITY_SERVER=${env.COVERITY_SERVER}",
                                         "-e COVERITY_USR=${USERNAME}",
                                         "-e COVERITY_PSW=${PASSWORD}",
                                         "-e WORKSPACE=${env.COVERITY_FOLDER}",
                                         "-e coverity_image_${BUILD_ID}",
                                         "-e COVERITY_SERVER=${env.COVERITY_SERVER}",
                                         "-t coverity_image_${BUILD_ID} ",].join(" ")
                        sh """
                            docker run ${docker_runtime_args} 
                            sudo chown -R \${USER}:\${USER} Logs
                            sudo chown -R \${USER}:\${USER} OWRBuild
                        """
                    }
                }
            }
        }
        stage("Coverity reports"){
            steps{
                script{
                    withCredentials([usernamePassword(
                        credentialsId: 'bbff5d12-094b-4009-9dce-b464d51f96d4',
                        usernameVariable: 'USERNAME',
                        passwordVariable: 'PASSWORD')]){
                        sh """
                            export no_proxy="\$no_proxy,.intel.com"
                            export USERNAME=${USERNAME}
                            export PASSWORD=${PASSWORD}
                            ./jenkins/scripts/generate_coverity_reports.sh
                            
                        """
                    }
                }
            }
        }
    }
    post{
        failure{
            script{
                //abi create files as root,  
                //it changes those files ownersip to $USER only if it executes succesfully
                // otherwise they are accessible only by root
                sh """ sudo chown \${USER}:\${USER} -R ${WORKSPACE} """
            }
        }
        always{
            script{
                    archiveArtifacts allowEmptyArchive: true, artifacts: "OWRBuild/*"
                    archiveArtifacts allowEmptyArchive: true, artifacts: "cov_report*"
                    archiveArtifacts allowEmptyArchive: true, artifacts: "*.pdf"

            }
        }
    }   
}
