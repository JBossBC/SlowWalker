pipeline{
    agent any
    environment {
        // 设置Go环境变量
        GO_HOME = "/usr/local/go"
        PATH = "${GO_HOME}/bin:${PATH}"
    }
    stages{
        stage("pull"){
            steps{
                git branch: 'test', credentialsId:"jenkins",url:'git@github.com:JBossBC/repliteWeb.git'
            }
        }
         stage('Set golang proxy') {
            steps {
                // 设置Go代理
                script {
                    env.GO111MODULE = 'on'
                    env.GOPROXY = 'https://goproxy.cn,direct'
                }
            }
        }
        // stage("test"){
        //     steps{
        //         sh 
        //     }
        // }
        stage("build"){
            steps{
                sh 'go build -o backend ./cmd/main.go'
                withCredentials([sshUserPrivateKey(credentialsId: 'jenkins', keyFileVariable: 'SSH_KEY')]) {
                    // 通过SCP命令将文件复制到节点
                    //TODO should use node to connect with ip
                    sh """
                        ssh-agent bash -c 'ssh-add $SSH_KEY; scp backend root@112.124.53.234:/opt/'
                    """
                    sh """
                        ssh-agent bash -c 'ssh-add $SSH_KEY; scp -r configs/ root@112.124.53.234:/opt/'
                    """    
                } 
            }
        }
        stage("report"){
            agent{
                node{
                    label 'backend'
                }
            }
            steps{
                 dir('/opt') {
                 sh './backend'
                }
        }
        }
}
}