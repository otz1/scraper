pipeline {
	agent { dockerfile true }
	stages {
		stage('build') {
			steps { sh 'docker build -t scraper .' }
		}
		stage('run') {
			steps { sh 'docker run scraper' }
		}
	}
}
