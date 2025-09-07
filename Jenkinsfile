node {
	def app

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Build image') {
		echo 'Building image...'
		app = docker.build("gorrent_image/latest")
		echo 'Image built'
	}

	stage('Deploy gorrent superserver') {
		
	}
}
