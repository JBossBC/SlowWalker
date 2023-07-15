pipeline{
    agent any
    stages{
        stage("pull"){
            steps{
                git branch 'test',credentialsId:"jenkins",url:"git@github.com:JBossBC/repliteWeb.git"
                withCredentials([sshUserPrivateKey(credentialsId: 'jenkins', keyFileVariable: 'SSH_KEY')]){
                    sh """
                        ssh-agent bash -c 'ssh-add $SSH_KEY; scp deployments/frontend.dockerfile root@159.75.177.48:/opt/configs/'
                    """
                     sh """
                        ssh-agent bash -c 'ssh-add $SSH_KEY; scp -r web/ root@159.75.177.48:/opt/web/'
                    """
                }
            }
        }
        stage("build"){
                agent {
        node{
          label "frontend"
        }
    }
            steps{
                dir("/opt"){
                    sh 'docker build -f deployments/frontend.dockerfile  -n frontend . '
                    sh 'docker run -p 80:80 fontend'
                }
            }
        }
    }
}