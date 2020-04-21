import child from "child_process";
import lodash from "lodash";
import q from "q";

// Add the project's required go tools to the PATH.
const devPath = `${process.env.PATH}`;

/**
 * The environment needed for the execution of any go command.
 */
const env = lodash.merge(process.env, { PATH: devPath });

/**
 * Spawns a Go process after making sure all Go prerequisites are
 * present. Backend source files must be packaged with 'package-backend'
 * task before running
 * this command.
 *
 * @param {!Array<string>} args - Arguments of the go command.
 * @param {function(?Error=)} doneFn - Callback.
 * @param {!Object<string, string>=} [envOverride] optional environment variables overrides map.
 */
export default function goCommand(args, doneFn, envOverride) {
  spawnGoProcess(args, envOverride)
    .then(doneFn)
    .fail((error) => doneFn(error));
}

/**
 * Spawns a process.
 * Promises an error if the go command process fails.
 *
 * @param {string} processName - Process name to spawn.
 * @param {!Array<string>} args - Arguments of the go command.
 * @param {!Object<string, string>=} [envOverride] optional environment variables overrides map.
 * @return {Q.Promise} A promise object.
 */
function spawnProcess(processName, args, envOverride) {
  let deferred = q.defer();
  let envLocal = lodash.merge(env, envOverride);
  let goTask = child.spawn(processName, args, {
    env: envLocal,
    stdio: "inherit",
  });
  // Call Gulp callback on task exit. This has to be done to make Gulp dependency management
  // work.
  goTask.on("exit", function(code) {
    if (code !== 0) {
      deferred.reject(Error(`Go command error, code: ${code}`));
      return;
    }
    deferred.resolve();
  });
  return deferred.promise;
}

/**
 * Spawns go process.
 * Promises an error if the go command process fails.
 *
 * @param {!Array<string>} args - Arguments of the go command.
 * @param {!Object<string, string>=} [envOverride] optional environment variables overrides map.
 * @return {Q.Promise} A promise object.
 */
function spawnGoProcess(args, envOverride) {
  return spawnProcess("go", args, envOverride);
}
