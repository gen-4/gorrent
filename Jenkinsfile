node {
	def app

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Build image') {
		echo 'Building image...'
		ls -a
		retry(3) {
			app = docker.build("gorrent_image:latest")
		}
		echo 'Image built'
	}

	stage('Deploying gorrent superserver') {
		
	}
}
