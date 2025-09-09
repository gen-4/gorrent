node {
	def app

	script {
		System.setProperty("org.jenkinsci.plugins.durabletask.BourneShellScript.HEARTBEAT_CHECK_INTERVAL", "86400")
	}

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Build image') {
		echo 'Building image...'
		retry(3) {
			app = docker.build("gorrent_image:latest")
		}
		echo 'Image built'
	}

	stage('Deploying gorrent superserver...') {
		
	}
}
