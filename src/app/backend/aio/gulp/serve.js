import child from "child_process";
import gulp from "gulp";
import path from "path";

const basePath = path.join(__dirname, "../../");

/**
 * Currently running backend process object. Null if the backend is not running.
 *
 * @type {?child.ChildProcess}
 */
let runningBackendProcess = null;

/**
 * Kills running backend process (if any).
 */
gulp.task("kill-backend", (doneFn) => {
  if (runningBackendProcess) {
    runningBackendProcess.on("exit", () => {
      // Mark that there is no backend process running anymore.
      runningBackendProcess = null;
      // Finish the task only when the backend is actually killed.
      doneFn();
    });
    runningBackendProcess.kill();
  } else {
    doneFn();
  }
});

/**
 * Watches for changes in source files and runs Gulp tasks to rebuild them.
 */
gulp.task("watch", () => {
  gulp.watch(
    path.join(basePath, "**/*.go"),
    gulp.parallel("spawn-backend", "watch")
  );
});

/**
 * Spawns new backend application process and finishes the task immediately. Previously spawned
 * backend process is killed beforehand, if any. The frontend pages are served by BrowserSync.
 */
gulp.task(
  "spawn-backend",
  gulp.series(gulp.parallel("kill-backend", "backend"), () => {
    runningBackendProcess = child.spawn(
      path.join(basePath, ".tmp/serve/backend"),
      { stdio: "inherit", cwd: path.join(basePath, ".tmp/serve") }
    );

    runningBackendProcess.on("exit", () => {
      // Mark that there is no backend process running anymore.
      runningBackendProcess = null;
    });
  })
);

/**
 * Serves the application in development mode. Watches for changes in the source files to rebuild
 * development artifacts.
 */
gulp.task("serve", gulp.parallel("spawn-backend", "watch"));
