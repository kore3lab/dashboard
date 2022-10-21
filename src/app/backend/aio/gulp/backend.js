// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
 * @fileoverview Gulp tasks for compiling backend application.
 */
import gulp from "gulp";
import lodash from "lodash";
import path from "path";

import goCommand from "./gocommand.js";

const basePath = path.join(__dirname, "../../");
const mainPackageName = "github.com/kore3lab/dashboard";

/**
 * Compiles backend application in development mode and places the binary in the serve
 * directory.
 */
gulp.task(
  "backend",
  gulp.series((doneFn) => {
    goCommand(
      [
        "build",
        // Install dependencies to speed up subsequent compilations.
        // "-i",
        // for debugging backend
        "-gcflags=all=-N -l",
        "-o",
        path.join(basePath, ".tmp/serve/backend"),
        mainPackageName,
      ],
      doneFn
    );
  })
);

/**
 * Compiles backend application in production mode for the current architecture and places the
 * binary in the dist directory.
 *
 * The production binary difference from development binary is only that it contains all
 * dependencies inside it and is targeted for a specific architecture.
 */
gulp.task(
  "backend:prod",
  gulp.series(() => {
    let outputBinaryPath = path.join(basePath, "dist", "amd64", "backend");
    return backendProd([[outputBinaryPath, "amd64"]]);
  })
);

/**
 * @param {!Array<!Array<string>>} outputBinaryPathsAndArchs array of
 *    (output binary path, architecture) pairs
 * @return {!Promise}
 */
function backendProd(outputBinaryPathsAndArchs) {
  let promiseFn = (path, arch) => {
    return (resolve, reject) => {
      goCommand(
        ["build", "-a", "-installsuffix", "cgo", "-o", path, mainPackageName],
        (err) => {
          if (err) {
            reject(err);
          } else {
            resolve();
          }
        },
        {
          // Disable cgo package. Required to run on scratch docker image.
          CGO_ENABLED: "0",
          GOARCH: arch,
        }
      );
    };
  };

  let goCommandPromises = outputBinaryPathsAndArchs.map(
    (pathAndArch) => new Promise(promiseFn(pathAndArch[0], pathAndArch[1]))
  );

  return Promise.all(goCommandPromises);
}
