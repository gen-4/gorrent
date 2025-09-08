node {
	def app

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Build image') {
		echo 'Building image...'
		def  FILES_LIST = sh (script: "ls   '${workers_dir}'", returnStdout: true).trim()
		echo "FILES_LIST : ${FILES_LIST}"
		retry(3) {
			app = docker.build("gorrent_image:latest")
		}
		echo 'Image built'
	}

	stage('Deploying gorrent superserver...') {
		
	}
}
