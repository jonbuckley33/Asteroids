import subprocess

def main():
	error = True
	while error:
		error = False
		try:
			subprocess.check_call("./asteroids -hostAt=\":10034\" &> /dev/null", shell=True)
		except subprocess.CalledProcessError:
			error = True


if __name__ == "__main__":
	main()