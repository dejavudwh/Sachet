/*
 * @Author: dejavudwh
 * @Date: 2021-09-21 19:14:41
 * @LastEditTime: 2021-09-21 19:50:21
 */
package nsenter

/*
 #include <errno.h>
 #include <sched.h>
 #include <stdio.h>
 #include <stdlib.h>
 #include <string.h>
 #include <fcntl.h>

// Called once at the beginning of each program start
 __attribute__((constructor)) void enter_namespace(void) {
	 char *sachet_pid;
	 sachet_pid = getenv("sachet_pid");
	 if (sachet_pid) {
		fprintf(stdout, "got sachet_pid=%s\n", sachet_pid);
	} else {
		fprintf(stdout, "missing sachet_pid env skip nsenter");
		return;
	}
	char *sachet_cmd;
	sachet_cmd = getenv("sachet_cmd");
	if (sachet_cmd) {
		fprintf(stdout, "got sachet_cmd=%s\n", sachet_cmd);
	} else {
		fprintf(stdout, "missing sachet_cmd env skip nsenter");
		return;
	}
	 int i;
	 char nspath[1024];
	 char *namespaces[] = { "ipc", "uts", "net", "pid", "mnt" };

	 for (i = 0; i < 5; i++) {
		 sprintf(nspath, "/proc/%s/ns/%s", sachet_pid, namespaces[i]);
		 int fd = open(nspath, O_RDONLY);

		 if (setns(fd, 0) == -1) {
			 fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		 } else {
			 fprintf(stdout, "setns on %s namespace succeeded\n", namespaces[i]);
		 }
		 close(fd);
	 }
	 int res = system(sachet_cmd);
	 exit(0);
	 return;
 }
*/
import "C"
